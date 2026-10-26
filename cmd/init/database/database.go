package initdatabase

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/database"
)

var (
	overwrite bool
)

// InitDatabaseCmd creates the SQLite databases of the initiators and
// acceptors keeping their sessions in one.
var InitDatabaseCmd = &cobra.Command{
	Use:               "database",
	Short:             "Initialize fix database",
	Long:              "Initialize fix database.",
	RunE:              Execute,
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: database.ValidateOptions,
}

func init() {
	InitDatabaseCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing databases")
}

// sqliteDB is a SQLite database, and the initiators and acceptors sharing
// it.
type sqliteDB struct {
	path     string
	owners   []string
	sessions []*config.Session
}

func Execute(cmd *cobra.Command, args []string) error {
	conf := config.GetConfig()

	if len(conf.Acceptors) == 0 && len(conf.Initiators) == 0 {
		fmt.Println("No initiator nor acceptor configured: no database to create.")
		return nil
	}

	var dbs []*sqliteDB
	add := func(owner string, store config.Store, defaultDBName string, match func(*config.Context) bool) error {
		if err := store.Validate(); err != nil {
			return fmt.Errorf("%s: %w", owner, err)
		}

		if store.IsMemory() {
			fmt.Printf("%s keeps its sessions in memory (store.type: sqlite3 to keep them in a database), skipping.\n", owner)
			return nil
		}

		path, err := store.SQLitePath(defaultDBName)
		if err != nil {
			return fmt.Errorf("%s: %w", owner, err)
		}
		if path, err = canonical(path); err != nil {
			return fmt.Errorf("%s: %w", owner, err)
		}

		i := slices.IndexFunc(dbs, func(db *sqliteDB) bool { return db.path == path })
		if i < 0 {
			dbs = append(dbs, &sqliteDB{path: path})
			i = len(dbs) - 1
		}

		db := dbs[i]
		db.owners = append(db.owners, owner)
		for _, s := range contextSessions(match) {
			if !slices.ContainsFunc(db.sessions, func(o *config.Session) bool { return o.ID() == s.ID() }) {
				db.sessions = append(db.sessions, s)
			}
		}

		return nil
	}

	for _, acceptor := range conf.Acceptors {
		match := func(c *config.Context) bool { return c.Acceptor == acceptor.Name }
		if err := add("acceptor "+acceptor.Name, acceptor.Store, "acceptor.db", match); err != nil {
			return err
		}
	}

	for _, initiator := range conf.Initiators {
		match := func(c *config.Context) bool { return c.Initiator == initiator.Name }
		if err := add("initiator "+initiator.Name, initiator.Store, "initiator.db", match); err != nil {
			return err
		}
	}

	for _, db := range dbs {
		if err := createDatabase(db); err != nil {
			return err
		}
	}

	return nil
}

// canonical returns path absolute, with the symlinks of its longest existing
// prefix resolved, so that the aliases of a database, existing or not, name
// it alike.
func canonical(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for dir := path; ; dir = filepath.Dir(dir) {
		resolved, err := filepath.EvalSymlinks(dir)
		if err == nil {
			rest, err := filepath.Rel(dir, path)
			if err != nil {
				return "", err
			}

			return filepath.Join(resolved, rest), nil
		}

		if !errors.Is(err, fs.ErrNotExist) || dir == filepath.Dir(dir) {
			return "", err
		}
	}
}

// contextSessions returns the sessions of the contexts match accepts.
func contextSessions(match func(*config.Context) bool) []*config.Session {
	var sessions []*config.Session

	for _, c := range config.GetContexts() {
		if !match(c) {
			continue
		}

		for _, name := range c.Sessions {
			if s, err := config.GetSession(name); err == nil {
				sessions = append(sessions, s)
			}
		}
	}

	return sessions
}

// createDatabase creates the SQLite database db, with the rows of its
// sessions, which fix2's store creates when opened.
func createDatabase(db *sqliteDB) error {
	owners := strings.Join(db.owners, ", ")

	if len(db.sessions) == 0 {
		fmt.Printf("No session uses database %s of %s, skipping.\n", db.path, owners)
		return nil
	}

	if _, err := os.Stat(db.path); err == nil {
		if !overwrite {
			fmt.Printf("Database %s of %s already exists, skipping.\n", db.path, owners)
			return nil
		}

		fmt.Printf("Removing database %s of %s.\n", db.path, owners)

		if err := os.Remove(db.path); err != nil {
			return err
		}

		// A store not closed cleanly leaves its write-ahead log behind,
		// which SQLite would replay into the new database.
		for _, suffix := range []string{"-wal", "-shm"} {
			if err := os.Remove(db.path + suffix); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	ctx := context.Background()
	store := config.Store{Type: "sqlite3", Path: db.path}

	for _, s := range db.sessions {
		st, err := store.Open("", s.ID())
		if err != nil {
			return err
		}

		if _, err := st.Open(ctx); err != nil {
			st.Close()
			return err
		}

		if err := st.Close(); err != nil {
			return err
		}

		fmt.Printf("Database %s of %s: session %s initialized.\n", db.path, owners, s.ID())
	}

	return nil
}
