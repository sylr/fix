package cmd

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"sylr.dev/fix/v2/config"
	fixerrors "sylr.dev/fix/v2/errors"
)

const testConfig = `current-context: client
contexts:
- name: client
  initiator: client
  sessions: [client]
initiators:
- name: client
  addresses: [127.0.0.1:1]
sessions:
- name: client
  beginString: FIXT.1.1
  defaultApplVerID: FIX.5.0SP2
  senderCompID: smallcorp
  targetCompID: BIGCORP
  heartBtInt: 5s
`

// execute runs fix with args, and returns its error.
func execute(t *testing.T, args ...string) error {
	t.Helper()

	t.Cleanup(func() {
		options := config.GetOptions()
		options.Metrics, options.HTTPPort = false, 8080
	})

	FixCmd.SetArgs(args)
	FixCmd.SetErr(new(nopWriter))
	return FixCmd.ExecuteContext(context.Background())
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }

// TestHooks checks the hooks run from the root to the command: the
// configuration is loaded before the command validates its options.
func TestHooks(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(path, []byte(testConfig), 0o600); err != nil {
		t.Fatal(err)
	}

	err := execute(t, "status", "security", "--config", path, "--symbol", "EURUSD", "--subscription-type", "never")
	if !errors.Is(err, fixerrors.ErrOptions) {
		t.Errorf("Execute(--subscription-type never) = %v, want an options error", err)
	}

	err = execute(t, "status", "security", "--config", filepath.Join(t.TempDir(), "none"), "--symbol", "EURUSD", "--subscription-type", "never")
	if err == nil || errors.Is(err, fixerrors.ErrOptions) {
		t.Errorf("Execute(--config none) = %v, want the configuration error first", err)
	}

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = l.Close() }()
	port := strconv.Itoa(l.Addr().(*net.TCPAddr).Port)

	err = execute(t, "status", "security", "--config", path, "--symbol", "EURUSD", "--subscription-type", "never", "--metrics", "--port", port)
	if err == nil || errors.Is(err, fixerrors.ErrOptions) {
		t.Errorf("Execute(--port %s, taken) = %v, want the listen error first", port, err)
	}
}
