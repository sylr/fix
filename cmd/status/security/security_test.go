package statussecurity

import (
	"bytes"
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
func setOptions(t *testing.T, symbol, subType, reqID string) {
	t.Helper()

	optionSymbol, optionSubType, optionSecurityStatusReqID = symbol, subType, reqID
	t.Cleanup(func() {
		optionSymbol, optionSubType, optionSecurityStatusReqID = "", "snapshot", ""
	})
}

func TestValidate(t *testing.T) {
	for _, tc := range []struct {
		name            string
		symbol, subType string
		wantErr         bool
	}{
		{"ok", "EURUSD", "snapshot", false},
		{"case", "EURUSD", "Snapshot_Plus_Updates", false},
		{"no symbol", "", "snapshot", true},
		{"bad subscription type", "EURUSD", "never", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setOptions(t, tc.symbol, tc.subType, "")

			err := Validate(StatusSecurityCmd, nil)
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Fatalf("Validate() = %v, want error %v", err, tc.wantErr)
			}
			if err != nil && !errors.Is(err, fixerrors.ErrOptions) {
				t.Errorf("Validate() = %v, want an options error", err)
			}
			if err == nil && optionSecurityStatusReqID == "" {
				t.Error("Validate() did not generate the SecurityStatusReqID")
			}
		})
	}
}

func TestBuildMessageVersions(t *testing.T) {
	setOptions(t, "EURUSD", "snapshot", "REQ-1")

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			s := fixmsgtest.Session(version)

			b, err := buildMessage(s)
			if err != nil {
				t.Fatal(err)
			}

			m := fixmsgtest.Validate(t, s, b)
			if id := string(m.Get(fix50sp2.TagSecurityStatusReqID)); id != "REQ-1" {
				t.Errorf("SecurityStatusReqID = %q, want REQ-1", id)
			}
		})
	}
}

// venue answers each SecurityStatusRequest with a SecurityStatus of another
// request, then statuses of it, of SecurityTradingStatus status if set, or
// fails with err, then logs out if logout is set.
type venue struct {
	session.BaseHandler
	statuses int
	status   fix50sp2.SecurityTradingStatus
	err      error
	logout   bool
}

func (v venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	if v.err != nil {
		return v.err
	}

	reqID := string(m.Get(fix50sp2.TagSecurityStatusReqID))
	for i := range v.statuses + 1 {
		id := reqID
		if i == 0 {
			id = "another"
		}
		status := fix50sp2.BuildSecurityStatus(s.ID().BeginString).
			SetSecurityStatusReqID(id).
			SetSymbol(string(m.Get(fix50sp2.TagSymbol)))
		if i > 0 && v.status != "" {
			status.SetSecurityTradingStatus(v.status)
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
	setOptions(t, "EURUSD", subType, "")

	if err := Validate(StatusSecurityCmd, nil); err != nil {
		t.Fatal(err)
	}

	StatusSecurityCmd.SetContext(ctx)
	StatusSecurityCmd.SetOut(out)
	t.Cleanup(func() { StatusSecurityCmd.SetOut(nil) })
	return Execute(StatusSecurityCmd, nil)
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
		{"unknown security", venue{statuses: 1, status: fix50sp2.SecurityTradingStatusUnknownOrInvalid}, fixerrors.ErrFixRejected, 1},
		{"reject", venue{err: &session.BusinessRejectError{Reason: session.BusinessRejectUnknownSecurity, Text: "unknown"}}, fixerrors.ErrFixRejected, 0},
		{"logout", venue{logout: true}, fixerrors.ErrFixLogout, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := execute(t, context.Background(), tc.handler, "snapshot", &out)
			if !errors.Is(err, tc.want) || (tc.want == nil && err != nil) {
				t.Errorf("Execute() = %v, want %v", err, tc.want)
			}
			if got := strings.Count(out.String(), optionSecurityStatusReqID); got != tc.printed {
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
	if got := strings.Count(out.String(), optionSecurityStatusReqID); got != 3 {
		t.Errorf("printed %d statuses, want 3:\n%s", got, out.String())
	}
}

// TestExecuteSubscriptionEnds checks a subscription ends with its security
// unknown, or once the timeout elapsed before its first status, the status
// of another request aside.
func TestExecuteSubscriptionEnds(t *testing.T) {
	for _, tc := range []struct {
		name    string
		handler session.Handler
		want    error
	}{
		{"unknown security", venue{statuses: 1, status: fix50sp2.SecurityTradingStatusUnknownOrInvalid}, fixerrors.ErrFixRejected},
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
