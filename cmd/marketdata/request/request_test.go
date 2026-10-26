package marketdatarequest

import (
	"bytes"
	"context"
	"errors"
	"slices"
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
	"sylr.dev/fix/v2/marketdata"
)

// setOptions sets the options of the command as its flags would, restoring
// their defaults at the end of t.
func setOptions(t *testing.T, symbols []string, subType string) {
	t.Helper()

	optionSymbols, optionTypes, optionSubType, optionUpdateType = symbols, []string{"bid", "offer"}, subType, "incremental_refresh"
	optionMDReqID, optionPrintData, optionMarketDepth, optionNoLogout = "", true, 0, false
	t.Cleanup(func() {
		optionSymbols, optionTypes, optionSubType, optionUpdateType = nil, []string{"bid", "offer"}, "snapshot", "incremental_refresh"
		optionMDReqID, optionPrintData, optionMarketDepth, optionNoLogout = "", true, 0, false
	})
}

func TestValidate(t *testing.T) {
	for _, tc := range []struct {
		name    string
		set     func()
		wantErr bool
	}{
		{"ok", func() {}, false},
		{"no symbol", func() { optionSymbols = nil }, true},
		{"empty symbol", func() { optionSymbols = []string{"EURUSD", ""} }, true},
		{"no type", func() { optionTypes = nil }, true},
		{"bad type", func() { optionTypes = []string{"bid", "never"} }, true},
		{"bad subscription type", func() { optionSubType = "never" }, true},
		{"bad update type", func() { optionUpdateType = "never" }, true},
		{"negative depth", func() { optionMarketDepth = -1 }, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setOptions(t, []string{"EURUSD"}, "snapshot")
			tc.set()

			err := Validate(MarketDataRequestCmd, nil)
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Fatalf("Validate() = %v, want error %v", err, tc.wantErr)
			}
			if err != nil && !errors.Is(err, fixerrors.ErrOptions) {
				t.Errorf("Validate() = %v, want an options error", err)
			}
			if err == nil && optionMDReqID == "" {
				t.Error("Validate() did not generate the MDReqID")
			}
		})
	}
}

func TestBuildMessageVersions(t *testing.T) {
	setOptions(t, []string{"EURUSD", "GBPUSD"}, "snapshot_plus_updates")
	optionMDReqID, optionMarketDepth = "REQ-1", 5

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			s := fixmsgtest.Session(version)

			b, err := buildMessage(s)
			if err != nil {
				t.Fatal(err)
			}

			m := fixmsgtest.Validate(t, s, b)
			request := fix50sp2.WrapMarketDataRequest(m)

			if id, _ := request.MDReqID(); id != "REQ-1" {
				t.Errorf("MDReqID = %q, want REQ-1", id)
			}
			if typ, _ := request.SubscriptionRequestType(); typ != fix50sp2.SubscriptionRequestTypeSnapshotPlusUpdates {
				t.Errorf("SubscriptionRequestType = %q, want 1", typ)
			}
			if depth, _, _ := request.MarketDepth(); depth != 5 {
				t.Errorf("MarketDepth = %d, want 5", depth)
			}
			if typ, _ := request.MDUpdateType(); typ != fix50sp2.MDUpdateTypeIncrementalRefresh {
				t.Errorf("MDUpdateType = %q, want 1", typ)
			}

			var types []fix50sp2.MDEntryType
			for _, e := range request.NoMDEntryTypesEntries() {
				typ, _ := e.MDEntryType()
				types = append(types, typ)
			}
			if len(types) != 2 || types[0] != fix50sp2.MDEntryTypeBid || types[1] != fix50sp2.MDEntryTypeOffer {
				t.Errorf("MDEntryTypes = %v, want [0 1]", types)
			}

			var symbols []string
			for _, e := range request.NoRelatedSymEntries() {
				symbol, _ := e.Symbol()
				symbols = append(symbols, symbol)
			}
			if len(symbols) != 2 || symbols[0] != "EURUSD" || symbols[1] != "GBPUSD" {
				t.Errorf("Symbols = %v, want [EURUSD GBPUSD]", symbols)
			}
		})
	}
}

// venue answers each MarketDataRequest with a snapshot of another request,
// then a snapshot of each symbol requested but the last lost ones, then
// updates incremental refreshes. reject, if set, rejects the request with a
// MarketDataRequestReject instead; err, with a Reject or a
// BusinessMessageReject.
type venue struct {
	session.BaseHandler
	lost    int
	updates int
	reject  bool
	err     error
}

func (v venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	if v.err != nil {
		return v.err
	}

	request := fix50sp2.WrapMarketDataRequest(m)
	reqID, _ := request.MDReqID()
	beginString := s.ID().BeginString

	if v.reject {
		_, err := s.Send(ctx, fix50sp2.BuildMarketDataRequestReject(beginString).
			SetMDReqID(reqID).
			SetMDReqRejReason(fix50sp2.MDReqRejReasonUnknownSymbol).
			SetText("unknown symbol"))
		return err
	}

	var symbols []string
	for _, e := range request.NoRelatedSymEntries() {
		symbol, _ := e.Symbol()
		symbols = append(symbols, symbol)
	}

	snapshot := func(id, symbol string) *fix50sp2.MarketDataSnapshotFullRefreshBuilder {
		return fix50sp2.BuildMarketDataSnapshotFullRefresh(beginString).
			SetMDReqID(id).
			SetSymbol(symbol).
			SetLastUpdateTime(time.Now()).
			AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().SetMDEntryType(fix50sp2.MDEntryTypeBid).SetMDEntryPx(1).SetMDEntrySize(1))
	}

	if _, err := s.Send(ctx, snapshot("another", "OTHER")); err != nil {
		return err
	}
	for _, symbol := range symbols[:len(symbols)-v.lost] {
		if _, err := s.Send(ctx, snapshot(reqID, symbol)); err != nil {
			return err
		}
	}
	for range v.updates {
		update := fix50sp2.BuildMarketDataIncrementalRefresh(beginString).
			SetMDReqID(reqID).
			AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().
				SetMDUpdateAction(fix50sp2.MDUpdateActionNew).
				SetMDEntryType(fix50sp2.MDEntryTypeOffer).
				SetSymbol("UPDATE").
				SetMDEntryPx(2).
				SetMDEntrySize(1))
		if _, err := s.Send(ctx, update); err != nil {
			return err
		}
	}

	return nil
}

// execute runs the command against handler for symbols, and returns its
// error; out receives its output.
func execute(t *testing.T, ctx context.Context, handler session.Handler, symbols []string, subType string, out *syncBuffer) error {
	t.Helper()

	// The acceptor reads the timeout too: it is restored once it stopped.
	config.GetOptions().Timeout = 500 * time.Millisecond
	t.Cleanup(func() { config.GetOptions().Timeout = 0 })
	initiatortest.Serve(t, handler, "FIX.5.0SP2")
	setOptions(t, symbols, subType)

	if err := Validate(MarketDataRequestCmd, nil); err != nil {
		t.Fatal(err)
	}

	MarketDataRequestCmd.SetContext(ctx)
	MarketDataRequestCmd.SetOut(out)
	t.Cleanup(func() { MarketDataRequestCmd.SetOut(nil) })

	return Execute(MarketDataRequestCmd, nil)
}

func TestExecute(t *testing.T) {
	symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}

	for _, tc := range []struct {
		name    string
		handler venue
		want    error
		printed []string
	}{
		{"snapshots", venue{}, nil, symbols},
		{"lost snapshot", venue{lost: 1}, fixerrors.ErrResponseTimeout, symbols[:2]},
		{"request reject", venue{reject: true}, marketdata.ErrRequestRejected, nil},
		{"business reject", venue{err: &session.BusinessRejectError{Reason: session.BusinessRejectUnsupportedMessageType, Text: "no market data"}}, fixerrors.ErrFixRejected, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := &syncBuffer{}
			err := execute(t, context.Background(), tc.handler, symbols, "snapshot", out)
			if !errors.Is(err, tc.want) || (tc.want == nil && err != nil) {
				t.Errorf("Execute() = %v, want %v", err, tc.want)
			}
			for _, symbol := range symbols {
				if printed := strings.Contains(out.String(), symbol); printed != slices.Contains(tc.printed, symbol) {
					t.Errorf("%s printed: %v, want %v:\n%s", symbol, printed, !printed, out.String())
				}
			}
			if strings.Contains(out.String(), "OTHER") {
				t.Errorf("printed the snapshot of another request:\n%s", out.String())
			}
		})
	}
}

// TestExecuteSubscription checks a subscription prints the updates until
// interrupted, past the timeout.
func TestExecuteSubscription(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	out := &syncBuffer{}
	go func() {
		done <- execute(t, ctx, venue{updates: 3}, []string{"EURUSD"}, "snapshot_plus_updates", out)
	}()

	select {
	case err := <-done:
		t.Fatalf("Execute() = %v before being interrupted", err)
	case <-time.After(time.Second):
	}

	cancel()
	if err := <-done; err != nil {
		t.Errorf("Execute() = %v, want nil once interrupted", err)
	}
	if got := strings.Count(out.String(), "UPDATE"); got != 3 {
		t.Errorf("printed %d updates, want 3:\n%s", got, out.String())
	}
}

// TestExecuteSubscriptionTimeout checks a subscription fails once the
// timeout elapsed before its first snapshot, that of another request aside.
func TestExecuteSubscriptionTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := execute(t, ctx, venue{lost: 1}, []string{"EURUSD"}, "snapshot_plus_updates", &syncBuffer{})
	if !errors.Is(err, fixerrors.ErrResponseTimeout) {
		t.Errorf("Execute() = %v, want a timeout", err)
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
