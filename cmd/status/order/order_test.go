package status_order

import (
	"testing"

	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/fix50sp2"
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
