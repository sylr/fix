package cancelmass

import (
	"context"
	"io"
	"testing"

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

func TestValidateOptions(t *testing.T) {
	for _, tc := range []struct {
		name  string
		flags []string
		want  error
	}{
		{"ok", []string{"sides", "buy,sell", "symbols", "EURUSD"}, nil},
		{"unknown side", []string{"sides", "buy,up", "symbols", "EURUSD"}, errors.ErrOptionOrderSideUnknown},
		{"no side", []string{"symbols", "EURUSD"}, errors.ErrOptionOrderSideUnknown},
		{"no symbol", []string{"sides", "buy"}, errors.ErrOptionsNoSymbolGiven},
		{"empty symbol", []string{"sides", "buy", "symbols", "EURUSD,"}, errors.ErrOptionsNoSymbolGiven},
	} {
		t.Run(tc.name, func(t *testing.T) {
			optionstest.Set(t, MassCancelOrderCmd, tc.flags...)
			err := validateOptions()
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("validateOptions() = %v, want %v", err, tc.want)
			}
		})
	}
}

// TestRequests checks --id is used as is for one request, as a prefix for
// several, and a uuid is generated per request without it.
func TestRequests(t *testing.T) {
	optionstest.Set(t, MassCancelOrderCmd, "id", "MASS", "sides", "buy", "symbols", "EURUSD")
	if r := requests(); len(r) != 1 || r[0].clOrdID != "MASS" {
		t.Errorf("requests() = %v", r)
	}

	optionstest.Set(t, MassCancelOrderCmd, "id", "MASS", "sides", "buy,sell", "symbols", "EURUSD,GBPUSD")
	r := requests()
	want := []request{{"MASS-1", "EURUSD", "buy"}, {"MASS-2", "EURUSD", "sell"}, {"MASS-3", "GBPUSD", "buy"}, {"MASS-4", "GBPUSD", "sell"}}
	if len(r) != len(want) {
		t.Fatalf("requests() = %v, want %v", r, want)
	}
	for i := range want {
		if r[i] != want[i] {
			t.Errorf("requests()[%d] = %v, want %v", i, r[i], want[i])
		}
	}

	optionstest.Set(t, MassCancelOrderCmd, "sides", "buy,sell", "symbols", "EURUSD")
	if r := requests(); len(r) != 2 || r[0].clOrdID == "" || r[0].clOrdID == r[1].clOrdID {
		t.Errorf("requests() = %v, want distinct ClOrdIDs", r)
	}
}

func TestBuildMessage(t *testing.T) {
	optionstest.Set(t, MassCancelOrderCmd)

	b, err := buildMessage(fixmsgtest.Session("FIX.5.0SP2"), request{"MASS-1", "EURUSD", "sell"})
	if err != nil {
		t.Fatal(err)
	}

	m, err := fix.Parse(b.Build())
	if err != nil {
		t.Fatal(err)
	}

	for tag, want := range map[int]string{
		fix50sp2.TagMsgType:               string(fix50sp2.MsgTypeOrdermasscancelrequest),
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

func TestBuildMessageVersions(t *testing.T) {
	optionstest.Set(t, MassCancelOrderCmd)

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			b, err := buildMessage(session, request{"MASS-1", "EURUSD", "sell"})
			if err != nil {
				t.Fatal(err)
			}

			fixmsgtest.Validate(t, session, b)
		})
	}
}

func massReport(beginString, clOrdID string, response fix50sp2.MassCancelResponse) *fix50sp2.OrderMassCancelReportBuilder {
	return fix50sp2.BuildOrderMassCancelReport(beginString).
		SetClOrdID(clOrdID).
		SetOrderID("MASS-ORD").
		SetMassActionReportID("REPORT-" + clOrdID).
		SetMassCancelRequestType(fix50sp2.MassCancelRequestTypeCancelOrdersForASecurity).
		SetMassCancelResponse(response)
}

func executionReport(beginString string) *fix50sp2.ExecutionReportBuilder {
	return fix50sp2.BuildExecutionReport(beginString).
		SetOrderID("ORD-1").
		SetClOrdID("ORD-1").
		SetExecID("1").
		SetExecType(fix50sp2.ExecTypeCanceled).
		SetOrdStatus(fix50sp2.OrdStatusCanceled).
		SetSide(fix50sp2.SideBuy).
		SetLeavesQty(0).
		SetCumQty(0).
		SetAvgPx(0)
}

func TestHandle(t *testing.T) {
	accepted := fix50sp2.MassCancelResponseCancelOrdersForASecurity
	for _, tc := range []struct {
		name     string
		messages []session.Outgoing
		done     bool
		want     error
	}{
		{"one of two", []session.Outgoing{massReport("FIXT.1.1", "A", accepted)}, false, nil},
		{"both", []session.Outgoing{massReport("FIXT.1.1", "A", accepted), massReport("FIXT.1.1", "B", accepted)}, true, nil},
		{"same twice", []session.Outgoing{massReport("FIXT.1.1", "A", accepted), massReport("FIXT.1.1", "A", accepted)}, false, initiator.ErrSkip},
		{"other request and execution report", []session.Outgoing{massReport("FIXT.1.1", "C", accepted), executionReport("FIXT.1.1")}, false, initiator.ErrSkip},
		{"rejected", []session.Outgoing{massReport("FIXT.1.1", "B", fix50sp2.MassCancelResponseCancelRequestRejected)}, true, errors.ErrFixOrderRejected},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := &handler{out: io.Discard, logger: config.GetLogger(), pending: []string{"A", "B"}}

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

// venue answers each OrderMassCancelRequest with an ExecutionReport, then
// an OrderMassCancelReport.
type venue struct {
	session.BaseHandler
}

func (venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	for _, r := range []session.Outgoing{
		executionReport(s.ID().BeginString),
		massReport(s.ID().BeginString, string(m.Get(fix50sp2.TagClOrdID)), fix50sp2.MassCancelResponseCancelOrdersForASecurity),
	} {
		if _, err := s.Send(ctx, r); err != nil {
			return err
		}
	}

	return nil
}

func TestExecute(t *testing.T) {
	initiatortest.Serve(t, venue{}, "FIX.5.0SP2")
	optionstest.Set(t, MassCancelOrderCmd, "id", "MASS", "sides", "buy,sell", "symbols", "EURUSD")

	if _, err := optionstest.Run(context.Background(), MassCancelOrderCmd); err != nil {
		t.Fatal(err)
	}
}
