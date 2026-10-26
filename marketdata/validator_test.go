//go:build validator

package marketdata

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
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
	if n := v.orders.Stats()[fix50sp2.OrdTypeLimit][fix50sp2.MDEntryTypeBid]; n != 2 {
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
	if _, err := v.orders.Find("O2", "", ""); !errors.Is(err, ErrOrderNotFound) {
		t.Errorf("O2: err = %v, want %v", err, ErrOrderNotFound)
	}
	o1, err := v.orders.Find("O1", "", "")
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

// TestValidatorWithoutOrdType checks the book of a FIX 4.4 session, whose
// market data entries have no OrdType, is kept.
func TestValidatorWithoutOrdType(t *testing.T) {
	logger := zerolog.Nop()
	v := NewValidator(&logger, "EURUSD")

	snapshot := fix50sp2.BuildMarketDataSnapshotFullRefresh("FIX.4.4").
		SetSymbol("EURUSD").
		AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().SetMDEntryType(fix50sp2.MDEntryTypeBid).SetOrderID("O1").SetMDEntrySize(10)).
		AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().SetMDEntryType(fix50sp2.MDEntryTypeOffer).SetOrderID("O2").SetMDEntrySize(20))

	if err := v.Handle(parse(t, snapshot)); err != nil {
		t.Fatal(err)
	}

	if n := v.orders.Len(); n != 2 {
		t.Fatalf("orders = %d, want 2", n)
	}
	if n := v.orders.Stats()[""][fix50sp2.MDEntryTypeOffer]; n != 1 {
		t.Errorf("offers without OrdType = %d, want 1", n)
	}
}

// mdEntry returns a bid or offer entry of an incremental refresh, of action,
// identified by its MDEntryID, without the fields left empty.
func mdEntry(action fix50sp2.MDUpdateAction, side fix50sp2.MDEntryType, id string, orderType fix50sp2.OrdType, size float64) *fix50sp2.NoMDEntriesEntryBuilder {
	e := fix50sp2.NewNoMDEntriesEntry().
		SetMDUpdateAction(action).
		SetMDEntryID(id)
	if side != "" {
		e.SetMDEntryType(side)
	}
	if orderType != "" {
		e.SetOrdType(orderType)
	}
	if size != 0 {
		e.SetMDEntrySize(size)
	}

	return e
}

func refresh(entries ...*fix50sp2.NoMDEntriesEntryBuilder) *fix50sp2.MarketDataIncrementalRefreshBuilder {
	b := fix50sp2.BuildMarketDataIncrementalRefresh(session.FIXT11)
	for _, e := range entries {
		b.AddNoMDEntriesEntry(e)
	}

	return b
}

func snapshot(entries ...*fix50sp2.NoMDEntriesEntryBuilder) *fix50sp2.MarketDataSnapshotFullRefreshBuilder {
	b := fix50sp2.BuildMarketDataSnapshotFullRefresh(session.FIXT11).SetSymbol("EURUSD")
	for _, e := range entries {
		b.AddNoMDEntriesEntry(e)
	}

	return b
}

// value returns the value of m, a gauge or a counter.
func value(t *testing.T, m prometheus.Metric) float64 {
	t.Helper()

	var d dto.Metric
	if err := m.Write(&d); err != nil {
		t.Fatal(err)
	}

	return d.GetGauge().GetValue() + d.GetCounter().GetValue()
}

func gauge(t *testing.T, security string, orderType, side string) float64 {
	return value(t, metricMarketDataValidatorOrders.WithLabelValues(security, orderType, side))
}

func updates(t *testing.T, security, update, orderType, side string) float64 {
	return value(t, metricMarketDataValidatorOrderUpdates.WithLabelValues(security, update, orderType, side))
}

func errorCount(t *testing.T, security string, err error) float64 {
	return value(t, metricMarketDataValidatorErrors.WithLabelValues(security, errorLabel(err)))
}

// TestValidatorDeleteWithoutOrdType checks deleted orders are counted with
// the type and side they were added with, not the delete entry's.
func TestValidatorDeleteWithoutOrdType(t *testing.T) {
	const security = "DELETE"
	logger := zerolog.Nop()
	v := NewValidator(&logger, security)

	if err := v.Handle(parse(t, snapshot(order("", fix50sp2.MDEntryTypeBid, "O1", 10), order("", fix50sp2.MDEntryTypeOffer, "O2", 20)))); err != nil {
		t.Fatal(err)
	}

	del := fix50sp2.NewNoMDEntriesEntry().
		SetMDUpdateAction(fix50sp2.MDUpdateActionDelete).
		SetMDEntryType(fix50sp2.MDEntryTypeBid).
		SetOrderID("O1")
	if err := v.Handle(parse(t, refresh(del))); err != nil {
		t.Fatal(err)
	}

	if n := updates(t, security, "delete", "limit", "bid"); n != 1 {
		t.Errorf("limit bids deleted = %v, want 1", n)
	}
	if n := updates(t, security, "delete", "", "bid"); n != 0 {
		t.Errorf("bids without type deleted = %v, want 0", n)
	}

	// The gauge of a type and side the book no longer holds drops to 0.
	if n := gauge(t, security, "limit", "bid"); n != 0 {
		t.Errorf("limit bids gauge = %v, want 0", n)
	}
	if n := gauge(t, security, "limit", "offer"); n != 1 {
		t.Errorf("limit offers gauge = %v, want 1", n)
	}
}

// TestValidatorNotFound checks updates of orders the book lacks are
// counted as errors and leave it as it is.
func TestValidatorNotFound(t *testing.T) {
	const security = "NOTFOUND"
	logger := zerolog.Nop()
	v := NewValidator(&logger, security)

	if err := v.Handle(parse(t, snapshot(order("", fix50sp2.MDEntryTypeBid, "O1", 10)))); err != nil {
		t.Fatal(err)
	}

	// O1 is a bid: as an offer, it is not found.
	if err := v.Handle(parse(t, refresh(
		order(fix50sp2.MDUpdateActionDelete, fix50sp2.MDEntryTypeOffer, "O1", 10),
		order(fix50sp2.MDUpdateActionChange, fix50sp2.MDEntryTypeBid, "O9", 10),
		mdEntry(fix50sp2.MDUpdateActionDelete, "", "E9", "", 0),
		order(fix50sp2.MDUpdateActionNew, fix50sp2.MDEntryTypeBid, "O1", 10),
	))); err != nil {
		t.Fatal(err)
	}

	if n := errorCount(t, security, ErrOrderNotFound); n != 3 {
		t.Errorf("orders not found = %v, want 3", n)
	}
	if n := errorCount(t, security, ErrOrderAlreadyExists); n != 1 {
		t.Errorf("orders already existing = %v, want 1", n)
	}
	if n := v.orders.Len(); n != 1 {
		t.Errorf("orders = %d, want 1", n)
	}
	if n := gauge(t, security, "limit", "bid"); n != 1 {
		t.Errorf("limit bids gauge = %v, want 1", n)
	}
}

// TestValidatorSecondSnapshot checks a snapshot replaces the book.
func TestValidatorSecondSnapshot(t *testing.T) {
	const security = "SNAPSHOT"
	logger := zerolog.Nop()
	v := NewValidator(&logger, security)

	if err := v.Handle(parse(t, snapshot(order("", fix50sp2.MDEntryTypeBid, "O1", 10), order("", fix50sp2.MDEntryTypeOffer, "O2", 20)))); err != nil {
		t.Fatal(err)
	}
	if n := gauge(t, security, "limit", "offer"); n != 1 {
		t.Errorf("limit offers gauge = %v, want 1", n)
	}

	if err := v.Handle(parse(t, snapshot(order("", fix50sp2.MDEntryTypeBid, "O1", 10), order("", fix50sp2.MDEntryTypeBid, "O3", 30), order("", fix50sp2.MDEntryTypeBid, "O3", 30)))); err != nil {
		t.Fatal(err)
	}

	if n := v.orders.Len(); n != 2 {
		t.Errorf("orders = %d, want 2", n)
	}
	if n := errorCount(t, security, ErrOrderAlreadyExists); n != 1 {
		t.Errorf("orders already existing = %v, want 1, O3 repeated", n)
	}
	if n := gauge(t, security, "limit", "bid"); n != 2 {
		t.Errorf("limit bids gauge = %v, want 2", n)
	}
	if n := gauge(t, security, "limit", "offer"); n != 0 {
		t.Errorf("limit offers gauge = %v, want 0", n)
	}
}

// TestValidatorSparseRefresh checks incremental refreshes identifying
// entries by MDEntryID, which may omit what does not change, MDEntryType
// included.
func TestValidatorSparseRefresh(t *testing.T) {
	const security = "SPARSE"
	logger := zerolog.Nop()
	v := NewValidator(&logger, security)

	if err := v.Handle(parse(t, refresh(
		mdEntry(fix50sp2.MDUpdateActionNew, fix50sp2.MDEntryTypeBid, "1", fix50sp2.OrdTypeLimit, 10),
		mdEntry(fix50sp2.MDUpdateActionNew, fix50sp2.MDEntryTypeOffer, "1", fix50sp2.OrdTypeLimit, 20),
		mdEntry(fix50sp2.MDUpdateActionNew, fix50sp2.MDEntryTypeOffer, "2", fix50sp2.OrdTypeMarket, 30),
	))); err != nil {
		t.Fatal(err)
	}

	if err := v.Handle(parse(t, refresh(
		mdEntry(fix50sp2.MDUpdateActionChange, fix50sp2.MDEntryTypeOffer, "1", "", 25),
		mdEntry(fix50sp2.MDUpdateActionChange, "", "2", "", 35),
		mdEntry(fix50sp2.MDUpdateActionDelete, fix50sp2.MDEntryTypeBid, "1", "", 0),
	))); err != nil {
		t.Fatal(err)
	}

	if n := errorCount(t, security, ErrOrderNotFound) + errorCount(t, security, ErrOrderAlreadyExists); n != 0 {
		t.Errorf("errors = %v, want 0", n)
	}

	offer, err := v.orders.Find("", "1", fix50sp2.MDEntryTypeOffer)
	if err != nil {
		t.Fatal(err)
	}
	if !offer.Size.Equal(decimal.NewFromInt(25)) || offer.Type != fix50sp2.OrdTypeLimit {
		t.Errorf("offer 1 = %s %s, want 25 limit", offer.Size, offer.Type)
	}

	market, err := v.orders.Find("", "2", "")
	if err != nil {
		t.Fatal(err)
	}
	if !market.Size.Equal(decimal.NewFromInt(35)) || market.Type != fix50sp2.OrdTypeMarket || market.Side != fix50sp2.MDEntryTypeOffer {
		t.Errorf("offer 2 = %s %s %s, want 35 market offer", market.Size, market.Type, market.Side)
	}

	if _, err := v.orders.Find("", "1", fix50sp2.MDEntryTypeBid); !errors.Is(err, ErrOrderNotFound) {
		t.Errorf("bid 1: err = %v, want %v", err, ErrOrderNotFound)
	}

	// A delete without MDEntryType deletes the entry its MDEntryID names.
	if err := v.Handle(parse(t, refresh(mdEntry(fix50sp2.MDUpdateActionDelete, "", "2", "", 0)))); err != nil {
		t.Fatal(err)
	}
	if n := updates(t, security, "delete", "market", "offer"); n != 1 {
		t.Errorf("market offers deleted = %v, want 1", n)
	}
	if n := v.orders.Len(); n != 1 {
		t.Errorf("orders = %d, want 1", n)
	}
	if n := gauge(t, security, "market", "offer"); n != 0 {
		t.Errorf("market offers gauge = %v, want 0", n)
	}
}

// TestValidatorRename checks changes referring to an entry by MDEntryRefID
// rename it to their MDEntryID, unless another entry has it or the book
// lacks the entry.
func TestValidatorRename(t *testing.T) {
	const security = "RENAME"
	logger := zerolog.Nop()
	v := NewValidator(&logger, security)

	if err := v.Handle(parse(t, refresh(
		mdEntry(fix50sp2.MDUpdateActionNew, fix50sp2.MDEntryTypeBid, "old", fix50sp2.OrdTypeLimit, 10),
		mdEntry(fix50sp2.MDUpdateActionNew, fix50sp2.MDEntryTypeBid, "taken", fix50sp2.OrdTypeLimit, 30),
		mdEntry(fix50sp2.MDUpdateActionNew, fix50sp2.MDEntryTypeOffer, "offer", fix50sp2.OrdTypeLimit, 40),
	))); err != nil {
		t.Fatal(err)
	}

	if err := v.Handle(parse(t, refresh(
		mdEntry(fix50sp2.MDUpdateActionChange, "", "new", "", 20).SetMDEntryRefID("old"),
		mdEntry(fix50sp2.MDUpdateActionChange, fix50sp2.MDEntryTypeBid, "new", "", 25).SetMDEntryRefID("taken"),
		mdEntry(fix50sp2.MDUpdateActionChange, fix50sp2.MDEntryTypeBid, "other", "", 50).SetMDEntryRefID("unknown"),
		// MDEntryIDs are unique per MDEntryType: an offer may take a bid's.
		mdEntry(fix50sp2.MDUpdateActionChange, fix50sp2.MDEntryTypeOffer, "taken", "", 45).SetMDEntryRefID("offer"),
	))); err != nil {
		t.Fatal(err)
	}

	if n := errorCount(t, security, ErrOrderAlreadyExists); n != 1 {
		t.Errorf("orders already existing = %v, want 1", n)
	}
	if n := errorCount(t, security, ErrOrderNotFound); n != 1 {
		t.Errorf("orders not found = %v, want 1", n)
	}
	if n := updates(t, security, "change", "limit", "bid"); n != 1 {
		t.Errorf("limit bids changed = %v, want 1", n)
	}

	if _, err := v.orders.Find("", "old", ""); !errors.Is(err, ErrOrderNotFound) {
		t.Errorf("old: err = %v, want %v", err, ErrOrderNotFound)
	}

	for _, want := range []struct {
		id   string
		side fix50sp2.MDEntryType
		size int64
	}{
		{"new", fix50sp2.MDEntryTypeBid, 20},
		{"taken", fix50sp2.MDEntryTypeBid, 30},
		{"taken", fix50sp2.MDEntryTypeOffer, 45},
	} {
		o, err := v.orders.Find("", want.id, want.side)
		if err != nil {
			t.Fatalf("%s %s: %v", want.side, want.id, err)
		}
		if !o.Size.Equal(decimal.NewFromInt(want.size)) || o.Type != fix50sp2.OrdTypeLimit {
			t.Errorf("%s %s = %s %s, want %d limit", want.side, want.id, o.Size, o.Type, want.size)
		}
	}

	if n := v.orders.Len(); n != 3 {
		t.Errorf("orders = %d, want 3", n)
	}
}
