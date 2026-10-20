//go:build validator
// +build validator

package marketdata

import (
	"fmt"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/fixt11"
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
	ErrBusinessMessageReject = fmt.Errorf("business message reject received")
	ErrOrderAlreadyExists    = fmt.Errorf("order already exists")
	ErrOrderNotFound         = fmt.Errorf("order not found")
)

// Validator maintains the order book of a security from the market data it
// receives and reports its inconsistencies as metrics.
type Validator struct {
	fix50sp2.UnimplementedHandler

	security string
	orders   Orders
	logger   *zerolog.Logger
}

func NewValidator(logger *zerolog.Logger, security string) *Validator {
	v := Validator{
		orders: Orders{
			orders:      make([]*Order, 0),
			typesVolume: make(map[fix50sp2.OrdType]int64),
			sidesVolume: make(map[fix50sp2.MDEntryType]int64),
		},
		security: security,
		logger:   logger,
	}

	// Initialize error vectors so that we we have pre-existing 0 values allowing
	// to do operations such as delta() when first errors are reported
	metricMarketDataValidatorErrors.WithLabelValues(security, ErrOrderNotFound.Error()).Add(0)
	metricMarketDataValidatorErrors.WithLabelValues(security, ErrOrderAlreadyExists.Error()).Add(0)

	return &v
}

// Handle handles m, a message received from initiator.Client.Messages:
// ErrBusinessMessageReject is returned for BusinessMessageRejects.
func (v *Validator) Handle(m *fix.Message) error {
	if string(m.MsgType()) == string(fixt11.MsgTypeReject) {
		v.logger.Warn().Msgf("Reject received: %s", m.Get(fixt11.TagText))
		return nil
	}

	return fix50sp2.Dispatch(m, v)
}

func (v *Validator) OnBusinessMessageReject(msg *fix50sp2.BusinessMessageReject) error {
	text, _ := msg.Text()
	return fmt.Errorf("%w: %s", ErrBusinessMessageReject, text)
}

func (v *Validator) OnMarketDataSnapshotFullRefresh(msg *fix50sp2.MarketDataSnapshotFullRefresh) error {
	v.logger.Info().Msg("Received snapshot full refresh")

	if _, ok, _ := msg.NoMDEntries(); !ok {
		v.logger.Error().Msgf("Received snapshot full refresh without NoMDEntries")
		return nil
	}

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

			if err := v.orders.AddOrder(order); err != nil {
				v.logger.Error().Msgf("Error while adding order (%s): %s", order.Id, err)
				metricMarketDataValidatorErrors.WithLabelValues(v.security, err.Error()).Inc()
			}

		case fix50sp2.MDEntryTypeTrade:
			metricMarketDataValidatorTradeUpdates.WithLabelValues(v.security, "new").Inc()

		default:
			v.logger.Warn().Msgf("Entry type not implemented: %s", entryType)
		}
	}

	v.logger.Info().Msgf("Order book: types=%+v sides=%+v", v.orders.typesVolume, v.orders.sidesVolume)

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

			updateAction, ok := mdentry.MDUpdateAction()
			if !ok {
				v.logger.Error().Msgf("No order update action in entry %d", i)
				continue
			}

			typeStr := strings.ToLower(dict.OrderTypesReversed[order.Type])
			sideStr := strings.ToLower(dict.MDEntryTypesReversed[entryType])

			switch updateAction {
			case fix50sp2.MDUpdateActionNew:
				metricMarketDataValidatorOrderUpdates.WithLabelValues(v.security, "new", typeStr, sideStr).Inc()

				if err := v.orders.AddOrder(order); err != nil {
					v.logger.Error().Msgf("Error while adding order (%s): %s", order.Id, err)
					metricMarketDataValidatorErrors.WithLabelValues(v.security, err.Error()).Inc()
				}
			case fix50sp2.MDUpdateActionChange:
				metricMarketDataValidatorOrderUpdates.WithLabelValues(v.security, "change", typeStr, sideStr).Inc()

				if err := v.orders.UpdateOrder(order); err != nil {
					v.logger.Error().Msgf("Error while updating order (%s): %s", order.Id, err)
					metricMarketDataValidatorErrors.WithLabelValues(v.security, err.Error()).Inc()
				}
			case fix50sp2.MDUpdateActionDelete:
				metricMarketDataValidatorOrderUpdates.WithLabelValues(v.security, "delete", typeStr, sideStr).Inc()

				if err := v.orders.DeleteOrder(order); err != nil {
					v.logger.Error().Msgf("Error while deleting order (%s): %s", order.Id, err)
					metricMarketDataValidatorErrors.WithLabelValues(v.security, err.Error()).Inc()
				}
			}
		case fix50sp2.MDEntryTypeTrade:
			updateAction, ok := mdentry.MDUpdateAction()
			if !ok {
				v.logger.Error().Msgf("No trade update action in entry %d", i)
				continue
			}

			switch updateAction {
			case fix50sp2.MDUpdateActionNew:
				metricMarketDataValidatorTradeUpdates.WithLabelValues(v.security, "new").Inc()
			}

		default:
			v.logger.Warn().Msgf("Entry type not implemented: %s", entryType)
		}
	}

	v.logger.Info().Msgf("Order book: types=%#v sides=%#v", v.orders.typesVolume, v.orders.sidesVolume)

	stats := v.orders.Stats()
	for ty, sides := range stats {
		tyStr := strings.ToLower(dict.OrderTypesReversed[ty])
		for si, count := range sides {
			siStr := strings.ToLower(dict.MDEntryTypesReversed[si])
			metricMarketDataValidatorOrders.WithLabelValues(v.security, tyStr, siStr).Set(float64(count))
		}
	}

	return nil
}

// newOrder returns the order of mdentry, a bid or an offer, its size zero if
// missing or invalid.
func newOrder(mdentry *fix50sp2.NoMDEntriesEntry, side fix50sp2.MDEntryType) (*Order, error) {
	orderID, ok := mdentry.OrderID()
	if !ok {
		return nil, fmt.Errorf("no order ID found")
	}

	orderType, ok := mdentry.OrdType()
	if !ok {
		return nil, fmt.Errorf("no order type found")
	}

	var size decimal.Decimal
	if raw, ok := mdentry.MDEntrySizeBytes(); ok {
		size, _ = decimal.NewFromString(string(raw))
	}

	return &Order{
		Id:            orderID,
		Size:          size,
		RemainingSize: size,
		Type:          orderType,
		Side:          side,
	}, nil
}

type Order struct {
	Id            string
	Size          decimal.Decimal
	RemainingSize decimal.Decimal
	Type          fix50sp2.OrdType
	Side          fix50sp2.MDEntryType
}

type Orders struct {
	orders      []*Order
	typesVolume map[fix50sp2.OrdType]int64
	sidesVolume map[fix50sp2.MDEntryType]int64
	mux         sync.RWMutex
}

func (o *Orders) Stats() map[fix50sp2.OrdType]map[fix50sp2.MDEntryType]int64 {
	stats := make(map[fix50sp2.OrdType]map[fix50sp2.MDEntryType]int64)

	o.mux.RLock()
	defer o.mux.RUnlock()

	for _, o := range o.orders {
		ty := o.Type
		si := o.Side

		if _, ok := stats[ty]; !ok {
			stats[ty] = make(map[fix50sp2.MDEntryType]int64)
			stats[ty][si] = 0
		}

		stats[ty][si] = stats[ty][si] + 1
	}

	return stats
}

func (o *Orders) Len() int {
	return len(o.orders)
}

func (o *Orders) GetOrder(id string) (*Order, int, error) {
	o.mux.RLock()
	defer o.mux.RUnlock()

	return o.getOrder(id)
}

func (o *Orders) getOrder(id string) (*Order, int, error) {
	for i, order := range o.orders {
		if order.Id == id {
			return order, i, nil
		}
	}

	return nil, 0, ErrOrderNotFound
}

func (o *Orders) AddOrder(order *Order) error {
	o.mux.Lock()
	defer o.mux.Unlock()

	_, _, err := o.getOrder(order.Id)
	if err == nil {
		return ErrOrderAlreadyExists
	}

	o.orders = append(o.orders, order)

	if currVolume, ok := o.typesVolume[order.Type]; ok {
		o.typesVolume[order.Type] = currVolume + 1
	} else {
		o.typesVolume[order.Type] = 1
	}

	if currVolume, ok := o.sidesVolume[order.Side]; ok {
		o.sidesVolume[order.Side] = currVolume + 1
	} else {
		o.sidesVolume[order.Side] = 1
	}

	return nil
}

func (o *Orders) DeleteOrder(order *Order) error {
	o.mux.Lock()
	defer o.mux.Unlock()

	_, i, err := o.getOrder(order.Id)
	if err != nil {
		return err
	}

	o.orders = append(o.orders[:i], o.orders[i+1:]...)

	if currVolume, ok := o.typesVolume[order.Type]; ok {
		o.typesVolume[order.Type] = currVolume - 1
	} else {
		return fmt.Errorf("something wrong happened")
	}

	if currVolume, ok := o.sidesVolume[order.Side]; ok {
		o.sidesVolume[order.Side] = currVolume - 1
	} else {
		return fmt.Errorf("something wrong happened")
	}

	return nil
}

func (o *Orders) UpdateOrder(order *Order) error {
	o.mux.Lock()
	defer o.mux.Unlock()

	_, i, err := o.getOrder(order.Id)
	if err != nil {
		return err
	}

	o.orders[i] = order

	return nil
}

type Trade struct {
	Id            string
	Size          float32
	RemainingSize float32
}
