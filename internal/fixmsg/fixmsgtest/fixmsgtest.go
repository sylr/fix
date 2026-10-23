// Package fixmsgtest validates, in tests, the messages commands build for
// each FIX version fix supports.
package fixmsgtest

import (
	"testing"
	"time"

	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/internal/fixmsg"
)

// Versions are the FIX versions fix supports.
var Versions = []string{"FIX.5.0SP2", "FIX.4.4", "FIX.4.3"}

// Session returns the configuration of a session of FIX version version,
// with fix's built-in dictionaries.
func Session(version string) *config.Session {
	s := &config.Session{
		Name:         "test",
		BeginString:  version,
		SenderCompID: "SENDER",
		TargetCompID: "TARGET",
		HeartBtInt:   30 * time.Second,
		Username:     "user",
	}

	if version == "FIX.5.0SP2" {
		s.BeginString, s.DefaultApplVerID = session.FIXT11, version
	}

	return s
}

// Validate fails t unless m, sent on s, passes the validations a
// counterparty using s's dictionaries runs on the messages it receives,
// enumerated values included, and returns it parsed.
func Validate(t testing.TB, s *config.Session, m session.Outgoing) *fix.Message {
	t.Helper()

	_, d, err := s.FIXDictionaries()
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fixmsg.Message(m, s.BeginString, d)
	if err != nil {
		t.Fatalf("%s: %v", s.Version(), err)
	}

	if err := fixmsg.Check(msg, d, s.Checks()); err != nil {
		t.Fatalf("%s: %v", s.Version(), err)
	}

	if err := d.ValidateWith(msg, s.Checks()); err != nil {
		t.Fatalf("%s: %v", s.Version(), err)
	}

	return msg
}

// Refuses fails t unless fixmsg.Check refuses m, sent on s, against s's
// dictionary, and returns the error.
func Refuses(t testing.TB, s *config.Session, m session.Outgoing) error {
	t.Helper()

	_, d, err := s.FIXDictionaries()
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fixmsg.Message(m, s.BeginString, d)
	if err != nil {
		t.Fatalf("%s: %v", s.Version(), err)
	}

	err = fixmsg.Check(msg, d, s.Checks())
	if err == nil {
		t.Fatalf("%s: message accepted, want it refused", s.Version())
	}

	return err
}
