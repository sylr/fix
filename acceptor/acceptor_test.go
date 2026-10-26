package acceptor

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
)

// serveContext configures two sessions served on a free address, and
// returns their context and the address.
func serveContext(t *testing.T) (*config.Context, string) {
	t.Helper()

	logger := zerolog.Nop()
	config.SetLogger(&logger)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := l.Addr().String()
	l.Close()

	cfgContext := &config.Context{Name: "server", Acceptor: "server", Sessions: []string{"a", "b"}}
	conf := config.GetConfig()
	conf.Contexts = []*config.Context{cfgContext}
	conf.Acceptors = []*config.Acceptor{{Name: "server", Address: address}}
	conf.Initiators = nil
	conf.Sessions = []*config.Session{
		{Name: "a", BeginString: session.FIXT11, DefaultApplVerID: "FIX.5.0SP2", SenderCompID: "A", TargetCompID: "B"},
		{Name: "b", BeginString: session.FIXT11, DefaultApplVerID: "FIX.5.0SP2", SenderCompID: "A", TargetCompID: "C"},
	}
	if err := conf.Validate(); err != nil {
		t.Fatal(err)
	}

	return cfgContext, address
}

// TestServeUnwinds checks a Serve that fails to set a session up releases
// its address.
func TestServeUnwinds(t *testing.T) {
	cfgContext, address := serveContext(t)
	errHandler := errors.New("no handler")

	err := ServeContext(context.Background(), cfgContext, func(s *config.Session) (session.Handler, error) {
		if s.Name == "b" {
			return nil, errHandler
		}
		return session.BaseHandler{}, nil
	})
	if !errors.Is(err, errHandler) {
		t.Fatalf("ServeContext() = %v, want %v", err, errHandler)
	}

	l, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("address still in use: %v", err)
	}
	l.Close()
}

// TestServeAddressInUse checks a Serve that cannot listen fails before
// setting up any session.
func TestServeAddressInUse(t *testing.T) {
	cfgContext, address := serveContext(t)

	l, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	called := false
	err = ServeContext(context.Background(), cfgContext, func(*config.Session) (session.Handler, error) {
		called = true
		return session.BaseHandler{}, nil
	})
	if err == nil {
		t.Fatal("ServeContext() = nil, want the address in use")
	}
	if called {
		t.Error("sessions set up before listening")
	}
}
