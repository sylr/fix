package application

import (
	"github.com/rs/zerolog"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"

	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/utils"
)

func NewNewOrder() *NewOrder {
	sod := NewOrder{
		Connected:     make(chan interface{}),
		FromAdminChan: make(chan *quickfix.Message),
		FromAppChan:   make(chan *quickfix.Message),
	}

	return &sod
}

type NewOrder struct {
	utils.QuickFixAppMessageLogger

	Settings *quickfix.Settings

	Connected chan interface{}

	FromAdminChan chan *quickfix.Message
	FromAppChan   chan *quickfix.Message
}

// Notification of a session begin created.
func (app *NewOrder) OnCreate(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("New session: %s", sessionID)
}

// Notification of a session successfully logging on.
func (app *NewOrder) OnLogon(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logon: %s", sessionID)

	app.Connected <- struct{}{}
}

// Notification of a session logging off or disconnecting.
func (app *NewOrder) OnLogout(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logout: %s", sessionID)

	close(app.Connected)
	close(app.FromAdminChan)
	close(app.FromAppChan)
}

// Notification of admin message being sent to target.
func (app *NewOrder) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
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
func (app *NewOrder) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from admin")

	msgType, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	switch msgType {
	case string(enum.MsgType_LOGON):
	case string(enum.MsgType_LOGOUT):
	default:
		app.FromAdminChan <- message
	}

	return nil
}

// Notification of app message being sent to target.
func (app *NewOrder) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.Logger.Debug().Msgf("-> Sending message to app")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)

	return nil
}

// Notification of app message being received from target.
func (app *NewOrder) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from app")

	typ, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	switch enum.MsgType(typ) {
	case enum.MsgType_EXECUTION_REPORT:
		app.FromAppChan <- message
	default:
		typName, err := dict.SearchValue(dict.MessageTypes, enum.MsgType(typ))
		if err != nil {
			app.Logger.Info().Msgf("Received unexpected message type: %s", typ)
		} else {
			app.Logger.Info().Msgf("Received unexpected message type: %s(%s)", typ, typName)
		}
	}

	return nil
}
