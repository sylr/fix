package initiator_test

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/acceptor"
	"sylr.dev/fix/v2/config"
	fixerrors "sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator"
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
	side, _ := order.Side()

	for range max(b.reports, 1) {
		report := fix50sp2.BuildExecutionReport(s.ID().BeginString).
			SetOrderID(clOrdID).
			SetClOrdID(clOrdID).
			SetExecID("1").
			SetExecType(fix50sp2.ExecTypeNew).
			SetOrdStatus(fix50sp2.OrdStatusNew).
			SetSide(side).
			SetLeavesQty(1).
			SetCumQty(0).
			SetAvgPx(0)

		if _, err := s.Send(ctx, report); err != nil {
			return err
		}
	}

	if b.logout {
		go s.Logout(context.Background(), "bye")
	}

	return nil
}

// hasty sends reports ExecutionReports as soon as its session logs on, then
// logs out.
type hasty struct {
	session.BaseHandler
	reports int
}

func (h hasty) OnLogon(ctx context.Context, s *session.Session) {
	for range h.reports {
		report := fix50sp2.BuildExecutionReport(s.ID().BeginString).
			SetOrderID("ORD-1").
			SetExecID("1").
			SetExecType(fix50sp2.ExecTypeNew).
			SetOrdStatus(fix50sp2.OrdStatusNew).
			SetSide(fix50sp2.SideBuy).
			SetLeavesQty(1).
			SetCumQty(0)

		if _, err := s.Send(ctx, report); err != nil {
			return
		}
	}

	go s.Logout(context.Background(), "bye")
}

// gatekeeper refuses every Logon.
type gatekeeper struct {
	session.BaseHandler
}

func (gatekeeper) Authenticate(context.Context, *session.LogonRequest) error {
	return errors.New("unknown user")
}

func freeAddress(t *testing.T) string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}

// serve runs an acceptor of FIX 5.0 SP2 sessions using handler, and returns
// the initiator and session configurations to connect to it.
func serve(t *testing.T, handler session.Handler) (*config.Initiator, *config.Session) {
	return serveVersion(t, handler, "FIX.5.0SP2")
}

// serveVersion is serve for sessions of FIX version version: FIX.5.0SP2 over
// FIXT.1.1, or a FIX 4.x BeginString.
func serveVersion(t *testing.T, handler session.Handler, version string) (*config.Initiator, *config.Session) {
	t.Helper()

	logger := zerolog.Nop()
	config.SetLogger(&logger)

	beginString, applVerID := version, ""
	if version == "FIX.5.0SP2" {
		beginString, applVerID = session.FIXT11, version
	}

	address := freeAddress(t)
	server := &config.Session{
		Name:             "server",
		BeginString:      beginString,
		SenderCompID:     "BIGCORP",
		TargetCompID:     "smallcorp",
		TargetSubID:      "john",
		DefaultApplVerID: applVerID,
		HeartBtInt:       5 * time.Second,
	}
	client := &config.Session{
		Name:             "client",
		BeginString:      beginString,
		SenderCompID:     "smallcorp",
		SenderSubID:      "john",
		TargetCompID:     "BIGCORP",
		DefaultApplVerID: applVerID,
		HeartBtInt:       5 * time.Second,
		LogonTimeout:     10 * time.Second,
		Username:         "john",
		Password:         "secret",
	}
	initiatorConfig := &config.Initiator{Name: "client", Addresses: []string{address}}

	conf := config.GetConfig()
	conf.Contexts = []*config.Context{{Name: "server", Acceptor: "server", Sessions: []string{"server"}}}
	conf.Acceptors = []*config.Acceptor{{Name: "server", Address: address}}
	conf.Initiators = []*config.Initiator{initiatorConfig}
	conf.Sessions = []*config.Session{server, client}
	config.GetOptions().Context = "server"
	if err := config.GetConfig().Validate(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	served := make(chan error, 1)
	go func() {
		served <- acceptor.Serve(ctx, func(*config.Session) (session.Handler, error) { return handler, nil })
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-served; err != nil {
			t.Error(err)
		}
	})

	for deadline := time.Now().Add(5 * time.Second); ; time.Sleep(10 * time.Millisecond) {
		conn, err := net.Dial("tcp", address)
		if err == nil {
			conn.Close()
			break
		} else if time.Now().After(deadline) {
			t.Fatal(err)
		}
	}

	return initiatorConfig, client
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

func TestClient(t *testing.T) {
	initiatorConfig, client := serve(t, broker{})
	ctx := context.Background()

	c, err := initiator.ConnectSession(ctx, initiatorConfig, client, initiator.Options{
		MsgTypes: []string{string(fix50sp2.MsgTypeExecutionreport)},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(true)

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

// TestClientVersions checks an order and its ExecutionReport go through
// sessions of every FIX version fix supports, both sides validating what
// they receive with the version's dictionary.
func TestClientVersions(t *testing.T) {
	for _, version := range []string{"FIX.5.0SP2", "FIX.4.4", "FIX.4.3"} {
		t.Run(version, func(t *testing.T) {
			initiatorConfig, client := serveVersion(t, broker{}, version)
			ctx := context.Background()

			c, err := initiator.ConnectSession(ctx, initiatorConfig, client, initiator.Options{
				MsgTypes: []string{string(fix50sp2.MsgTypeExecutionreport)},
			})
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close(true)

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
	initiatorConfig, client := serveVersion(t, broker{}, "FIX.4.3")
	ctx := context.Background()

	c, err := initiator.ConnectSession(ctx, initiatorConfig, client, initiator.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(true)

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

	c, err := initiator.ConnectSession(ctx, initiatorConfig, client, initiator.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(true)

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

	c, err := initiator.ConnectSession(context.Background(), initiatorConfig, client, initiator.Options{})
	if err != nil {
		t.Fatalf("ConnectSession() = %v", err)
	}
	defer c.Close(true)

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
	if err := c.Err(); err == nil || strings.Contains(err.Error(), "refused") {
		t.Errorf("Err() = %v, want the counterparty's logout", err)
	}
}

// TestClientLogonRejected checks that a Logon answer the session itself
// rejects, here for its invalid NextExpectedMsgSeqNum, ends Connect rather
// than being taken for an accepted Logon and retried until the timeout.
func TestClientLogonRejected(t *testing.T) {
	logger := zerolog.Nop()
	config.SetLogger(&logger)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })

	// A fake acceptor answering every Logon with one whose
	// NextExpectedMsgSeqNum (789) is beyond what was sent.
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				if _, err := fix.NewFramer(conn, fix.FramerOptions{}).Next(); err != nil {
					return
				}
				answer := fix.NewBuilder(session.FIXT11, "A").
					SetString(49, "BIGCORP").
					SetString(56, "smallcorp").
					SetString(57, "john").
					SetInt(34, 1).
					SetString(52, time.Now().UTC().Format("20060102-15:04:05.000")).
					SetInt(98, 0).
					SetInt(108, 5).
					SetString(1137, "9").
					SetInt(789, 100)
				conn.Write(answer.Build())
				io.Copy(io.Discard, conn)
			}()
		}
	}()

	client := &config.Session{
		Name:                  "client",
		BeginString:           session.FIXT11,
		SenderCompID:          "smallcorp",
		SenderSubID:           "john",
		TargetCompID:          "BIGCORP",
		DefaultApplVerID:      "FIX.5.0SP2",
		HeartBtInt:            5 * time.Second,
		LogonTimeout:          10 * time.Second,
		NextExpectedMsgSeqNum: true,
	}
	initiatorConfig := &config.Initiator{Name: "client", Addresses: []string{l.Addr().String()}}

	start := time.Now()
	_, err = initiator.ConnectSession(context.Background(), initiatorConfig, client, initiator.Options{})
	if !errors.Is(err, fixerrors.FixLogout) || !strings.Contains(err.Error(), "logon refused") {
		t.Fatalf("ConnectSession() = %v, want a refused logon", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("ConnectSession() took %v, the logon timeout is 10s", elapsed)
	}
}

// TestClientLogonRefused checks that Connect fails at once when the
// counterparty refuses the Logon, rather than at the logon timeout.
func TestClientLogonRefused(t *testing.T) {
	initiatorConfig, client := serve(t, gatekeeper{})

	start := time.Now()
	_, err := initiator.ConnectSession(context.Background(), initiatorConfig, client, initiator.Options{})
	if !errors.Is(err, fixerrors.FixLogout) || !strings.Contains(err.Error(), "logon refused") {
		t.Fatalf("ConnectSession() = %v, want a refused logon", err)
	}
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Errorf("ConnectSession() took %v, the logon timeout is 10s", elapsed)
	}
}
