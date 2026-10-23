package config

import (
	"testing"
	"time"

	"github.com/sylr/fix2/fix/session/schedule"
)

func TestSessionSchedule(t *testing.T) {
	daily, err := (&Session{Schedule: &Schedule{
		Start: "17:00:00", End: "16:55:00", Days: []string{"sun", "Monday"}, Location: "America/New_York",
	}}).schedule()
	if err != nil {
		t.Fatal(err)
	}
	d, ok := daily.(schedule.Daily)
	if !ok || d.Start != schedule.Clock(17, 0, 0) || len(d.Days) != 2 || d.Days[1] != time.Monday || d.Location.String() != "America/New_York" {
		t.Errorf("daily schedule = %#v", daily)
	}

	weekly, err := (&Session{Schedule: &Schedule{
		Start: "07:00:00", End: "18:00:00", StartDay: "monday", EndDay: "fri",
	}}).schedule()
	if err != nil {
		t.Fatal(err)
	}
	w, ok := weekly.(schedule.Weekly)
	if !ok || w.StartDay != time.Monday || w.EndDay != time.Friday || w.End != schedule.Clock(18, 0, 0) {
		t.Errorf("weekly schedule = %#v", weekly)
	}

	for name, s := range map[string]*Schedule{
		"bad clock":          {Start: "7am", End: "18:00:00"},
		"bad day":            {Start: "07:00:00", End: "18:00:00", Days: []string{"someday"}},
		"weekly without end": {Start: "07:00:00", End: "18:00:00", StartDay: "mon"},
		"weekly with days":   {Start: "07:00:00", End: "18:00:00", StartDay: "mon", EndDay: "fri", Days: []string{"mon"}},
		"bad location":       {Start: "07:00:00", End: "18:00:00", Location: "Mars/Olympus"},
	} {
		if _, err := (&Session{Schedule: s}).schedule(); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}

func TestSessionFIXConfig(t *testing.T) {
	s := &Session{
		BeginString:      "FIXT.1.1",
		SenderCompID:     "A",
		TargetCompID:     "B",
		DefaultApplVerID: "FIX.5.0SP2",
		HeartBtInt:       5 * time.Second,
		Tolerations:      []string{"quickfix", "finerTimestamps"},
		SkipChecks:       []string{"groupOrder"},
	}
	if err := s.validate(); err != nil {
		t.Fatal(err)
	}

	cfg, err := s.FIXConfig(nil, nil, nil)
	if err == nil {
		t.Fatal("FIXConfig() without a handler: no error")
	}
	if cfg.Dictionaries.DefaultApplVerID != "9" || cfg.Dictionaries.Apps["9"] == nil || cfg.Dictionaries.Session == nil {
		t.Errorf("Dictionaries = %+v", cfg.Dictionaries)
	}
	if !cfg.TolerateNoHeartbeatAwaitingTestRequest || !cfg.TolerateFinerTimestamps {
		t.Errorf("Tolerations = %+v", cfg.Tolerations)
	}
	if cfg.LogonTimeout != defaultTimeout {
		t.Errorf("LogonTimeout = %v, want %v", cfg.LogonTimeout, defaultTimeout)
	}

	if err := (&Session{Tolerations: []string{"everything"}}).validate(); err == nil {
		t.Error("unknown toleration: no error")
	}
	if err := (&Session{SkipChecks: []string{"nothing"}}).validate(); err == nil {
		t.Error("unknown check: no error")
	}
}

func TestConfigValidateStoreType(t *testing.T) {
	for name, conf := range map[string]fixConfig{
		"initiator": {Initiators: []*Initiator{{Name: "i", Store: Store{Type: "sqlit3"}}}},
		"acceptor":  {Acceptors: []*Acceptor{{Name: "a", Store: Store{Type: "postgres"}}}},
	} {
		if err := conf.Validate(); err == nil {
			t.Errorf("%s: unknown store type: no error", name)
		}
	}

	ok := fixConfig{
		Initiators: []*Initiator{{Name: "i", Store: Store{Type: "sqlite3"}}},
		Acceptors:  []*Acceptor{{Name: "a"}, {Name: "b", Store: Store{Type: "memory"}}},
	}
	if err := ok.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}
