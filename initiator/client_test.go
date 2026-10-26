package initiator_test

import (
	"bytes"
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
	fixerrors "sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/initiator/initiatortest"
)

// broker answers every NewOrderSingle with reports ExecutionReports, then
// logs out if logout is set.
type broker struct {
	session.BaseHandler
	reports int
	logout  bool
}

func (b broker) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	order := fix50sp2.WrapNewOrderSingle(m)
	clOrdID, _ := order.ClOrdID()

	for range max(b.reports, 1) {
		if _, err := s.Send(ctx, executionReport(s, clOrdID)); err != nil {
			return err
		}
	}

	if b.logout {
		go func() { _ = s.Logout(context.Background(), "bye") }()
	}

	return nil
}

func executionReport(s *session.Session, clOrdID string) *fix50sp2.ExecutionReportBuilder {
	return fix50sp2.BuildExecutionReport(s.ID().BeginString).
		SetOrderID(clOrdID).
		SetClOrdID(clOrdID).
		SetExecID("1").
		SetExecType(fix50sp2.ExecTypeNew).
		SetOrdStatus(fix50sp2.OrdStatusNew).
		SetSide(fix50sp2.SideBuy).
		SetLeavesQty(1).
		SetCumQty(0).
		SetAvgPx(0)
}

// hasty sends reports ExecutionReports as soon as its session logs on, then
// logs out.
type hasty struct {
	session.BaseHandler
	reports int
}

func (h hasty) OnLogon(ctx context.Context, s *session.Session) {
	for range h.reports {
		if _, err := s.Send(ctx, executionReport(s, "ORD-1")); err != nil {
			return
		}
	}

	go func() { _ = s.Logout(context.Background(), "bye") }()
}

// gatekeeper refuses every Logon.
type gatekeeper struct {
	session.BaseHandler
}

func (gatekeeper) Authenticate(context.Context, *session.LogonRequest) error {
	return errors.New("unknown user")
}

// refuser rejects every application message with err.
type refuser struct {
	session.BaseHandler
	err error
}

func (r refuser) FromApp(context.Context, *session.Session, *fix.Message, session.Inbound) error {
	return r.err
}

// streamer answers a NewOrderSingle with reports ExecutionReports, one
// every period.
type streamer struct {
	session.BaseHandler
	reports int
	period  time.Duration
}

func (st streamer) FromApp(_ context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	clOrdID, _ := fix50sp2.WrapNewOrderSingle(m).ClOrdID()

	go func() {
		for range st.reports {
			time.Sleep(st.period)
			if _, err := s.Send(context.Background(), executionReport(s, clOrdID)); err != nil {
				return
			}
		}
	}()

	return nil
}

// stranger answers a NewOrderSingle with a BusinessMessageReject of a
// message it was not sent, then an ExecutionReport.
type stranger struct {
	session.BaseHandler
}

func (stranger) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	clOrdID, _ := fix50sp2.WrapNewOrderSingle(m).ClOrdID()

	reject := fix50sp2.BuildBusinessMessageReject(s.ID().BeginString).
		SetRefSeqNum(999).
		SetRefMsgType(fix50sp2.RefMsgType(fix50sp2.MsgTypeNewordersingle)).
		SetBusinessRejectReason(fix50sp2.BusinessRejectReasonOther).
		SetText("not yours")
	if _, err := s.Send(ctx, reject); err != nil {
		return err
	}

	_, err := s.Send(ctx, executionReport(s, clOrdID))
	return err
}

// recorder tells of the Logouts its session receives, "logout", and of its
// end, "end".
type recorder struct {
	session.BaseHandler
	events chan string
}

func (r recorder) FromAdmin(_ context.Context, _ *session.Session, m *fix.Message, _ session.Inbound) {
	if string(m.MsgType()) == "5" {
		r.events <- "logout"
	}
}

func (r recorder) OnLogout(context.Context, *session.Session, error) {
	r.events <- "end"
}

func serve(t *testing.T, handler session.Handler) (*config.Initiator, *config.Session) {
	return initiatortest.Serve(t, handler, "FIX.5.0SP2")
}

func newOrder(beginString string) *fix50sp2.NewOrderSingleBuilder {
	return fix50sp2.BuildNewOrderSingle(beginString).
		SetClOrdID("ORD-1").
		SetSide(fix50sp2.SideBuy).
		SetTransactTime(time.Now()).
		SetOrdType(fix50sp2.OrdTypeMarket).
		SetSymbol("EURUSD").
		SetOrderQty(1)
}

// connect connects to the acceptor of initiatorConfig, closing the client
// at the end of t.
func connect(t *testing.T, initiatorConfig *config.Initiator, client *config.Session, opts initiator.Options) *initiator.Client {
	t.Helper()

	c, err := initiator.ConnectSession(context.Background(), initiatorConfig, client, opts)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close(true) })

	return c
}

func TestClient(t *testing.T) {
	initiatorConfig, client := serve(t, broker{})
	ctx := context.Background()

	c := connect(t, initiatorConfig, client, initiator.Options{
		MsgTypes: []string{string(fix50sp2.MsgTypeExecutionreport)},
	})

	if _, err := c.Send(ctx, newOrder(c.BeginString())); err != nil {
		t.Fatal(err)
	}

	select {
	case m := <-c.Messages():
		report := fix50sp2.WrapExecutionReport(m)
		if id, _ := report.ClOrdID(); id != "ORD-1" {
			t.Errorf("ClOrdID = %q, want ORD-1", id)
		}
		if status, _ := report.OrdStatus(); status != fix50sp2.OrdStatusNew {
			t.Errorf("OrdStatus = %q, want %q", status, fix50sp2.OrdStatusNew)
		}
	case <-c.Done():
		t.Fatalf("logged out: %v", c.Err())
	case <-time.After(5 * time.Second):
		t.Fatal("no ExecutionReport")
	}
}

// TestConnect checks Connect connects the session of the current context,
// and that a negative Buffer is taken for the default.
func TestConnect(t *testing.T) {
	serve(t, broker{})

	c, err := initiator.Connect(context.Background(), initiator.Options{Buffer: -1})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(true)

	if err := c.Await(context.Background(), initiator.AwaitOptions{Timeout: 5 * time.Second}, sendThen(t, c, newOrder(c.BeginString()), 1)); err != nil {
		t.Fatal(err)
	}
}

// TestConnectNoSession checks Connect refuses a context without sessions.
func TestConnectNoSession(t *testing.T) {
	serve(t, broker{})
	config.GetConfig().Contexts[1].Sessions = nil

	if _, err := initiator.Connect(context.Background(), initiator.Options{}); !errors.Is(err, fixerrors.ErrConfigContextNoSession) {
		t.Errorf("Connect() = %v, want %v", err, fixerrors.ErrConfigContextNoSession)
	}
}

// sendThen sends m, then returns a handler done after n messages.
func sendThen(t *testing.T, c *initiator.Client, m session.Outgoing, n int) func(*fix.Message) (bool, error) {
	t.Helper()

	if _, err := c.Send(context.Background(), m); err != nil {
		t.Fatal(err)
	}

	return func(*fix.Message) (bool, error) {
		n--
		return n == 0, nil
	}
}

// TestClientVersions checks an order and its ExecutionReport go through
// sessions of every FIX version fix supports, both sides validating what
// they receive with the version's dictionary.
func TestClientVersions(t *testing.T) {
	for _, version := range []string{"FIX.5.0SP2", "FIX.4.4", "FIX.4.3"} {
		t.Run(version, func(t *testing.T) {
			initiatorConfig, client := initiatortest.Serve(t, broker{}, version)
			ctx := context.Background()

			c := connect(t, initiatorConfig, client, initiator.Options{
				MsgTypes: []string{string(fix50sp2.MsgTypeExecutionreport)},
			})

			order := newOrder(c.BeginString())
			if client.Requires(string(fix50sp2.MsgTypeNewordersingle), fix50sp2.TagHandlInst) {
				order.SetHandlInst(fix50sp2.HandlInstAutomatedExecutionOrderPrivateNoBrokerIntervention)
			}

			if _, err := c.Send(ctx, order); err != nil {
				t.Fatal(err)
			}

			select {
			case m := <-c.Messages():
				if string(m.MsgType()) != string(fix50sp2.MsgTypeExecutionreport) {
					t.Fatalf("received MsgType %s, want an ExecutionReport: %s", m.MsgType(), m.Get(58))
				}
				if got := string(m.BeginString()); got != client.BeginString {
					t.Errorf("BeginString = %s, want %s", got, client.BeginString)
				}
			case <-c.Done():
				t.Fatalf("logged out: %v", c.Err())
			case <-time.After(5 * time.Second):
				t.Fatal("no ExecutionReport")
			}
		})
	}
}

// TestClientSendChecks checks Send refuses a message the session's FIX
// version does not define, or that lacks what it requires.
func TestClientSendChecks(t *testing.T) {
	initiatorConfig, client := initiatortest.Serve(t, broker{}, "FIX.4.3")
	ctx := context.Background()

	c := connect(t, initiatorConfig, client, initiator.Options{})

	// FIX 4.3 requires HandlInst (21).
	if _, err := c.Send(ctx, newOrder(c.BeginString())); err == nil || !strings.Contains(err.Error(), "HandlInst(21)") {
		t.Errorf("Send(NewOrderSingle without HandlInst) = %v, want HandlInst required", err)
	}

	// OrderOrigination (1724) came with FIX 5.0 SP1.
	order := newOrder(c.BeginString()).SetHandlInst(fix50sp2.HandlInstAutomatedExecutionOrderPrivateNoBrokerIntervention)
	order.SetString(fix50sp2.TagOrderOrigination, "1")
	if _, err := c.Send(ctx, order); err == nil || !strings.Contains(err.Error(), "OrderOrigination") {
		t.Errorf("Send(NewOrderSingle with OrderOrigination) = %v, want OrderOrigination undefined", err)
	}
}

// TestClientMessagesBeforeLogout checks that a reader selecting on Messages
// and Done gets every message received before the counterparty logged out,
// however slow it is.
func TestClientMessagesBeforeLogout(t *testing.T) {
	const reports = 5

	initiatorConfig, client := serve(t, broker{reports: reports, logout: true})
	ctx := context.Background()

	c := connect(t, initiatorConfig, client, initiator.Options{})

	if _, err := c.Send(ctx, newOrder(c.BeginString())); err != nil {
		t.Fatal(err)
	}

	received := 0
LOOP:
	for {
		select {
		case <-c.Messages():
			received++
			time.Sleep(20 * time.Millisecond)
		case <-c.Done():
			break LOOP
		case <-time.After(5 * time.Second):
			t.Fatal("no logout")
		}
	}

	if received != reports {
		t.Errorf("received %d ExecutionReports before the logout, want %d", received, reports)
	}
	if c.Err() == nil {
		t.Error("Err() = nil after the counterparty logged out")
	}
}

// TestClientLogoutRightAfterLogon checks that a counterparty logging out
// right after logging on, before OnLogon may have run, is not taken for a
// refused Logon, and that what it sent in between is delivered.
func TestClientLogoutRightAfterLogon(t *testing.T) {
	const reports = 3

	initiatorConfig, client := serve(t, hasty{reports: reports})

	c := connect(t, initiatorConfig, client, initiator.Options{})

	received := 0
LOOP:
	for {
		select {
		case <-c.Messages():
			received++
		case <-c.Done():
			break LOOP
		case <-time.After(5 * time.Second):
			t.Fatal("no logout")
		}
	}

	if received != reports {
		t.Errorf("received %d ExecutionReports before the logout, want %d", received, reports)
	}
	if err := c.Err(); err == nil || errors.Is(err, fixerrors.ErrFixLogonRefused) {
		t.Errorf("Err() = %v, want the counterparty's logout", err)
	}
}

// TestClientLogonRefused checks that Connect fails at once when the
// counterparty refuses the Logon, rather than at the logon timeout.
func TestClientLogonRefused(t *testing.T) {
	initiatorConfig, client := serve(t, gatekeeper{})

	start := time.Now()
	_, err := initiator.ConnectSession(context.Background(), initiatorConfig, client, initiator.Options{})
	if !errors.Is(err, fixerrors.ErrFixLogonRefused) {
		t.Fatalf("ConnectSession() = %v, want a refused logon", err)
	}
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Errorf("ConnectSession() took %v, the logon timeout is 10s", elapsed)
	}
}

// peer is a fake acceptor's end of a connection, talking FIXT.1.1 by hand.
type peer struct {
	t    *testing.T
	conn net.Conn
	in   chan fix.Message
	seq  int
}

// field is a tag and its value.
type field struct {
	tag   int
	value string
}

// send sends a message of msgType with fields, numbered after the last.
func (p *peer) send(msgType string, fields ...field) {
	p.seq++
	b := fix.NewBuilder(session.FIXT11, msgType).
		SetString(49, "BIGCORP").
		SetString(56, "smallcorp").
		SetString(57, "john").
		SetInt(34, p.seq).
		SetString(52, time.Now().UTC().Format("20060102-15:04:05.000"))
	for _, f := range fields {
		b.SetString(f.tag, f.value)
	}
	if _, err := p.conn.Write(b.Build()); err != nil {
		p.t.Error(err)
	}
}

// expect returns the next message of msgType received, skipping others.
func (p *peer) expect(msgType string) fix.Message {
	for {
		select {
		case m, ok := <-p.in:
			if !ok {
				p.t.Fatalf("connection closed waiting for MsgType %s", msgType)
			}
			if string(m.MsgType()) == msgType {
				return m
			}
		case <-time.After(5 * time.Second):
			p.t.Fatalf("no MsgType %s received", msgType)
		}
	}
}

// logon answers the client's Logon.
func (p *peer) logon() {
	p.expect("A")
	p.send("A", field{98, "0"}, field{108, "5"}, field{1137, "9"})
}

// executionReport sends an ExecutionReport of ORD-1.
func (p *peer) executionReport() {
	p.send("8", field{37, "ORD-1"}, field{17, "E1"}, field{150, "0"}, field{39, "0"}, field{54, "1"}, field{151, "1"}, field{14, "0"})
}

// fake runs a fake acceptor whose connections script handles, and returns
// the initiator and session configurations to connect to it.
func fake(t *testing.T, script func(*peer)) (*config.Initiator, *config.Session) {
	logger := zerolog.Nop()
	config.SetLogger(&logger)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			t.Cleanup(func() { conn.Close() })

			p := &peer{t: t, conn: conn, in: make(chan fix.Message, 64)}
			go func() {
				defer close(p.in)
				framer := fix.NewFramer(conn, fix.FramerOptions{})
				for {
					frame, err := framer.Next()
					if err != nil {
						return
					}
					m, err := fix.Parse(bytes.Clone(frame))
					if err != nil {
						continue
					}
					p.in <- m
				}
			}()
			go script(p)
		}
	}()

	client := &config.Session{
		Name:             "client",
		BeginString:      session.FIXT11,
		SenderCompID:     "smallcorp",
		SenderSubID:      "john",
		TargetCompID:     "BIGCORP",
		DefaultApplVerID: "FIX.5.0SP2",
		HeartBtInt:       5 * time.Second,
		LogonTimeout:     10 * time.Second,
		// The fake never answers the Logout.
		LogoutTimeout: 100 * time.Millisecond,
	}
	initiatorConfig := &config.Initiator{Name: "client", Addresses: []string{l.Addr().String()}, ReconnectMin: 30 * time.Second}

	return initiatorConfig, client
}

// TestClientLogonRejected checks that a Logon answer the session itself
// rejects, here for its invalid NextExpectedMsgSeqNum, ends Connect rather
// than being taken for an accepted Logon and retried until the timeout.
func TestClientLogonRejected(t *testing.T) {
	initiatorConfig, client := fake(t, func(p *peer) {
		p.expect("A")
		p.send("A", field{98, "0"}, field{108, "5"}, field{1137, "9"}, field{789, "100"})
	})
	client.NextExpectedMsgSeqNum = true

	start := time.Now()
	_, err := initiator.ConnectSession(context.Background(), initiatorConfig, client, initiator.Options{})
	if !errors.Is(err, fixerrors.ErrFixLogonRefused) {
		t.Fatalf("ConnectSession() = %v, want a refused logon", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("ConnectSession() took %v, the logon timeout is 10s", elapsed)
	}
}

// TestClientLogonRejectedByReject checks that a Logon the counterparty
// answers with a Reject ends Connect at once with the Reject's Text.
func TestClientLogonRejectedByReject(t *testing.T) {
	initiatorConfig, client := fake(t, func(p *peer) {
		p.expect("A")
		p.send("3", field{45, "1"}, field{373, "99"}, field{58, "bad credentials"})
	})

	start := time.Now()
	_, err := initiator.ConnectSession(context.Background(), initiatorConfig, client, initiator.Options{})
	if !errors.Is(err, fixerrors.ErrFixLogonRefused) || !strings.Contains(err.Error(), "bad credentials") {
		t.Fatalf("ConnectSession() = %v, want a logon refused with its Text", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("ConnectSession() took %v, the logon timeout is 10s", elapsed)
	}
}

// TestClientLogonConnectionClosed checks that a counterparty closing the
// connection on a Logon ends Connect at once, not at the next dial.
func TestClientLogonConnectionClosed(t *testing.T) {
	initiatorConfig, client := fake(t, func(p *peer) {
		p.expect("A")
		p.conn.Close()
	})

	start := time.Now()
	_, err := initiator.ConnectSession(context.Background(), initiatorConfig, client, initiator.Options{})
	if !errors.Is(err, fixerrors.ErrFixLogonRefused) || !strings.Contains(err.Error(), "closed the connection") {
		t.Fatalf("ConnectSession() = %v, want a logon refused by a closed connection", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("ConnectSession() took %v, reconnecting waits 30s", elapsed)
	}
}

// TestClientDialError checks that a Connect that could not dial tells why.
func TestClientDialError(t *testing.T) {
	logger := zerolog.Nop()
	config.SetLogger(&logger)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := l.Addr().String()
	l.Close()

	client := &config.Session{
		Name:             "client",
		BeginString:      session.FIXT11,
		SenderCompID:     "smallcorp",
		TargetCompID:     "BIGCORP",
		DefaultApplVerID: "FIX.5.0SP2",
		LogonTimeout:     500 * time.Millisecond,
	}
	initiatorConfig := &config.Initiator{Name: "client", Addresses: []string{address}}

	_, err = initiator.ConnectSession(context.Background(), initiatorConfig, client, initiator.Options{})
	if !errors.Is(err, fixerrors.ErrConnectionTimeout) || !strings.Contains(err.Error(), "refused") {
		t.Fatalf("ConnectSession() = %v, want a timeout telling the connection was refused", err)
	}
}

// TestClientClose checks Close logs out or drops the connection, as asked,
// and stops the session either way.
func TestClientClose(t *testing.T) {
	for _, logout := range []bool{true, false} {
		t.Run(map[bool]string{true: "logout", false: "drop"}[logout], func(t *testing.T) {
			events := make(chan string, 4)
			initiatorConfig, client := serve(t, recorder{events: events})

			c, err := initiator.ConnectSession(context.Background(), initiatorConfig, client, initiator.Options{})
			if err != nil {
				t.Fatal(err)
			}

			c.Close(logout)
			select {
			case <-c.Done():
			default:
				t.Fatal("Done not closed once Close returned")
			}
			if err := c.Err(); err != nil {
				t.Errorf("Err() = %v after Close", err)
			}

			select {
			case event := <-events:
				if want := map[bool]string{true: "logout", false: "end"}[logout]; event != want {
					t.Errorf("the acceptor's session got %q first, want %q", event, want)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("the acceptor's session did not end")
			}

			// Closing again returns at once.
			c.Close(true)
		})
	}
}

// TestAwait checks Await hands the responses, and returns once the handler
// is done.
func TestAwait(t *testing.T) {
	initiatorConfig, client := serve(t, broker{reports: 3})
	c := connect(t, initiatorConfig, client, initiator.Options{})

	if _, err := c.Send(context.Background(), newOrder(c.BeginString())); err != nil {
		t.Fatal(err)
	}

	received := 0
	err := c.Await(context.Background(), initiator.AwaitOptions{What: "execution report", Timeout: 5 * time.Second}, func(m *fix.Message) (bool, error) {
		if string(m.MsgType()) != string(fix50sp2.MsgTypeExecutionreport) {
			t.Errorf("MsgType %s, want an ExecutionReport", m.MsgType())
		}
		received++
		return received == 3, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if received != 3 {
		t.Errorf("received %d, want 3", received)
	}
}

// TestAwaitReject checks a Reject or BusinessMessageReject of the request
// fails Await, telling why.
func TestAwaitReject(t *testing.T) {
	for _, tc := range []struct {
		name    string
		err     error
		msgType string
		reason  int
	}{
		{"Reject", &session.RejectError{Reason: dictionary.ReasonRequiredTagMissing, Tag: 55, Text: "no symbol"}, "3", int(dictionary.ReasonRequiredTagMissing)},
		{"BusinessMessageReject", &session.BusinessRejectError{Reason: session.BusinessRejectUnknownSecurity, RefID: "ORD-1", Text: "no symbol"}, "j", session.BusinessRejectUnknownSecurity},
	} {
		t.Run(tc.name, func(t *testing.T) {
			initiatorConfig, client := serve(t, refuser{err: tc.err})
			c := connect(t, initiatorConfig, client, initiator.Options{
				MsgTypes: []string{string(fix50sp2.MsgTypeExecutionreport)},
			})

			receipt, err := c.Send(context.Background(), newOrder(c.BeginString()))
			if err != nil {
				t.Fatal(err)
			}

			err = c.Await(context.Background(), initiator.AwaitOptions{Timeout: 5 * time.Second}, func(m *fix.Message) (bool, error) {
				t.Errorf("handled MsgType %s", m.MsgType())
				return true, nil
			})

			var reject *initiator.RejectError
			if !errors.As(err, &reject) || !errors.Is(err, fixerrors.ErrFixRejected) {
				t.Fatalf("Await() = %v, want a RejectError", err)
			}
			if reject.MsgType != tc.msgType || reject.Reason != tc.reason || reject.RefSeqNum != receipt.Seq() || reject.RefMsgType != "D" {
				t.Errorf("RejectError = %+v, want a %s of MsgSeqNum %d, reason %d", reject, tc.msgType, receipt.Seq(), tc.reason)
			}
			if !strings.Contains(err.Error(), "no symbol") {
				t.Errorf("Await() = %v, want the reject's Text", err)
			}
		})
	}
}

// TestAwaitRejectOfAnother checks a BusinessMessageReject of a message not
// sent is skipped.
func TestAwaitRejectOfAnother(t *testing.T) {
	initiatorConfig, client := serve(t, stranger{})
	c := connect(t, initiatorConfig, client, initiator.Options{})

	if _, err := c.Send(context.Background(), newOrder(c.BeginString())); err != nil {
		t.Fatal(err)
	}

	err := c.Await(context.Background(), initiator.AwaitOptions{Timeout: 5 * time.Second}, func(m *fix.Message) (bool, error) {
		if string(m.MsgType()) != string(fix50sp2.MsgTypeExecutionreport) {
			t.Errorf("handled MsgType %s, want an ExecutionReport", m.MsgType())
		}
		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestAwaitSessionRejectOfAnother checks a Reject of a message not sent is
// skipped.
func TestAwaitSessionRejectOfAnother(t *testing.T) {
	initiatorConfig, client := fake(t, func(p *peer) {
		p.logon()
		p.expect("D")
		p.send("3", field{45, "999"}, field{373, "1"}, field{58, "not yours"})
		p.executionReport()
	})
	c := connect(t, initiatorConfig, client, initiator.Options{})

	err := c.Await(context.Background(), initiator.AwaitOptions{Timeout: 5 * time.Second}, sendThen(t, c, newOrder(c.BeginString()), 1))
	if err != nil {
		t.Fatal(err)
	}
}

// TestAwaitRejectOrder checks a Reject is handled after the application
// messages received before it, however slow the handler, though the session
// handles Rejects on a goroutine of their own.
func TestAwaitRejectOrder(t *testing.T) {
	const reports = 5

	initiatorConfig, client := fake(t, func(p *peer) {
		p.logon()
		order := p.expect("D")
		for range reports {
			p.executionReport()
		}
		p.send("3", field{45, string(order.Get(34))}, field{373, "5"}, field{58, "too late"})
	})
	c := connect(t, initiatorConfig, client, initiator.Options{Buffer: 1})

	if _, err := c.Send(context.Background(), newOrder(c.BeginString())); err != nil {
		t.Fatal(err)
	}

	received := 0
	err := c.Await(context.Background(), initiator.AwaitOptions{Timeout: 5 * time.Second}, func(*fix.Message) (bool, error) {
		received++
		time.Sleep(20 * time.Millisecond)
		return false, nil
	})
	if !errors.Is(err, fixerrors.ErrFixRejected) || !strings.Contains(err.Error(), "too late") {
		t.Fatalf("Await() = %v, want the Reject", err)
	}
	if received != reports {
		t.Errorf("handled %d ExecutionReports before the Reject, want %d", received, reports)
	}
}

// TestAwaitTimeout checks Await fails once its timeout elapsed.
func TestAwaitTimeout(t *testing.T) {
	initiatorConfig, client := serve(t, session.BaseHandler{})
	c := connect(t, initiatorConfig, client, initiator.Options{})

	start := time.Now()
	err := c.Await(context.Background(), initiator.AwaitOptions{What: "execution report", Timeout: 200 * time.Millisecond}, sendThen(t, c, newOrder(c.BeginString()), 1))
	if !errors.Is(err, fixerrors.ErrResponseTimeout) || !strings.Contains(err.Error(), "execution report") {
		t.Fatalf("Await() = %v, want a timeout", err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("Await() took %v, the timeout is 200ms", elapsed)
	}
}

// TestAwaitLogout checks Await fails once the counterparty logged out, after
// handling what it sent before.
func TestAwaitLogout(t *testing.T) {
	initiatorConfig, client := serve(t, broker{reports: 3, logout: true})
	c := connect(t, initiatorConfig, client, initiator.Options{})

	if _, err := c.Send(context.Background(), newOrder(c.BeginString())); err != nil {
		t.Fatal(err)
	}

	received := 0
	err := c.Await(context.Background(), initiator.AwaitOptions{Timeout: 5 * time.Second}, func(*fix.Message) (bool, error) {
		received++
		return false, nil
	})
	if !errors.Is(err, fixerrors.ErrFixLogout) || c.Err() == nil || !strings.Contains(err.Error(), c.Err().Error()) {
		t.Fatalf("Await() = %v, want the logout and Err() (%v)", err, c.Err())
	}
	if received != 3 {
		t.Errorf("handled %d ExecutionReports before the logout, want 3", received)
	}
	if _, err := c.Send(context.Background(), newOrder(c.BeginString())); err == nil {
		t.Error("Send() after the logout succeeded")
	}
}

// TestAwaitCancel checks Await returns nil once ctx is done.
func TestAwaitCancel(t *testing.T) {
	initiatorConfig, client := serve(t, session.BaseHandler{})
	c := connect(t, initiatorConfig, client, initiator.Options{})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := c.Await(ctx, initiator.AwaitOptions{Timeout: 5 * time.Second}, sendThen(t, c, newOrder(c.BeginString()), 1)); err != nil {
		t.Fatalf("Await() = %v, want nil", err)
	}
}

// TestAwaitTimeoutModes checks what the timeout bounds: a stream of 6
// responses, one every 50ms, outlasts a timeout of 150ms, unless it bounds
// the first response or each one only.
func TestAwaitTimeoutModes(t *testing.T) {
	const reports = 6

	for _, tc := range []struct {
		name     string
		mode     initiator.TimeoutMode
		received int
	}{
		{"total", initiator.TimeoutTotal, 2},
		{"idle", initiator.TimeoutIdle, reports},
		{"first", initiator.TimeoutFirst, reports},
	} {
		t.Run(tc.name, func(t *testing.T) {
			initiatorConfig, client := serve(t, streamer{reports: reports, period: 50 * time.Millisecond})
			c := connect(t, initiatorConfig, client, initiator.Options{})

			// The stream is cancelled after the last response.
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			send := sendThen(t, c, newOrder(c.BeginString()), reports)
			received := 0
			err := c.Await(ctx, initiator.AwaitOptions{Timeout: 150 * time.Millisecond, Mode: tc.mode}, func(m *fix.Message) (bool, error) {
				received++
				if done, _ := send(m); done && tc.mode == initiator.TimeoutFirst {
					cancel()
				}
				return false, nil
			})

			switch tc.mode {
			case initiator.TimeoutFirst:
				if err != nil {
					t.Errorf("Await() = %v, want nil once cancelled", err)
				}
			default:
				if !errors.Is(err, fixerrors.ErrResponseTimeout) {
					t.Errorf("Await() = %v, want a timeout", err)
				}
			}
			if tc.mode == initiator.TimeoutTotal && received >= reports || tc.mode != initiator.TimeoutTotal && received != reports {
				t.Errorf("handled %d responses", received)
			}
		})
	}
}

// TestAwaitTick checks Await calls OnTick at every tick, from the start and
// however frequent the responses, and fails on a Reject of what it sent.
func TestAwaitTick(t *testing.T) {
	initiatorConfig, client := serve(t, broker{reports: 20})
	c := connect(t, initiatorConfig, client, initiator.Options{})

	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	var ticks atomic.Int32
	received := 0
	err := c.Await(context.Background(), initiator.AwaitOptions{
		Timeout: 5 * time.Second,
		Tick:    ticker.C,
		OnTick: func() error {
			ticks.Add(1)
			_, err := c.Send(context.Background(), newOrder(c.BeginString()))
			return err
		},
	}, func(*fix.Message) (bool, error) {
		received++
		time.Sleep(5 * time.Millisecond)
		return received == 50, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n := ticks.Load(); n < 3 {
		t.Errorf("%d ticks, want 3 at least", n)
	}
}

// TestAwaitTickReject checks a Reject of what OnTick sent fails Await.
func TestAwaitTickReject(t *testing.T) {
	initiatorConfig, client := serve(t, refuser{err: &session.BusinessRejectError{Reason: session.BusinessRejectOther, Text: "no"}})
	c := connect(t, initiatorConfig, client, initiator.Options{})

	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	sent := false
	err := c.Await(context.Background(), initiator.AwaitOptions{
		Timeout: 5 * time.Second,
		Tick:    ticker.C,
		OnTick: func() error {
			if sent {
				return nil
			}
			sent = true
			_, err := c.Send(context.Background(), newOrder(c.BeginString()))
			return err
		},
	}, func(*fix.Message) (bool, error) {
		return true, nil
	})
	if !errors.Is(err, fixerrors.ErrFixRejected) {
		t.Fatalf("Await() = %v, want the BusinessMessageReject", err)
	}
}

// TestAwaitSkip checks the messages the handler skips with ErrSkip neither
// stop nor restart the timeout, whatever it bounds, nor count as received.
func TestAwaitSkip(t *testing.T) {
	for _, mode := range []initiator.TimeoutMode{initiator.TimeoutTotal, initiator.TimeoutIdle, initiator.TimeoutFirst} {
		t.Run(map[initiator.TimeoutMode]string{initiator.TimeoutTotal: "total", initiator.TimeoutIdle: "idle", initiator.TimeoutFirst: "first"}[mode], func(t *testing.T) {
			// Unrelated responses, one every 30ms, for 600ms.
			initiatorConfig, client := serve(t, streamer{reports: 20, period: 30 * time.Millisecond})
			c := connect(t, initiatorConfig, client, initiator.Options{})

			if _, err := c.Send(context.Background(), newOrder(c.BeginString())); err != nil {
				t.Fatal(err)
			}

			skipped := 0
			start := time.Now()
			err := c.Await(context.Background(), initiator.AwaitOptions{What: "response", Timeout: 200 * time.Millisecond, Mode: mode}, func(*fix.Message) (bool, error) {
				skipped++
				return false, initiator.ErrSkip
			})
			if !errors.Is(err, fixerrors.ErrResponseTimeout) || !strings.Contains(err.Error(), "(0 received)") {
				t.Fatalf("Await() = %v, want a timeout with none received", err)
			}
			if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
				t.Errorf("Await() took %v, the timeout is 200ms", elapsed)
			}
			if skipped == 0 {
				t.Error("no message skipped")
			}
		})
	}
}

// gapFill sends a PossDup SequenceReset-GapFill numbered seq, to newSeqNo.
func (p *peer) gapFill(seq, newSeqNo int) {
	p.seq = seq - 1
	p.send("4", field{43, "Y"}, field{122, time.Now().UTC().Format("20060102-15:04:05.000")}, field{123, "Y"}, field{36, strconv.Itoa(newSeqNo)})
}

// TestAwaitRejectAfterGapFill checks a Reject is not held for messages
// received ahead of a gap that a GapFill skipped, which are never
// delivered, but still is for those received before it.
func TestAwaitRejectAfterGapFill(t *testing.T) {
	for _, tc := range []struct {
		name string
		// script sends ExecutionReports, of which the client handles
		// handled, and a GapFill, numbered from 2; the Reject follows.
		script  func(*peer)
		handled int
	}{
		{"in sequence", func(p *peer) {
			p.executionReport() // 2
			p.seq = 3
			p.executionReport() // 4, ahead of the gap
			p.gapFill(3, 5)     // skips 4
			p.seq = 4
		}, 1},
		{"ahead of the gap", func(p *peer) {
			p.executionReport() // 2
			p.seq = 4
			p.executionReport() // 5, ahead of the gap
			p.gapFill(4, 6)     // ahead of the gap too, skips 5
			p.seq = 2
			p.executionReport() // 3, which the GapFill follows
			p.seq = 5
		}, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			initiatorConfig, client := fake(t, func(p *peer) {
				p.logon()
				order := p.expect("D")
				tc.script(p)
				p.send("3", field{45, string(order.Get(34))}, field{373, "5"}, field{58, "too late"})
			})
			c := connect(t, initiatorConfig, client, initiator.Options{Buffer: 1})

			if _, err := c.Send(context.Background(), newOrder(c.BeginString())); err != nil {
				t.Fatal(err)
			}

			handled := 0
			err := c.Await(context.Background(), initiator.AwaitOptions{Timeout: 3 * time.Second}, func(*fix.Message) (bool, error) {
				handled++
				time.Sleep(20 * time.Millisecond)
				return false, nil
			})
			if !errors.Is(err, fixerrors.ErrFixRejected) || !strings.Contains(err.Error(), "too late") {
				t.Fatalf("Await() = %v, want the Reject", err)
			}
			if handled != tc.handled {
				t.Errorf("handled %d ExecutionReports before the Reject, want %d", handled, tc.handled)
			}
		})
	}
}

// TestAwaitRejectAfterReset checks the Rejects are ordered and matched
// within the sequence number generation the counterparty's Logon with
// ResetSeqNumFlag started: neither the MsgSeqNums of the messages it
// delivered before nor those of the messages sent before count.
func TestAwaitRejectAfterReset(t *testing.T) {
	const reports = 4

	initiatorConfig, client := fake(t, func(p *peer) {
		p.logon()
		p.expect("D")
		p.expect("D")
		for range 10 {
			p.executionReport() // 2 to 11
		}

		p.seq = 0
		p.send("A", field{98, "0"}, field{108, "5"}, field{141, "Y"}, field{1137, "9"})
		p.expect("A")
		p.executionReport() // 2

		order := p.expect("D")
		for range reports {
			p.executionReport() // 3 to 6
		}
		// Our 3, sent before the reset only: not one of the client's.
		p.send("3", field{45, "3"}, field{373, "5"}, field{58, "stale"})
		p.send("3", field{45, string(order.Get(34))}, field{373, "5"}, field{58, "too late"})
	})
	c := connect(t, initiatorConfig, client, initiator.Options{Buffer: 1})

	// MsgSeqNums 2 and 3.
	for range 2 {
		if _, err := c.Send(context.Background(), newOrder(c.BeginString())); err != nil {
			t.Fatal(err)
		}
	}
	// The reports before the reset, then the first after it.
	for _, n := range []int{10, 1} {
		if err := c.Await(context.Background(), initiator.AwaitOptions{Timeout: 5 * time.Second}, func(*fix.Message) (bool, error) {
			n--
			return n == 0, nil
		}); err != nil {
			t.Fatal(err)
		}
	}

	receipt, err := c.Send(context.Background(), newOrder(c.BeginString()))
	if err != nil {
		t.Fatal(err)
	}

	handled := 0
	err = c.Await(context.Background(), initiator.AwaitOptions{Timeout: 5 * time.Second}, func(*fix.Message) (bool, error) {
		handled++
		time.Sleep(20 * time.Millisecond)
		return false, nil
	})
	var reject *initiator.RejectError
	if !errors.As(err, &reject) || reject.Text != "too late" || reject.RefSeqNum != receipt.Seq() {
		t.Fatalf("Await() = %v, want the Reject of MsgSeqNum %d", err, receipt.Seq())
	}
	if handled != reports {
		t.Errorf("handled %d ExecutionReports before the Reject, want %d", handled, reports)
	}
}
