package application

import (
	"bytes"
	"text/template"

	natsd "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fixt11"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

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

	sideString, _ := dict.OrderSides[enum.Side(side)]
	typeString, _ := dict.OrderTypes[enum.OrdType(ordType)]

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
	message.Body.Set(field.NewOrderID(utils.MustNot(order.Body.GetString(tag.OrderID))))
	message.Body.Set(field.NewExecID("0"))
	message.Body.Set(field.NewExecType(enum.ExecType(status)))
	message.Body.Set(field.NewOrdStatus(status))
	message.Body.Set(field.NewSide(enum.Side(utils.MustNot(order.Body.GetString(tag.Side)))))
	message.Body.Set(field.NewLeavesQty(utils.MustNot(decimal.NewFromString(utils.MustNot(order.Body.GetString(tag.LeavesQty)))), 2))
	message.Body.Set(field.NewCumQty(utils.MustNot(decimal.NewFromString(utils.MustNot(order.Body.GetString(tag.CumQty)))), 2))

	message.Body.Set(field.NewOrderQty(utils.MustNot(decimal.NewFromString(utils.MustNot(order.Body.GetString(tag.OrderQty)))), 2))
	message.Body.Set(field.NewClOrdID(utils.MustNot(order.Body.GetString(tag.ClOrdID))))

	utils.QuickFixMessagePartSet(&message.Header, utils.MustNot(order.Header.GetString(tag.SenderCompID)), field.NewTargetCompID)
	utils.QuickFixMessagePartSet(&message.Header, utils.MustNot(order.Header.GetString(tag.SenderSubID)), field.NewTargetSubID)
	utils.QuickFixMessagePartSet(&message.Header, utils.MustNot(order.Header.GetString(tag.TargetCompID)), field.NewSenderCompID)
	utils.QuickFixMessagePartSet(&message.Header, utils.MustNot(order.Header.GetString(tag.TargetSubID)), field.NewSenderSubID)

	return quickfix.Send(message)
}
