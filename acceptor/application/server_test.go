package application

import (
	"context"
	"errors"
	"testing"
	"text/template"
	"time"

	natsd "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"

	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/initiator/initiatortest"
	"sylr.dev/fix/v2/internal/fixmsg"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

func TestBuildExecutionReportVersions(t *testing.T) {
	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			order := fix50sp2.BuildNewOrderSingle(session.BeginString).
				SetClOrdID("order-1").
				SetSide(fix50sp2.SideBuy).
				SetTransactTime(time.Now()).
				SetOrdType(fix50sp2.OrdTypeLimit).
				SetSymbol("EURUSD")
			order.SetString(fix50sp2.TagOrderQty, "10.125")
			order.SetString(fix50sp2.TagPrice, "1.08345")

			msg, err := fixmsg.Message(order, session.BeginString, nil)
			if err != nil {
				t.Fatal(err)
			}

			report := buildExecutionReport(session.BeginString, fix50sp2.WrapNewOrderSingle(msg), fix50sp2.OrdStatusNew)
			fixmsgtest.Validate(t, session, report)

			reportMsg, err := fixmsg.Message(report, session.BeginString, nil)
			if err != nil {
				t.Fatal(err)
			}
			for tag, want := range map[int]string{
				fix50sp2.TagSymbol:    "EURUSD",
				fix50sp2.TagOrdType:   string(fix50sp2.OrdTypeLimit),
				fix50sp2.TagPrice:     "1.08345",
				fix50sp2.TagOrderQty:  "10.125",
				fix50sp2.TagLeavesQty: "10.125",
			} {
				if got := string(reportMsg.Get(tag)); got != want {
					t.Errorf("tag %d = %q, want %q", tag, got, want)
				}
			}

			other := buildExecutionReport(session.BeginString, fix50sp2.WrapNewOrderSingle(msg), fix50sp2.OrdStatusNew)
			otherMsg, err := fixmsg.Message(other, session.BeginString, nil)
			if err != nil {
				t.Fatal(err)
			}
			if id := reportMsg.Get(fix50sp2.TagExecID); string(id) == string(otherMsg.Get(fix50sp2.TagExecID)) {
				t.Errorf("two ExecutionReports share ExecID %s", id)
			}
		})
	}
}

func TestSubjectToken(t *testing.T) {
	for symbol, want := range map[string]bool{
		"EURUSD":   true,
		"EUR/USD":  true,
		"":         false,
		"BTC.USD":  false,
		"*":        false,
		">":        false,
		"EUR USD":  false,
		"EUR\tUSD": false,
	} {
		if got := subjectToken(symbol); got != want {
			t.Errorf("subjectToken(%q) = %v, want %v", symbol, got, want)
		}
	}
}

// TestAcceptor checks the acceptor publishes the orders it receives on NATS
// and acknowledges them, and rejects those it cannot.
func TestAcceptor(t *testing.T) {
	server, err := natsd.NewServer(&natsd.Options{Host: "127.0.0.1", Port: -1})
	if err != nil {
		t.Fatal(err)
	}
	server.Start()
	if !server.ReadyForConnections(5 * time.Second) {
		t.Fatal("NATS server not ready")
	}
	conn, err := nats.Connect(server.ClientURL())
	if err != nil {
		t.Fatal(err)
	}

	logger := zerolog.Nop()
	app := &Acceptor{
		natsConn:         conn,
		natsServer:       server,
		NatsOrderSubject: template.Must(template.New("fix").Parse("orders.{{.Symbol}}.{{.Side}}.{{.Type}}")),
		Logger:           &logger,
	}
	defer app.Close()

	published, err := conn.SubscribeSync("orders.>")
	if err != nil {
		t.Fatal(err)
	}

	initiatorConfig, client := initiatortest.Serve(t, app, "FIX.5.0SP2")
	c, err := initiator.ConnectSession(context.Background(), initiatorConfig, client, initiator.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(true)

	order := func(symbol string) *fix50sp2.NewOrderSingleBuilder {
		o := fix50sp2.BuildNewOrderSingle(c.BeginString()).
			SetClOrdID("order-" + symbol).
			SetSide(fix50sp2.SideBuy).
			SetTransactTime(time.Now()).
			SetOrdType(fix50sp2.OrdTypeLimit)
		if symbol != "" {
			o.SetSymbol(symbol)
		}
		o.SetString(fix50sp2.TagOrderQty, "10.125")
		o.SetString(fix50sp2.TagPrice, "1.08345")
		return o
	}
	await := func() (*fix.Message, error) {
		var report *fix.Message
		err := c.Await(context.Background(), initiator.AwaitOptions{Timeout: 5 * time.Second}, func(m *fix.Message) (bool, error) {
			report = m
			return true, nil
		})
		return report, err
	}

	if _, err := c.Send(context.Background(), order("EURUSD")); err != nil {
		t.Fatal(err)
	}
	report, err := await()
	if err != nil {
		t.Fatal(err)
	}
	for tag, want := range map[int]string{
		fix50sp2.TagSymbol:  "EURUSD",
		fix50sp2.TagOrdType: string(fix50sp2.OrdTypeLimit),
		fix50sp2.TagPrice:   "1.08345",
	} {
		if got := string(report.Get(tag)); got != want {
			t.Errorf("ExecutionReport tag %d = %q, want %q", tag, got, want)
		}
	}
	if msg, err := published.NextMsg(5 * time.Second); err != nil {
		t.Error(err)
	} else if msg.Subject != "orders.EURUSD.BUY.LIMIT" {
		t.Errorf("published on %s, want orders.EURUSD.BUY.LIMIT", msg.Subject)
	}

	var reject *initiator.RejectError

	if _, err := c.Send(context.Background(), order("")); err != nil {
		t.Fatal(err)
	}
	if _, err := await(); !errors.As(err, &reject) || reject.MsgType != "3" || reject.Reason != int(dictionary.ReasonRequiredTagMissing) || reject.RefTagID != fix50sp2.TagSymbol {
		t.Errorf("order without Symbol: %v, want a Reject of the missing Symbol", err)
	}

	if _, err := c.Send(context.Background(), order("EUR.USD")); err != nil {
		t.Fatal(err)
	}
	if _, err := await(); !errors.As(err, &reject) || reject.MsgType != "j" {
		t.Errorf("order of EUR.USD: %v, want a BusinessMessageReject", err)
	}
}
