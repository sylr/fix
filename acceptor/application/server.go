package application

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/template"

	natsd "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/fix50sp2"
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
	}

	if options.NATSEmbeded {
		s.natsServer, err = natsd.NewServer(&natsd.Options{})
		if err != nil {
			return nil, err
		}

		s.natsServer.Start()
	}

	natsOptions := []nats.Option{
		nats.DontRandomize(),
		nats.RetryOnFailedConnect(true),
	}
	s.natsConn, err = nats.Connect(options.NATSURL, natsOptions...)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

type NewOrderSingleNatsSubject struct {
	Symbol string
	Side   string
	Type   string
}

// Acceptor is the session.Handler forwarding the NewOrderSingle it receives to
// NATS and answering them with a new ExecutionReport.
type Acceptor struct {
	session.BaseHandler

	natsConn   *nats.Conn
	natsServer *natsd.Server

	NatsOrderSubject *template.Template
	Logger           *zerolog.Logger
}

func (app *Acceptor) Close() {
	app.natsConn.Close()

	if app.natsServer != nil {
		app.natsServer.Shutdown()
	}
}

// Notification of a session successfully logging on.
func (app *Acceptor) OnLogon(ctx context.Context, s *session.Session) {
	app.Logger.Debug().Msgf("Logon: %s", s.ID())
}

// Notification of a session logging off or disconnecting.
func (app *Acceptor) OnLogout(ctx context.Context, s *session.Session, reason error) {
	app.Logger.Debug().Msgf("Logout: %s", s.ID())
}

// Notification of app message being received from target.
func (app *Acceptor) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	app.Logger.Debug().Msgf("<- Message received from app %s", s.ID())

	return fix50sp2.Dispatch(m, &dispatcher{Acceptor: app, ctx: ctx, session: s, message: m})
}

// dispatcher handles the message being received on a session.
type dispatcher struct {
	fix50sp2.UnimplementedHandler
	*Acceptor

	ctx     context.Context
	session *session.Session
	message *fix.Message
}

func (d *dispatcher) OnNewOrderSingle(order *fix50sp2.NewOrderSingle) error {
	clOrdID, _ := order.ClOrdID()

	symbol, ok := order.Symbol()
	if !ok {
		return missingField(clOrdID, fix50sp2.TagSymbol)
	}

	side, ok := order.Side()
	if !ok {
		return missingField(clOrdID, fix50sp2.TagSide)
	}

	ordType, ok := order.OrdType()
	if !ok {
		return missingField(clOrdID, fix50sp2.TagOrdType)
	}

	sideString, _ := dict.Search(dict.OrderSides, string(side))
	typeString, _ := dict.Search(dict.OrderTypes, string(ordType))

	buf := bytes.NewBuffer([]byte{})
	subj := NewOrderSingleNatsSubject{
		Symbol: symbol,
		Side:   string(sideString),
		Type:   string(typeString),
	}

	if err := d.NatsOrderSubject.Execute(buf, subj); err != nil {
		return applicationError(clOrdID, err)
	}

	if err := d.natsConn.Publish(buf.String(), rawMessage(d.message)); err != nil {
		return applicationError(clOrdID, err)
	}

	if err := d.sendExecutionReport(order, side, fix50sp2.OrdStatusNew); err != nil {
		return applicationError(clOrdID, err)
	}

	return nil
}

func (d *dispatcher) sendExecutionReport(order *fix50sp2.NewOrderSingle, side fix50sp2.Side, status fix50sp2.OrdStatus) error {
	clOrdID, _ := order.ClOrdID()

	orderQty := decimal.Zero
	if v, ok := order.OrderQtyBytes(); ok {
		var err error
		if orderQty, err = decimal.NewFromString(string(v)); err != nil {
			return fmt.Errorf("OrderQty: %w", err)
		}
	}

	report := fix50sp2.BuildExecutionReport(d.session.ID().BeginString).
		SetOrderID(clOrdID).
		SetExecID("0").
		SetExecType(fix50sp2.ExecType(status)).
		SetOrdStatus(status).
		SetSide(side).
		SetClOrdID(clOrdID)
	report.SetString(fix50sp2.TagLeavesQty, orderQty.StringFixed(2))
	report.SetString(fix50sp2.TagCumQty, decimal.Zero.StringFixed(2))
	report.SetString(fix50sp2.TagOrderQty, orderQty.StringFixed(2))

	_, err := d.session.Send(d.ctx, report)
	return err
}

func missingField(clOrdID string, tag int) error {
	return &session.BusinessRejectError{
		Reason: session.BusinessRejectConditionallyRequiredMissing,
		RefID:  clOrdID,
		Text:   "Required tag missing: " + strconv.Itoa(tag),
	}
}

func applicationError(clOrdID string, err error) error {
	return &session.BusinessRejectError{
		Reason: session.BusinessRejectApplicationNotAvailable,
		RefID:  clOrdID,
		Text:   err.Error(),
	}
}

// rawMessage rebuilds the tag=value SOH-delimited frame of m.
func rawMessage(m *fix.Message) []byte {
	var b []byte
	for _, f := range m.Fields() {
		b = strconv.AppendInt(b, int64(f.Tag), 10)
		b = append(b, '=')
		b = append(b, f.Value...)
		b = append(b, 0x01)
	}
	return b
}
