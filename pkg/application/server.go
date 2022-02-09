package application

import (
	"bytes"
	"text/template"

	"github.com/nats-io/nats.go"
	"github.com/quickfixgo/fix50sp1/newordersingle"
	"github.com/quickfixgo/quickfix"

	"sylr.dev/fix/pkg/utils"
)

func NewServer() (*Server, error) {
	tpl := template.New("orders")
	tpl.Parse("orders.{{.Symbol}}.{{.Side}}.{{.type}}")

	s := Server{
		NatsSubject: tpl,
		router:      quickfix.NewMessageRouter(),
	}

	s.router.AddRoute(newordersingle.Route(s.onNewOrderSingle))

	return &s, nil
}

type NewOrderSingleNatsSubject struct {
	Symbol string
	Side   string
	Type   string
}

type Server struct {
	utils.AppMessageLogger

	natsConnetion *nats.Conn
	NatsSubject   *template.Template
	router        *quickfix.MessageRouter
	Settings      *quickfix.Settings
}

func (app *Server) NatsConnect(connString string) error {
	var err error
	app.natsConnetion, err = nats.Connect(connString)
	return err
}

// Notification of a session begin created.
func (app *Server) OnCreate(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("New session: %s", sessionID)
}

// Notification of a session successfully logging on.
func (app *Server) OnLogon(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logon: %s", sessionID)
}

// Notification of a session logging off or disconnecting.
func (app *Server) OnLogout(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logout: %s", sessionID)
}

// Notification of admin message being sent to target.
func (app *Server) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("-> Sending message to admin")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(message, sessionID, true)
}

// Notification of admin message being received from target.
func (app *Server) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from admin")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(message, sessionID, false)

	return nil
}

// Notification of app message being sent to target.
func (app *Server) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.Logger.Debug().Msgf("-> Sending message to app")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(message, sessionID, true)

	return nil
}

// Notification of app message being received from target.
func (app *Server) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from app")

	return app.router.Route(message, sessionID)
}

func (app *Server) onNewOrderSingle(order newordersingle.NewOrderSingle, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	symbol, err := order.GetSymbol()
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})
	app.NatsSubject.Execute(
		buf,
		NewOrderSingleNatsSubject{
			Symbol: symbol,
			Side:   symbol,
			Type:   symbol,
		},
	)
	if err != nil {
		return err
	}
	app.natsConnetion.Publish(
		buf.String(),
		[]byte(order.ToMessage().String()),
	)
	return nil
}
