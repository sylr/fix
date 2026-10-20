package databaseinit

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/database"
	"sylr.dev/fix/v2/utils"
)

var (
	overwrite bool
)

// ConfigCmd represents the buy command
var InitDatabaseCmd = &cobra.Command{
	Use:               "database",
	Short:             "Initialize fix database",
	Long:              "Initialize fix database.",
	RunE:              Execute,
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := utils.ValidateRequiredFlags(cmd); err != nil {
			return err
		}

		if err := database.ValidateOptions(cmd, args); err != nil {
			return err
		}

		if cmd.HasParent() {
			parent := cmd.Parent()
			if parent.PersistentPreRunE != nil {
				return parent.PersistentPreRunE(cmd, args)
			}
		}
		return nil
	},
}

func init() {
	InitDatabaseCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing databases")
}

func Execute(cmd *cobra.Command, args []string) error {
	conf := config.GetConfig()

	for _, acceptor := range conf.Acceptors {
		sessions := contextSessions(func(c *config.Context) bool { return c.Acceptor == acceptor.Name })
		if err := createDatabase(acceptor.Store, "acceptor.db", sessions); err != nil {
			return err
		}
	}

	for _, initiator := range conf.Initiators {
		sessions := contextSessions(func(c *config.Context) bool { return c.Initiator == initiator.Name })
		if err := createDatabase(initiator.Store, "initiator.db", sessions); err != nil {
			return err
		}
	}

	return nil
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

// createDatabase creates the SQLite database of a sqlite3 store, and the
// tables of its sessions, which fix2's store creates when opened.
func createDatabase(store config.Store, defaultDBName string, sessions []*config.Session) error {
	if store.Type != "sqlite3" {
		return nil
	}

	file, err := store.SQLitePath(defaultDBName)
	if err != nil {
		return err
	}

	if _, err := os.Stat(file); err == nil {
		if !overwrite {
			fmt.Printf("Database %s already exists, skipping.\n", file)
			return nil
		}

		fmt.Printf("Database %s already exists but will be overwritten.\n", file)

		if err := os.Remove(file); err != nil {
			return err
		}

		// A store not closed cleanly leaves its write-ahead log behind,
		// which SQLite would replay into the new database.
		for _, suffix := range []string{"-wal", "-shm"} {
			if err := os.Remove(file + suffix); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	ctx := context.Background()

	// A database without sessions is created by opening one with a fixed ID.
	ids := make([]session.SessionID, 0, len(sessions)+1)
	for _, s := range sessions {
		ids = append(ids, s.ID())
	}
	if len(ids) == 0 {
		ids = append(ids, session.SessionID{})
	}

	for _, id := range ids {
		st, err := store.Open(defaultDBName, id)
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
	}

	return nil
}
