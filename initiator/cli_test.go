package initiator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
)

const testConfig = `
contexts:
  - name: default
    initiator: venue-a
    sessions: [a]
  - name: other
    initiator: venue-b
    sessions: [b]
  - name: two
    initiator: venue-a
    sessions: [a, b]
  - name: served
    acceptor: server
    sessions: [a]
  # Named as --initiator venue-b --session a would name its context.
  - name: (venue-b/a)
    initiator: venue-a
    sessions: [b]
  - name: (venue-b/a)#2
    initiator: venue-a
    sessions: [a]
initiators:
  - name: venue-a
    addresses: [127.0.0.1:5001]
  - name: venue-b
    addresses: [127.0.0.1:5002]
acceptors:
  - name: server
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
		options.Config, options.Context, options.Initiator, options.Session = "", "", "", ""
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
		name                        string
		file                        string
		context, initiator, session string
		// wantInitiator and wantSession are those of the context resolved,
		// wantErr what the error holds otherwise.
		wantInitiator, wantSession string
		wantErr                    error
	}{
		{name: "current-context", file: withCurrent, wantInitiator: "venue-a", wantSession: "a"},
		{name: "--context", file: withCurrent, context: "other", wantInitiator: "venue-b", wantSession: "b"},
		{name: "--initiator --session over current-context", file: withCurrent, initiator: "venue-b", session: "b", wantInitiator: "venue-b", wantSession: "b"},
		{name: "--initiator --session", file: withoutCurrent, initiator: "venue-a", session: "b", wantInitiator: "venue-a", wantSession: "b"},
		{name: "--initiator --session named as a context", file: withCurrent, initiator: "venue-b", session: "a", wantInitiator: "venue-b", wantSession: "a"},
		{name: "--initiator alone", file: withCurrent, initiator: "venue-b", wantErr: errors.ErrOptions},
		{name: "--session alone", file: withCurrent, session: "b", wantErr: errors.ErrOptions},
		{name: "--context --session", file: withCurrent, context: "other", session: "b", wantErr: errors.ErrOptions},
		{name: "nothing", file: withoutCurrent, wantErr: errors.ErrOptions},
		{name: "unknown initiator", file: withCurrent, initiator: "venue-c", session: "a", wantErr: errors.ErrConfigInitiatorNotFound},
		{name: "unknown session", file: withCurrent, initiator: "venue-a", session: "c", wantErr: errors.ErrConfigSessionNotFound},
		{name: "two sessions", file: withCurrent, context: "two", wantErr: errors.ErrConfigContextMultipleSessions},
		{name: "no initiator", file: withCurrent, context: "served", wantErr: errors.ErrConfig},
	} {
		t.Run(tc.name, func(t *testing.T) {
			options := config.GetOptions()
			options.Config, options.Interactive = tc.file, false
			options.Context, options.Initiator, options.Session = tc.context, tc.initiator, tc.session

			err := ValidateOptions(nil, nil)
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
			if cfgContext.Initiator != tc.wantInitiator || strings.Join(cfgContext.Sessions, ",") != tc.wantSession {
				t.Errorf("context %s: initiator %s, sessions %v, want %s and %s", cfgContext.Name, cfgContext.Initiator, cfgContext.Sessions, tc.wantInitiator, tc.wantSession)
			}

			// Resolving again, as a second PersistentPreRunE would, agrees.
			if err := ValidateOptions(nil, nil); err != nil {
				t.Errorf("ValidateOptions() again = %v", err)
			}
		})
	}
}
