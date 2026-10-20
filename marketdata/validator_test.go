//go:build validator
// +build validator

package marketdata

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
)

func order(action fix50sp2.MDUpdateAction, side fix50sp2.MDEntryType, id string, size float64) *fix50sp2.NoMDEntriesEntryBuilder {
	e := fix50sp2.NewNoMDEntriesEntry().
		SetMDEntryType(side).
		SetOrderID(id).
		SetOrdType(fix50sp2.OrdTypeLimit).
		SetMDEntrySize(size)
	if action != "" {
		e.SetMDUpdateAction(action)
	}

	return e
}

func TestValidator(t *testing.T) {
	logger := zerolog.Nop()
	v := NewValidator(&logger, "EURUSD")

	snapshot := fix50sp2.BuildMarketDataSnapshotFullRefresh(session.FIXT11).
		SetSymbol("EURUSD").
		AddNoMDEntriesEntry(order("", fix50sp2.MDEntryTypeBid, "O1", 10)).
		AddNoMDEntriesEntry(order("", fix50sp2.MDEntryTypeBid, "O2", 20)).
		AddNoMDEntriesEntry(order("", fix50sp2.MDEntryTypeOffer, "O3", 30)).
		AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().SetMDEntryType(fix50sp2.MDEntryTypeTrade).SetTradeID("T1"))

	if err := v.Handle(parse(t, snapshot)); err != nil {
		t.Fatal(err)
	}

	if n := v.orders.Len(); n != 3 {
		t.Fatalf("orders = %d, want 3", n)
	}
	if n := v.orders.sidesVolume[fix50sp2.MDEntryTypeBid]; n != 2 {
		t.Errorf("bids = %d, want 2", n)
	}

	refresh := fix50sp2.BuildMarketDataIncrementalRefresh(session.FIXT11).
		AddNoMDEntriesEntry(order(fix50sp2.MDUpdateActionNew, fix50sp2.MDEntryTypeOffer, "O4", 40)).
		AddNoMDEntriesEntry(order(fix50sp2.MDUpdateActionChange, fix50sp2.MDEntryTypeBid, "O1", 5.5)).
		AddNoMDEntriesEntry(order(fix50sp2.MDUpdateActionDelete, fix50sp2.MDEntryTypeBid, "O2", 20)).
		AddNoMDEntriesEntry(order(fix50sp2.MDUpdateActionDelete, fix50sp2.MDEntryTypeBid, "O9", 1)).
		AddNoMDEntriesEntry(order(fix50sp2.MDUpdateActionNew, fix50sp2.MDEntryTypeOffer, "O3", 30))

	if err := v.Handle(parse(t, refresh)); err != nil {
		t.Fatal(err)
	}

	if n := v.orders.Len(); n != 3 {
		t.Fatalf("orders = %d, want 3", n)
	}
	if _, _, err := v.orders.GetOrder("O2"); !errors.Is(err, ErrOrderNotFound) {
		t.Errorf("O2: err = %v, want %v", err, ErrOrderNotFound)
	}
	o1, _, err := v.orders.GetOrder("O1")
	if err != nil {
		t.Fatal(err)
	}
	if !o1.Size.Equal(decimal.RequireFromString("5.5")) {
		t.Errorf("O1 size = %s, want 5.5", o1.Size)
	}

	stats := v.orders.Stats()
	if n := stats[fix50sp2.OrdTypeLimit][fix50sp2.MDEntryTypeOffer]; n != 2 {
		t.Errorf("limit offers = %d, want 2", n)
	}
	if n := stats[fix50sp2.OrdTypeLimit][fix50sp2.MDEntryTypeBid]; n != 1 {
		t.Errorf("limit bids = %d, want 1", n)
	}
}

func TestValidatorBusinessMessageReject(t *testing.T) {
	logger := zerolog.Nop()
	v := NewValidator(&logger, "EURUSD")

	reject := fix50sp2.BuildBusinessMessageReject(session.FIXT11).
		SetRefMsgType(fix50sp2.RefMsgTypeMarketdatarequest).
		SetBusinessRejectReason(fix50sp2.BusinessRejectReasonUnsupportedMessageType).
		SetText("unsupported")

	if err := v.Handle(parse(t, reject)); !errors.Is(err, ErrBusinessMessageReject) {
		t.Fatalf("err = %v, want %v", err, ErrBusinessMessageReject)
	}
}
