package amendorder

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/cli/options/optionstest"
	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/initiator/initiatortest"
	"sylr.dev/fix/v2/internal/fixmsg"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

func init() {
	logger := zerolog.Nop()
	config.SetLogger(&logger)
}

// amendment are the flags of the amendment of a limit order.
var amendment = []string{"id", "oid", "origclordid", "cl-1", "clordid", "cl-2", "side", "buy", "type", "limit", "symbol", "EURUSD", "quantity", "10", "price", "1.08345"}

func with(flags []string, more ...string) []string {
	return append(flags[:len(flags):len(flags)], more...)
}

func set(t *testing.T, flags ...string) {
	t.Helper()

	optionstest.Set(t, AmendOrderCmd, flags...)
	if err := validateOptions(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateOptions(t *testing.T) {
	for _, tc := range []struct {
		name  string
		flags []string
		want  error
	}{
		{"limit", amendment, nil},
		{"only id", with(amendment, "origclordid", ""), nil},
		{"no identifier", with(amendment, "id", "", "origclordid", ""), errors.ErrOptions},
		{"market with price", with(amendment, "type", "market"), errors.ErrOptionsInvalidMarketPrice},
		{"limit without price", with(amendment, "price", ""), errors.ErrOptionsNoPriceGiven},
		{"zero quantity", with(amendment, "quantity", "0"), errors.ErrOptions},
		{"unknown side", with(amendment, "side", "up"), errors.ErrOptionOrderSideUnknown},
		{"unknown attribute", with(amendment, "attribute-type", "nope"), errors.ErrOptionOrderAttributeTypeUnknown},
	} {
		t.Run(tc.name, func(t *testing.T) {
			optionstest.Set(t, AmendOrderCmd, tc.flags...)
			err := validateOptions()
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("validateOptions() = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestBuildMessage(t *testing.T) {
	set(t, with(amendment, "attribute-type", "aggregated_order", "attribute-value", "x")...)

	builder, err := buildMessage(fixmsgtest.Session("FIX.5.0SP2"))
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fix.Parse(builder.Build())
	if err != nil {
		t.Fatal(err)
	}

	for tag, want := range map[int]string{
		fix50sp2.TagMsgType:            string(fix50sp2.MsgTypeOrdercancelreplacerequest),
		fix50sp2.TagOrderID:            "oid",
		fix50sp2.TagOrigClOrdID:        "cl-1",
		fix50sp2.TagClOrdID:            "cl-2",
		fix50sp2.TagSymbol:             "EURUSD",
		fix50sp2.TagOrderQty:           "10",
		fix50sp2.TagPrice:              "1.08345",
		fix50sp2.TagOrderAttributeType: string(fix50sp2.OrderAttributeTypeAggregatedOrder),
	} {
		if got := string(msg.Get(tag)); got != want {
			t.Errorf("tag %d: got %q, want %q", tag, got, want)
		}
	}
}

// TestBuildMessageMarket checks a market order is amended without a price.
func TestBuildMessageMarket(t *testing.T) {
	set(t, with(amendment, "type", "market", "price", "")...)

	builder, err := buildMessage(fixmsgtest.Session("FIX.5.0SP2"))
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fix.Parse(builder.Build())
	if err != nil {
		t.Fatal(err)
	}

	if price := msg.Get(fix50sp2.TagPrice); price != nil {
		t.Errorf("Price = %q, want none", price)
	}
}

func TestBuildMessageVersions(t *testing.T) {
	set(t, amendment...)

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			builder, err := buildMessage(session)
			if err != nil {
				t.Fatal(err)
			}

			fixmsgtest.Validate(t, session, builder)
		})
	}
}

// TestBuildMessageOrigClOrdIDRequired checks FIX 4.x, which requires
// OrigClOrdID, refuses an amendment by OrderID alone.
func TestBuildMessageOrigClOrdIDRequired(t *testing.T) {
	set(t, with(amendment, "origclordid", "")...)

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			_, err := buildMessage(fixmsgtest.Session(version))
			if version == "FIX.5.0SP2" {
				if err != nil {
					t.Fatal(err)
				}
			} else if !errors.Is(err, errors.ErrOptions) || !strings.Contains(err.Error(), "--origclordid") {
				t.Errorf("buildMessage() = %v, want --origclordid required", err)
			}
		})
	}
}

// TestBuildMessageAttributesFix5Only checks FIX 4.4 refuses
// --attribute-type, which only FIX 5.0 SP2 defines.
func TestBuildMessageAttributesFix5Only(t *testing.T) {
	set(t, with(amendment, "attribute-type", "aggregated_order")...)
	session := fixmsgtest.Session("FIX.4.4")

	builder, err := buildMessage(session)
	if err != nil {
		t.Fatal(err)
	}

	if err := fixmsgtest.Refuses(t, session, builder); !strings.Contains(err.Error(), "NoOrderAttributes") {
		t.Errorf("err = %v, want it to name NoOrderAttributes", err)
	}
}

func report(beginString, clOrdID string, ordStatus fix50sp2.OrdStatus) *fix50sp2.ExecutionReportBuilder {
	return fix50sp2.BuildExecutionReport(beginString).
		SetOrderID("oid").
		SetClOrdID(clOrdID).
		SetOrigClOrdID("cl-1").
		SetExecID("1").
		SetExecType(fix50sp2.ExecType(ordStatus)).
		SetOrdStatus(ordStatus).
		SetSide(fix50sp2.SideBuy).
		SetLeavesQty(1).
		SetCumQty(0).
		SetAvgPx(0)
}

func cancelReject(beginString, clOrdID string) *fix50sp2.OrderCancelRejectBuilder {
	return fix50sp2.BuildOrderCancelReject(beginString).
		SetOrderID("oid").
		SetClOrdID(clOrdID).
		SetOrigClOrdID("cl-1").
		SetOrdStatus(fix50sp2.OrdStatusNew).
		SetCxlRejResponseTo(fix50sp2.CxlRejResponseToOrderCancelReplaceRequest).
		SetText("too late")
}

func TestHandle(t *testing.T) {
	for _, tc := range []struct {
		name     string
		messages []session.Outgoing
		done     bool
		want     error
	}{
		{"replaced", []session.Outgoing{report("FIXT.1.1", "cl-2", fix50sp2.OrdStatusReplaced)}, true, nil},
		{"pending replace", []session.Outgoing{report("FIXT.1.1", "cl-2", fix50sp2.OrdStatusPendingReplace)}, false, nil},
		{"pending then new", []session.Outgoing{
			report("FIXT.1.1", "cl-2", fix50sp2.OrdStatusPendingReplace),
			report("FIXT.1.1", "cl-2", fix50sp2.OrdStatusNew),
		}, true, nil},
		{"original order", []session.Outgoing{report("FIXT.1.1", "cl-1", fix50sp2.OrdStatusNew)}, false, initiator.ErrSkip},
		{"rejected", []session.Outgoing{report("FIXT.1.1", "cl-2", fix50sp2.OrdStatusRejected)}, true, errors.ErrFixOrderRejected},
		{"cancel reject", []session.Outgoing{cancelReject("FIXT.1.1", "cl-2")}, true, errors.ErrFixOrderRejected},
		{"cancel reject of another", []session.Outgoing{cancelReject("FIXT.1.1", "cl-3")}, false, initiator.ErrSkip},
	} {
		t.Run(tc.name, func(t *testing.T) {
			set(t, amendment...)
			h := &handler{out: io.Discard, logger: config.GetLogger()}

			var done bool
			var err error
			for _, m := range tc.messages {
				msg, perr := fixmsg.Message(m, "FIXT.1.1", nil)
				if perr != nil {
					t.Fatal(perr)
				}
				if done, err = h.handle(msg); done || (err != nil && !errors.Is(err, initiator.ErrSkip)) {
					break
				}
			}

			if done != tc.done || (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("handle() = %v, %v, want %v, %v", done, err, tc.done, tc.want)
			}
		})
	}
}

// venue answers each OrderCancelReplaceRequest with a pending replace
// report then a replaced one, an OrderCancelReject if reject is set, or the
// report of another order if other is.
type venue struct {
	session.BaseHandler
	reject, other bool
}

func (v venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	clOrdID := string(m.Get(fix50sp2.TagClOrdID))
	beginString := s.ID().BeginString

	// FIX 4.4 reports a replacement with ExecType Replaced.
	replaced := report(beginString, clOrdID, fix50sp2.OrdStatusNew).SetExecType(fix50sp2.ExecTypeReplaced)
	responses := []session.Outgoing{report(beginString, clOrdID, fix50sp2.OrdStatusPendingReplace), replaced}
	switch {
	case v.reject:
		responses = []session.Outgoing{cancelReject(beginString, clOrdID)}
	case v.other:
		responses = []session.Outgoing{report(beginString, "cl-9", fix50sp2.OrdStatusNew)}
	}

	for _, r := range responses {
		if _, err := s.Send(ctx, r); err != nil {
			return err
		}
	}

	return nil
}

func TestExecute(t *testing.T) {
	for _, tc := range []struct {
		name  string
		venue venue
		want  error
	}{
		{"replaced", venue{}, nil},
		{"rejected", venue{reject: true}, errors.ErrFixOrderRejected},
	} {
		t.Run(tc.name, func(t *testing.T) {
			initiatortest.Serve(t, tc.venue, "FIX.4.4")
			optionstest.Set(t, AmendOrderCmd, amendment...)

			out, err := optionstest.Run(context.Background(), AmendOrderCmd)
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Fatalf("Run() = %v, want %v", err, tc.want)
			}
			if !strings.Contains(out, "cl-2") {
				t.Errorf("output lacks the response:\n%s", out)
			}
		})
	}
}

// TestExecuteOther checks the report of another order does not defeat the
// timeout of the first report.
func TestExecuteOther(t *testing.T) {
	initiatortest.Serve(t, venue{other: true}, "FIX.4.4")
	optionstest.Set(t, AmendOrderCmd, with(amendment, "exec-reports", "0", "exec-reports-timeout", "300ms")...)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := optionstest.Run(ctx, AmendOrderCmd); !errors.Is(err, errors.ErrResponseTimeout) {
		t.Errorf("Run() = %v, want a timeout", err)
	}
}
