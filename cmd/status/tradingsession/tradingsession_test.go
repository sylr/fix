package statustradingsession

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
	fixerrors "sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator/initiatortest"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

// setOptions sets the options of the command as its flags would, restoring
// their defaults at the end of t.
func setOptions(t *testing.T, reqID, tradingSessionID, subType string) {
	t.Helper()

	optionTradSesReqID, optionTradingSessionID, optionSubType = reqID, tradingSessionID, subType
	t.Cleanup(func() {
		optionTradSesReqID, optionTradingSessionID, optionSubType = "", "", "snapshot"
	})
}

func TestValidate(t *testing.T) {
	setOptions(t, "", "", "snapshot")
	if err := Validate(StatusTradingSessionCmd, nil); err != nil {
		t.Fatal(err)
	}
	if optionTradSesReqID == "" {
		t.Error("Validate() did not generate the TradSesReqID")
	}

	setOptions(t, "", "", "never")
	if err := Validate(StatusTradingSessionCmd, nil); !errors.Is(err, fixerrors.ErrOptions) {
		t.Errorf("Validate(--subscription-type never) = %v, want an options error", err)
	}
}

func TestBuildMessageVersions(t *testing.T) {
	setOptions(t, "REQ-1", "1", "snapshot")

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			s := fixmsgtest.Session(version)

			b, err := buildMessage(s)
			if err != nil {
				t.Fatal(err)
			}

			m := fixmsgtest.Validate(t, s, b)
			if id := string(m.Get(fix50sp2.TagTradSesReqID)); id != "REQ-1" {
				t.Errorf("TradSesReqID = %q, want REQ-1", id)
			}
			if id := string(m.Get(fix50sp2.TagTradingSessionID)); id != "1" {
				t.Errorf("TradingSessionID = %q, want 1", id)
			}
		})
	}
}

// venue answers each TradingSessionStatusRequest with a status of another
// request, then statuses of it, Open or status, or fails with err, then
// logs out if logout is set.
type venue struct {
	session.BaseHandler
	statuses int
	status   fix50sp2.TradSesStatus
	err      error
	logout   bool
}

func (v venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	if v.err != nil {
		return v.err
	}

	reqID := string(m.Get(fix50sp2.TagTradSesReqID))
	for i := range v.statuses + 1 {
		id, st := reqID, cmp.Or(v.status, fix50sp2.TradSesStatusOpen)
		if i == 0 {
			id, st = "another", fix50sp2.TradSesStatusOpen
		}
		status := fix50sp2.BuildTradingSessionStatus(s.ID().BeginString).
			SetTradSesReqID(id).
			SetTradingSessionID("1").
			SetTradSesStatus(st)
		if st == fix50sp2.TradSesStatusRequestRejected {
			status.SetTradSesStatusRejReason(fix50sp2.TradSesStatusRejReasonUnknownOrInvalidTradingsessionid)
		}
		if _, err := s.Send(ctx, status); err != nil {
			return err
		}
	}

	if v.logout {
		go func() { _ = s.Logout(context.Background(), "bye") }()
	}

	return nil
}

// execute runs the command against handler with the options given, and
// returns its error; out receives its output.
func execute(t *testing.T, ctx context.Context, handler session.Handler, subType string, out io.Writer) error {
	t.Helper()

	// The acceptor reads the timeout too: it is restored once it stopped.
	config.GetOptions().Timeout = 500 * time.Millisecond
	t.Cleanup(func() { config.GetOptions().Timeout = 0 })
	initiatortest.Serve(t, handler, "FIX.5.0SP2")
	setOptions(t, "", "", subType)

	if err := Validate(StatusTradingSessionCmd, nil); err != nil {
		t.Fatal(err)
	}

	StatusTradingSessionCmd.SetContext(ctx)
	StatusTradingSessionCmd.SetOut(out)
	t.Cleanup(func() { StatusTradingSessionCmd.SetOut(nil) })

	return Execute(StatusTradingSessionCmd, nil)
}

func TestExecute(t *testing.T) {
	for _, tc := range []struct {
		name    string
		handler session.Handler
		want    error
		printed int
	}{
		{"status", venue{statuses: 2}, nil, 1},
		{"timeout", venue{}, fixerrors.ErrResponseTimeout, 0},
		{"rejected", venue{statuses: 1, status: fix50sp2.TradSesStatusRequestRejected}, fixerrors.ErrFixRejected, 1},
		{"reject", venue{err: &session.RejectError{Reason: 5, Tag: 336, Text: "unknown session"}}, fixerrors.ErrFixRejected, 0},
		{"logout", venue{logout: true}, fixerrors.ErrFixLogout, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := execute(t, context.Background(), tc.handler, "snapshot", &out)
			if !errors.Is(err, tc.want) || (tc.want == nil && err != nil) {
				t.Errorf("Execute() = %v, want %v", err, tc.want)
			}
			if got := strings.Count(out.String(), optionTradSesReqID); got != tc.printed {
				t.Errorf("printed %d statuses, want %d:\n%s", got, tc.printed, out.String())
			}
		})
	}
}

// TestExecuteSubscription checks a subscription prints updates until
// interrupted, past the timeout.
func TestExecuteSubscription(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	out := &syncBuffer{}
	go func() { done <- execute(t, ctx, venue{statuses: 3}, "snapshot_plus_updates", out) }()

	select {
	case err := <-done:
		t.Fatalf("Execute() = %v before being interrupted", err)
	case <-time.After(time.Second):
	}

	cancel()
	if err := <-done; err != nil {
		t.Errorf("Execute() = %v, want nil once interrupted", err)
	}
	if got := strings.Count(out.String(), optionTradSesReqID); got != 3 {
		t.Errorf("printed %d statuses, want 3:\n%s", got, out.String())
	}
}

// TestExecuteSubscriptionEnds checks a subscription ends with its
// rejection, or once the timeout elapsed before its first status, the
// status of another request aside.
func TestExecuteSubscriptionEnds(t *testing.T) {
	for _, tc := range []struct {
		name    string
		handler session.Handler
		want    error
	}{
		{"rejected", venue{statuses: 1, status: fix50sp2.TradSesStatusRequestRejected}, fixerrors.ErrFixRejected},
		{"timeout", venue{}, fixerrors.ErrResponseTimeout},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			if err := execute(t, ctx, tc.handler, "snapshot_plus_updates", io.Discard); !errors.Is(err, tc.want) {
				t.Errorf("Execute() = %v, want %v", err, tc.want)
			}
		})
	}
}

// syncBuffer is a bytes.Buffer safe for concurrent use.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}
