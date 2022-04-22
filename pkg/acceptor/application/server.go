package application

import (
	"bytes"
	"text/template"

	natsd "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fixt11"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"

	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/utils"
)

type AcceptorOptions struct {
	NATSEmbeded      bool
	NATSURL          string
	NATSOrderSubject string
}

func NewAcceptor(options *AcceptorOptions) (*Acceptor, error) {
	var err error

	tpl := template.New("fix")
	_, err = tpl.Parse(options.NATSOrderSubject)
	if err != nil {
		return nil, err
	}

	s := Acceptor{
		NatsOrderSubject: tpl,
		router:           quickfix.NewMessageRouter(),
	}

	if options.NATSEmbeded {
		s.natsServer, err = natsd.NewServer(&natsd.Options{})
		s.natsServer.Start()

		if err != nil {
			return nil, err
		}
	}

	natsOptions := []nats.Option{
		nats.DontRandomize(),
		nats.RetryOnFailedConnect(true),
	}
	s.natsConn, err = nats.Connect(options.NATSURL, natsOptions...)
	if err != nil {
		return nil, err
	}

	//s.router.AddRoute(fix50sp2nos.Route(s.onNewOrderSingle))
	s.router.AddRoute(quickfix.ApplVerIDFIX50SP2, string(enum.MsgType_ORDER_SINGLE), s.onNewOrderSingle)

	return &s, nil
}

type NewOrderSingleNatsSubject struct {
	Symbol string
	Side   string
	Type   string
}

type Acceptor struct {
	utils.QuickFixAppMessageLogger

	natsConn   *nats.Conn
	natsServer *natsd.Server

	NatsOrderSubject *template.Template
	router           *quickfix.MessageRouter
	Settings         *quickfix.Settings
}

func (app *Acceptor) Close() {
	app.natsConn.Close()
	app.natsServer.Shutdown()
}

// Notification of a session begin created.
func (app *Acceptor) OnCreate(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("New session: %s", sessionID)
}

// Notification of a session successfully logging on.
func (app *Acceptor) OnLogon(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logon: %s", sessionID)
}

// Notification of a session logging off or disconnecting.
func (app *Acceptor) OnLogout(sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("Logout: %s", sessionID)
}

// Notification of admin message being sent to target.
func (app *Acceptor) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
	app.Logger.Debug().Msgf("-> Sending message to admin")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)
}

// Notification of admin message being received from target.
func (app *Acceptor) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from admin %s", sessionID)

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, false)

	return nil
}

// Notification of app message being sent to target.
func (app *Acceptor) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	app.Logger.Debug().Msgf("-> Sending message to app")

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)

	return nil
}

// Notification of app message being received from target.
func (app *Acceptor) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	app.Logger.Debug().Msgf("<- Message received from app %s", sessionID)

	_, err := message.MsgType()
	if err != nil {
		app.Logger.Error().Msgf("Message type error: %s", err)
	}

	app.LogMessage(zerolog.TraceLevel, message, sessionID, true)

	return app.router.Route(message, sessionID)
}

func (app *Acceptor) onNewOrderSingle(order *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	symbol, ferr := order.Body.GetString(tag.Symbol)
	if ferr != nil {
		return ferr
	}

	side, ferr := order.Body.GetString(tag.Side)
	if ferr != nil {
		return ferr
	}

	ordType, ferr := order.Body.GetString(tag.OrdType)
	if ferr != nil {
		return ferr
	}

	sideString := dict.OrderSides[enum.Side(side)]
	typeString := dict.OrderTypes[enum.OrdType(ordType)]

	buf := bytes.NewBuffer([]byte{})
	subj := NewOrderSingleNatsSubject{
		Symbol: symbol,
		Side:   sideString,
		Type:   typeString,
	}

	err := app.NatsOrderSubject.Execute(buf, subj)
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

func (a *Acceptor) sendExecutionReport(order *quickfix.Message, status enum.OrdStatus) error {
	message := quickfix.NewMessage()
	header := fixt11.NewHeader(&message.Header)

	header.SetField(tag.MsgType, field.NewMsgType(enum.MsgType_EXECUTION_REPORT))
	utils.QuickFixMessagePartSetString(&header, utils.MustNot(order.Header.GetString(tag.SenderCompID)), field.NewTargetCompID)
	utils.QuickFixMessagePartSetString(&header, utils.MustNot(order.Header.GetString(tag.SenderSubID)), field.NewTargetSubID)
	utils.QuickFixMessagePartSetString(&header, utils.MustNot(order.Header.GetString(tag.TargetCompID)), field.NewSenderCompID)
	utils.QuickFixMessagePartSetString(&header, utils.MustNot(order.Header.GetString(tag.TargetSubID)), field.NewSenderSubID)

	utils.QuickFixMessagePartSetString(&message.Body, (utils.MustNot(order.Body.GetString(tag.ClOrdID))), field.NewOrderID)
	utils.QuickFixMessagePartSetString(&message.Body, "0", field.NewExecID)
	utils.QuickFixMessagePartSetString(&message.Body, enum.ExecType(status), field.NewExecType)
	utils.QuickFixMessagePartSetString(&message.Body, status, field.NewOrdStatus)
	utils.QuickFixMessagePartSetString(&message.Body, enum.Side(utils.MustNot(order.Body.GetString(tag.Side))), field.NewSide)
	utils.QuickFixMessagePartSetString(&message.Body, utils.MustNot(order.Body.GetString(tag.ClOrdID)), field.NewClOrdID)
	utils.QuickFixMessagePartSetDecimal(&message.Body, utils.MustNot(order.Body.GetString(tag.LeavesQty)), field.NewLeavesQty, 2)
	utils.QuickFixMessagePartSetDecimal(&message.Body, utils.MustNot(order.Body.GetString(tag.CumQty)), field.NewCumQty, 2)
	utils.QuickFixMessagePartSetDecimal(&message.Body, utils.MustNot(order.Body.GetString(tag.OrderQty)), field.NewOrderQty, 2)

	return quickfix.Send(message)
}
