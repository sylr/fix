package newquote

import (
	"cmp"
	"context"
	"io"
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

var quoteFlags = []string{"id", "quote-1", "symbol", "EURUSD", "buy-quantities", "10", "buy-prices", "1.08345", "sell-quantities", "20.5", "sell-prices", "1.2"}

func with(flags []string, more ...string) []string {
	return append(flags[:len(flags):len(flags)], more...)
}

func set(t *testing.T, flags ...string) {
	t.Helper()

	optionstest.Set(t, NewQuoteCmd, flags...)
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

func parse(t *testing.T, m session.Outgoing) *fix.Message {
	t.Helper()

	msg, err := fixmsg.Message(m, "FIXT.1.1", nil)
	if err != nil {
		t.Fatal(err)
	}

	return msg
}

func TestValidateOptions(t *testing.T) {
	for _, tc := range []struct {
		name  string
		flags []string
		want  error
	}{
		{"ok", quoteFlags, nil},
		{"no symbol", with(quoteFlags, "symbol", ""), errors.ErrOptionsNoSymbolGiven},
		{"prices without quantities", []string{"symbol", "EURUSD", "buy-prices", "1,2", "buy-quantities", "1"}, errors.ErrOptionsInconsistentValues},
		{"zero price", with(quoteFlags, "sell-prices", "0", "sell-quantities", "1"), nil},
		{"zero quantity", with(quoteFlags, "sell-prices", "1", "sell-quantities", "0"), errors.ErrOptions},
		{"negative quantity", []string{"symbol", "EURUSD", "buy-prices", "1", "buy-quantities", "-1"}, errors.ErrOptions},
		{"auto price with several prices", with(quoteFlags, "buy-prices", "1.1", "buy-quantities", "1", "auto-price-update", "true"), errors.ErrOptions},
		{"party sub id types alone", with(quoteFlags, "party-sub-id-types", "firm"), errors.ErrOptions},
	} {
		t.Run(tc.name, func(t *testing.T) {
			optionstest.Set(t, NewQuoteCmd, tc.flags...)
			err := validateOptions()
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("validateOptions() = %v, want %v", err, tc.want)
			}
		})
	}
}

// TestNoOrderFlags checks the flags a Quote can't carry are gone.
func TestNoOrderFlags(t *testing.T) {
	for _, name := range []string{"origination", "attribute-type", "attribute-value"} {
		if NewQuoteCmd.Flags().Lookup(name) != nil {
			t.Errorf("--%s exists", name)
		}
	}
}

func TestBuildMessage(t *testing.T) {
	set(t, quoteFlags...)

	builder, err := buildMessage(fixmsgtest.Session("FIX.5.0SP2"), 0, "msg-1")
	if err != nil {
		t.Fatal(err)
	}

	checkFields(t, parse(t, builder), map[int]string{
		fix50sp2.TagQuoteID:    "quote-1",
		fix50sp2.TagQuoteMsgID: "msg-1",
		fix50sp2.TagSymbol:     "EURUSD",
		fix50sp2.TagBidSize:    "10",
		fix50sp2.TagBidPx:      "1.08345",
		fix50sp2.TagOfferSize:  "20.5",
		fix50sp2.TagOfferPx:    "1.2",
	})
}

func TestBuildMessageAutoPriceUpdate(t *testing.T) {
	set(t, with(quoteFlags, "auto-price-update", "true")...)

	builder, err := buildMessage(fixmsgtest.Session("FIX.5.0SP2"), 3, "msg-1")
	if err != nil {
		t.Fatal(err)
	}

	checkFields(t, parse(t, builder), map[int]string{
		fix50sp2.TagBidPx:   "1.05345",
		fix50sp2.TagOfferPx: "1.23",
	})
}

func TestBuildCancelMessage(t *testing.T) {
	set(t, quoteFlags...)

	builder, err := buildCancelMessage(fixmsgtest.Session("FIX.5.0SP2"), "msg-2")
	if err != nil {
		t.Fatal(err)
	}

	checkFields(t, parse(t, builder), map[int]string{
		fix50sp2.TagQuoteID:         "quote-1",
		fix50sp2.TagQuoteMsgID:      "msg-2",
		fix50sp2.TagQuoteCancelType: string(fix50sp2.QuoteCancelTypeCancelForOneOrMoreSecurities),
		fix50sp2.TagNoQuoteEntries:  "1",
		fix50sp2.TagSymbol:          "EURUSD",
	})
}

func TestBuildMessageVersions(t *testing.T) {
	set(t, quoteFlags...)

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			builder, err := buildMessage(session, 0, "msg-1")
			if err != nil {
				t.Fatal(err)
			}
			fixmsgtest.Validate(t, session, builder)

			cancel, err := buildCancelMessage(session, "msg-2")
			if err != nil {
				t.Fatal(err)
			}
			fixmsgtest.Validate(t, session, cancel)
		})
	}
}

// TestUpdate checks the updates keep the QuoteID, give each message its
// own QuoteMsgID, and end with the cancellation.
func TestUpdate(t *testing.T) {
	set(t, "id", "quote-1", "symbol", "EURUSD", "buy-quantities", "1,2", "buy-prices", "1.1,1.2", "cancel-quote-after-updates", "true")

	var sent []*fix.Message
	q := newQuoter(io.Discard, nil, fixmsgtest.Session("FIX.5.0SP2"), func(m session.Outgoing) error {
		sent = append(sent, parse(t, m))
		return nil
	})

	if err := q.quote(); err != nil {
		t.Fatal(err)
	}
	for range 3 {
		if err := q.update(); err != nil {
			t.Fatal(err)
		}
	}

	if len(sent) != 3 {
		t.Fatalf("%d messages sent, want 3", len(sent))
	}
	checkFields(t, sent[1], map[int]string{fix50sp2.TagMsgType: string(fix50sp2.MsgTypeQuote), fix50sp2.TagBidPx: "1.2", fix50sp2.TagBidSize: "2"})
	checkFields(t, sent[2], map[int]string{fix50sp2.TagMsgType: string(fix50sp2.MsgTypeQuotecancel)})

	msgIDs := map[string]bool{}
	for _, m := range sent {
		if got := string(m.Get(fix50sp2.TagQuoteID)); got != "quote-1" {
			t.Errorf("QuoteID = %q, want quote-1", got)
		}
		msgIDs[string(m.Get(fix50sp2.TagQuoteMsgID))] = true
	}
	if len(msgIDs) != 3 {
		t.Errorf("QuoteMsgIDs %v, want 3 distinct", msgIDs)
	}
}

func statusReport(quoteID, msgID string, status fix50sp2.QuoteStatus) *fix50sp2.QuoteStatusReportBuilder {
	b := fix50sp2.BuildQuoteStatusReport("FIXT.1.1").
		SetQuoteID(quoteID).
		SetSymbol("EURUSD").
		SetQuoteStatus(status)
	if msgID != "" {
		b.SetQuoteMsgID(msgID)
	}

	return b
}

// fill is the ExecutionReport of an EURUSD order filled, with QuoteRespID
// respID if not empty.
func fill(beginString, respID string) *fix50sp2.ExecutionReportBuilder {
	b := fix50sp2.BuildExecutionReport(beginString).
		SetOrderID("ORD-1").
		SetClOrdID("CL-1").
		SetExecID("1").
		SetExecType(fix50sp2.ExecTypeTrade).
		SetOrdStatus(fix50sp2.OrdStatusFilled).
		SetSymbol("EURUSD").
		SetSide(fix50sp2.SideBuy).
		SetLeavesQty(0).
		SetCumQty(1).
		SetAvgPx(1)
	if respID != "" {
		b.SetQuoteRespID(respID)
	}

	return b
}

func TestHandle(t *testing.T) {
	for _, tc := range []struct {
		name     string
		flags    []string
		messages func(msgID string) []session.Outgoing
		done     bool
		want     error
	}{
		{"accepted", nil, func(id string) []session.Outgoing {
			return []session.Outgoing{statusReport("quote-1", id, fix50sp2.QuoteStatusAccepted)}
		}, true, nil},
		{"accepted by QuoteID", nil, func(string) []session.Outgoing {
			return []session.Outgoing{statusReport("quote-1", "", fix50sp2.QuoteStatusAccepted)}
		}, true, nil},
		{"rejected", nil, func(id string) []session.Outgoing {
			return []session.Outgoing{statusReport("quote-1", id, fix50sp2.QuoteStatusRejected)}
		}, true, errors.ErrFixQuoteRejected},
		{"quote not found", nil, func(id string) []session.Outgoing {
			return []session.Outgoing{statusReport("quote-1", id, fix50sp2.QuoteStatusQuoteNotFound)}
		}, true, errors.ErrFixQuoteRejected},
		{"another quote", nil, func(string) []session.Outgoing {
			return []session.Outgoing{statusReport("quote-2", "", fix50sp2.QuoteStatusRejected), statusReport("quote-1", "other", fix50sp2.QuoteStatusRejected)}
		}, false, initiator.ErrSkip},
		{"order of the symbol", nil, func(string) []session.Outgoing {
			return []session.Outgoing{fill("FIXT.1.1", ""), fill("FIXT.1.1", "resp-1")}
		}, false, initiator.ErrSkip},
		{"order by QuoteMsgID", nil, func(id string) []session.Outgoing {
			return []session.Outgoing{fill("FIXT.1.1", "").SetQuoteMsgID(id)}
		}, true, nil},
		{"order by QuoteRespID", []string{"exec-reports", "2"}, func(id string) []session.Outgoing {
			return []session.Outgoing{statusReport("quote-1", "", fix50sp2.QuoteStatusAccepted).SetQuoteRespID("resp-1"), fill("FIXT.1.1", "resp-1")}
		}, true, nil},
		{"final state", []string{"exec-reports", "0", "stop-on-final-state", "true"}, func(id string) []session.Outgoing {
			return []session.Outgoing{statusReport("quote-1", id, fix50sp2.QuoteStatusAccepted), statusReport("quote-1", id, fix50sp2.QuoteStatusCanceled)}
		}, true, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			set(t, with(quoteFlags, tc.flags...)...)
			q := newQuoter(io.Discard, nil, fixmsgtest.Session("FIX.5.0SP2"), func(session.Outgoing) error { return nil })
			if err := q.quote(); err != nil {
				t.Fatal(err)
			}

			var done bool
			var err error
			for _, m := range tc.messages(q.quoteMsgIDs[0]) {
				if done, err = q.handle(parse(t, m)); done || (err != nil && !errors.Is(err, initiator.ErrSkip)) {
					break
				}
			}

			if done != tc.done || (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("handle() = %v, %v, want %v, %v", done, err, tc.done, tc.want)
			}
		})
	}
}

// venue answers each Quote and QuoteCancel with a QuoteStatusReport, or
// with what respond returns if not nil.
type venue struct {
	session.BaseHandler
	cancelStatus fix50sp2.QuoteStatus
	respond      func(beginString string) session.Outgoing

	mu       sync.Mutex
	received []*fix.Message
}

func (v *venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	v.mu.Lock()
	v.received = append(v.received, m.Clone())
	v.mu.Unlock()

	if v.respond != nil {
		_, err := s.Send(ctx, v.respond(s.ID().BeginString))
		return err
	}

	status := fix50sp2.QuoteStatusAccepted
	if string(m.MsgType()) == string(fix50sp2.MsgTypeQuotecancel) {
		status = cmp.Or(v.cancelStatus, fix50sp2.QuoteStatusCanceled)
	}

	report := fix50sp2.BuildQuoteStatusReport(s.ID().BeginString).
		SetQuoteID(string(m.Get(fix50sp2.TagQuoteID))).
		SetSymbol("EURUSD").
		SetQuoteStatus(status)
	if id := m.Get(fix50sp2.TagQuoteMsgID); id != nil {
		report.SetQuoteMsgID(string(id))
	}

	_, err := s.Send(ctx, report)
	return err
}

// TestExecute checks the updates, then the cancellation, are sent, and the
// command stops at the cancellation's report.
func TestExecute(t *testing.T) {
	v := &venue{}
	initiatortest.Serve(t, v, "FIX.5.0SP2")
	optionstest.Set(t, NewQuoteCmd, with(quoteFlags, "buy-quantities", "1", "buy-prices", "1.1", "update-period", "30ms",
		"cancel-quote-after-updates", "true", "exec-reports", "0", "stop-on-final-state", "true")...)

	if _, err := optionstest.Run(context.Background(), NewQuoteCmd); err != nil {
		t.Fatal(err)
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	if len(v.received) != 3 {
		t.Fatalf("%d messages received, want 3", len(v.received))
	}
	if got := string(v.received[2].MsgType()); got != string(fix50sp2.MsgTypeQuotecancel) {
		t.Errorf("last message %s, want a QuoteCancel", got)
	}
}

func TestExecuteTimeout(t *testing.T) {
	initiatortest.Serve(t, session.BaseHandler{}, "FIX.4.4")
	optionstest.Set(t, NewQuoteCmd, with(quoteFlags, "exec-reports-timeout", "200ms", "update-period", "0s")...)

	if _, err := optionstest.Run(context.Background(), NewQuoteCmd); !errors.Is(err, errors.ErrResponseTimeout) {
		t.Errorf("Run() = %v, want a timeout", err)
	}
}

// TestExecuteCancelNotFound checks a quote not found by its cancellation
// fails the command, rather than leaving it waiting.
func TestExecuteCancelNotFound(t *testing.T) {
	initiatortest.Serve(t, &venue{cancelStatus: fix50sp2.QuoteStatusQuoteNotFound}, "FIX.5.0SP2")
	optionstest.Set(t, NewQuoteCmd, with(quoteFlags, "update-period", "30ms", "cancel-quote-after-updates", "true",
		"exec-reports", "0", "stop-on-final-state", "true")...)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := optionstest.Run(ctx, NewQuoteCmd); !errors.Is(err, errors.ErrFixQuoteRejected) {
		t.Errorf("Run() = %v, want a rejection", err)
	}
}

// TestExecuteUnrelated checks the report of an order of the quote's symbol
// it cannot tie to the quote neither completes it nor defeats its timeout,
// whatever the timeout bounds.
func TestExecuteUnrelated(t *testing.T) {
	for _, tc := range []struct {
		name  string
		flags []string
	}{
		{"total", nil},
		{"idle", []string{"exec-reports-timeout-reset", "true", "auto-price-update", "true", "update-period", "50ms"}},
		{"first", []string{"exec-reports", "0"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			initiatortest.Serve(t, &venue{respond: func(b string) session.Outgoing { return fill(b, "") }}, "FIX.4.4")
			optionstest.Set(t, NewQuoteCmd, with(with(quoteFlags, "exec-reports-timeout", "300ms", "update-period", "0s"), tc.flags...)...)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			if _, err := optionstest.Run(ctx, NewQuoteCmd); !errors.Is(err, errors.ErrResponseTimeout) {
				t.Errorf("Run() = %v, want a timeout", err)
			}
		})
	}
}
