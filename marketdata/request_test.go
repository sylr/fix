package marketdata

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/fixt11"
)

func parse(t *testing.T, b interface{ Build() []byte }) *fix.Message {
	t.Helper()

	m, err := fix.Parse(b.Build())
	if err != nil {
		t.Fatal(err)
	}

	return &m
}

func TestRequestHandlerSnapshot(t *testing.T) {
	snapshot := fix50sp2.BuildMarketDataSnapshotFullRefresh(session.FIXT11).
		SetSymbol("EURUSD").
		SetLastUpdateTime(time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)).
		AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().
			SetMDEntryType(fix50sp2.MDEntryTypeBid).
			SetOrderID("O1").
			SetOrdType(fix50sp2.OrdTypeLimit).
			SetMDEntryPx(1.1).
			SetMDEntrySize(1000)).
		AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().
			SetMDEntryType(fix50sp2.MDEntryTypeTrade).
			SetTradeID("T1").
			SetMDEntryPx(1.2).
			SetMDEntrySize(10).
			SetMDEntryDate("20261020"))

	var out bytes.Buffer
	h := &RequestHandler{Writer: &out, Dictionary: fix50sp2.Dictionary()}
	if err := h.Handle(parse(t, snapshot)); err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"+ EURUSD", "O1", "Bid (Limit)", "1.1", "1000", "= EURUSD", "T1", "Trade", "2026-10-20", "2026-10-20T12:00:00"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output does not contain %q:\n%s", want, out.String())
		}
	}
}

func TestRequestHandlerIncrementalRefresh(t *testing.T) {
	refresh := fix50sp2.BuildMarketDataIncrementalRefresh(session.FIXT11).
		AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().
			SetMDUpdateAction(fix50sp2.MDUpdateActionDelete).
			SetMDEntryType(fix50sp2.MDEntryTypeOffer).
			SetSymbol("EURUSD").
			SetOrderID("O2").
			SetMDEntryPx(1.3))

	var out bytes.Buffer
	h := &RequestHandler{Writer: &out, Dictionary: fix50sp2.Dictionary()}
	if err := h.Handle(parse(t, refresh)); err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"- EURUSD", "O2", "Delete", "Offer", "1.3", nilstr} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output does not contain %q:\n%s", want, out.String())
		}
	}
}

func TestRequestHandlerWithoutWriter(t *testing.T) {
	refresh := fix50sp2.BuildMarketDataIncrementalRefresh(session.FIXT11).
		AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().SetMDEntryType(fix50sp2.MDEntryTypeBid))

	h := &RequestHandler{}
	if err := h.Handle(parse(t, refresh)); err != nil {
		t.Fatal(err)
	}
}

func TestRequestHandlerRejects(t *testing.T) {
	h := &RequestHandler{Dictionary: fix50sp2.Dictionary()}

	reject := fix50sp2.BuildMarketDataRequestReject(session.FIXT11).
		SetMDReqID("1").
		SetMDReqRejReason(fix50sp2.MDReqRejReasonUnknownSymbol).
		SetText("no such symbol")

	err := h.Handle(parse(t, reject))
	if !errors.Is(err, ErrRequestRejected) {
		t.Fatalf("err = %v, want %v", err, ErrRequestRejected)
	}
	if !strings.Contains(err.Error(), "UnknownSymbol") || !strings.Contains(err.Error(), "no such symbol") {
		t.Errorf("err = %v, want reason and text", err)
	}

	sessionReject := fixt11.BuildReject(session.FIXT11).
		SetRefSeqNum(2).
		SetText("invalid tag")

	err = h.Handle(parse(t, sessionReject))
	if !errors.Is(err, ErrRejected) {
		t.Fatalf("err = %v, want %v", err, ErrRejected)
	}
	if !strings.Contains(err.Error(), "invalid tag") {
		t.Errorf("err = %v, want text", err)
	}
}
