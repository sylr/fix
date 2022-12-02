package application

import (
	"sync"

	"github.com/rs/zerolog"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"

	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/utils"
)

func NewCancelOrder() *CancelOrder {
	o := CancelOrder{
		Connected:       make(chan quickfix.SessionID),
		FromAppMessages: make(chan *quickfix.Message, 1),
	}

	return &o
}

type CancelOrder struct {
	utils.QuickFixAppMessageLogger

	Settings        *quickfix.Settings
	Connected       chan quickfix.SessionID
	FromAppMessages chan *quickfix.Message
	stopped         bool
	mux             sync.RWMutex
}

// Stop ensures the app chans are emptied so that quickfix can carry on with
// the LOGOUT process correctly.
func (app *CancelOrder) Stop() {
	app.Logger.Debug().Msgf("Stopping NewOrder application")

	app.mux.Lock()
	defer app.mux.Unlock()

	app.stopped = true

	// Empty the channel to avoid blocking
	for len(app.FromAppMessages) > 0 {
		<-app.FromAppMessages
	}
}

// Notification of a session begin created.
func (app *CancelOrder) OnCreate(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("New session: %s", sessionID)
}

// Notification of a session successfully logging on.
func (app *CancelOrder) OnLogon(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logon: %s", sessionID)

	app.Connected <- sessionID
}

// Notification of a session logging off or disconnecting.
func (app *CancelOrder) OnLogout(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logout: %s", sessionID)

	close(app.Connected)
	close(app.FromAppMessages)
}

// Notification of admin message being sent to target.
func (app *CancelOrder) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("-> Sending message to admin")

	typ, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	// Logon
	if err == nil && typ == string(enum.MsgType_LOGON) {
		sets := app.Settings.SessionSettings()
		if session, ok := sets[sessionID]; ok {
			if username, err := session.Setting("Username"); err == nil && len(username) > 0 {
				app.Logger.Debug().Msg("Username injected in logon message")
				message.Header.SetField(tag.Username, quickfix.FIXString(username))
			}
			if password, err := session.Setting("Password"); err == nil && len(password) > 0 {
				app.Logger.Debug().Msg("Password injected in logon message")
				message.Header.SetField(tag.Password, quickfix.FIXString(password))
			}
		}
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)
}

func (app *CancelOrder) treatMessageByType(message *quickfix.Message, f func(enum.MsgType, *quickfix.Message)) {
	msgType, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	f(enum.MsgType(msgType), message)
}

// Notification of admin message being received from target.
func (app *CancelOrder) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from admin")
	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	app.treatMessageByType(message, func(msgType enum.MsgType, _ *quickfix.Message) {
		if msgType == enum.MsgType_REJECT {
			app.FromAppMessages <- message
		}
	})

	return nil
}

// Notification of app message being sent to target.
func (app *CancelOrder) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.Logger.Debug().Msgf("-> Sending message to app")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)

	return nil
}

// Notification of app message being received from target.
func (app *CancelOrder) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from app")

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	app.mux.RLock()
	if app.stopped {
		app.mux.RUnlock()
		return nil
	}
	app.mux.RUnlock()

	app.treatMessageByType(message, func(msgType enum.MsgType, _ *quickfix.Message) {
		switch msgType {
		case enum.MsgType_EXECUTION_REPORT:
			fallthrough
		case enum.MsgType_ORDER_CANCEL_REJECT:
			fallthrough
		case enum.MsgType_ORDER_MASS_CANCEL_REPORT:
			app.FromAppMessages <- message
		default:
			typeName, err := dict.SearchValue(dict.MessageTypes, msgType)
			if err != nil {
				app.Logger.Info().Msgf("Received unexpected message type: %s", msgType)
			} else {
				app.Logger.Info().Msgf("Received unexpected message type: %s(%s)", msgType, typeName)
			}
		}
	})

	return nil
}
