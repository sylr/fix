//go:build cgo

package initdatabase

import (
	"database/sql"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"sylr.dev/fix/v2/config"
)

// sessionKeys returns the sessions the SQLite database path holds.
func sessionKeys(t *testing.T, path string) []string {
	t.Helper()
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT session_key FROM fix2_sessions ORDER BY session_key")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			t.Fatal(err)
		}
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return keys
}

func TestExecuteSharedDatabase(t *testing.T) {
	dir := t.TempDir()
	shared := filepath.Join(dir, "sub", "shared.db")
	unused := filepath.Join(dir, "unused.db")
	store := config.Store{Type: "sqlite3", Path: shared}

	client := &config.Session{Name: "client", BeginString: "FIXT.1.1", DefaultApplVerID: "FIX.5.0SP2", SenderCompID: "smallcorp", TargetCompID: "BIGCORP"}
	server := &config.Session{Name: "server", BeginString: "FIX.4.4", SenderCompID: "BIGCORP", TargetCompID: "smallcorp"}

	setConfig(t,
		[]*config.Initiator{{Name: "i", Store: store}, {Name: "idle", Store: config.Store{Type: "sqlite3", Path: unused}}},
		[]*config.Acceptor{{Name: "a", Store: store}},
		[]*config.Context{
			{Name: "client", Initiator: "i", Sessions: []string{"client"}},
			{Name: "client2", Initiator: "i", Sessions: []string{"client"}},
			{Name: "server", Acceptor: "a", Sessions: []string{"server"}},
		},
		[]*config.Session{client, server},
		false)

	want := []string{client.ID().String(), server.ID().String()}
	slices.Sort(want)

	if err := Execute(nil, nil); err != nil {
		t.Fatal(err)
	}
	if keys := sessionKeys(t, shared); !slices.Equal(keys, want) {
		t.Errorf("sessions = %q, want %q", keys, want)
	}
	if _, err := os.Stat(unused); !os.IsNotExist(err) {
		t.Errorf("database of an initiator without sessions: %v", err)
	}

	// Existing databases are skipped without --overwrite, and replaced with
	// it: a broken one would not open.
	if err := os.WriteFile(shared, []byte("broken"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Execute(nil, nil); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(shared); string(b) != "broken" {
		t.Errorf("database replaced without --overwrite")
	}

	overwrite = true
	if err := Execute(nil, nil); err != nil {
		t.Fatal(err)
	}
	if keys := sessionKeys(t, shared); !slices.Equal(keys, want) {
		t.Errorf("sessions after --overwrite = %q, want %q", keys, want)
	}
}

// TestExecuteDatabaseAliases checks the stores naming one database through
// a symlinked directory share it: with --overwrite, the second does not
// remove what the first created.
func TestExecuteDatabaseAliases(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "real"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("real", filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}

	client := &config.Session{Name: "client", BeginString: "FIXT.1.1", DefaultApplVerID: "FIX.5.0SP2", SenderCompID: "smallcorp", TargetCompID: "BIGCORP"}
	server := &config.Session{Name: "server", BeginString: "FIX.4.4", SenderCompID: "BIGCORP", TargetCompID: "smallcorp"}

	want := []string{client.ID().String(), server.ID().String()}
	slices.Sort(want)

	// The database, and the directory holding it, need not exist.
	for _, name := range []string{"shared.db", filepath.Join("sub", "shared.db")} {
		setConfig(t,
			[]*config.Initiator{{Name: "i", Store: config.Store{Type: "sqlite3", Path: filepath.Join(dir, "link", name)}}},
			[]*config.Acceptor{{Name: "a", Store: config.Store{Type: "sqlite3", Path: filepath.Join(dir, "real", name)}}},
			[]*config.Context{
				{Name: "client", Initiator: "i", Sessions: []string{"client"}},
				{Name: "server", Acceptor: "a", Sessions: []string{"server"}},
			},
			[]*config.Session{client, server},
			true)

		if err := Execute(nil, nil); err != nil {
			t.Fatal(err)
		}
		if keys := sessionKeys(t, filepath.Join(dir, "real", name)); !slices.Equal(keys, want) {
			t.Errorf("%s: sessions = %q, want %q", name, keys, want)
		}

		// Existing, the database is overwritten once.
		if err := Execute(nil, nil); err != nil {
			t.Fatal(err)
		}
		if keys := sessionKeys(t, filepath.Join(dir, "real", name)); !slices.Equal(keys, want) {
			t.Errorf("%s: sessions after --overwrite = %q, want %q", name, keys, want)
		}
	}
}
