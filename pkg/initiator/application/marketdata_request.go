package application

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/datadictionary"
	"github.com/quickfixgo/tag"

	"sylr.dev/fix/pkg/utils"
)

const nilstr = "<nil>"

func NewMarketDataRequest(printData bool) *MarketDataRequest {
	mdr := MarketDataRequest{
		Connected:       make(chan interface{}),
		FromAppMessages: make(chan quickfix.Messagable, 1),
		router:          quickfix.NewMessageRouter(),
		printData:       printData,
	}

	mdr.router.AddRoute(quickfix.ApplVerIDFIX50SP2, string(enum.MsgType_MARKET_DATA_INCREMENTAL_REFRESH), mdr.onMarketDataIncrementalRefresh)
	mdr.router.AddRoute(quickfix.ApplVerIDFIX50SP2, string(enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH), mdr.onMarketDataSnapshotFullRefresh)

	return &mdr
}

type MarketDataRequest struct {
	utils.QuickFixAppMessageLogger

	Settings        *quickfix.Settings
	Connected       chan interface{}
	FromAppMessages chan quickfix.Messagable
	stopped         bool
	mux             sync.RWMutex
	router          *quickfix.MessageRouter
	printData       bool
}

var _ quickfix.Application = (*MarketDataRequest)(nil)

// Stop ensures the app chans are emptied so that quickfix can carry on with
// the LOGOUT process correctly.
func (app *MarketDataRequest) Stop() {
	app.Logger.Debug().Msgf("Stopping MarketDataRequest application")

	app.mux.Lock()
	defer app.mux.Unlock()

	app.stopped = true

	// Empty the channel to avoid blocking
	for len(app.FromAppMessages) > 0 {
		<-app.FromAppMessages
	}
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
	close(app.FromAppMessages)
}

// Notification of admin message being sent to target.
func (app *MarketDataRequest) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
	app.LogMessageType(message, sessionID, "-> Sending message to admin:    ")

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
	app.LogMessageType(message, sessionID, "<- Message received from admin: ")

	typ, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	switch typ {
	case string(enum.MsgType_REJECT):
		app.FromAppMessages <- message
	}

	return nil
}

// Notification of app message being sent to target.
func (app *MarketDataRequest) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.LogMessageType(message, sessionID, "-> Sending message to app:      ")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)

	return nil
}

// Notification of app message being received from target.
func (app *MarketDataRequest) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.LogMessageType(message, sessionID, "<- Message received from app:   ")

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
			quickfix.GroupElement(tag.OrdType),
			quickfix.GroupElement(tag.TradeID),
			quickfix.GroupElement(tag.MDEntryDate),
			quickfix.GroupElement(tag.MDEntryTime),
			quickfix.GroupElement(tag.Text),
		},
	)
	msg.Body.GetGroup(group)

	if app.printData {
		printFIX50NoMDEntriesFull(group, msg, app.AppDataDictionary)
	}

	app.mux.RLock()
	if app.stopped {
		app.mux.RUnlock()
		return nil
	}
	app.mux.RUnlock()

	app.FromAppMessages <- msg

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
			quickfix.GroupElement(tag.OrdType),
			quickfix.GroupElement(tag.Text),
			quickfix.GroupElement(tag.TradeID),
			quickfix.GroupElement(tag.MDEntryDate),
			quickfix.GroupElement(tag.MDEntryTime),
			quickfix.GroupElement(tag.TradeCondition),
			quickfix.GroupElement(tag.OpenCloseSettlFlag),
			quickfix.GroupElement(tag.Symbol),
			quickfix.GroupElement(tag.Text),
		},
	)
	msg.Body.GetGroup(group)

	if app.printData {
		printFIX50NoMDEntriesInc(group, app.AppDataDictionary)
	}

	app.mux.RLock()
	if app.stopped {
		app.mux.RUnlock()
		return nil
	}
	app.mux.RUnlock()

	app.FromAppMessages <- msg

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
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"SYMBOL", "ORDER ID", "TYPE", "PRICE", "SIZE", "TIME"})
	table.SetBorders(tablewriter.Border{false, false, false, true})
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("-")

	symbol, err := msg.Body.GetString(tag.Symbol)
	if err != nil {
		symbol = nilstr
	}

	for i := 0; i < group.Len(); i++ {
		var typSign, typ, orderTradeID, price, size, tim string

		s := group.Get(i)

		entryType, err := s.GetString(tag.MDEntryType)
		if err != nil {
			typ = nilstr
		} else {
			tagField := dict.FieldTypeByTag[int(tag.MDEntryType)]
			typ = strcase.ToCamel(strings.ToLower(tagField.Enums[entryType].Description))
			typSign = type2sym[typ]
		}

		orderType, err := s.GetString(tag.OrdType)
		if err == nil {
			tagField := dict.FieldTypeByTag[int(tag.OrdType)]
			ordTyp := strcase.ToCamel(strings.ToLower(tagField.Enums[orderType].Description))
			typ = fmt.Sprintf("%s (%s)", typ, ordTyp)
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

		price, err = s.GetString(tag.MDEntryPx)
		if err != nil {
			price = nilstr
		}

		size, err = s.GetString(tag.MDEntrySize)
		if err != nil {
			size = nilstr
		}

		stringDate, errDate := s.GetString(tag.MDEntryDate)
		stringTime, errTime := s.GetString(tag.MDEntryTime)

		if errDate == nil && errTime == nil {
			timeDate, _ := time.Parse("20060102", stringDate)
			timeTime, _ := time.Parse("15:04:05.999999999", stringTime)
			tim = utils.CombineDateAndTime(timeDate, timeTime).Format(time.RFC3339Nano)
		} else if errDate == nil {
			timeDate, _ := time.Parse("20060102", stringDate)
			tim = timeDate.Format("2006-01-02")
		} else if errTime == nil {
			timeTime, _ := time.Parse("15:04:05.999999999", stringTime)
			tim = timeTime.Format("15:04:05.999")
		}

		table.Append([]string{fmt.Sprintf("%s %s", typSign, symbol), orderTradeID, typ, price, size, tim})
	}

	last, err := msg.Body.GetTime(tag.LastUpdateTime)
	if err == nil {
		table.SetFooter([]string{"", "", "", "", "Last Time", last.Format(time.RFC3339Nano)})
	}

	table.Render()
}

func printFIX50NoMDEntriesInc(group *quickfix.RepeatingGroup, dict *datadictionary.DataDictionary) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"SYMBOL", "ID", "ACTION", "TYPE", "PRICE", "SIZE", "TIME"})
	table.SetBorders(tablewriter.Border{false, false, false, true})
	table.SetColMinWidth(2, 8)
	table.SetColMinWidth(3, 14)
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("-")

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

		orderType, err := s.GetString(tag.OrdType)
		if err == nil {
			tagField := dict.FieldTypeByTag[int(tag.OrdType)]
			ordTyp := strcase.ToCamel(strings.ToLower(tagField.Enums[orderType].Description))
			typ = fmt.Sprintf("%s (%s)", typ, ordTyp)
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

		stringDate, errDate := s.GetString(tag.MDEntryDate)
		stringTime, errTime := s.GetString(tag.MDEntryTime)

		if errDate == nil && errTime == nil {
			timeDate, _ := time.Parse("20060102", stringDate)
			timeTime, _ := time.Parse("15:04:05.999999999", stringTime)
			tim = utils.CombineDateAndTime(timeDate, timeTime).Format(time.RFC3339Nano)
		} else if errDate == nil {
			timeDate, _ := time.Parse("20060102", stringDate)
			tim = timeDate.Format("2006-01-02")
		} else if errTime == nil {
			timeTime, _ := time.Parse("15:04:05.999999999", stringTime)
			tim = timeTime.Format("15:04:05.999")
		}

		table.Append([]string{fmt.Sprintf("%s %s", typSign, symbol), orderTradeID, action, typ, price, size, tim})
	}

	table.Render()
}
