package probe

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/utils"
)

var (
	optionProbeInterval time.Duration

	status = make(chan result)
)

// ConfigCmd represents the buy command
var ProbeCmd = &cobra.Command{
	Use:               "probe",
	Short:             "Test fix session",
	Long:              "Test fix session.",
	RunE:              Execute,
	PersistentPreRunE: utils.MakePersistentPreRunE(Validate),
}

func init() {
	ProbeCmd.Flags().DurationVar(&optionProbeInterval, "probe-interval", 10*time.Second, "Interval between each probing")

	prometheus.MustRegister(metricFixSessionFailuresTotal)
	prometheus.MustRegister(metricFixSessionSuccessesTotal)
	prometheus.MustRegister(metricFixSessionStatus)
}

func Validate(_ *cobra.Command, _ []string) error {
	return nil
}

func Execute(_ *cobra.Command, _ []string) error {
	logger := config.GetLogger()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	time.Sleep(time.Until(time.Now().Truncate(optionProbeInterval).Add(optionProbeInterval)))
	ticker := time.NewTicker(optionProbeInterval)

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

type result struct {
	context   string
	session   string
	connected bool
}

func probe() {
	logger := config.GetLogger()

	logger.Info().Msgf("Start probing")

	contexts := config.GetContexts()

	for _, cfgContext := range contexts {
		if len(cfgContext.Initiator) == 0 {
			continue
		}

		initiatorConfig, err := cfgContext.GetInitiator()
		if err != nil {
			logger.Error().Err(err).Str("context", cfgContext.GetName()).Msg("Unable to get initiator")
			continue
		}

		sessions, err := cfgContext.GetSessions()
		if err != nil {
			logger.Error().Err(err).Str("context", cfgContext.GetName()).Msg("Unable to get sessions")
			continue
		}

		for _, session := range sessions {
			go func(cfgContext *config.Context, session *config.Session) {
				name := session.Name
				logger.Debug().Str("context", cfgContext.GetName()).Str("session", name).Msgf("Start probing session")

				c, err := initiator.ConnectSession(context.Background(), initiatorConfig, session, initiator.Options{})
				if err != nil {
					logger.Error().Err(err).Str("context", cfgContext.GetName()).Str("session", name).Msg("Unable to connect")
					status <- result{context: cfgContext.GetName(), session: name, connected: false}
					return
				}
				defer c.Close(true)

				logger.Debug().Msgf("connected")
				status <- result{context: cfgContext.GetName(), session: name, connected: true}
			}(cfgContext, session)
		}
	}
}
