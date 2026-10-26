package statusorder

import (
	"context"
	"io"
	"strings"
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

var statusFlags = []string{"id", "ORD-1", "clordid", "CL-1", "side", "buy", "symbol", "EURUSD"}

func with(flags []string, more ...string) []string {
	return append(flags[:len(flags):len(flags)], more...)
}

func set(t *testing.T, flags ...string) {
	t.Helper()

	optionstest.Set(t, StatusOrderCmd, flags...)
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
		{"all", statusFlags, nil},
		{"only id", with(statusFlags, "clordid", ""), nil},
		{"only clordid", with(statusFlags, "id", ""), nil},
		{"deprecated origclordid", []string{"origclordid", "CL-1", "side", "buy", "symbol", "EURUSD"}, nil},
		{"no identifier", with(statusFlags, "id", "", "clordid", ""), errors.ErrOptions},
		{"empty symbol", with(statusFlags, "symbol", ""), errors.ErrOptionsNoSymbolGiven},
		{"unknown side", with(statusFlags, "side", "up"), errors.ErrOptionOrderSideUnknown},
	} {
		t.Run(tc.name, func(t *testing.T) {
			optionstest.Set(t, StatusOrderCmd, tc.flags...)
			err := validateOptions()
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("validateOptions() = %v, want %v", err, tc.want)
			}
		})
	}
}

// TestBuildMessageVersions checks ClOrdID names the order, OrigClOrdID is
// never sent, and OrdStatusReqID is where the version defines it.
func TestBuildMessageVersions(t *testing.T) {
	set(t, statusFlags...)

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			b, err := buildMessage(session)
			if err != nil {
				t.Fatal(err)
			}

			msg := fixmsgtest.Validate(t, session, b)

			for tag, want := range map[int]string{
				fix50sp2.TagMsgType: string(fix50sp2.MsgTypeOrderstatusrequest),
				fix50sp2.TagClOrdID: "CL-1",
				fix50sp2.TagOrderID: "ORD-1",
				fix50sp2.TagSide:    string(fix50sp2.SideBuy),
				fix50sp2.TagSymbol:  "EURUSD",
			} {
				if got := string(msg.Get(tag)); got != want {
					t.Errorf("tag %d = %q, want %q", tag, got, want)
				}
			}
			if got := msg.Get(fix50sp2.TagOrigClOrdID); got != nil {
				t.Errorf("OrigClOrdID = %q, want none", got)
			}

			reqID := string(msg.Get(fix50sp2.TagOrdStatusReqID))
			if want := version != "FIX.4.3"; (reqID == statusReqID) != want {
				t.Errorf("OrdStatusReqID = %q, want it sent: %v", reqID, want)
			}
		})
	}
}

// TestBuildMessageClOrdIDRequired checks FIX 4.x, which requires ClOrdID,
// refuses a request by OrderID alone.
func TestBuildMessageClOrdIDRequired(t *testing.T) {
	set(t, with(statusFlags, "clordid", "")...)

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)
			b, err := buildMessage(session)
			if version == "FIX.5.0SP2" {
				if err != nil {
					t.Fatal(err)
				}
				fixmsgtest.Validate(t, session, b)
			} else if !errors.Is(err, errors.ErrOptions) || !strings.Contains(err.Error(), "--clordid") {
				t.Errorf("buildMessage() = %v, want --clordid required", err)
			}
		})
	}
}

func report(beginString, clOrdID, reqID string, ordStatus fix50sp2.OrdStatus) *fix50sp2.ExecutionReportBuilder {
	b := fix50sp2.BuildExecutionReport(beginString).
		SetOrderID("ORD-1").
		SetClOrdID(clOrdID).
		SetExecID("1").
		SetExecType(fix50sp2.ExecTypeOrderStatus).
		SetOrdStatus(ordStatus).
		SetSide(fix50sp2.SideBuy).
		SetLeavesQty(1).
		SetCumQty(0).
		SetAvgPx(0)
	if reqID != "" {
		b.SetOrdStatusReqID(reqID)
	}

	return b
}

func TestHandle(t *testing.T) {
	set(t, statusFlags...)

	unknown := report("FIXT.1.1", "CL-1", statusReqID, fix50sp2.OrdStatusRejected).SetOrdRejReason(fix50sp2.OrdRejReasonUnknownOrder)
	for _, tc := range []struct {
		name          string
		byStatusReqID bool
		message       session.Outgoing
		done          bool
		want          error
	}{
		{"by request id", true, report("FIXT.1.1", "CL-1", statusReqID, fix50sp2.OrdStatusNew), true, nil},
		{"another request", true, report("FIXT.1.1", "CL-1", "other", fix50sp2.OrdStatusNew), false, initiator.ErrSkip},
		{"request id not echoed", true, report("FIXT.1.1", "CL-1", "", fix50sp2.OrdStatusFilled), true, nil},
		{"by clordid", false, report("FIXT.1.1", "CL-1", "", fix50sp2.OrdStatusNew), true, nil},
		{"another order", false, report("FIXT.1.1", "CL-2", "", fix50sp2.OrdStatusNew), false, initiator.ErrSkip},
		{"rejected order", true, report("FIXT.1.1", "CL-1", statusReqID, fix50sp2.OrdStatusRejected), true, nil},
		{"unknown order", true, unknown, true, errors.ErrFixOrderRejected},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := fixmsg.Message(tc.message, "FIXT.1.1", nil)
			if err != nil {
				t.Fatal(err)
			}

			h := &handler{out: io.Discard, logger: config.GetLogger(), byStatusReqID: tc.byStatusReqID}
			done, err := h.handle(msg)
			if done != tc.done || (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("handle() = %v, %v, want %v, %v", done, err, tc.done, tc.want)
			}
		})
	}
}

// venue answers each OrderStatusRequest with a report of another request,
// then its own.
type venue struct {
	session.BaseHandler
}

func (venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	clOrdID, reqID := string(m.Get(fix50sp2.TagClOrdID)), string(m.Get(fix50sp2.TagOrdStatusReqID))
	for _, r := range []session.Outgoing{
		report(s.ID().BeginString, clOrdID, "other", fix50sp2.OrdStatusFilled),
		report(s.ID().BeginString, clOrdID, reqID, fix50sp2.OrdStatusNew),
	} {
		if _, err := s.Send(ctx, r); err != nil {
			return err
		}
	}

	return nil
}

func TestExecute(t *testing.T) {
	for _, version := range []string{"FIX.5.0SP2", "FIX.4.4"} {
		t.Run(version, func(t *testing.T) {
			initiatortest.Serve(t, venue{}, version)
			optionstest.Set(t, StatusOrderCmd, statusFlags...)

			out, err := optionstest.Run(context.Background(), StatusOrderCmd)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out, statusReqID) || strings.Contains(out, "other") {
				t.Errorf("output is not the report of the request:\n%s", out)
			}
		})
	}
}

func TestExecuteTimeout(t *testing.T) {
	initiatortest.Serve(t, session.BaseHandler{}, "FIX.5.0SP2")
	optionstest.Set(t, StatusOrderCmd, with(statusFlags, "exec-reports-timeout", "200ms")...)

	if _, err := optionstest.Run(context.Background(), StatusOrderCmd); !errors.Is(err, errors.ErrResponseTimeout) {
		t.Errorf("Run() = %v, want a timeout", err)
	}
}
