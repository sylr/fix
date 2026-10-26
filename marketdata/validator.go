//go:build validator

package marketdata

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
)

var (
	metricMarketDataValidatorIncrementalRefreshes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "fix",
			Subsystem: "marketdata_validator",
			Name:      "incremental_refreshes_total",
			Help:      "Number of incremental refresh messages received",
		},
		[]string{"security"},
	)
	metricMarketDataValidatorOrderUpdates = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "fix",
			Subsystem: "marketdata_validator",
			Name:      "order_updates_total",
			Help:      "Number of order updates",
		},
		[]string{"security", "update", "type", "side"},
	)
	metricMarketDataValidatorTradeUpdates = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "fix",
			Subsystem: "marketdata_validator",
			Name:      "trade_updates_total",
			Help:      "Number of trade updates",
		},
		[]string{"security", "type"},
	)
	metricMarketDataValidatorErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "fix",
			Subsystem: "marketdata_validator",
			Name:      "errors_total",
			Help:      "Number validator errors",
		},
		[]string{"security", "error"},
	)
	metricMarketDataValidatorOrders = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "fix",
			Subsystem: "marketdata_validator",
			Name:      "orders",
			Help:      "Current orders in the book",
		},
		[]string{"security", "type", "side"},
	)
)

func init() {
	prometheus.MustRegister(metricMarketDataValidatorIncrementalRefreshes)
	prometheus.MustRegister(metricMarketDataValidatorOrderUpdates)
	prometheus.MustRegister(metricMarketDataValidatorTradeUpdates)
	prometheus.MustRegister(metricMarketDataValidatorErrors)
	prometheus.MustRegister(metricMarketDataValidatorOrders)
}

var (
	ErrOrderAlreadyExists = fmt.Errorf("%w: order already exists", errors.ErrFix)
	ErrOrderNotFound      = fmt.Errorf("%w: order not found", errors.ErrFix)
)

// Validator maintains the order book of a security from the market data it
// receives and reports its inconsistencies as metrics.
type Validator struct {
	fix50sp2.UnimplementedHandler

	security string
	orders   Orders
	logger   *zerolog.Logger
	// gauged holds the type and side labels the orders gauge was set for.
	gauged map[[2]string]bool
}

func NewValidator(logger *zerolog.Logger, security string) *Validator {
	v := Validator{
		security: security,
		logger:   logger,
		gauged:   map[[2]string]bool{},
	}

	// Initialize error vectors so that we we have pre-existing 0 values allowing
	// to do operations such as delta() when first errors are reported
	metricMarketDataValidatorErrors.WithLabelValues(security, errorLabel(ErrOrderNotFound)).Add(0)
	metricMarketDataValidatorErrors.WithLabelValues(security, errorLabel(ErrOrderAlreadyExists)).Add(0)

	return &v
}

// Handle handles m, a message initiator.Client.Await hands.
func (v *Validator) Handle(m *fix.Message) error {
	return fix50sp2.Dispatch(m, v)
}

// OnMarketDataSnapshotFullRefresh replaces the book with the snapshot's.
func (v *Validator) OnMarketDataSnapshotFullRefresh(msg *fix50sp2.MarketDataSnapshotFullRefresh) error {
	v.logger.Info().Msg("Received snapshot full refresh")

	if _, ok, _ := msg.NoMDEntries(); !ok {
		v.logger.Error().Msgf("Received snapshot full refresh without NoMDEntries")
		return nil
	}

	var orders []*Order
	for i, mdentry := range msg.NoMDEntriesEntries() {
		entryType, ok := mdentry.MDEntryType()
		if !ok {
			v.logger.Error().Msgf("No entry type in entry %d", i)
			continue
		}

		switch entryType {
		case fix50sp2.MDEntryTypeBid, fix50sp2.MDEntryTypeOffer:
			order, err := newOrder(&mdentry, entryType)
			if err != nil {
				v.logger.Error().Msgf("%s in entry %d", err, i)
				continue
			}

			if slices.ContainsFunc(orders, func(o *Order) bool { return o.matches(order.Id, order.EntryID, order.Side) }) {
				v.fail(fmt.Sprintf("adding order (%s)", order), ErrOrderAlreadyExists)
				continue
			}

			orders = append(orders, order)

		case fix50sp2.MDEntryTypeTrade:
			metricMarketDataValidatorTradeUpdates.WithLabelValues(v.security, "new").Inc()

		default:
			v.logger.Warn().Msgf("Entry type not implemented: %s", entryType)
		}
	}

	v.orders.Replace(orders)
	v.publish()

	return nil
}

func (v *Validator) OnMarketDataIncrementalRefresh(msg *fix50sp2.MarketDataIncrementalRefresh) error {
	metricMarketDataValidatorIncrementalRefreshes.WithLabelValues(v.security).Inc()

	if _, ok, _ := msg.NoMDEntries(); !ok {
		v.logger.Error().Msgf("Received incremental refresh without NoMDEntries")
		return nil
	}

	mdentries := msg.NoMDEntriesEntries()

	v.logger.Info().Msgf("Received incremental refresh with %d entries", len(mdentries))

	for i, mdentry := range mdentries {
		updateAction, ok := mdentry.MDUpdateAction()
		if !ok {
			v.logger.Error().Msgf("No update action in entry %d", i)
			continue
		}

		entryType, ok := mdentry.MDEntryType()
		if !ok {
			// FIX requires the MDEntryType of new entries only: the others
			// take the type of the entry they refer to.
			if updateAction == fix50sp2.MDUpdateActionNew {
				v.logger.Error().Msgf("No entry type in entry %d", i)
				continue
			}

			orderID, entryID, err := refs(&mdentry, updateAction)
			if err != nil {
				v.logger.Error().Msgf("%s in entry %d", err, i)
				continue
			}

			stored, err := v.orders.Find(orderID, entryID, "")
			if err != nil {
				v.fail(fmt.Sprintf("finding entry %d (%s)", i, cmp.Or(orderID, entryID)), err)
				continue
			}

			entryType = stored.Side
		}

		switch entryType {
		case fix50sp2.MDEntryTypeBid, fix50sp2.MDEntryTypeOffer:
			v.updateOrder(i, &mdentry, updateAction, entryType)

		case fix50sp2.MDEntryTypeTrade:
			switch updateAction {
			case fix50sp2.MDUpdateActionNew:
				metricMarketDataValidatorTradeUpdates.WithLabelValues(v.security, "new").Inc()
			}

		default:
			v.logger.Warn().Msgf("Entry type not implemented: %s", entryType)
		}
	}

	v.publish()

	return nil
}

// updateOrder applies to the book mdentry, the i-th entry of an incremental
// refresh, of side, whose update action is action. The entries changing or
// deleting an order are counted with its type and side, those it was given.
func (v *Validator) updateOrder(i int, mdentry *fix50sp2.NoMDEntriesEntry, action fix50sp2.MDUpdateAction, side fix50sp2.MDEntryType) {
	orderID, entryID, err := refs(mdentry, action)
	if err != nil {
		v.logger.Error().Msgf("%s in entry %d", err, i)
		return
	}

	switch action {
	case fix50sp2.MDUpdateActionNew:
		order, err := newOrder(mdentry, side)
		if err != nil {
			v.logger.Error().Msgf("%s in entry %d", err, i)
			return
		}

		v.countUpdate("new", order)

		if err := v.orders.AddOrder(order); err != nil {
			v.fail(fmt.Sprintf("adding order (%s)", order), err)
		}

	case fix50sp2.MDUpdateActionChange:
		order, err := v.orders.UpdateOrder(orderID, entryID, side, mdentry)
		if err != nil {
			v.fail(fmt.Sprintf("updating order (%s)", cmp.Or(orderID, entryID)), err)
			return
		}

		v.countUpdate("change", order)

	case fix50sp2.MDUpdateActionDelete:
		order, err := v.orders.DeleteOrder(orderID, entryID, side)
		if err != nil {
			v.fail(fmt.Sprintf("deleting order (%s)", cmp.Or(orderID, entryID)), err)
			return
		}

		v.countUpdate("delete", order)

	default:
		v.logger.Warn().Msgf("Update action not implemented: %s in entry %d", action, i)
	}
}

// countUpdate counts the update of order.
func (v *Validator) countUpdate(update string, order *Order) {
	metricMarketDataValidatorOrderUpdates.WithLabelValues(v.security, update, typeLabel(order.Type), sideLabel(order.Side)).Inc()
}

// fail logs err, which happened while doing what, and counts it.
func (v *Validator) fail(what string, err error) {
	v.logger.Error().Msgf("Error while %s: %s", what, err)
	metricMarketDataValidatorErrors.WithLabelValues(v.security, errorLabel(err)).Inc()
}

// publish sets the orders gauge to the book's counts: to zero for the types
// and sides it no longer holds.
func (v *Validator) publish() {
	stats := v.orders.Stats()

	v.logger.Info().Msgf("Order book: %v", stats)

	current := map[[2]string]bool{}
	for ty, sides := range stats {
		for si, count := range sides {
			labels := [2]string{typeLabel(ty), sideLabel(si)}
			current[labels], v.gauged[labels] = true, true
			metricMarketDataValidatorOrders.WithLabelValues(v.security, labels[0], labels[1]).Set(float64(count))
		}
	}

	for labels := range v.gauged {
		if !current[labels] {
			metricMarketDataValidatorOrders.WithLabelValues(v.security, labels[0], labels[1]).Set(0)
		}
	}
}

// errorLabel returns the error label of err, one of the book's errors, in
// the errors metric: its text without the "FIX: " prefix.
func errorLabel(err error) string {
	return strings.TrimPrefix(err.Error(), errors.ErrFix.Error()+": ")
}

func typeLabel(t fix50sp2.OrdType) string {
	return strings.ToLower(dict.OrderTypesReversed[t])
}

func sideLabel(s fix50sp2.MDEntryType) string {
	return strings.ToLower(dict.MDEntryTypesReversed[s])
}

// ids returns the OrderID and MDEntryID of mdentry, either of which
// identifies it, an error if it has neither.
func ids(mdentry *fix50sp2.NoMDEntriesEntry) (orderID, entryID string, err error) {
	orderID, _ = mdentry.OrderID()
	entryID, _ = mdentry.MDEntryID()
	if orderID == "" && entryID == "" {
		return "", "", fmt.Errorf("no order ID or MDEntryID")
	}

	return orderID, entryID, nil
}

// refs returns the OrderID and MDEntryID of the entry mdentry, of action,
// refers to: the MDEntryRefID of a change, which renames that entry to the
// MDEntryID of mdentry, mdentry's own IDs otherwise.
func refs(mdentry *fix50sp2.NoMDEntriesEntry, action fix50sp2.MDUpdateAction) (orderID, entryID string, err error) {
	if refID, _ := mdentry.MDEntryRefID(); action == fix50sp2.MDUpdateActionChange && refID != "" {
		return "", refID, nil
	}

	return ids(mdentry)
}

// entrySize returns the MDEntrySize of mdentry, false if missing or invalid.
func entrySize(mdentry *fix50sp2.NoMDEntriesEntry) (decimal.Decimal, bool) {
	raw, ok := mdentry.MDEntrySizeBytes()
	if !ok {
		return decimal.Decimal{}, false
	}

	size, err := decimal.NewFromString(string(raw))
	if err != nil {
		return decimal.Decimal{}, false
	}

	return size, true
}

// newOrder returns the order of mdentry, a bid or an offer, its size zero if
// missing or invalid.
func newOrder(mdentry *fix50sp2.NoMDEntriesEntry, side fix50sp2.MDEntryType) (*Order, error) {
	orderID, entryID, err := ids(mdentry)
	if err != nil {
		return nil, err
	}

	// FIX 4.3 and 4.4 market data entries have no OrdType: such orders
	// are counted with an empty type.
	orderType, _ := mdentry.OrdType()
	size, _ := entrySize(mdentry)

	return &Order{
		Id:            orderID,
		EntryID:       entryID,
		Size:          size,
		RemainingSize: size,
		Type:          orderType,
		Side:          side,
	}, nil
}

// Order is a bid or an offer of the book, identified by its OrderID or its
// MDEntryID.
type Order struct {
	Id            string
	EntryID       string
	Size          decimal.Decimal
	RemainingSize decimal.Decimal
	Type          fix50sp2.OrdType
	Side          fix50sp2.MDEntryType
}

func (o *Order) String() string {
	return cmp.Or(o.Id, o.EntryID)
}

// matches reports whether o is identified by orderID or entryID and, unless
// side is empty, of side: MDEntryIDs may be unique per MDEntryType only.
func (o *Order) matches(orderID, entryID string, side fix50sp2.MDEntryType) bool {
	if side != "" && o.Side != side {
		return false
	}

	return orderID != "" && o.Id == orderID || entryID != "" && o.EntryID == entryID
}

// changed returns a copy of o changed by mdentry: what mdentry omits, FIX
// allowing it to send only what changed, is o's. The MDEntryID of mdentry
// renames o if mdentry refers to it by MDEntryRefID.
func (o *Order) changed(mdentry *fix50sp2.NoMDEntriesEntry) *Order {
	c := *o

	if refID, _ := mdentry.MDEntryRefID(); refID != "" {
		if entryID, _ := mdentry.MDEntryID(); entryID != "" {
			c.EntryID = entryID
		}
	}

	if orderType, ok := mdentry.OrdType(); ok {
		c.Type = orderType
	}

	if size, ok := entrySize(mdentry); ok {
		c.Size, c.RemainingSize = size, size
	}

	return &c
}

// Orders is a book of orders, safe for concurrent use.
type Orders struct {
	orders []*Order
	mux    sync.RWMutex
}

// Stats returns the number of orders of each type and side.
func (o *Orders) Stats() map[fix50sp2.OrdType]map[fix50sp2.MDEntryType]int64 {
	stats := make(map[fix50sp2.OrdType]map[fix50sp2.MDEntryType]int64)

	o.mux.RLock()
	defer o.mux.RUnlock()

	for _, o := range o.orders {
		if _, ok := stats[o.Type]; !ok {
			stats[o.Type] = make(map[fix50sp2.MDEntryType]int64)
		}

		stats[o.Type][o.Side]++
	}

	return stats
}

func (o *Orders) Len() int {
	o.mux.RLock()
	defer o.mux.RUnlock()

	return len(o.orders)
}

// Find returns the order identified by orderID or entryID, of side unless
// empty.
func (o *Orders) Find(orderID, entryID string, side fix50sp2.MDEntryType) (*Order, error) {
	o.mux.RLock()
	defer o.mux.RUnlock()

	i, err := o.find(orderID, entryID, side)
	if err != nil {
		return nil, err
	}

	return o.orders[i], nil
}

func (o *Orders) find(orderID, entryID string, side fix50sp2.MDEntryType) (int, error) {
	i := slices.IndexFunc(o.orders, func(order *Order) bool { return order.matches(orderID, entryID, side) })
	if i < 0 {
		return 0, ErrOrderNotFound
	}

	return i, nil
}

// AddOrder adds order, unless the book holds one identified as it is.
func (o *Orders) AddOrder(order *Order) error {
	o.mux.Lock()
	defer o.mux.Unlock()

	if _, err := o.find(order.Id, order.EntryID, order.Side); err == nil {
		return ErrOrderAlreadyExists
	}

	o.orders = append(o.orders, order)

	return nil
}

// DeleteOrder deletes the order identified by orderID or entryID, of side
// unless empty, and returns it.
func (o *Orders) DeleteOrder(orderID, entryID string, side fix50sp2.MDEntryType) (*Order, error) {
	o.mux.Lock()
	defer o.mux.Unlock()

	i, err := o.find(orderID, entryID, side)
	if err != nil {
		return nil, err
	}

	order := o.orders[i]
	o.orders = slices.Delete(o.orders, i, i+1)

	return order, nil
}

// UpdateOrder changes by mdentry the order identified by orderID or entryID,
// of side unless empty, and returns it changed, unless another order is
// identified as it is then.
func (o *Orders) UpdateOrder(orderID, entryID string, side fix50sp2.MDEntryType, mdentry *fix50sp2.NoMDEntriesEntry) (*Order, error) {
	o.mux.Lock()
	defer o.mux.Unlock()

	i, err := o.find(orderID, entryID, side)
	if err != nil {
		return nil, err
	}

	order := o.orders[i].changed(mdentry)
	for j, other := range o.orders {
		if j != i && other.matches(order.Id, order.EntryID, order.Side) {
			return nil, ErrOrderAlreadyExists
		}
	}

	o.orders[i] = order

	return order, nil
}

// Replace replaces the book's orders with orders.
func (o *Orders) Replace(orders []*Order) {
	o.mux.Lock()
	defer o.mux.Unlock()

	o.orders = orders
}
