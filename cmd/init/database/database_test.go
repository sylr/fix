package databaseinit

import (
	"os"
	"path/filepath"
	"testing"

	"sylr.dev/fix/v2/config"
)

func TestCreateDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "initiator.db")
	store := config.Store{Type: "sqlite3", Path: path}
	sessions := []*config.Session{{
		Name:         "client",
		BeginString:  "FIXT.1.1",
		SenderCompID: "smallcorp",
		TargetCompID: "BIGCORP",
	}}

	if err := createDatabase("initiator test", store, "initiator.db", sessions); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Fatal("database is empty")
	}

	// Existing databases are skipped without --overwrite.
	if err := createDatabase("initiator test", store, "initiator.db", sessions); err != nil {
		t.Fatal(err)
	}

	overwrite = true
	defer func() { overwrite = false }()

	if err := createDatabase("initiator test", store, "initiator.db", sessions); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDatabaseMemory(t *testing.T) {
	if err := createDatabase("initiator test", config.Store{}, "initiator.db", nil); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDatabaseUnknownType(t *testing.T) {
	path := filepath.Join(t.TempDir(), "initiator.db")
	if err := os.WriteFile(path, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}

	overwrite = true
	defer func() { overwrite = false }()

	if err := createDatabase("initiator test", config.Store{Type: "sqlit3", Path: path}, "initiator.db", nil); err == nil {
		t.Fatal("unknown store type: no error")
	}
	if b, err := os.ReadFile(path); err != nil || string(b) != "keep" {
		t.Errorf("an unknown store type touched %s: %q, %v", path, b, err)
	}
}
