package neworder

import (
	"testing"

	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/fix50sp2"
)

func testSession() *config.Session {
	return &config.Session{BeginString: "FIXT.1.1", DefaultApplVerID: "FIX.5.0SP2"}
}

func checkFields(t *testing.T, msg fix.Message, want map[int]string) {
	t.Helper()

	for tag, value := range want {
		if got := string(msg.Get(tag)); got != value {
			t.Errorf("tag %d: got %q, want %q", tag, got, value)
		}
	}
}

func TestBuildMessage(t *testing.T) {
	optionOrderID = "order-1"
	optionOrderSide = "buy"
	optionOrderType = "limit"
	optionOrderExpiry = "day"
	optionOrderSymbol = "EURUSD"
	optionOrderQuantity = 10
	optionOrderPrice = 1.5

	builder, err := buildMessage("FIXT.1.1", testSession())
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fix.Parse(builder.Build())
	if err != nil {
		t.Fatal(err)
	}

	checkFields(t, msg, map[int]string{
		fix50sp2.TagMsgType:  string(fix50sp2.MsgTypeNewordersingle),
		fix50sp2.TagClOrdID:  "order-1",
		fix50sp2.TagSymbol:   "EURUSD",
		fix50sp2.TagOrderQty: "10.00",
		fix50sp2.TagPrice:    "1.50",
	})
}

func TestBuildMessageMarketHasNoPrice(t *testing.T) {
	optionOrderID = "order-2"
	optionOrderSide = "sell"
	optionOrderType = "market"
	optionOrderExpiry = "day"
	optionOrderSymbol = "EURUSD"
	optionOrderQuantity = 1
	optionOrderPrice = 0

	builder, err := buildMessage("FIXT.1.1", testSession())
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fix.Parse(builder.Build())
	if err != nil {
		t.Fatal(err)
	}

	if msg.Get(fix50sp2.TagPrice) != nil {
		t.Errorf("market order must not carry a price")
	}
}

func TestBuildCancelReplaceMessage(t *testing.T) {
	optionOrderExpiry = "day"
	optionOrderSymbol = "EURUSD"
	optionUpdateOrderQuantity = 5
	optionUpdateOrderPrice = 0.25

	raw := fix.NewBuilder("FIXT.1.1", string(fix50sp2.MsgTypeExecutionreport)).
		SetString(fix50sp2.TagOrderID, "oid").
		SetString(fix50sp2.TagClOrdID, "cl-1").
		SetString(fix50sp2.TagOrdType, string(fix50sp2.OrdTypeLimit)).
		SetString(fix50sp2.TagSide, string(fix50sp2.SideBuy)).
		SetString(fix50sp2.TagOrderQty, "10.00").
		SetString(fix50sp2.TagPrice, "1.50").
		Build()
	er, err := fix.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}

	builder, err := buildCancelReplaceMessage("FIXT.1.1", testSession(), &er)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fix.Parse(builder.Build())
	if err != nil {
		t.Fatal(err)
	}

	checkFields(t, msg, map[int]string{
		fix50sp2.TagMsgType:     string(fix50sp2.MsgTypeOrdercancelreplacerequest),
		fix50sp2.TagOrderID:     "oid",
		fix50sp2.TagOrigClOrdID: "cl-1",
		fix50sp2.TagOrderQty:    "15.00",
		fix50sp2.TagPrice:       "1.75",
	})

	if _, err := buildCancelReplaceMessage("FIXT.1.1", testSession(), nil); err == nil {
		t.Errorf("expected error without execution report")
	}
}
