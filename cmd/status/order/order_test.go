package status_order

import (
	"testing"

	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

func TestBuildMessage(t *testing.T) {
	optionOrderID = "ORD-1"
	optionOrigClOrdID = "ORIG-1"
	optionSide = "buy"
	optionSymbol = "EURUSD"

	session := &config.Session{BeginString: "FIXT.1.1", DefaultApplVerID: "FIX.5.0SP2"}
	b, err := buildMessage(session)
	if err != nil {
		t.Fatal(err)
	}

	m, err := fix.Parse(b.Build())
	if err != nil {
		t.Fatal(err)
	}

	if got := string(m.MsgType()); got != string(fix50sp2.MsgTypeOrderstatusrequest) {
		t.Errorf("MsgType = %q", got)
	}

	for tag, want := range map[int]string{
		fix50sp2.TagOrderID:     "ORD-1",
		fix50sp2.TagOrigClOrdID: "ORIG-1",
		fix50sp2.TagSide:        string(fix50sp2.SideBuy),
		fix50sp2.TagSymbol:      "EURUSD",
	} {
		if got := string(m.Get(tag)); got != want {
			t.Errorf("tag %d = %q, want %q", tag, got, want)
		}
	}

	if len(m.Get(fix50sp2.TagClOrdID)) == 0 {
		t.Error("ClOrdID not set")
	}
}

func TestBuildMessageVersions(t *testing.T) {
	optionOrderID = "ORD-1"
	optionOrigClOrdID = "ORIG-1"
	optionSide = "buy"
	optionSymbol = "EURUSD"

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			b, err := buildMessage(session)
			if err != nil {
				t.Fatal(err)
			}

			msg := fixmsgtest.Validate(t, session, b)

			// FIX 5.0 SP2 names the order with OrigClOrdID, FIX 4.x with ClOrdID.
			tag := fix50sp2.TagOrigClOrdID
			if version != "FIX.5.0SP2" {
				tag = fix50sp2.TagClOrdID
			}
			if got := string(msg.Get(tag)); got != "ORIG-1" {
				t.Errorf("tag %d = %q, want ORIG-1", tag, got)
			}
		})
	}
}

// TestBuildMessageSkipUnknownTags checks a FIX 4.x session skipping the
// unknown tags check still names the order with ClOrdID, its dictionary
// being whole.
func TestBuildMessageSkipUnknownTags(t *testing.T) {
	optionOrderID = "ORD-1"
	optionOrigClOrdID = "ORIG-1"
	optionSide = "buy"
	optionSymbol = "EURUSD"

	session := fixmsgtest.Session("FIX.4.4")
	session.SkipChecks = []string{"unknownTags"}

	b, err := buildMessage(session)
	if err != nil {
		t.Fatal(err)
	}

	msg := fixmsgtest.Validate(t, session, b)
	if got := string(msg.Get(fix50sp2.TagClOrdID)); got != "ORIG-1" {
		t.Errorf("ClOrdID = %q, want ORIG-1", got)
	}
	if got := msg.Get(fix50sp2.TagOrigClOrdID); got != nil {
		t.Errorf("OrigClOrdID = %q, want none", got)
	}
}
