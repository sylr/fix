package cancelorder

import (
	"testing"

	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/fix50sp2"
)

func TestBuildMessage(t *testing.T) {
	optionQuoteID = "Q-1"
	optionOrderSymbols = []string{"EURUSD", "GBPUSD"}

	session := &config.Session{BeginString: "FIXT.1.1", DefaultApplVerID: "FIX.5.0SP2"}
	b, err := buildMessage(session)
	if err != nil {
		t.Fatal(err)
	}

	m, err := fix.Parse(b.Build())
	if err != nil {
		t.Fatal(err)
	}

	if got := string(m.MsgType()); got != string(fix50sp2.MsgTypeQuotecancel) {
		t.Errorf("MsgType = %q", got)
	}

	for tag, want := range map[int]string{
		fix50sp2.TagQuoteID:         "Q-1",
		fix50sp2.TagQuoteMsgID:      "msg_Q-1",
		fix50sp2.TagQuoteCancelType: string(fix50sp2.QuoteCancelTypeCancelForOneOrMoreSecurities),
		fix50sp2.TagNoQuoteEntries:  "2",
	} {
		if got := string(m.Get(tag)); got != want {
			t.Errorf("tag %d = %q, want %q", tag, got, want)
		}
	}

	var symbols []string
	for _, f := range m.Fields() {
		if f.Tag == fix50sp2.TagSymbol {
			symbols = append(symbols, string(f.Value))
		}
	}
	if len(symbols) != 2 || symbols[0] != "EURUSD" || symbols[1] != "GBPUSD" {
		t.Errorf("symbols = %v", symbols)
	}
}
