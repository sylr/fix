package initdatabase

import (
	"os"
	"path/filepath"
	"testing"

	"sylr.dev/fix/v2/config"
)

// setConfig sets the initiators, acceptors, contexts and sessions of the
// configuration, and --overwrite, for the test.
func setConfig(t *testing.T, initiators []*config.Initiator, acceptors []*config.Acceptor, contexts []*config.Context, sessions []*config.Session, overwriteFlag bool) {
	t.Helper()
	conf := config.GetConfig()
	old, oldOverwrite := *conf, overwrite
	t.Cleanup(func() { *conf, overwrite = old, oldOverwrite })

	conf.Initiators, conf.Acceptors, conf.Contexts, conf.Sessions = initiators, acceptors, contexts, sessions
	overwrite = overwriteFlag
}

func TestExecuteMemory(t *testing.T) {
	setConfig(t, []*config.Initiator{{Name: "i"}}, []*config.Acceptor{{Name: "a", Store: config.Store{Type: "memory"}}}, nil, nil, true)

	if err := Execute(nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestExecuteRejectedStore(t *testing.T) {
	types := []string{"sqlit3"}
	if !config.SQLiteStore {
		types = append(types, "sqlite3")
	}

	for _, typ := range types {
		path := filepath.Join(t.TempDir(), "initiator.db")
		if err := os.WriteFile(path, []byte("keep"), 0o600); err != nil {
			t.Fatal(err)
		}

		setConfig(t,
			[]*config.Initiator{{Name: "i", Store: config.Store{Type: typ, Path: path}}}, nil,
			[]*config.Context{{Name: "c", Initiator: "i", Sessions: []string{"s"}}},
			[]*config.Session{{Name: "s", BeginString: "FIX.4.4", SenderCompID: "A", TargetCompID: "B"}},
			true)

		if err := Execute(nil, nil); err == nil {
			t.Errorf("store type %s: no error", typ)
		}
		if b, err := os.ReadFile(path); err != nil || string(b) != "keep" {
			t.Errorf("store type %s touched %s: %q, %v", typ, path, b, err)
		}
	}
}
