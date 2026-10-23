package neworder

import (
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/internal/fixmsg"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
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

func TestBuildMessageVersions(t *testing.T) {
	optionOrderID = "order-1"
	optionOrderSide = "buy"
	optionOrderType = "limit"
	optionOrderExpiry = "day"
	optionOrderSymbol = "EURUSD"
	optionOrderQuantity = 10
	optionOrderPrice = 1.5
	optionOrderOrigination = ""

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

func TestBuildCancelReplaceMessageVersions(t *testing.T) {
	optionOrderExpiry = "day"
	optionOrderSymbol = "EURUSD"
	optionUpdateOrderQuantity = 5
	optionUpdateOrderPrice = 0.25

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			report := fix50sp2.BuildExecutionReport(session.BeginString).
				SetOrderID("oid").
				SetClOrdID("cl-1").
				SetOrdType(fix50sp2.OrdTypeLimit).
				SetSide(fix50sp2.SideBuy)
			report.SetString(fix50sp2.TagOrderQty, "10.00")
			report.SetString(fix50sp2.TagPrice, "1.50")

			er, err := fixmsg.Message(report, session.BeginString, nil)
			if err != nil {
				t.Fatal(err)
			}

			builder, err := buildCancelReplaceMessage(session.BeginString, session, er)
			if err != nil {
				t.Fatal(err)
			}

			fixmsgtest.Validate(t, session, builder)
		})
	}
}

// setSliceFlag sets a repeatable flag of NewOrderCmd for the rest of the test.
func setSliceFlag(t *testing.T, name string, values ...string) {
	t.Helper()

	f := NewOrderCmd.Flags().Lookup(name)
	slice := f.Value.(pflag.SliceValue)
	if err := slice.Replace(values); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { slice.Replace(nil) })
}

// TestBuildMessageFix5Only checks FIX 4.4 refuses what only FIX 5.0 SP2
// defines: --origination, --attribute-type and --party-role-qualifier.
func TestBuildMessageFix5Only(t *testing.T) {
	optionOrderID = "order-1"
	optionOrderSide = "buy"
	optionOrderType = "limit"
	optionOrderExpiry = "day"
	optionOrderSymbol = "EURUSD"
	optionOrderQuantity = 10
	optionOrderPrice = 1.5

	session := fixmsgtest.Session("FIX.4.4")

	t.Run("origination", func(t *testing.T) {
		optionOrderOrigination = "customer"
		t.Cleanup(func() { optionOrderOrigination = "" })

		builder, err := buildMessage(session.BeginString, session)
		if err != nil {
			t.Fatal(err)
		}

		if err := fixmsgtest.Refuses(t, session, builder); !strings.Contains(err.Error(), "OrderOrigination") {
			t.Errorf("err = %v, want it to name OrderOrigination", err)
		}
	})

	t.Run("attribute", func(t *testing.T) {
		setSliceFlag(t, "attribute-type", "aggregated_order")
		setSliceFlag(t, "attribute-value", "x")

		builder, err := buildMessage(session.BeginString, session)
		if err != nil {
			t.Fatal(err)
		}

		if err := fixmsgtest.Refuses(t, session, builder); !strings.Contains(err.Error(), "NoOrderAttributes") {
			t.Errorf("err = %v, want it to name NoOrderAttributes", err)
		}
	})

	t.Run("party-role-qualifier", func(t *testing.T) {
		setSliceFlag(t, "party-id", "P1")
		setSliceFlag(t, "party-id-source", "proprietary_custom_code")
		setSliceFlag(t, "party-role", "executing_firm")
		setSliceFlag(t, "party-role-qualifier", "agency")

		builder, err := buildMessage(session.BeginString, session)
		if err != nil {
			t.Fatal(err)
		}

		if err := fixmsgtest.Refuses(t, session, builder); !strings.Contains(err.Error(), "PartyRoleQualifier") {
			t.Errorf("err = %v, want it to name PartyRoleQualifier", err)
		}
	})
}
