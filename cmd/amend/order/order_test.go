package amendorder

import (
	"testing"

	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

func TestBuildMessage(t *testing.T) {
	optionOrderID = "oid"
	optionOrigClOrdID = "cl-1"
	optionClOrdID = "cl-2"
	optionOrderSide = "buy"
	optionOrderType = "limit"
	optionOrderExpiry = "day"
	optionOrderSymbol = "EURUSD"
	optionOrderQuantity = 10
	optionOrderPrice = 1.5

	session := &config.Session{BeginString: "FIXT.1.1", DefaultApplVerID: "FIX.5.0SP2"}
	builder, err := buildMessage("FIXT.1.1", session)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fix.Parse(builder.Build())
	if err != nil {
		t.Fatal(err)
	}

	for tag, want := range map[int]string{
		fix50sp2.TagMsgType:     string(fix50sp2.MsgTypeOrdercancelreplacerequest),
		fix50sp2.TagOrderID:     "oid",
		fix50sp2.TagOrigClOrdID: "cl-1",
		fix50sp2.TagClOrdID:     "cl-2",
		fix50sp2.TagSymbol:      "EURUSD",
		fix50sp2.TagOrderQty:    "10.00",
		fix50sp2.TagPrice:       "1.50",
	} {
		if got := string(msg.Get(tag)); got != want {
			t.Errorf("tag %d: got %q, want %q", tag, got, want)
		}
	}
}

func TestBuildMessageVersions(t *testing.T) {
	optionOrderID = "oid"
	optionOrigClOrdID = "cl-1"
	optionClOrdID = "cl-2"
	optionOrderSide = "buy"
	optionOrderType = "limit"
	optionOrderExpiry = "day"
	optionOrderSymbol = "EURUSD"
	optionOrderQuantity = 10
	optionOrderPrice = 1.5

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			builder, err := buildMessage(session.BeginString, session)
			if err != nil {
				t.Fatal(err)
			}

			fixmsgtest.Validate(t, session, builder)
		})
	}
}
