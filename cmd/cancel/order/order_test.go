package cancelorder

import (
	"testing"

	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
)

var testSession = &config.Session{BeginString: "FIXT.1.1", DefaultApplVerID: "FIX.5.0SP2"}

func TestBuildMessage(t *testing.T) {
	optionClientOrderID = "CL-1"
	optionOrigClientOrderID = "ORIG-1"
	optionOrderID = "ORD-1"
	optionOrderSide = "buy"
	optionOrderSymbol = "EURUSD"

	b, err := buildMessage(testSession)
	if err != nil {
		t.Fatal(err)
	}

	m, err := fix.Parse(b.Build())
	if err != nil {
		t.Fatal(err)
	}

	if got := string(m.MsgType()); got != string(fix50sp2.MsgTypeOrdercancelrequest) {
		t.Errorf("MsgType = %q", got)
	}

	for tag, want := range map[int]string{
		fix50sp2.TagClOrdID:     "CL-1",
		fix50sp2.TagOrigClOrdID: "ORIG-1",
		fix50sp2.TagOrderID:     "ORD-1",
		fix50sp2.TagSide:        string(fix50sp2.SideBuy),
		fix50sp2.TagSymbol:      "EURUSD",
	} {
		if got := string(m.Get(tag)); got != want {
			t.Errorf("tag %d = %q, want %q", tag, got, want)
		}
	}

	if len(m.Get(fix50sp2.TagTransactTime)) == 0 {
		t.Error("TransactTime not set")
	}
}

func TestBuildMessageUnsupportedVersion(t *testing.T) {
	optionOrderSide = "buy"

	_, err := buildMessage(&config.Session{BeginString: "FIX.4.4"})
	if !errors.Is(err, errors.FixVersionNotImplemented) {
		t.Errorf("err = %v, want %v", err, errors.FixVersionNotImplemented)
	}
}
