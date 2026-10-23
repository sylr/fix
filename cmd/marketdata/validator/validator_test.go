//go:build validator
// +build validator

package marketdatavalidator

import (
	"testing"

	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/internal/fixmsg"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

// The validator deliberately subscribes to MDEntryType "101", a non-standard
// value (Trade History): the request is checked for its structure only, as
// full validation would refuse the value.
func TestBuildMessageVersions(t *testing.T) {
	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			b, err := buildMessage(&initiator.Client{Session: session}, "EURUSD")
			if err != nil {
				t.Fatal(err)
			}

			_, d, err := session.FIXDictionaries()
			if err != nil {
				t.Fatal(err)
			}

			msg, err := fixmsg.Message(b, session.BeginString, d)
			if err != nil {
				t.Fatal(err)
			}

			if err := fixmsg.Check(msg, d, session.Checks()); err != nil {
				t.Fatalf("%s: %v", session.Version(), err)
			}
		})
	}
}
