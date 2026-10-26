package neworder

import (
	"context"
	"io"
	"strings"
	"sync"
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

// limitOrder are the flags of a limit order.
var limitOrder = []string{"id", "order-1", "side", "buy", "type", "limit", "symbol", "EURUSD", "quantity", "10", "price", "1.08345"}

func with(flags []string, more ...string) []string {
	return append(flags[:len(flags):len(flags)], more...)
}

// set sets the flags of NewOrderCmd for the test, and validates them.
func set(t *testing.T, flags ...string) {
	t.Helper()

	optionstest.Set(t, NewOrderCmd, flags...)
	if err := validateOptions(); err != nil {
		t.Fatal(err)
	}
}

func checkFields(t *testing.T, msg *fix.Message, want map[int]string) {
	t.Helper()

	for tag, value := range want {
		if got := string(msg.Get(tag)); got != value {
			t.Errorf("tag %d: got %q, want %q", tag, got, value)
		}
	}
}

func TestValidateOptions(t *testing.T) {
	for _, tc := range []struct {
		name  string
		flags []string
		want  error
	}{
		{"limit", limitOrder, nil},
		{"market", []string{"side", "sell", "type", "market", "symbol", "EURUSD", "quantity", "1"}, nil},
		{"market with price", []string{"side", "sell", "type", "market", "symbol", "EURUSD", "quantity", "1", "price", "1"}, errors.ErrOptionsInvalidMarketPrice},
		{"limit without price", []string{"side", "sell", "type", "limit", "symbol", "EURUSD", "quantity", "1"}, errors.ErrOptionsNoPriceGiven},
		{"zero price", with(limitOrder, "price", "0"), nil},
		{"negative price", with(limitOrder, "price", "-0.25"), nil},
		{"zero quantity", with(limitOrder, "quantity", "0"), errors.ErrOptions},
		{"negative quantity", with(limitOrder, "quantity", "-1"), errors.ErrOptions},
		{"unknown side", with(limitOrder, "side", "up"), errors.ErrOptionOrderSideUnknown},
		{"unknown type", with(limitOrder, "type", "nope"), errors.ErrOptionOrderTypeUnknown},
		{"unknown expiry", with(limitOrder, "expiry", "never"), errors.ErrOptionOrderExpiryUnknown},
		{"unknown origination", with(limitOrder, "origination", "nope"), errors.ErrOptionOrderOriginationUnknown},
		{"empty symbol", with(limitOrder, "symbol", ""), errors.ErrOptionsNoSymbolGiven},
		{"market price update", []string{"side", "sell", "type", "market", "symbol", "EURUSD", "quantity", "1", "update-period", "1s", "update-order-price", "1"}, errors.ErrOptions},
		{"party sub id types alone", with(limitOrder, "party-sub-id-types", "firm"), errors.ErrOptions},
	} {
		t.Run(tc.name, func(t *testing.T) {
			optionstest.Set(t, NewOrderCmd, tc.flags...)
			err := validateOptions()
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("validateOptions() = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestDecimalFlags(t *testing.T) {
	for _, v := range []string{"NaN", "Inf", "1,5", "abc"} {
		if err := NewOrderCmd.Flags().Set("price", v); err == nil {
			t.Errorf("--price %s accepted", v)
		}
	}
	optionstest.Set(t, NewOrderCmd)
}

func TestValidateGeneratesClOrdID(t *testing.T) {
	set(t, with(limitOrder, "id", "")...)
	if len(optionOrderID) == 0 {
		t.Error("ClOrdID not generated")
	}
}

func TestBuildMessage(t *testing.T) {
	set(t, limitOrder...)

	builder, err := buildMessage(fixmsgtest.Session("FIX.5.0SP2"))
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fix.Parse(builder.Build())
	if err != nil {
		t.Fatal(err)
	}

	checkFields(t, &msg, map[int]string{
		fix50sp2.TagMsgType:  string(fix50sp2.MsgTypeNewordersingle),
		fix50sp2.TagClOrdID:  "order-1",
		fix50sp2.TagSymbol:   "EURUSD",
		fix50sp2.TagOrderQty: "10",
		fix50sp2.TagPrice:    "1.08345",
	})
}

func TestBuildMessageMarketHasNoPrice(t *testing.T) {
	set(t, "side", "sell", "type", "market", "symbol", "EURUSD", "quantity", "1.5")

	builder, err := buildMessage(fixmsgtest.Session("FIX.5.0SP2"))
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
	checkFields(t, &msg, map[int]string{fix50sp2.TagOrderQty: "1.5"})
}

func TestBuildMessageVersions(t *testing.T) {
	set(t, limitOrder...)

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

func TestBuildReplaceMessageVersions(t *testing.T) {
	set(t, limitOrder...)

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			builder, err := buildReplaceMessage(session, "order-1", "oid", replacement{clOrdID: "cl-2", quantity: optionOrderQuantity.Value(), price: optionOrderPrice.Value()})
			if err != nil {
				t.Fatal(err)
			}

			msg := fixmsgtest.Validate(t, session, builder)
			checkFields(t, msg, map[int]string{
				fix50sp2.TagClOrdID:     "cl-2",
				fix50sp2.TagOrigClOrdID: "order-1",
				fix50sp2.TagOrderID:     "oid",
				fix50sp2.TagOrderQty:    "10",
				fix50sp2.TagPrice:       "1.08345",
			})
		})
	}
}

// TestBuildMessageFix5Only checks FIX 4.4 refuses what only FIX 5.0 SP2
// defines: --origination, --attribute-type and --party-role-qualifier.
func TestBuildMessageFix5Only(t *testing.T) {
	session := fixmsgtest.Session("FIX.4.4")
	party := []string{"party-id", "P1", "party-id-source", "proprietary", "party-role", "executing_firm"}

	for _, tc := range []struct {
		name  string
		flags []string
		want  string
	}{
		{"origination", []string{"origination", "customer"}, "OrderOrigination"},
		{"attribute", []string{"attribute-type", "aggregated_order", "attribute-value", "x"}, "NoOrderAttributes"},
		{"party-role-qualifier", with(party, "party-role-qualifier", "agency"), "PartyRoleQualifier"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			set(t, with(limitOrder, tc.flags...)...)

			builder, err := buildMessage(session)
			if err != nil {
				t.Fatal(err)
			}

			if err := fixmsgtest.Refuses(t, session, builder); !strings.Contains(err.Error(), tc.want) {
				t.Errorf("err = %v, want it to name %s", err, tc.want)
			}
		})
	}
}

// report builds an ExecutionReport of order clOrdID, of status ordStatus.
func report(beginString, clOrdID string, ordStatus fix50sp2.OrdStatus) *fix50sp2.ExecutionReportBuilder {
	return fix50sp2.BuildExecutionReport(beginString).
		SetOrderID("OID-" + clOrdID).
		SetClOrdID(clOrdID).
		SetExecID("1").
		SetExecType(fix50sp2.ExecType(ordStatus)).
		SetOrdStatus(ordStatus).
		SetSide(fix50sp2.SideBuy).
		SetLeavesQty(1).
		SetCumQty(0).
		SetAvgPx(0)
}

func parse(t *testing.T, m session.Outgoing) *fix.Message {
	t.Helper()

	msg, err := fixmsg.Message(m, "FIXT.1.1", nil)
	if err != nil {
		t.Fatal(err)
	}

	return msg
}

func cancelReject(clOrdID string) *fix50sp2.OrderCancelRejectBuilder {
	return fix50sp2.BuildOrderCancelReject("FIXT.1.1").
		SetOrderID("OID").
		SetClOrdID(clOrdID).
		SetOrigClOrdID("order-1").
		SetOrdStatus(fix50sp2.OrdStatusNew).
		SetCxlRejResponseTo(fix50sp2.CxlRejResponseToOrderCancelReplaceRequest).
		SetText("too late")
}

func TestHandle(t *testing.T) {
	for _, tc := range []struct {
		name     string
		flags    []string
		messages []session.Outgoing
		done     bool
		want     error
	}{
		{"new", nil, []session.Outgoing{report("FIXT.1.1", "order-1", fix50sp2.OrdStatusNew)}, true, nil},
		{"pending new", nil, []session.Outgoing{report("FIXT.1.1", "order-1", fix50sp2.OrdStatusPendingNew)}, false, nil},
		{"pending new then new", nil, []session.Outgoing{
			report("FIXT.1.1", "order-1", fix50sp2.OrdStatusPendingNew),
			report("FIXT.1.1", "order-1", fix50sp2.OrdStatusNew),
		}, true, nil},
		{"other order", nil, []session.Outgoing{report("FIXT.1.1", "order-2", fix50sp2.OrdStatusNew)}, false, initiator.ErrSkip},
		{"rejected", nil, []session.Outgoing{report("FIXT.1.1", "order-1", fix50sp2.OrdStatusRejected)}, true, errors.ErrFixOrderRejected},
		{"canceled", nil, []session.Outgoing{report("FIXT.1.1", "order-1", fix50sp2.OrdStatusCanceled)}, true, errors.ErrFixOrderCanceled},
		{"cancel reject", nil, []session.Outgoing{cancelReject("cl-2")}, true, errors.ErrFixOrderRejected},
		{"two reports", []string{"exec-reports", "2"}, []session.Outgoing{report("FIXT.1.1", "order-1", fix50sp2.OrdStatusNew)}, false, nil},
		{"final state", []string{"exec-reports", "0", "stop-on-final-state", "true"}, []session.Outgoing{
			report("FIXT.1.1", "order-1", fix50sp2.OrdStatusNew),
			report("FIXT.1.1", "order-1", fix50sp2.OrdStatusFilled),
		}, true, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			set(t, with(limitOrder, tc.flags...)...)
			h := newHandler(io.Discard, nil, fixmsgtest.Session("FIX.5.0SP2"), nil)

			var done bool
			var err error
			for _, m := range tc.messages {
				if done, err = h.handle(parse(t, m)); done || (err != nil && !errors.Is(err, initiator.ErrSkip)) {
					break
				}
			}

			if done != tc.done || (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("handle() = %v, %v, want %v, %v", done, err, tc.done, tc.want)
			}
		})
	}
}

// TestUpdate checks updates wait for the order to be acknowledged and for
// the last replacement to be confirmed, and replace the order last
// confirmed.
func TestUpdate(t *testing.T) {
	set(t, with(limitOrder, "exec-reports", "0", "update-period", "1s", "update-order-quantity", "1", "update-order-price", "-0.00001")...)

	var sent []*fix.Message
	h := newHandler(io.Discard, nil, fixmsgtest.Session("FIX.5.0SP2"), func(m session.Outgoing) error {
		sent = append(sent, parse(t, m))
		return nil
	})

	handle := func(m session.Outgoing) {
		t.Helper()
		if done, err := h.handle(parse(t, m)); done || err != nil {
			t.Fatalf("handle() = %v, %v", done, err)
		}
	}
	update := func(want int) {
		t.Helper()
		if err := h.update(); err != nil {
			t.Fatal(err)
		}
		if len(sent) != want {
			t.Fatalf("%d replacements sent, want %d", len(sent), want)
		}
	}

	update(0)
	handle(report("FIXT.1.1", "order-1", fix50sp2.OrdStatusNew))
	update(1)
	checkFields(t, sent[0], map[int]string{
		fix50sp2.TagOrigClOrdID: "order-1",
		fix50sp2.TagOrderID:     "OID-order-1",
		fix50sp2.TagOrderQty:    "11",
		fix50sp2.TagPrice:       "1.08344",
	})
	first := string(sent[0].Get(fix50sp2.TagClOrdID))

	// Pending: not confirmed yet.
	update(1)
	handle(report("FIXT.1.1", first, fix50sp2.OrdStatusPendingReplace))
	update(1)

	handle(report("FIXT.1.1", first, fix50sp2.OrdStatusNew))
	update(2)
	checkFields(t, sent[1], map[int]string{
		fix50sp2.TagOrigClOrdID: first,
		fix50sp2.TagOrderQty:    "12",
		fix50sp2.TagPrice:       "1.08343",
	})
}

// venue answers each NewOrderSingle with the ExecutionReports of the
// statuses given, those of another order if other is set, or with err.
type venue struct {
	session.BaseHandler
	statuses []fix50sp2.OrdStatus
	other    bool
	err      error

	mu       sync.Mutex
	received []*fix.Message
}

func (v *venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	v.mu.Lock()
	v.received = append(v.received, m.Clone())
	v.mu.Unlock()

	if v.err != nil {
		return v.err
	}

	clOrdID := string(m.Get(fix50sp2.TagClOrdID))
	if v.other {
		clOrdID = "order-9"
	}
	for _, status := range v.statuses {
		if _, err := s.Send(ctx, report(s.ID().BeginString, clOrdID, status)); err != nil {
			return err
		}
	}

	return nil
}

func TestExecute(t *testing.T) {
	for _, tc := range []struct {
		name  string
		venue *venue
		flags []string
		want  error
	}{
		{"new", &venue{statuses: []fix50sp2.OrdStatus{fix50sp2.OrdStatusPendingNew, fix50sp2.OrdStatusNew}}, nil, nil},
		{"rejected", &venue{statuses: []fix50sp2.OrdStatus{fix50sp2.OrdStatusRejected}}, nil, errors.ErrFixOrderRejected},
		{"business reject", &venue{err: &session.BusinessRejectError{Reason: session.BusinessRejectOther, Text: "no"}}, nil, errors.ErrFixRejected},
		{"timeout", &venue{}, []string{"exec-reports-timeout", "200ms"}, errors.ErrResponseTimeout},
		{"another order", &venue{statuses: []fix50sp2.OrdStatus{fix50sp2.OrdStatusNew}, other: true}, []string{"exec-reports-timeout", "200ms"}, errors.ErrResponseTimeout},
		{"another order, first timeout", &venue{statuses: []fix50sp2.OrdStatus{fix50sp2.OrdStatusNew}, other: true}, []string{"exec-reports", "0", "exec-reports-timeout", "200ms"}, errors.ErrResponseTimeout},
	} {
		t.Run(tc.name, func(t *testing.T) {
			initiatortest.Serve(t, tc.venue, "FIX.5.0SP2")
			optionstest.Set(t, NewOrderCmd, with(limitOrder, tc.flags...)...)

			// A timeout defeated would wait until ctx is done.
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			out, err := optionstest.Run(ctx, NewOrderCmd)
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Fatalf("Run() = %v, want %v", err, tc.want)
			}

			tc.venue.mu.Lock()
			defer tc.venue.mu.Unlock()
			if len(tc.venue.received) != 1 {
				t.Fatalf("%d messages received, want 1", len(tc.venue.received))
			}
			checkFields(t, tc.venue.received[0], map[int]string{fix50sp2.TagPrice: "1.08345", fix50sp2.TagOrderQty: "10"})
			if tc.want == nil && !strings.Contains(out, "order-1") {
				t.Errorf("output lacks the report:\n%s", out)
			}
		})
	}
}

// TestExecuteUpdates checks --update-period replaces the order, the first
// update not waiting for a response beyond the acknowledgement.
func TestExecuteUpdates(t *testing.T) {
	v := &venue{statuses: []fix50sp2.OrdStatus{fix50sp2.OrdStatusNew}}
	initiatortest.Serve(t, v, "FIX.4.4")
	optionstest.Set(t, NewOrderCmd, with(limitOrder, "exec-reports", "3", "update-period", "50ms", "update-order-quantity", "1")...)

	if _, err := optionstest.Run(context.Background(), NewOrderCmd); err != nil {
		t.Fatal(err)
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	if len(v.received) < 3 {
		t.Fatalf("%d messages received, want 3", len(v.received))
	}
	checkFields(t, v.received[1], map[int]string{
		fix50sp2.TagMsgType:     string(fix50sp2.MsgTypeOrdercancelreplacerequest),
		fix50sp2.TagOrigClOrdID: "order-1",
		fix50sp2.TagOrderQty:    "11",
	})
	checkFields(t, v.received[2], map[int]string{
		fix50sp2.TagOrigClOrdID: string(v.received[1].Get(fix50sp2.TagClOrdID)),
		fix50sp2.TagOrderQty:    "12",
	})
}

// TestValidateBeforeConnect checks Validate refuses, from the session's
// dictionary, what its FIX version can't carry.
func TestValidateBeforeConnect(t *testing.T) {
	initiatortest.Serve(t, &venue{}, "FIX.4.4")
	optionstest.Set(t, NewOrderCmd, with(limitOrder, "origination", "customer")...)

	if err := Validate(NewOrderCmd, nil); err == nil || !strings.Contains(err.Error(), "OrderOrigination") {
		t.Errorf("Validate() = %v, want OrderOrigination refused", err)
	}
}
