package cancelmassorder

import (
	"testing"

	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/fix50sp2"
)

func TestBuildMessage(t *testing.T) {
	optionOrderID = "MASS-1"

	session := &config.Session{BeginString: "FIXT.1.1", DefaultApplVerID: "FIX.5.0SP2"}
	b, err := buildMessage(session, "EURUSD", "sell")
	if err != nil {
		t.Fatal(err)
	}

	m, err := fix.Parse(b.Build())
	if err != nil {
		t.Fatal(err)
	}

	if got := string(m.MsgType()); got != string(fix50sp2.MsgTypeOrdermasscancelrequest) {
		t.Errorf("MsgType = %q", got)
	}

	for tag, want := range map[int]string{
		fix50sp2.TagClOrdID:               "MASS-1",
		fix50sp2.TagMassCancelRequestType: string(fix50sp2.MassCancelRequestTypeCancelOrdersForASecurity),
		fix50sp2.TagSide:                  string(fix50sp2.SideSell),
		fix50sp2.TagSymbol:                "EURUSD",
	} {
		if got := string(m.Get(tag)); got != want {
			t.Errorf("tag %d = %q, want %q", tag, got, want)
		}
	}
}
