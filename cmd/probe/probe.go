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

var (
	metricFixSessionErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "fix",
			Subsystem: "session",
			Name:      "errors_total",
			Help:      "Number fix session errors",
		},
		[]string{},
	)
	metricFixSessionStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "fix",
			Subsystem: "session",
			Name:      "status",
			Help:      "Status of the fix session ",
		},
		[]string{},
	)
)

var (
	status = make(chan bool)
)

func init() {
	prometheus.MustRegister(metricFixSessionErrors)
	prometheus.MustRegister(metricFixSessionStatus)

	metricFixSessionErrors.WithLabelValues().Add(0)
}

func connect() {
	options := config.GetOptions()
	logger := config.GetLogger()

	context, err := config.GetCurrentContext()
	if err != nil {
		status <- false
		return
	}

	settings, err := context.ToQuickFixInitiatorSettings()
	if err != nil {
		status <- false
		return
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
		status <- false
		return
	}

	// Start session
	if err = init.Start(); err != nil {
		status <- false
		return
	}

	defer func() {
		app.Stop()
		init.Stop()
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
		quickfix.UnregisterSession(app.SessionID)
		status <- false

	case _, ok := <-app.Connected:
		if !ok {
			status <- false
		} else {
			logger.Debug().Msgf("connected")
			quickfix.UnregisterSession(app.SessionID)
			status <- true
		}
	}
}

func Execute(_ *cobra.Command, _ []string) error {
	logger := config.GetLogger()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go connect()

LOOP:
	for {
		select {
		case signal := <-interrupt:
			logger.Debug().Msgf("Received signal: %s", signal)
			break LOOP

		case <-time.After(30 * time.Second):
			go connect()

		case ok := <-status:
			if ok {
				metricFixSessionStatus.WithLabelValues().Set(1.0)
			} else {
				metricFixSessionErrors.WithLabelValues().Inc()
				metricFixSessionStatus.WithLabelValues().Set(0.0)
			}
		}
	}

	return nil
}
