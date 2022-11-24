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

func NewSecurityList() *SecurityList {
	sl := SecurityList{
		Connected:       make(chan interface{}),
		FromAppMessages: make(chan *quickfix.Message, 1),
	}

	return &sl
}

type SecurityList struct {
	utils.QuickFixAppMessageLogger

	Settings        *quickfix.Settings
	Connected       chan interface{}
	FromAppMessages chan *quickfix.Message
	stopped         bool
	mux             sync.RWMutex
}

// Stop ensures the app chans are emptied so that quickfix can carry on with
// the LOGOUT process correctly.
func (app *SecurityList) Stop() {
	app.Logger.Debug().Msgf("Stopping SecurityList application")

	app.mux.Lock()
	defer app.mux.Unlock()

	app.stopped = true

	// Empty the channel to avoid blocking
	for len(app.FromAppMessages) > 0 {
		<-app.FromAppMessages
	}
}

// Notification of a session begin created.
func (app *SecurityList) OnCreate(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("New session: %s", sessionID)
}

// Notification of a session successfully logging on.
func (app *SecurityList) OnLogon(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logon: %s", sessionID)

	app.Connected <- struct{}{}
}

// Notification of a session logging off or disconnecting.
func (app *SecurityList) OnLogout(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logout: %s", sessionID)

	close(app.Connected)
	close(app.FromAppMessages)
}

// Notification of admin message being sent to target.
func (app *SecurityList) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
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
func (app *SecurityList) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from admin")

	typ, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	app.mux.RLock()
	if app.stopped {
		app.mux.RUnlock()
		return nil
	}
	app.mux.RUnlock()

	switch typ {
	case string(enum.MsgType_REJECT):
		app.FromAppMessages <- utils.QuickFixCopyMessage(message)
	}

	return nil
}

// Notification of app message being sent to target.
func (app *SecurityList) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.Logger.Debug().Msgf("-> Sending message to app")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)

	return nil
}

// Notification of app message being received from target.
func (app *SecurityList) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from app")

	typ, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	app.mux.RLock()
	if app.stopped {
		app.mux.RUnlock()
		return nil
	}
	app.mux.RUnlock()

	switch enum.MsgType(typ) {
	case enum.MsgType_SECURITY_LIST:
		app.FromAppMessages <- message
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
