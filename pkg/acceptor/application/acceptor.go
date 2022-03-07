package application

import (
	"bytes"
	"text/template"

	natsd "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	er50sp2 "github.com/quickfixgo/fix50sp2/executionreport"
	nos50sp2 "github.com/quickfixgo/fix50sp2/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/rs/zerolog"

	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/utils"
)

func NewServer(natsdOptions *natsd.Options) (*Server, error) {
	var err error

	tpl := template.New("fix")
	tpl.Parse("orders.{{.Symbol}}.{{.Side}}.{{.Type}}")

	s := Server{
		NatsSubject: tpl,
		router:      quickfix.NewMessageRouter(),
	}

	if natsdOptions != nil {
		s.natsServer, err = natsd.NewServer(natsdOptions)
		s.natsServer.Start()

		if err != nil {
			return nil, err
		}
	}

	s.natsConn, err = nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		return nil, err
	}

	s.router.AddRoute(nos50sp2.Route(s.onNewOrderSingle))

	return &s, nil
}

type NewOrderSingleNatsSubject struct {
	Symbol string
	Side   string
	Type   string
}

type Server struct {
	utils.QuickFixAppMessageLogger

	natsConn   *nats.Conn
	natsServer *natsd.Server

	NatsSubject *template.Template
	router      *quickfix.MessageRouter
	Settings    *quickfix.Settings
}

func (app *Server) Close() {
	app.natsConn.Close()
	app.natsServer.Shutdown()
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

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)
}

// Notification of admin message being received from target.
func (app *Server) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from admin %s", sessionID)

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	return nil
}

// Notification of app message being sent to target.
func (app *Server) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.Logger.Debug().Msgf("-> Sending message to app")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)

	return nil
}

// Notification of app message being received from target.
func (app *Server) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from app %s", sessionID)

	return app.router.Route(message, sessionID)
}

func (app *Server) onNewOrderSingle(order nos50sp2.NewOrderSingle, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	symbol, ferr := order.GetSymbol()
	if ferr != nil {
		return ferr
	}

	side, ferr := order.GetSide()
	if ferr != nil {
		return ferr
	}

	ordType, ferr := order.GetOrdType()
	if ferr != nil {
		return ferr
	}

	sideString, _ := dict.OrderSides[side]
	typeString, _ := dict.OrderTypes[ordType]

	buf := bytes.NewBuffer([]byte{})
	subj := NewOrderSingleNatsSubject{
		Symbol: symbol,
		Side:   sideString,
		Type:   typeString,
	}

	err := app.NatsSubject.Execute(buf, subj)
	if err != nil {
		return quickfix.NewMessageRejectError(err.Error(), int(tag.ApplResponseError), nil)
	}

	err = app.natsConn.Publish(buf.String(), []byte(order.ToMessage().String()))
	if err != nil {
		return quickfix.NewMessageRejectError(err.Error(), int(tag.ApplResponseError), nil)
	}

	err = app.sendExecutionReport(order, enum.OrdStatus(enum.OrdStatus_NEW))
	if err != nil {
		return quickfix.NewMessageRejectError(err.Error(), int(tag.ApplResponseError), nil)
	}

	return nil
}

func (a *Server) sendExecutionReport(order nos50sp2.NewOrderSingle, status enum.OrdStatus) error {
	execReport := er50sp2.New(
		field.NewOrderID(utils.MustNot(order.GetClOrdID())),
		field.NewExecID("0"),
		field.NewExecType(enum.ExecType(status)),
		field.NewOrdStatus(status),
		field.NewSide(utils.MustNot(order.GetSide())),
		field.NewLeavesQty(utils.MustNot(order.GetOrderQty()), 2),
		field.NewCumQty(utils.MustNot(order.GetOrderQty()), 2),
	)

	execReport.SetOrderQty(utils.MustNot(order.GetOrderQty()), 2)
	execReport.SetClOrdID(utils.MustNot(order.GetClOrdID()))

	utils.QuickFixMessagePartSet(&execReport.Header, utils.MustNot(order.GetSenderCompID()), field.NewTargetCompID)
	utils.QuickFixMessagePartSet(&execReport.Header, utils.MustNot(order.GetSenderSubID()), field.NewTargetSubID)
	utils.QuickFixMessagePartSet(&execReport.Header, utils.MustNot(order.GetTargetCompID()), field.NewSenderCompID)
	utils.QuickFixMessagePartSet(&execReport.Header, utils.MustNot(order.GetTargetSubID()), field.NewSenderSubID)

	return quickfix.Send(execReport)
}
