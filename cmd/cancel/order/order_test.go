package cancelorder

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

var cancelFlags = []string{"clordid", "CL-1", "origclordid", "ORIG-1", "id", "ORD-1", "side", "buy", "symbol", "EURUSD"}

func with(flags []string, more ...string) []string {
	return append(flags[:len(flags):len(flags)], more...)
}

func set(t *testing.T, flags ...string) {
	t.Helper()

	optionstest.Set(t, CancelOrderCmd, flags...)
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
		{"all", cancelFlags, nil},
		{"only id", with(cancelFlags, "origclordid", ""), nil},
		{"only origclordid", with(cancelFlags, "id", ""), nil},
		{"no identifier", with(cancelFlags, "id", "", "origclordid", ""), errors.ErrOptions},
		{"unknown side", with(cancelFlags, "side", "up"), errors.ErrOptionOrderSideUnknown},
	} {
		t.Run(tc.name, func(t *testing.T) {
			optionstest.Set(t, CancelOrderCmd, tc.flags...)
			err := validateOptions()
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("validateOptions() = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestValidateGeneratesClOrdID(t *testing.T) {
	set(t, with(cancelFlags, "clordid", "")...)
	if len(optionClientOrderID) == 0 {
		t.Error("ClOrdID not generated")
	}
}

func TestBuildMessage(t *testing.T) {
	set(t, cancelFlags...)

	b, err := buildMessage(fixmsgtest.Session("FIX.5.0SP2"))
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
	set(t, cancelFlags...)

	_, err := buildMessage(&config.Session{BeginString: "FIX.4.2"})
	if !errors.Is(err, errors.ErrFixVersionNotImplemented) {
		t.Errorf("err = %v, want %v", err, errors.ErrFixVersionNotImplemented)
	}
}

func TestBuildMessageVersions(t *testing.T) {
	set(t, cancelFlags...)

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			b, err := buildMessage(session)
			if err != nil {
				t.Fatal(err)
			}

			fixmsgtest.Validate(t, session, b)
		})
	}
}

// TestBuildMessageOrigClOrdIDRequired checks FIX 4.x, which requires
// OrigClOrdID, refuses a cancel request by OrderID alone.
func TestBuildMessageOrigClOrdIDRequired(t *testing.T) {
	set(t, with(cancelFlags, "origclordid", "")...)

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

func report(beginString, clOrdID string, ordStatus fix50sp2.OrdStatus) *fix50sp2.ExecutionReportBuilder {
	return fix50sp2.BuildExecutionReport(beginString).
		SetOrderID("ORD-1").
		SetClOrdID(clOrdID).
		SetOrigClOrdID("ORIG-1").
		SetExecID("1").
		SetExecType(fix50sp2.ExecType(ordStatus)).
		SetOrdStatus(ordStatus).
		SetSide(fix50sp2.SideBuy).
		SetLeavesQty(0).
		SetCumQty(0).
		SetAvgPx(0)
}

func cancelReject(beginString, clOrdID string) *fix50sp2.OrderCancelRejectBuilder {
	return fix50sp2.BuildOrderCancelReject(beginString).
		SetOrderID("ORD-1").
		SetClOrdID(clOrdID).
		SetOrigClOrdID("ORIG-1").
		SetOrdStatus(fix50sp2.OrdStatusFilled).
		SetCxlRejResponseTo(fix50sp2.CxlRejResponseToOrderCancelRequest).
		SetText("too late to cancel")
}

func TestHandle(t *testing.T) {
	for _, tc := range []struct {
		name     string
		flags    []string
		messages []session.Outgoing
		done     bool
		want     error
	}{
		{"canceled", nil, []session.Outgoing{report("FIXT.1.1", "CL-1", fix50sp2.OrdStatusCanceled)}, true, nil},
		{"pending cancel", nil, []session.Outgoing{report("FIXT.1.1", "CL-1", fix50sp2.OrdStatusPendingCancel)}, false, nil},
		{"pending then canceled", nil, []session.Outgoing{
			report("FIXT.1.1", "CL-1", fix50sp2.OrdStatusPendingCancel),
			report("FIXT.1.1", "CL-1", fix50sp2.OrdStatusCanceled),
		}, true, nil},
		{"canceled by the order's id", nil, []session.Outgoing{report("FIXT.1.1", "ORIG-1", fix50sp2.OrdStatusCanceled)}, true, nil},
		{"rejected", nil, []session.Outgoing{report("FIXT.1.1", "CL-1", fix50sp2.OrdStatusRejected)}, true, errors.ErrFixOrderRejected},
		{"cancel reject", nil, []session.Outgoing{cancelReject("FIXT.1.1", "CL-1")}, true, errors.ErrFixOrderRejected},
		{"another order", []string{"origclordid", "", "id", "ORD-2"}, []session.Outgoing{
			report("FIXT.1.1", "CL-2", fix50sp2.OrdStatusCanceled),
			cancelReject("FIXT.1.1", "CL-2"),
		}, false, initiator.ErrSkip},
	} {
		t.Run(tc.name, func(t *testing.T) {
			set(t, with(cancelFlags, tc.flags...)...)
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

// venue answers each OrderCancelRequest with the responses respond returns.
type venue struct {
	session.BaseHandler
	respond func(beginString, clOrdID string) []session.Outgoing
}

func (v venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	for _, r := range v.respond(s.ID().BeginString, string(m.Get(fix50sp2.TagClOrdID))) {
		if _, err := s.Send(ctx, r); err != nil {
			return err
		}
	}

	return nil
}

func TestExecute(t *testing.T) {
	for _, tc := range []struct {
		name    string
		respond func(beginString, clOrdID string) []session.Outgoing
		flags   []string
		want    error
	}{
		{"pending then canceled", func(b, id string) []session.Outgoing {
			return []session.Outgoing{report(b, id, fix50sp2.OrdStatusPendingCancel), report(b, id, fix50sp2.OrdStatusCanceled)}
		}, nil, nil},
		{"rejected", func(b, id string) []session.Outgoing {
			return []session.Outgoing{cancelReject(b, id)}
		}, nil, errors.ErrFixOrderRejected},
		{"pending only", func(b, id string) []session.Outgoing {
			return []session.Outgoing{report(b, id, fix50sp2.OrdStatusPendingCancel)}
		}, []string{"exec-reports-timeout", "300ms"}, errors.ErrResponseTimeout},
	} {
		t.Run(tc.name, func(t *testing.T) {
			initiatortest.Serve(t, venue{respond: tc.respond}, "FIX.5.0SP2")
			optionstest.Set(t, CancelOrderCmd, with(cancelFlags, tc.flags...)...)

			start := time.Now()
			_, err := optionstest.Run(context.Background(), CancelOrderCmd)
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Fatalf("Run() = %v, want %v", err, tc.want)
			}
			if time.Since(start) > 3*time.Second {
				t.Errorf("Run() took %s", time.Since(start))
			}
		})
	}
}
