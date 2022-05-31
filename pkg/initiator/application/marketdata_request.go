package application

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/datadictionary"
	"github.com/quickfixgo/tag"

	"sylr.dev/fix/pkg/utils"
)

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

	app.FromAppChan <- msg

	return nil
}

func (app *MarketDataRequest) onMarketDataIncrementalRefresh(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
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

	printFIX50NoMDEntriesInc(group, app.AppDataDictionary)

	app.FromAppChan <- msg

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
	tw.Write([]byte("   SYMBOL\tTYPE\tPRICE\tSIZE\n"))

	for i := 0; i < group.Len(); i++ {
		s := group.Get(i)
		v, _ := s.GetString(tag.MDEntryType)
		tagField := dict.FieldTypeByTag[int(tag.MDEntryType)]

		ty := strcase.ToCamel(strings.ToLower(tagField.Enums[string(v)].Description))
		sy, _ := msg.Body.GetString(tag.Symbol)
		px, _ := decimal.NewFromString(utils.MustNot(s.GetString(tag.MDEntryPx)))
		si, _ := decimal.NewFromString(utils.MustNot(s.GetString(tag.MDEntrySize)))
		tw.Write([]byte(fmt.Sprintf("%s  %s\t%s\t%s\t%s\t\n", type2sym[ty], sy, ty, px.StringFixed(2), si.StringFixed(2))))
	}

	tw.Flush()
}

func printFIX50NoMDEntriesInc(group *quickfix.RepeatingGroup, dict *datadictionary.DataDictionary) {
	tw := tabwriter.NewWriter(os.Stdout, 15, 0, 2, ' ', 0)
	tw.Write([]byte(fmt.Sprintf("   SYMBOL\tTYPE\tPRICE\tSIZE\n")))

	for i := 0; i < group.Len(); i++ {
		s := group.Get(i)
		v, _ := s.GetString(tag.MDEntryType)
		tagField := dict.FieldTypeByTag[int(tag.MDEntryType)]

		ty := strcase.ToCamel(strings.ToLower(tagField.Enums[string(v)].Description))
		sy, _ := s.GetString(tag.Symbol)
		px, _ := decimal.NewFromString(utils.MustNot(s.GetString(tag.MDEntryPx)))
		si, _ := decimal.NewFromString(utils.MustNot(s.GetString(tag.MDEntrySize)))
		tw.Write([]byte(fmt.Sprintf("%s  %s\t%s\t%s\t%s\t\n", type2sym[ty], sy, ty, px.StringFixed(2), si.StringFixed(2))))
	}

	tw.Flush()
}
