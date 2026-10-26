package probe

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/utils"
)

var (
	optionProbeInterval time.Duration
)

// ProbeCmd probes the sessions of every context with an initiator.
var ProbeCmd = &cobra.Command{
	Use:               "probe",
	Short:             "Test fix session",
	Long:              "Test fix session: log on to the sessions of every context with an initiator, and log out, at every interval.",
	RunE:              Execute,
	PersistentPreRunE: utils.MakePersistentPreRunE(Validate),
}

func init() {
	ProbeCmd.Flags().DurationVar(&optionProbeInterval, "probe-interval", 10*time.Second, "Interval between each probing")

	prometheus.MustRegister(metricFixSessionFailuresTotal)
	prometheus.MustRegister(metricFixSessionSuccessesTotal)
	prometheus.MustRegister(metricFixSessionStatus)
}

// Validate checks the interval and loads the configuration.
func Validate(_ *cobra.Command, _ []string) error {
	if optionProbeInterval <= 0 {
		return fmt.Errorf("%w: --probe-interval must be positive", errors.ErrOptions)
	}

	return config.Load()
}

func Execute(cmd *cobra.Command, _ []string) error {
	targets, err := probeTargets()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// A second signal kills.
	context.AfterFunc(ctx, stop)

	run(ctx, optionProbeInterval, targets, connect)

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
	metricFixSessionSuccessesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
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

// target is a session to probe.
type target struct {
	context   string
	initiator *config.Initiator
	session   *config.Session
	// running is set while it is probed.
	running atomic.Bool
}

// probeTargets returns the sessions of the contexts with an initiator.
func probeTargets() ([]*target, error) {
	var targets []*target

	for _, cfgContext := range config.GetContexts() {
		if len(cfgContext.Initiator) == 0 {
			continue
		}

		initiatorConfig, err := cfgContext.GetInitiator()
		if err != nil {
			return nil, fmt.Errorf("context %s: %w", cfgContext.Name, err)
		}

		sessions, err := cfgContext.GetSessions()
		if err != nil {
			return nil, fmt.Errorf("context %s: %w", cfgContext.Name, err)
		}

		for _, s := range sessions {
			targets = append(targets, &target{context: cfgContext.Name, initiator: initiatorConfig, session: s})
		}
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("%w: no context with an initiator and sessions to probe", errors.ErrConfig)
	}

	return targets, nil
}

// run probes targets with probe at every interval, aligned on it, until ctx
// is done, then waits for the probes running. A target whose last probe
// still runs is skipped.
func run(ctx context.Context, interval time.Duration, targets []*target, probe func(context.Context, *target) error) {
	logger := config.GetLogger()

	var wg sync.WaitGroup
	defer wg.Wait()

	timer := time.NewTimer(time.Until(time.Now().Truncate(interval).Add(interval)))
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		logger.Info().Msg("Start probing")

		for _, t := range targets {
			if !t.running.CompareAndSwap(false, true) {
				logger.Warn().Str("context", t.context).Str("session", t.session.Name).Msg("Previous probe still running, skipped")
				continue
			}

			wg.Go(func() {
				defer t.running.Store(false)

				err := probe(ctx, t)
				if ctx.Err() != nil {
					return
				}
				record(t, err)
			})
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// record records the result of a probe of t.
func record(t *target, err error) {
	logger := config.GetLogger()

	if err != nil {
		logger.Error().Err(err).Str("context", t.context).Str("session", t.session.Name).Msg("Unable to connect")
		metricFixSessionFailuresTotal.WithLabelValues(t.context, t.session.Name).Inc()
		metricFixSessionStatus.WithLabelValues(t.context, t.session.Name).Set(0)
		return
	}

	logger.Debug().Str("context", t.context).Str("session", t.session.Name).Msg("Connected")
	metricFixSessionSuccessesTotal.WithLabelValues(t.context, t.session.Name).Inc()
	metricFixSessionStatus.WithLabelValues(t.context, t.session.Name).Set(1)
}

// connect logs on to t's session, then out.
func connect(ctx context.Context, t *target) error {
	config.GetLogger().Debug().Str("context", t.context).Str("session", t.session.Name).Msg("Start probing session")

	c, err := initiator.ConnectSession(ctx, t.initiator, t.session, initiator.Options{})
	if err != nil {
		return err
	}
	c.Close(true)

	return nil
}
