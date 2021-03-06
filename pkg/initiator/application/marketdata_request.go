package application

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/datadictionary"
	"github.com/quickfixgo/tag"

	"sylr.dev/fix/pkg/utils"
)

const nilstr = "<nil>"

func NewMarketDataRequest() *MarketDataRequest {
	mdr := MarketDataRequest{
		Connected:   make(chan interface{}),
		FromAppChan: make(chan quickfix.Messagable),
		router:      quickfix.NewMessageRouter(),
	}

	mdr.router.AddRoute(quickfix.ApplVerIDFIX50SP2, string(enum.MsgType_MARKET_DATA_INCREMENTAL_REFRESH), mdr.onMarketDataIncrementalRefresh)
	mdr.router.AddRoute(quickfix.ApplVerIDFIX50SP2, string(enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH), mdr.onMarketDataSnapshotFullRefresh)

	return &mdr
}

type MarketDataRequest struct {
	utils.QuickFixAppMessageLogger

	Settings *quickfix.Settings
	router   *quickfix.MessageRouter

	Connected   chan interface{}
	FromAppChan chan quickfix.Messagable

	mux     sync.Mutex
	stopped bool
}

// Stopped
func (app *MarketDataRequest) Stop() {
	app.mux.Lock()
	defer app.mux.Unlock()

	app.stopped = true
}

// IsStopped
func (app *MarketDataRequest) IsStopped() bool {
	app.mux.Lock()
	defer app.mux.Unlock()

	return app.stopped
}

// Notification of a session begin created.
func (app *MarketDataRequest) OnCreate(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("New session: %s", sessionID)
}

// Notification of a session successfully logging on.
func (app *MarketDataRequest) OnLogon(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logon: %s", sessionID)

	app.Connected <- struct{}{}
}

// Notification of a session logging off or disconnecting.
func (app *MarketDataRequest) OnLogout(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logout: %s", sessionID)

	close(app.Connected)
	close(app.FromAppChan)
}

// Notification of admin message being sent to target.
func (app *MarketDataRequest) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
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
func (app *MarketDataRequest) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from admin")

	typ, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	switch typ {
	case string(enum.MsgType_REJECT):
		app.FromAppChan <- message
	}

	return nil
}

// Notification of app message being sent to target.
func (app *MarketDataRequest) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.Logger.Debug().Msgf("-> Sending message to app")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)

	return nil
}

// Notification of app message being received from target.
func (app *MarketDataRequest) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from app")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	return app.router.Route(message, sessionID)
}

func (app *MarketDataRequest) onMarketDataSnapshotFullRefresh(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	group := quickfix.NewRepeatingGroup(
		tag.NoMDEntries,
		quickfix.GroupTemplate{
			quickfix.GroupElement(tag.MDEntryType),
			quickfix.GroupElement(tag.MDEntryPx),
			quickfix.GroupElement(tag.MDEntrySize),
			quickfix.GroupElement(tag.OrderID),
		},
	)
	msg.Body.GetGroup(group)

	printFIX50NoMDEntriesFull(group, msg, app.AppDataDictionary)

	if !app.IsStopped() {
		app.FromAppChan <- msg
	}

	return nil
}

func (app *MarketDataRequest) onMarketDataIncrementalRefresh(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	group := quickfix.NewRepeatingGroup(
		tag.NoMDEntries,
		quickfix.GroupTemplate{
			quickfix.GroupElement(tag.MDUpdateAction),
			quickfix.GroupElement(tag.MDEntryType),
			quickfix.GroupElement(tag.MDEntryPx),
			quickfix.GroupElement(tag.MDEntrySize),
			quickfix.GroupElement(tag.OrderID),
			quickfix.GroupElement(tag.Text),
			quickfix.GroupElement(tag.TradeID),
			quickfix.GroupElement(tag.MDEntryTime),
			quickfix.GroupElement(tag.MDEntryDate),
			quickfix.GroupElement(tag.TradeCondition),
			quickfix.GroupElement(tag.OpenCloseSettlFlag),
			quickfix.GroupElement(tag.Symbol),
		},
	)
	msg.Body.GetGroup(group)

	printFIX50NoMDEntriesInc(group, app.AppDataDictionary)

	if !app.IsStopped() {
		app.FromAppChan <- msg
	}

	return nil
}

type Messager interface {
	GetSymbol() (string, quickfix.MessageRejectError)
}

var (
	type2sym = map[string]string{
		"Bid":   "+",
		"Trade": "=",
		"Offer": "-",
	}
)

func printFIX50NoMDEntriesFull(group *quickfix.RepeatingGroup, msg *quickfix.Message, dict *datadictionary.DataDictionary) {
	tw := tabwriter.NewWriter(os.Stdout, 15, 0, 2, ' ', 0)
	tw.Write([]byte(fmt.Sprintf("    SYMBOL\t ORDER ID\t TYPE\t PRICE\t SIZE\n")))
	tw.Write([]byte("   " + strings.Repeat("-", 62) + "\n"))

	for i := 0; i < group.Len(); i++ {
		var typSign, typ, symbol, orderID, price, size string

		s := group.Get(i)

		entryType, err := s.GetString(tag.MDEntryType)
		if err != nil {
			typ = nilstr
		} else {
			tagField := dict.FieldTypeByTag[int(tag.MDEntryType)]
			typ = strcase.ToCamel(strings.ToLower(tagField.Enums[entryType].Description))
			typSign = type2sym[typ]
		}

		symbol, err = msg.Body.GetString(tag.Symbol)
		if err != nil {
			symbol = nilstr
		}

		orderID, err = s.GetString(tag.OrderID)
		if err != nil {
			symbol = nilstr
		}

		price, err = s.GetString(tag.MDEntryPx)
		if err != nil {
			price = nilstr
		}

		size, err = s.GetString(tag.MDEntrySize)
		if err != nil {
			size = nilstr
		}

		tw.Write([]byte(fmt.Sprintf("%s  %s\t%s\t%s\t%10s\t%10s\n", typSign, symbol, orderID, typ, price, size)))
	}

	tw.Write([]byte{10})

	tw.Flush()
}

func printFIX50NoMDEntriesInc(group *quickfix.RepeatingGroup, dict *datadictionary.DataDictionary) {
	tw := tabwriter.NewWriter(os.Stdout, 15, 0, 2, ' ', 0)
	tw.Write([]byte("    SYMBOL\t   ORDER/TRADE ID     \t ACTION\t TYPE\t PRICE\t SIZE\t       TIME\n"))
	tw.Write([]byte("   " + strings.Repeat("-", 115) + "\n"))

	for i := 0; i < group.Len(); i++ {
		var typSign, typ, orderTradeID, action, symbol, price, size, tim string

		s := group.Get(i)

		updateAction, err := s.GetString(tag.MDUpdateAction)
		if err != nil {
			action = nilstr
		} else {
			tagField := dict.FieldTypeByTag[int(tag.MDUpdateAction)]
			action = strcase.ToCamel(strings.ToLower(tagField.Enums[updateAction].Description))
		}

		entryType, err := s.GetString(tag.MDEntryType)
		if err != nil {
			typ = nilstr
		} else {
			tagField := dict.FieldTypeByTag[int(tag.MDEntryType)]
			typ = strcase.ToCamel(strings.ToLower(tagField.Enums[entryType].Description))
			typSign = type2sym[typ]
		}

		orderTradeID, err = s.GetString(tag.OrderID)
		if err != nil {
			orderTradeID = nilstr
		}

		if len(orderTradeID) == 0 || orderTradeID == nilstr {
			orderTradeID, err = s.GetString(tag.TradeID)
			if err != nil {
				orderTradeID = nilstr
			}
		}

		symbol, err = s.GetString(tag.Symbol)
		if err != nil {
			symbol = nilstr
		}

		price, err = s.GetString(tag.MDEntryPx)
		if err != nil {
			price = nilstr
		}

		size, err = s.GetString(tag.MDEntrySize)
		if err != nil {
			size = nilstr
		}

		tim, err = s.GetString(tag.MDEntryTime)
		if err != nil {
			size = nilstr
		}

		tw.Write([]byte(fmt.Sprintf("%s  %s\t%-21s\t%s\t%s\t%10s\t%10s\t%s\n", typSign, symbol, orderTradeID, action, typ, price, size, tim)))
	}

	tw.Write([]byte{10})

	tw.Flush()
}
