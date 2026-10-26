package probe

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rs/zerolog"

	"sylr.dev/fix/v2/config"
	fixerrors "sylr.dev/fix/v2/errors"
)

func init() {
	logger := zerolog.Nop()
	config.SetLogger(&logger)
}

func TestValidateInterval(t *testing.T) {
	defer func(v time.Duration) { optionProbeInterval = v }(optionProbeInterval)

	for _, interval := range []time.Duration{0, -time.Second} {
		optionProbeInterval = interval
		if err := Validate(nil, nil); !errors.Is(err, fixerrors.ErrOptions) {
			t.Errorf("Validate() with interval %v = %v, want %v", interval, err, fixerrors.ErrOptions)
		}
	}
}

func TestProbeTargets(t *testing.T) {
	conf := config.GetConfig()
	conf.Contexts = []*config.Context{
		{Name: "one", Initiator: "venue", Sessions: []string{"a"}},
		{Name: "served", Acceptor: "server", Sessions: []string{"a"}},
		{Name: "two", Initiator: "venue", Sessions: []string{"a", "b"}},
	}
	conf.Initiators = []*config.Initiator{{Name: "venue", Addresses: []string{"127.0.0.1:1"}}}
	conf.Sessions = []*config.Session{{Name: "a"}, {Name: "b"}}

	targets, err := probeTargets()
	if err != nil {
		t.Fatal(err)
	}

	var got []string
	for _, target := range targets {
		got = append(got, target.context+"/"+target.session.Name)
	}
	if want := []string{"one/a", "two/a", "two/b"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Errorf("probeTargets() = %v, want %v", got, want)
	}

	conf.Contexts = conf.Contexts[1:2]
	if _, err := probeTargets(); !errors.Is(err, fixerrors.ErrConfig) {
		t.Errorf("probeTargets() without initiator = %v, want %v", err, fixerrors.ErrConfig)
	}
}

// TestRunCancelledWhileWaiting checks run returns when cancelled before the
// first probe.
func TestRunCancelledWhileWaiting(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	run(ctx, time.Hour, []*target{{context: "c", session: &config.Session{Name: "s"}}}, func(context.Context, *target) error {
		t.Error("probed")
		return nil
	})
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("run() took %v once cancelled", elapsed)
	}
}

// TestRun checks run probes at every interval, skips a target whose probe
// still runs, counts the results, and waits for the probes on return.
func TestRun(t *testing.T) {
	slow := &target{context: "run", session: &config.Session{Name: "slow"}}
	failing := &target{context: "run", session: &config.Session{Name: "failing"}}

	var slowProbes, failingProbes atomic.Int32
	failures := value(t, metricFixSessionFailuresTotal.WithLabelValues("run", "failing"))
	successes := value(t, metricFixSessionSuccessesTotal.WithLabelValues("run", "slow"))
	release := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	run(ctx, 20*time.Millisecond, []*target{slow, failing}, func(ctx context.Context, tg *target) error {
		if tg == failing {
			failingProbes.Add(1)
			return errors.New("refused")
		}
		if slowProbes.Add(1) == 1 {
			<-ctx.Done()
			close(release)
		}
		return nil
	})

	select {
	case <-release:
	default:
		t.Fatal("run() returned before its probe")
	}
	if n := slowProbes.Load(); n != 1 {
		t.Errorf("slow target probed %d times, want once: its probe never ended", n)
	}
	if n := failingProbes.Load(); n < 3 {
		t.Errorf("failing target probed %d times, want 3 at least", n)
	}
	// The last probe may be cancelled, and not counted.
	if got := value(t, metricFixSessionFailuresTotal.WithLabelValues("run", "failing")) - failures; got < 3 || got > float64(failingProbes.Load()) {
		t.Errorf("failures_total grew by %v, want 3 at least and %d at most", got, failingProbes.Load())
	}
	if got := value(t, metricFixSessionSuccessesTotal.WithLabelValues("run", "slow")) - successes; got != 0 {
		t.Errorf("successes_total grew by %v for a probe cancelled, want 0", got)
	}
}

func TestRecord(t *testing.T) {
	ok := &target{context: "record", session: &config.Session{Name: "s"}}

	successes := value(t, metricFixSessionSuccessesTotal.WithLabelValues("record", "s"))
	record(ok, nil)
	record(ok, nil)
	if got := value(t, metricFixSessionSuccessesTotal.WithLabelValues("record", "s")) - successes; got != 2 {
		t.Errorf("successes_total grew by %v, want 2", got)
	}
	if got := value(t, metricFixSessionStatus.WithLabelValues("record", "s")); got != 1 {
		t.Errorf("status = %v, want 1", got)
	}

	record(ok, errors.New("refused"))
	if got := value(t, metricFixSessionStatus.WithLabelValues("record", "s")); got != 0 {
		t.Errorf("status = %v, want 0", got)
	}
}

// value returns the value of counter or gauge m.
func value(t *testing.T, m prometheus.Metric) float64 {
	t.Helper()

	var d dto.Metric
	if err := m.Write(&d); err != nil {
		t.Fatal(err)
	}
	if d.Counter != nil {
		return d.Counter.GetValue()
	}

	return d.Gauge.GetValue()
}
