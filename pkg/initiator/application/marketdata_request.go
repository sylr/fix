package application

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/iancoleman/strcase"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/datadictionary"
	"github.com/quickfixgo/tag"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	fix50sp1mdsfr "github.com/quickfixgo/fix50sp1/marketdatasnapshotfullrefresh"
	fix50sp2mdsfr "github.com/quickfixgo/fix50sp2/marketdatasnapshotfullrefresh"

	"sylr.dev/fix/pkg/utils"
)

func NewMarketDataRequest() *MarketDataRequest {
	mdr := MarketDataRequest{
		Connected:   make(chan interface{}),
		FromAppChan: make(chan quickfix.Messagable),
		router:      quickfix.NewMessageRouter(),
	}

	mdr.router.AddRoute(fix50sp1mdsfr.Route(mdr.onFIX50SP1MarketDataSnapshotFullRefresh))
	mdr.router.AddRoute(fix50sp2mdsfr.Route(mdr.onFIX50SP2MarketDataSnapshotFullRefresh))

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

func (app *MarketDataRequest) onFIX50SP1MarketDataSnapshotFullRefresh(msg fix50sp1mdsfr.MarketDataSnapshotFullRefresh, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	group := fix50sp1mdsfr.NewNoMDEntriesRepeatingGroup()
	msg.GetGroup(group)

	printFIX50NoMDEntries[NoMDEntriesRG[fix50sp1mdsfr.NoMDEntries], fix50sp1mdsfr.NoMDEntries](group, msg, app.AppDataDictionary)

	return nil
}

func (app *MarketDataRequest) onFIX50SP2MarketDataSnapshotFullRefresh(msg fix50sp2mdsfr.MarketDataSnapshotFullRefresh, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	group := fix50sp2mdsfr.NewNoMDEntriesRepeatingGroup()
	msg.GetGroup(group)

	printFIX50NoMDEntries[NoMDEntriesRG[fix50sp2mdsfr.NoMDEntries], fix50sp2mdsfr.NoMDEntries](group, msg, app.AppDataDictionary)

	app.FromAppChan <- msg

	return nil
}

type NoMDEntries interface {
	GetMDEntryType() (enum.MDEntryType, quickfix.MessageRejectError)
	GetMDEntryPx() (decimal.Decimal, quickfix.MessageRejectError)
	GetMDEntrySize() (decimal.Decimal, quickfix.MessageRejectError)
}

type NoMDEntriesRG[T NoMDEntries] interface {
	Len() int
	Get(int) T
}

type Messager interface {
	GetSymbol() (string, quickfix.MessageRejectError)
}

func printFIX50NoMDEntries[T NoMDEntriesRG[U], U NoMDEntries](group T, msg Messager, dict *datadictionary.DataDictionary) {
	tw := tabwriter.NewWriter(os.Stdout, 15, 0, 2, ' ', 0)
	tw.Write([]byte(fmt.Sprintf("   SYMBOL\tTYPE\tPRICE\tSIZE\n")))

	for i := 0; i < group.Len(); i++ {
		s := group.Get(i)
		v, _ := s.GetMDEntryType()
		tagField := dict.FieldTypeByTag[int(tag.MDEntryType)]

		type2sym := map[string]string{
			"Bid":   "+",
			"Trade": "=",
			"Offer": "-",
		}
		ty := strcase.ToCamel(strings.ToLower(tagField.Enums[string(v)].Description))
		sy, _ := msg.GetSymbol()
		px, _ := s.GetMDEntryPx()
		si, _ := s.GetMDEntrySize()
		tw.Write([]byte(fmt.Sprintf("%s  %s\t%s\t%s\t%s\t\n", type2sym[ty], sy, ty, px.StringFixed(2), si.StringFixed(2))))
	}

	tw.Flush()
}
