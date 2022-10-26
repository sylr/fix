package application

import (
	"sync"

	"github.com/rs/zerolog"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"

	"sylr.dev/fix/pkg/utils"
)

func NewInitiator() *Initiator {
	sl := Initiator{
		Connected: make(chan interface{}),
		Messages:  make(chan *quickfix.Message, 10),
	}

	return &sl
}

type Initiator struct {
	utils.QuickFixAppMessageLogger

	Settings  *quickfix.Settings
	SessionID quickfix.SessionID

	Connected chan interface{}
	Messages  chan *quickfix.Message

	mux sync.RWMutex
}

// Notification of a session begin created.
func (app *Initiator) OnCreate(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("New session: %s", sessionID)
	app.SessionID = sessionID
}

// Notification of a session successfully logging on.
func (app *Initiator) OnLogon(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logon: %s", sessionID)

	app.Connected <- struct{}{}
}

// Notification of a session logging off or disconnecting.
func (app *Initiator) OnLogout(sessionID quickfix.SessionID) {
	app.mux.Lock()
	defer app.mux.Unlock()

	app.Logger.Debug().Msgf("Logout: %s", sessionID)

	close(app.Connected)
	close(app.Messages)
}

// Notification of admin message being sent to target.
func (app *Initiator) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
	app.mux.Lock()
	defer app.mux.Unlock()

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
func (app *Initiator) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.mux.Lock()
	defer app.mux.Unlock()

	app.LogMessageType(message, sessionID, "<- Message received from admin: ")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	return nil
}

// Notification of app message being sent to target.
func (app *Initiator) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.mux.RLock()
	defer app.mux.RUnlock()

	app.LogMessageType(message, sessionID, "-> Sending message to app:      ")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.InfoLevel, message, sessionID, true)

	app.Messages <- message

	return nil
}

// Notification of app message being received from target.
func (app *Initiator) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.mux.RLock()
	defer app.mux.RUnlock()

	app.LogMessageType(message, sessionID, "<- Message received from app:   ")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.InfoLevel, message, sessionID, false)

	app.Messages <- message

	return nil
}
