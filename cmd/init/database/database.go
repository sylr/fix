package databaseinit

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	qsql "github.com/quickfixgo/quickfix/_sql"
	"github.com/spf13/cobra"
	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/database"
)

var (
	overwrite bool
)

// ConfigCmd represents the buy command
var InitDatabaseCmd = &cobra.Command{
	Use:               "database",
	Short:             "Initialize quickfix database",
	Long:              "Initialize quickfix database.",
	RunE:              Execute,
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := database.ValidateOptions(cmd, args)
		if err != nil {
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
	acceptors := config.GetConfig().Acceptors
	initiators := config.GetConfig().Initiators

	// Acceptors
	if err := createDatabase(acceptors, "acceptor.db"); err != nil {
		return err
	}

	// Initiators
	if err := createDatabase(initiators, "initiator.db"); err != nil {
		return err
	}

	return nil
}

func createDatabase[T config.SQLStoreConfig](values []T, defaultDBName string) error {
	for _, value := range values {
		if len(value.GetSQLStoreDriver()) == 0 {
			continue
		}

		var file string
		switch value.GetSQLStoreDriver() {
		case "sqlite3":
			file = value.GetSQLStoreDataSourceName()
			if len(file) == 0 {
				file = os.ExpandEnv(filepath.Join("$HOME", ".fix", defaultDBName))
			}

			databaseDir := path.Dir(file)
			err := os.MkdirAll(databaseDir, 0700)
			if err != nil {
				return err
			}

			if _, err := os.Stat(file); err == nil {
				if !overwrite {
					fmt.Printf("Database %s already exists, skipping.\n", file)
					continue
				} else {
					fmt.Printf("Database %s already exists but will be overwritten.\n", file)
				}
			}
		default:
			file = value.GetSQLStoreDataSourceName()
		}

		dbConn, err := sql.Open(value.GetSQLStoreDriver(), file)

		if err != nil {
			return err
		}

		defer dbConn.Close()

		dir, _ := qsql.FS.ReadDir(value.GetSQLStoreDriver())

		for _, f := range dir {
			if !f.Type().IsRegular() {
				continue
			}
			fd, _ := qsql.FS.Open(filepath.Join(value.GetSQLStoreDriver(), f.Name()))
			sqlBytes, _ := io.ReadAll(fd)

			if _, err := dbConn.Exec(string(sqlBytes)); err != nil {
				return err
			}
		}
	}

	return nil
}
