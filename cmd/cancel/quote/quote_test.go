package cancelquote

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

var cancelFlags = []string{"id", "Q-1", "symbols", "EURUSD,GBPUSD"}

func set(t *testing.T, flags ...string) {
	t.Helper()

	optionstest.Set(t, CancelQuoteCmd, flags...)
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
		{"ok", cancelFlags, nil},
		{"no symbol", []string{"id", "Q-1"}, errors.ErrOptionsNoSymbolGiven},
		{"empty symbol", []string{"symbols", ",EURUSD"}, errors.ErrOptionsNoSymbolGiven},
	} {
		t.Run(tc.name, func(t *testing.T) {
			optionstest.Set(t, CancelQuoteCmd, tc.flags...)
			err := validateOptions()
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("validateOptions() = %v, want %v", err, tc.want)
			}
		})
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

	if got := string(m.MsgType()); got != string(fix50sp2.MsgTypeQuotecancel) {
		t.Errorf("MsgType = %q", got)
	}

	for tag, want := range map[int]string{
		fix50sp2.TagQuoteID:         "Q-1",
		fix50sp2.TagQuoteMsgID:      quoteMsgID,
		fix50sp2.TagQuoteCancelType: string(fix50sp2.QuoteCancelTypeCancelForOneOrMoreSecurities),
		fix50sp2.TagNoQuoteEntries:  "2",
	} {
		if got := string(m.Get(tag)); got != want {
			t.Errorf("tag %d = %q, want %q", tag, got, want)
		}
	}

	var symbols []string
	for _, f := range m.Fields() {
		if f.Tag == fix50sp2.TagSymbol {
			symbols = append(symbols, string(f.Value))
		}
	}
	if len(symbols) != 2 || symbols[0] != "EURUSD" || symbols[1] != "GBPUSD" {
		t.Errorf("symbols = %v", symbols)
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

func statusReport(beginString, quoteID, msgID string, status fix50sp2.QuoteStatus) *fix50sp2.QuoteStatusReportBuilder {
	b := fix50sp2.BuildQuoteStatusReport(beginString).
		SetQuoteID(quoteID).
		SetSymbol("EURUSD").
		SetQuoteStatus(status)
	if msgID != "" {
		b.SetQuoteMsgID(msgID)
	}

	return b
}

// massQuoteAck builds a MassQuoteAcknowledgement, which fix50sp2 lacks.
func massQuoteAck(quoteID string, status fix50sp2.QuoteStatus) *fix.Builder {
	return fix.NewBuilder("FIXT.1.1", string(fix50sp2.MsgTypeMassquoteack)).
		SetString(fix50sp2.TagQuoteID, quoteID).
		SetString(fix50sp2.TagQuoteStatus, string(status))
}

func TestHandle(t *testing.T) {
	set(t, cancelFlags...)

	for _, tc := range []struct {
		name     string
		messages []session.Outgoing
		done     bool
		want     error
	}{
		{"canceled", []session.Outgoing{statusReport("FIXT.1.1", "Q-1", quoteMsgID, fix50sp2.QuoteStatusCanceled)}, true, nil},
		{"canceled for securities", []session.Outgoing{statusReport("FIXT.1.1", "Q-1", "", fix50sp2.QuoteStatusCanceledForSpecificSecurities)}, true, nil},
		{"pending then canceled", []session.Outgoing{
			statusReport("FIXT.1.1", "Q-1", quoteMsgID, fix50sp2.QuoteStatusPending),
			statusReport("FIXT.1.1", "Q-1", quoteMsgID, fix50sp2.QuoteStatusCanceled),
		}, true, nil},
		{"rejected", []session.Outgoing{statusReport("FIXT.1.1", "Q-1", quoteMsgID, fix50sp2.QuoteStatusRejected)}, true, errors.ErrFixQuoteRejected},
		{"reject reason", []session.Outgoing{statusReport("FIXT.1.1", "Q-1", quoteMsgID, fix50sp2.QuoteStatusPending).SetQuoteRejectReason(fix50sp2.QuoteRejectReasonUnknownQuote)}, true, errors.ErrFixQuoteRejected},
		{"not found", []session.Outgoing{statusReport("FIXT.1.1", "Q-1", quoteMsgID, fix50sp2.QuoteStatusQuoteNotFound)}, true, errors.ErrFixQuoteRejected},
		{"another quote", []session.Outgoing{statusReport("FIXT.1.1", "Q-2", "", fix50sp2.QuoteStatusRejected)}, false, initiator.ErrSkip},
		{"another message", []session.Outgoing{statusReport("FIXT.1.1", "Q-1", "other", fix50sp2.QuoteStatusRejected)}, false, initiator.ErrSkip},
		{"mass quote ack", []session.Outgoing{massQuoteAck("Q-1", fix50sp2.QuoteStatusCanceledAll)}, true, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
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

// venue answers each QuoteCancel with a QuoteStatusReport of status.
type venue struct {
	session.BaseHandler
	status fix50sp2.QuoteStatus
}

func (v venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	_, err := s.Send(ctx, statusReport(s.ID().BeginString, string(m.Get(fix50sp2.TagQuoteID)), string(m.Get(fix50sp2.TagQuoteMsgID)), v.status))
	return err
}

func TestExecute(t *testing.T) {
	for _, tc := range []struct {
		name    string
		version string
		status  fix50sp2.QuoteStatus
		want    error
	}{
		{"canceled", "FIX.5.0SP2", fix50sp2.QuoteStatusCanceled, nil},
		{"canceled FIX.4.4", "FIX.4.4", fix50sp2.QuoteStatusCanceledForSpecificSecurities, nil},
		{"rejected", "FIX.5.0SP2", fix50sp2.QuoteStatusRejected, errors.ErrFixQuoteRejected},
	} {
		t.Run(tc.name, func(t *testing.T) {
			initiatortest.Serve(t, venue{status: tc.status}, tc.version)
			optionstest.Set(t, CancelQuoteCmd, cancelFlags...)

			_, err := optionstest.Run(context.Background(), CancelQuoteCmd)
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Fatalf("Run() = %v, want %v", err, tc.want)
			}
		})
	}
}
