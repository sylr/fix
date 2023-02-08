package probe

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/quickfixgo/quickfix"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/initiator/application"
	"sylr.dev/fix/pkg/utils"
)

// Metrics
var (
	metricFixSessionFailuresTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "fix",
			Subsystem: "session",
			Name:      "failures_total",
			Help:      "Total number of session errors",
		},
		[]string{"context", "session"},
	)
	metricFixSessionSuccessesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "fix",
			Subsystem: "session",
			Name:      "successes_total",
			Help:      "Total number of session successes",
		},
		[]string{"context", "session"},
	)
	metricFixSessionStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "fix",
			Subsystem: "session",
			Name:      "status",
			Help:      "Status of the fix session",
		},
		[]string{"context", "session"},
	)
)

// ConfigCmd represents the buy command
var ProbeCmd = &cobra.Command{
	Use:   "probe",
	Short: "Test fix session",
	Long:  "Test fix session.",
	RunE:  Execute,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := utils.ValidateRequiredFlags(cmd); err != nil {
			return err
		}

		if err := initiator.ValidateOptions(cmd, args); err != nil {
			return err
		}

		if cmd.HasParent() {
			parent := cmd.Parent()
			if parent.PersistentPreRunE != nil {
				return parent.PersistentPreRunE(parent, args)
			}
		}

		return nil
	},
}

func init() {
	prometheus.MustRegister(metricFixSessionFailuresTotal)
	prometheus.MustRegister(metricFixSessionSuccessesTotal)
	prometheus.MustRegister(metricFixSessionStatus)
}

type result struct {
	context   string
	session   string
	connected bool
}

var (
	status = make(chan result)
)

func probe() {
	options := config.GetOptions()
	logger := config.GetLogger()

	logger.Info().Msgf("Start probing")

	contexts := config.GetContexts()

	for i, context := range contexts {
		if len(context.Initiator) == 0 {
			continue
		}

		for j, session := range context.Sessions {
			go func(context *config.Context, i int, session string) {
				logger.Debug().Str("context", context.GetName()).Str("session", session).Msgf("Start probing session")

				// Make a copy of the context which has only one session.
				contextSingleSession := *context
				contextSingleSession.Sessions = nil
				contextSingleSession.Sessions = context.Sessions[i : i+1]

				logger.Debug().Str("context", context.GetName()).Str("session", session).Strs("sessions", contextSingleSession.Sessions).Msgf("Copied context")

				settings, err := contextSingleSession.ToQuickFixInitiatorSettings()
				if err != nil {
					panic(err)
				}

				app := application.NewInitiator()
				app.Settings = settings
				app.Logger = logger

				var quickfixLogger *zerolog.Logger
				if options.QuickFixLogging {
					quickfixLogger = logger
				}

				init, err := initiator.Initiate(app, settings, quickfixLogger)
				if err != nil {
					logger.Error().Err(err).Str("context", context.GetName()).Str("session", session).Msg("Unable to initiate initiator")
					status <- result{context: context.GetName(), session: session, connected: false}
					return
				}

				// Start session
				if err = init.Start(); err != nil {
					logger.Error().Err(err).Str("context", context.GetName()).Str("session", session).Msg("Unable to start initiator")
					status <- result{context: context.GetName(), session: session, connected: false}
					return
				}

				defer func() {
					app.Stop()
					init.Stop()
					_ = quickfix.UnregisterSession(app.SessionID)
				}()

				// Choose right timeout cli option > config > default value (5s)
				var timeout time.Duration
				if options.Timeout != time.Duration(0) {
					timeout = options.Timeout
				} else {
					timeout = 5 * time.Second
				}

				// Wait for session connection
				select {
				case <-time.After(timeout):
					status <- result{context: context.GetName(), session: session, connected: false}

				case _, ok := <-app.Connected:
					if !ok {
						status <- result{context: context.GetName(), session: session, connected: false}
					} else {
						logger.Debug().Msgf("connected")
						status <- result{context: context.GetName(), session: session, connected: true}
					}
				}
			}(contexts[i], j, session)
		}
	}
}

func Execute(_ *cobra.Command, _ []string) error {
	logger := config.GetLogger()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	time.Sleep(time.Until(time.Now().Truncate(10 * time.Second).Add(10 * time.Second)))
	ticker := time.NewTicker(10 * time.Second)

LOOP:
	for {
		select {
		case signal := <-interrupt:
			logger.Debug().Msgf("Received signal: %s", signal)
			break LOOP

		case <-ticker.C:
			go probe()

		case result := <-status:
			if result.connected {
				metricFixSessionSuccessesTotal.WithLabelValues(result.context, result.session).Inc()
				metricFixSessionStatus.WithLabelValues(result.context, result.session).Set(1.0)
			} else {
				metricFixSessionFailuresTotal.WithLabelValues(result.context, result.session).Inc()
				metricFixSessionStatus.WithLabelValues(result.context, result.session).Set(0.0)
			}
		}
	}

	return nil
}
