package application

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/google/uuid"
	natsd "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/fix50sp2"
)

type AcceptorOptions struct {
	NATSEmbedded     bool
	NATSURL          string
	NATSOrderSubject string
}

// NewAcceptor returns an Acceptor connected to NATS, starting the embedded
// NATS server first if asked to.
func NewAcceptor(options *AcceptorOptions) (*Acceptor, error) {
	tpl, err := template.New("fix").Parse(options.NATSOrderSubject)
	if err != nil {
		return nil, err
	}

	s := Acceptor{
		NatsOrderSubject: tpl,
	}

	if options.NATSEmbedded {
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
		if s.natsServer != nil {
			s.natsServer.Shutdown()
		}
		return nil, err
	}

	return &s, nil
}

// NewOrderSingleNatsSubject is what the NATS order subject template is
// executed with: tokens of a subject, without '.', '*', '>' or spaces.
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
	app.Logger.Debug().Msgf("Logout: %s: %v", s.ID(), reason)
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
		return missingField(fix50sp2.TagSymbol)
	}

	side, ok := order.Side()
	if !ok {
		return missingField(fix50sp2.TagSide)
	}

	ordType, ok := order.OrdType()
	if !ok {
		return missingField(fix50sp2.TagOrdType)
	}

	subj := NewOrderSingleNatsSubject{
		Symbol: symbol,
		Side:   name(dict.OrderSidesReversed, side),
		Type:   name(dict.OrderTypesReversed, ordType),
	}

	if !subjectToken(subj.Symbol) {
		return &session.BusinessRejectError{
			Reason: session.BusinessRejectUnknownSecurity,
			RefID:  clOrdID,
			Text:   fmt.Sprintf("Symbol %q can't be a NATS subject token", subj.Symbol),
		}
	}

	buf := bytes.NewBuffer([]byte{})
	if err := d.NatsOrderSubject.Execute(buf, subj); err != nil {
		return applicationError(clOrdID, err)
	}

	if err := d.natsConn.Publish(buf.String(), rawMessage(d.message)); err != nil {
		return applicationError(clOrdID, err)
	}

	report := buildExecutionReport(d.session.ID().BeginString, order, fix50sp2.OrdStatusNew)
	if _, err := d.session.Send(d.ctx, report); err != nil {
		return applicationError(clOrdID, err)
	}

	return nil
}

// name returns the name names gives v, v itself if none.
func name[T ~string](names map[T]string, v T) string {
	if n, ok := names[v]; ok {
		return n
	}

	return string(v)
}

// subjectToken reports whether s can be a token of a NATS subject.
func subjectToken(s string) bool {
	return len(s) > 0 && strings.IndexFunc(s, func(r rune) bool {
		return r == '.' || r == '*' || r == '>' || unicode.IsSpace(r) || unicode.IsControl(r)
	}) < 0
}

// buildExecutionReport returns the ExecutionReport acknowledging order with
// status: it echoes its identifiers, side, symbol, type, quantity and price,
// the latter two as they were sent, and has an ExecID of its own.
func buildExecutionReport(beginString string, order *fix50sp2.NewOrderSingle, status fix50sp2.OrdStatus) *fix50sp2.ExecutionReportBuilder {
	clOrdID, _ := order.ClOrdID()
	side, _ := order.Side()
	symbol, _ := order.Symbol()
	ordType, _ := order.OrdType()

	orderQty := []byte("0")
	if v, ok := order.OrderQtyBytes(); ok {
		orderQty = v
	}

	report := fix50sp2.BuildExecutionReport(beginString).
		SetOrderID(clOrdID).
		SetExecID(uuid.NewString()).
		SetExecType(fix50sp2.ExecType(status)).
		SetOrdStatus(status).
		SetSide(side).
		SetClOrdID(clOrdID).
		SetSymbol(symbol).
		SetOrdType(ordType)
	report.SetString(fix50sp2.TagOrderQty, string(orderQty))
	report.SetString(fix50sp2.TagLeavesQty, string(orderQty))
	report.SetString(fix50sp2.TagCumQty, "0")
	// FIX 4.3 and 4.4 require AvgPx: nothing was executed yet.
	report.SetString(fix50sp2.TagAvgPx, "0")

	if price, ok := order.PriceBytes(); ok {
		report.SetString(fix50sp2.TagPrice, string(price))
	}

	return report
}

// missingField is the Reject of a NewOrderSingle lacking tag: Symbol, which
// the dictionaries do not require, or another when the session skips the
// requiredTags check.
func missingField(tag int) error {
	return &session.RejectError{
		Reason: dictionary.ReasonRequiredTagMissing,
		Tag:    tag,
		Text:   "Required tag missing",
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
