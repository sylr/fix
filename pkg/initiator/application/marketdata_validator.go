//go:build validator
// +build validator

package application

import (
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/artex-io/quickfixgo-fix50sp2/marketdataincrementalrefresh"
	"github.com/artex-io/quickfixgo-fix50sp2/marketdatasnapshotfullrefresh"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"

	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/utils"
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
	metricMarketDataValidatorCrossedUpdates = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "fix",
			Subsystem: "marketdata_validator",
			Name:      "crossed_updates_totals",
			Help:      "Number of updates that resulted in a crossed book",
		},
		[]string{"security"},
	)
	metricMarketDataValidatorBookCrossed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "fix",
			Subsystem: "marketdata_validator",
			Name:      "book_crossed",
			Help:      "Book crossed for security",
		},
		[]string{"security"},
	)
)

func init() {
	prometheus.MustRegister(metricMarketDataValidatorIncrementalRefreshes)
	prometheus.MustRegister(metricMarketDataValidatorOrderUpdates)
	prometheus.MustRegister(metricMarketDataValidatorTradeUpdates)
	prometheus.MustRegister(metricMarketDataValidatorErrors)
	prometheus.MustRegister(metricMarketDataValidatorOrders)
	prometheus.MustRegister(metricMarketDataValidatorBookCrossed)
	prometheus.MustRegister(metricMarketDataValidatorCrossedUpdates)
}

func NewMarketDataValidator(logger *zerolog.Logger, securities []string) *MarketDataValidator {
	secs := make(map[string]*Orders, len(securities))
	for _, security := range securities {
		secs[security] = &Orders{
			orders:        make([]*Order, 0),
			typesVolume:   make(map[enum.OrdType]int64),
			sidesVolume:   make(map[enum.MDEntryType]int64),
			bestBuyOrder:  &Order{},
			bestSellOrder: &Order{},
		}
	}

	mdr := MarketDataValidator{
		Connected: make(chan interface{}),
		Validator: &Validator{
			secList: secs,
			logger:  logger,
		},
		router: quickfix.NewMessageRouter(),
	}
	mdr.Logger = logger

	mdr.router.AddRoute(marketdataincrementalrefresh.Route(mdr.onMarketDataIncrementalRefresh))
	mdr.router.AddRoute(marketdatasnapshotfullrefresh.Route(mdr.onMarketDataSnapshotFullRefresh))

	// Initialize error vectors so that we we have pre-existing 0 values allowing
	// to do operations such as delta() when first errors are reported
	for _, security := range securities {
		metricMarketDataValidatorErrors.WithLabelValues(security, ErrOrderNotFound.Error()).Add(0)
		metricMarketDataValidatorErrors.WithLabelValues(security, ErrOrderAlreadyExists.Error()).Add(0)
		metricMarketDataValidatorCrossedUpdates.WithLabelValues(security).Add(0)
	}

	return &mdr
}

type MarketDataValidator struct {
	utils.QuickFixAppMessageLogger

	Settings        *quickfix.Settings
	Connected       chan interface{}
	FromAppMessages chan quickfix.Messagable
	stopped         bool
	mux             sync.RWMutex
	router          *quickfix.MessageRouter

	Validator *Validator
}

var _ quickfix.Application = (*MarketDataValidator)(nil)

// Stop ensures the app chans are emptied so that quickfix can carry on with
// the LOGOUT process correctly.
func (app *MarketDataValidator) Stop() {
	app.Logger.Debug().Msgf("Stopping MarketDataValidator application")

	app.mux.Lock()
	defer app.mux.Unlock()

	app.stopped = true

	// Empty the channel to avoid blocking
	for len(app.FromAppMessages) > 0 {
		<-app.FromAppMessages
	}
}

// Notification of a session begin created.
func (app *MarketDataValidator) OnCreate(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("New session: %s", sessionID)
}

// Notification of a session successfully logging on.
func (app *MarketDataValidator) OnLogon(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logon: %s", sessionID)

	app.Connected <- struct{}{}
}

// Notification of a session logging off or disconnecting.
func (app *MarketDataValidator) OnLogout(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logout: %s", sessionID)

	close(app.Connected)
}

// Notification of admin message being sent to target.
func (app *MarketDataValidator) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("-> Sending message to admin")

	typ, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	// Logon
	if err == nil && typ == string(enum.MsgType_LOGON) {
		sets := app.Settings.SessionSettings()
		if session, ok := sets[sessionID]; ok {
			if session.HasSetting("Username") {
				username, err := session.Setting("Username")
				if err == nil && len(username) > 0 {
					app.Logger.Debug().Msg("Username injected in logon message")
					message.Header.SetField(tag.Username, quickfix.FIXString(username))
				}
			}
			if session.HasSetting("Password") {
				password, err := session.Setting("Password")
				if err == nil && len(password) > 0 {
					app.Logger.Debug().Msg("Password injected in logon message")
					message.Header.SetField(tag.Password, quickfix.FIXString(password))
				}
			}
		}
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)
}

// Notification of admin message being received from target.
func (app *MarketDataValidator) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from admin")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	return nil
}

// Notification of app message being sent to target.
func (app *MarketDataValidator) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.Logger.Debug().Msgf("-> Sending message to app")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)

	return nil
}

// Notification of app message being received from target.
func (app *MarketDataValidator) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from app")

	msgType, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
		return err
	}

	switch msgType {
	case string(enum.MsgType_BUSINESS_MESSAGE_REJECT):
		app.LogMessage(zerolog.TraceLevel, message, sessionID, false)
		app.Connected <- struct{}{}
		return nil
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	return app.router.Route(message, sessionID)
}

func (app *MarketDataValidator) onMarketDataSnapshotFullRefresh(msg marketdatasnapshotfullrefresh.MarketDataSnapshotFullRefresh, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Info().Msg("Received snapshot full refresh")

	security, err := msg.GetSymbol()
	if err != nil {
		app.Logger.Error().Err(err).Msgf("NoSymbol")
		return err
	}

	var orders *Orders
	var ok bool

	if orders, ok = app.Validator.secList[security]; !ok {
		reason := fmt.Sprintf("symbol not found internaly : %s", security)
		err = quickfix.NewMessageRejectError(reason, 0, nil)
		app.Logger.Error().Err(err).Msgf(reason)
		return err
	}

	mdentries, err := msg.GetNoMDEntries()
	if err != nil {
		app.Logger.Error().Err(err).Msgf("NoMDEntries")
		return err
	}

	for i := 0; i < mdentries.Len(); i++ {
		mdentry := mdentries.Get(i)

		entryType, err := mdentry.GetMDEntryType()
		if err != nil {
			app.Logger.Error().Msgf("No entry type: %v", mdentry.FieldMap)
			continue
		}

		switch entryType {
		case enum.MDEntryType_BID, enum.MDEntryType_OFFER:
			orderID, err := mdentry.GetOrderID()
			if err != nil {
				app.Logger.Error().Msgf("No order ID found: %v", mdentry.FieldMap)
				continue
			}

			orderType, err := mdentry.GetOrdType()
			if err != nil {
				app.Logger.Error().Msgf("No order type found: %v", mdentry.FieldMap)
				continue
			}

			order := Order{
				Id:            orderID,
				Size:          utils.MustNot(mdentry.GetMDEntrySize()),
				RemainingSize: utils.MustNot(mdentry.GetMDEntrySize()),
				Type:          orderType,
				Side:          entryType,
			}

			if err := orders.AddOrder(&order); err != nil {
				app.Logger.Error().Msgf("Error while adding order (%s): %s", order.Id, err)
				metricMarketDataValidatorErrors.WithLabelValues(security, err.Error()).Inc()
			}

		case enum.MDEntryType_TRADE:
			metricMarketDataValidatorTradeUpdates.WithLabelValues(security, "new").Inc()

		default:
			app.Logger.Warn().Msgf("Entry type not implemented: %s", entryType)
		}
	}

	app.Logger.Info().Str("security", security).Any("types", orders.typesVolume).Any("sides", orders.sidesVolume).Msgf("Order book:")

	if orders.isOrderBookCrossed(security) {
		app.Logger.Error().Str("security", security).Any("Best BUY order", orders.bestBuyOrder).Any("Best SELL order", orders.bestSellOrder).Msgf("Order book is crossed")
	}

	return nil
}

func (app *MarketDataValidator) onMarketDataIncrementalRefresh(msg marketdataincrementalrefresh.MarketDataIncrementalRefresh, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	mdentries, err := msg.GetNoMDEntries()
	if err != nil {
		app.Logger.Error().Err(err).Msgf("Received incremental refresh without NoMDEntries")
		return err
	}

	if mdentries.Len() == 0 {
		reason := "MDEntries seems emtpty"
		app.Logger.Error().Err(err).Msgf(reason)
		return quickfix.NewMessageRejectError(reason, 0, nil)
	}

	var security string
	security, err = mdentries.Get(0).GetSymbol()
	if err != nil {
		app.Logger.Error().Err(err).Msgf("No security found in MDEntries")
		return err
	}

	var orders *Orders
	var ok bool
	if orders, ok = app.Validator.secList[security]; !ok {
		reason := fmt.Sprintf("security not found: %s", security)
		app.Logger.Error().Err(err).Msgf(reason)
		return quickfix.NewMessageRejectError(reason, 0, nil)
	}

	metricMarketDataValidatorIncrementalRefreshes.WithLabelValues(security).Inc()

	app.Logger.Info().Msgf("Received incremental refresh with %d entries", mdentries.Len())

	for i := 0; i < mdentries.Len(); i++ {
		mdentry := mdentries.Get(i)
		entryType, err := mdentry.GetMDEntryType()
		if err != nil {
			app.Logger.Error().Err(err).Msgf("MDEntryType error")
			continue
		}

		switch entryType {
		case enum.MDEntryType_BID, enum.MDEntryType_OFFER:
			orderID, err := mdentry.GetOrderID()
			if err != nil {
				app.Logger.Error().Msgf("No order ID found: %v", mdentry.FieldMap)
				continue
			}

			orderType, err := mdentry.GetOrdType()
			if err != nil {
				app.Logger.Error().Msgf("No order type found: %v", mdentry.FieldMap)
				continue
			}

			px, err := mdentry.GetMDEntryPx()
			if err != nil {
				app.Logger.Error().Msgf("No px found: %v", mdentry.FieldMap)
				continue
			}

			order := Order{
				Id:            orderID,
				Size:          utils.MustNot(mdentry.GetMDEntrySize()),
				RemainingSize: utils.MustNot(mdentry.GetMDEntrySize()),
				Type:          orderType,
				Side:          entryType,
				Price:         px,
			}

			updateAction, err := mdentry.GetMDUpdateAction()
			if err != nil {
				app.Logger.Error().Err(err).Msgf("Order GetMDUpdateAction")
				continue
			}

			typeStr := strings.ToLower(dict.OrderTypesReversed[orderType])
			sideStr := strings.ToLower(dict.MDEntryTypesReversed[entryType])

			switch updateAction {
			case enum.MDUpdateAction_NEW:
				metricMarketDataValidatorOrderUpdates.WithLabelValues(security, "new", typeStr, sideStr).Inc()

				if err := orders.AddOrder(&order); err != nil {
					app.Logger.Error().Msgf("Error while adding order (%s): %s", order.Id, err)
					metricMarketDataValidatorErrors.WithLabelValues(security, err.Error()).Inc()
				}

			case enum.MDUpdateAction_CHANGE:
				metricMarketDataValidatorOrderUpdates.WithLabelValues(security, "change", typeStr, sideStr).Inc()

				if err := orders.UpdateOrder(&order); err != nil {
					app.Logger.Error().Msgf("Error while updating order (%s): %s", order.Id, err)
					metricMarketDataValidatorErrors.WithLabelValues(security, err.Error()).Inc()
				}

			case enum.MDUpdateAction_DELETE:
				metricMarketDataValidatorOrderUpdates.WithLabelValues(security, "delete", typeStr, sideStr).Inc()

				if err := orders.DeleteOrder(&order); err != nil {
					app.Logger.Error().Msgf("Error while deleting order (%s): %s", order.Id, err)
					metricMarketDataValidatorErrors.WithLabelValues(security, err.Error()).Inc()
				}
			}

		case enum.MDEntryType_TRADE:
			updateAction, err := mdentry.GetMDUpdateAction()
			if err != nil {
				app.Logger.Error().Err(err).Msgf("Trade GetMDUpdateAction")
				continue
			}

			switch updateAction {
			case enum.MDUpdateAction_NEW:
				metricMarketDataValidatorTradeUpdates.WithLabelValues(security, "new").Inc()
			}

		default:
			app.Logger.Warn().Msgf("Entry type not implemented: %s", entryType)
		}
	}
	app.Logger.Info().Str("security", security).Any("types", orders.typesVolume).Any("sides", orders.sidesVolume).Msgf("Order book:")

	stats := orders.Stats()
	for ty, sides := range stats {
		tyStr := strings.ToLower(dict.OrderTypesReversed[ty])
		for si, count := range sides {
			siStr := strings.ToLower(dict.MDEntryTypesReversed[si])
			metricMarketDataValidatorOrders.WithLabelValues(security, tyStr, siStr).Set(float64(count))
		}
	}

	if orders.isOrderBookCrossed(security) {
		app.Logger.Error().Str("security", security).Any("Best BUY order", orders.bestBuyOrder).Any("Best SELL order", orders.bestSellOrder).Msgf("Order book is crossed")
	}

	return nil
}

var (
	ErrOrderAlreadyExists = fmt.Errorf("order already exists")
	ErrOrderNotFound      = fmt.Errorf("order not found")
)

type Order struct {
	Id            string
	Size          decimal.Decimal
	RemainingSize decimal.Decimal
	Type          enum.OrdType
	Side          enum.MDEntryType
	Price         decimal.Decimal
}

type Orders struct {
	orders        []*Order
	typesVolume   map[enum.OrdType]int64
	sidesVolume   map[enum.MDEntryType]int64
	mux           sync.RWMutex
	bestBuyOrder  *Order
	bestSellOrder *Order
	isCrossed     bool
}

func (o *Orders) Stats() map[enum.OrdType]map[enum.MDEntryType]int64 {
	stats := make(map[enum.OrdType]map[enum.MDEntryType]int64)

	o.mux.RLock()
	defer o.mux.RUnlock()

	for _, o := range o.orders {
		ty := o.Type
		si := o.Side

		if _, ok := stats[ty]; !ok {
			stats[ty] = make(map[enum.MDEntryType]int64)
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

	o.fillBestOrder(order)

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

	o.updateBestPrice(order, enum.MDUpdateAction_DELETE)

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

	o.updateBestPrice(order, enum.MDUpdateAction_CHANGE)

	return nil
}

type Trade struct {
	Id            string
	Size          float32
	RemainingSize float32
}

type Validator struct {
	secList map[string]*Orders
	logger  *zerolog.Logger
}

func (o *Orders) fillBestOrder(order *Order) {
	if order.Side == enum.MDEntryType_BID {
		if o.bestBuyOrder == nil || order.Price.GreaterThan(o.bestBuyOrder.Price) {
			o.bestBuyOrder = order
		}
	} else {
		if o.bestSellOrder == nil || order.Price.LessThan(o.bestSellOrder.Price) {
			o.bestSellOrder = order
		}
	}
}

func (o *Orders) updateBestPrice(order *Order, action enum.MDUpdateAction) {
	if (order.Side == enum.MDEntryType_BID && o.bestBuyOrder == nil) || (order.Side == enum.MDEntryType_OFFER && o.bestSellOrder == nil) {
		goto UPDATE_BEST_ORDER
	}

	if action == enum.MDUpdateAction_CHANGE {
		o.fillBestOrder(order)
		return
	}

	if action == enum.MDUpdateAction_DELETE {
		if order.Side == enum.MDEntryType_BID && order.Id == o.bestBuyOrder.Id {
			o.bestBuyOrder = &Order{}
		} else if order.Side == enum.MDEntryType_OFFER && order.Id == o.bestSellOrder.Id {
			o.bestSellOrder = &Order{}
		} else {
			return
		}
	}

UPDATE_BEST_ORDER:
	clo := o.getOrdersBySide(order.Side)

	for _, co := range clo {
		o.fillBestOrder(co)
	}
}

func (o *Orders) getOrdersBySide(side enum.MDEntryType) []*Order {
	var lo []*Order

	for _, o := range o.orders {
		if o.Side != side {
			continue
		}
		lo = append(lo, o)
	}

	return lo
}

func (o *Orders) isOrderBookCrossed(security string) bool {
	o.mux.Lock()
	defer o.mux.Unlock()

	if o.bestBuyOrder != nil && o.bestSellOrder != nil && o.bestBuyOrder.Price.GreaterThanOrEqual(o.bestSellOrder.Price) {
		// Book is crossed
		if !o.isCrossed {
			metricMarketDataValidatorCrossedUpdates.WithLabelValues(security).Add(1)
			metricMarketDataValidatorBookCrossed.WithLabelValues(security).Set(1)
			o.isCrossed = true
		}
		return true
	}

	if o.isCrossed {
		metricMarketDataValidatorBookCrossed.WithLabelValues(security).Set(0)
		o.isCrossed = false
	}

	return false
}
