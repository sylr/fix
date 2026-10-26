// Package initiatortest runs, in tests, an acceptor for initiator.Clients to
// connect to.
package initiatortest

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/acceptor"
	"sylr.dev/fix/v2/config"
)

// Serve runs an acceptor of sessions of FIX version version, FIX.5.0SP2 over
// FIXT.1.1 or a FIX 4.x BeginString, using handler, until t ends. It makes
// the "client" context, whose initiator and session connect to it, the
// current context, and returns them. The loggers log nothing.
func Serve(t testing.TB, handler session.Handler, version string) (*config.Initiator, *config.Session) {
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
	serverContext := &config.Context{Name: "server", Acceptor: "server", Sessions: []string{"server"}}

	conf := config.GetConfig()
	conf.Contexts = []*config.Context{serverContext, {Name: "client", Initiator: "client", Sessions: []string{"client"}}}
	conf.Acceptors = []*config.Acceptor{{Name: "server", Address: address}}
	conf.Initiators = []*config.Initiator{initiatorConfig}
	conf.Sessions = []*config.Session{server, client}
	conf.CurrentContext = "client"
	config.GetOptions().Context = ""
	if err := conf.Validate(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	served := make(chan error, 1)
	go func() {
		served <- acceptor.ServeContext(ctx, serverContext, func(*config.Session) (session.Handler, error) { return handler, nil })
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

func freeAddress(t testing.TB) string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}
