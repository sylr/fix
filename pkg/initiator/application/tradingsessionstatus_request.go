package application

import (
	"github.com/rs/zerolog"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"

	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/utils"
)

func NewTradingSessionStatusRequest() *TradingSessionStatusRequest {
	sod := TradingSessionStatusRequest{
		Connected:     make(chan interface{}),
		FromAdminChan: make(chan *quickfix.Message),
		FromAppChan:   make(chan *quickfix.Message),
	}

	return &sod
}

type TradingSessionStatusRequest struct {
	utils.QuickFixAppMessageLogger

	Settings *quickfix.Settings

	Connected chan interface{}

	FromAdminChan chan *quickfix.Message
	FromAppChan   chan *quickfix.Message
}

// Notification of a session begin created.
func (app *TradingSessionStatusRequest) OnCreate(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("New session: %s", sessionID)
}

// Notification of a session successfully logging on.
func (app *TradingSessionStatusRequest) OnLogon(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logon: %s", sessionID)

	app.Connected <- struct{}{}
}

// Notification of a session logging off or disconnecting.
func (app *TradingSessionStatusRequest) OnLogout(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logout: %s", sessionID)

	close(app.Connected)
	close(app.FromAdminChan)
	close(app.FromAppChan)
}

// Notification of admin message being sent to target.
func (app *TradingSessionStatusRequest) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
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
func (app *TradingSessionStatusRequest) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
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
func (app *TradingSessionStatusRequest) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.Logger.Debug().Msgf("-> Sending message to app")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)

	return nil
}

// Notification of app message being received from target.
func (app *TradingSessionStatusRequest) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from app")

	typ, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	switch enum.MsgType(typ) {
	case enum.MsgType_TRADING_SESSION_STATUS:
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
