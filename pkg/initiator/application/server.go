package application

import (
	"bytes"
	"fmt"
	"text/template"

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

func NewServer() (*Server, error) {
	tpl := template.New("orders")
	tpl.Parse("orders.{{.Symbol}}.{{.Side}}.{{.Type}}")

	s := Server{
		NatsSubject: tpl,
		router:      quickfix.NewMessageRouter(),
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

	//clOrdID, ferr := order.GetClOrdID()
	//if ferr != nil {
	//	return ferr
	//}

	//senderCompID, ferr := order.Header.GetSenderCompID()
	//if ferr != nil {
	//	return ferr
	//}

	//targetCompID, ferr := order.Header.GetTargetCompID()
	//if ferr != nil {
	//	return ferr
	//}

	side, ferr := order.GetSide()
	if ferr != nil {
		return ferr
	}

	ordType, ferr := order.GetOrdType()
	if ferr != nil {
		return ferr
	}

	//price, ferr := order.GetPrice()
	//if ferr != nil {
	//	return ferr
	//}

	//orderQty, ferr := order.GetOrderQty()
	//if ferr != nil {
	//	return ferr
	//}

	sideString, _ := dict.OrderSides[side]
	typeString, _ := dict.OrderTypes[ordType]

	buf := bytes.NewBuffer([]byte{})
	err := app.NatsSubject.Execute(
		buf,
		NewOrderSingleNatsSubject{
			Symbol: symbol,
			Side:   sideString,
			Type:   typeString,
		},
	)
	if err != nil {
		return quickfix.NewMessageRejectError(err.Error(), int(tag.SessionRejectReason), nil)
	}

	err = app.natsConnetion.Publish(
		buf.String(),
		[]byte(order.ToMessage().String()),
	)
	if err != nil {
		return quickfix.NewMessageRejectError(err.Error(), int(tag.SessionRejectReason), nil)
	}

	app.updateOrder(order, enum.OrdStatus(enum.OrdStatus_NEW))

	return nil
}

func (a *Server) updateOrder(order nos50sp2.NewOrderSingle, status enum.OrdStatus) {
	execReport := er50sp2.New(
		field.NewOrderID(MustNot(order.GetClOrdID())),
		field.NewExecID("aze"),
		field.NewExecType(enum.ExecType(status)),
		field.NewOrdStatus(status),
		field.NewSide(MustNot(order.GetSide())),
		field.NewLeavesQty(MustNot(order.GetOrderQty()), 2),
		field.NewCumQty(MustNot(order.GetOrderQty()), 2),
	)
	execReport.SetOrderQty(MustNot(order.GetOrderQty()), 2)
	execReport.SetClOrdID(MustNot(order.GetClOrdID()))

	execReport.Header.SetSenderCompID(MustNot(order.GetTargetCompID()))
	execReport.Header.SetSenderSubID(MustNot(order.GetTargetSubID()))
	execReport.Header.SetTargetCompID(MustNot(order.GetSenderCompID()))
	execReport.Header.SetTargetSubID(MustNot(order.GetSenderSubID()))

	sendErr := quickfix.Send(execReport)
	if sendErr != nil {
		fmt.Println(sendErr)
	}
}

func MustNot[T any](v T, err quickfix.MessageRejectError) T {
	return v
}
