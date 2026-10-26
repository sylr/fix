package acceptor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
)

const testConfig = `
contexts:
  - name: default
    acceptor: server-a
    sessions: [a]
  - name: other
    acceptor: server-b
    sessions: [a, b]
  - name: client
    initiator: venue
    sessions: [a]
  # Named as --acceptor server-b --session a would name its context.
  - name: (server-b/a)
    acceptor: server-a
    sessions: [b]
  - name: (server-b/a)#2
    acceptor: server-a
    sessions: [a]
initiators:
  - name: venue
    addresses: [127.0.0.1:5001]
acceptors:
  - name: server-a
    address: 127.0.0.1:5002
  - name: server-b
    address: 127.0.0.1:5003
sessions:
  - name: a
    beginString: FIXT.1.1
    defaultApplVerID: FIX.5.0SP2
    senderCompID: A
    targetCompID: B
  - name: b
    beginString: FIXT.1.1
    defaultApplVerID: FIX.5.0SP2
    senderCompID: A
    targetCompID: C
`

func TestValidateOptions(t *testing.T) {
	// No SSH key: Load must not read the user's.
	t.Setenv("HOME", t.TempDir())
	t.Cleanup(func() {
		options := config.GetOptions()
		options.Config, options.Context, options.Acceptor, options.Session = "", "", "", ""
	})

	dir := t.TempDir()
	withCurrent := filepath.Join(dir, "current")
	withoutCurrent := filepath.Join(dir, "none")
	if err := os.WriteFile(withCurrent, []byte("current-context: default\n"+testConfig), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(withoutCurrent, []byte(testConfig), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name                       string
		file                       string
		context, acceptor, session string
		// wantAcceptor and wantSessions are those of the context resolved,
		// wantErr what the error holds otherwise.
		wantAcceptor, wantSessions string
		wantErr                    error
	}{
		{name: "current-context", file: withCurrent, wantAcceptor: "server-a", wantSessions: "a"},
		{name: "--context", file: withCurrent, context: "other", wantAcceptor: "server-b", wantSessions: "a,b"},
		{name: "--acceptor --session over current-context", file: withCurrent, acceptor: "server-b", session: "b", wantAcceptor: "server-b", wantSessions: "b"},
		{name: "--acceptor --session", file: withoutCurrent, acceptor: "server-a", session: "b", wantAcceptor: "server-a", wantSessions: "b"},
		{name: "--acceptor --session named as a context", file: withCurrent, acceptor: "server-b", session: "a", wantAcceptor: "server-b", wantSessions: "a"},
		{name: "--acceptor alone", file: withCurrent, acceptor: "server-b", wantErr: errors.ErrOptions},
		{name: "--session alone", file: withCurrent, session: "b", wantErr: errors.ErrOptions},
		{name: "--context --acceptor", file: withCurrent, context: "other", acceptor: "server-a", wantErr: errors.ErrOptions},
		{name: "nothing", file: withoutCurrent, wantErr: errors.ErrOptions},
		{name: "unknown acceptor", file: withCurrent, acceptor: "server-c", session: "a", wantErr: errors.ErrConfigAcceptorNotFound},
		{name: "no acceptor", file: withCurrent, context: "client", wantErr: errors.ErrConfig},
	} {
		t.Run(tc.name, func(t *testing.T) {
			options := config.GetOptions()
			options.Config, options.Interactive = tc.file, false
			options.Context, options.Acceptor, options.Session = tc.context, tc.acceptor, tc.session

			err := ValidateOptions(&cobra.Command{}, nil)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("ValidateOptions() = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			cfgContext, err := config.GetCurrentContext()
			if err != nil {
				t.Fatal(err)
			}
			if cfgContext.Acceptor != tc.wantAcceptor || strings.Join(cfgContext.Sessions, ",") != tc.wantSessions {
				t.Errorf("context %s: acceptor %s, sessions %v, want %s and %s", cfgContext.Name, cfgContext.Acceptor, cfgContext.Sessions, tc.wantAcceptor, tc.wantSessions)
			}
		})
	}
}
