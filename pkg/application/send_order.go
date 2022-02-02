package application

import (
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/rs/zerolog"
	"sylr.dev/fix/pkg/utils"
)

func NewSendOrder() *SendOrder {
	sod := SendOrder{
		Connected:   make(chan interface{}),
		FromAppChan: make(chan *quickfix.Message),
	}

	return &sod
}

type SendOrder struct {
	Logger      *zerolog.Logger
	Settings    *quickfix.Settings
	FixVersion  string
	Connected   chan interface{}
	FromAppChan chan *quickfix.Message
}

// Notification of a session begin created.
func (app *SendOrder) OnCreate(sessionID quickfix.SessionID) {
	app.Logger.Trace().Msgf("*SendOrder.OnCreate new session: %s", sessionID)
}

// Notification of a session successfully logging on.
func (app *SendOrder) OnLogon(sessionID quickfix.SessionID) {
	app.Logger.Trace().Msgf("*SendOrder.OnLogon: %s", sessionID)

	app.Connected <- struct{}{}
}

// Notification of a session logging off or disconnecting.
func (app *SendOrder) OnLogout(sessionID quickfix.SessionID) {
	app.Logger.Trace().Msg("*SendOrder.OnLogout")
}

// Notification of admin message being sent to target.
func (app *SendOrder) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
	typ, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("*SendOrder.ToAdmin: %s", err)
	}

	app.Logger.Trace().Msgf("*SendOrder.ToAdmin")
	utils.LogMessage(*app.Logger, zerolog.TraceLevel, message, app.FixVersion)

	// Logon
	if err == nil && typ == string(enum.MsgType_LOGON) {
		sets := app.Settings.SessionSettings()
		if session, ok := sets[sessionID]; ok {
			if session.HasSetting("Username") {
				username, err := session.Setting("Username")
				if err == nil && len(username) > 0 {
					app.Logger.Trace().Msg("*SendOrder.ToAdmin Username injected")
					message.Header.SetField(tag.Username, quickfix.FIXString(username))
				}
			}
		}
	}
}

// Notification of admin message being received from target.
func (app *SendOrder) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("*SendOrder.FromAdmin: %s", err)
	}

	app.Logger.Trace().Msgf("*SendOrder.FromAdmin")
	utils.LogMessage(*app.Logger, zerolog.TraceLevel, message, app.FixVersion)

	return nil
}

// Notification of app message being sent to target.
func (app *SendOrder) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("*SendOrder.ToApp: %s", err)
	}

	app.Logger.Trace().Msgf("*SendOrder.ToApp")
	utils.LogMessage(*app.Logger, zerolog.TraceLevel, message, app.FixVersion)

	return nil
}

// Notification of app message being received from target.
func (app *SendOrder) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("SendOrder.FromApp: %s", err)
	}

	app.Logger.Trace().Msgf("*SendOrder.FromApp")
	utils.LogMessage(*app.Logger, zerolog.DebugLevel, message, app.FixVersion)

	app.FromAppChan <- message

	return nil
}
