package databaseinit

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	qsql "github.com/quickfixgo/quickfix/_sql"
	"github.com/spf13/cobra"
	"sylr.dev/fix/config"
)

var (
	overwrite bool
)

// ConfigCmd represents the buy command
var DatabaseInitCmd = &cobra.Command{
	Use:               "init",
	Short:             "Initialize quickfix database",
	Long:              "Initialize quickfix database.",
	RunE:              Execute,
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
	DatabaseInitCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing databases")
}

func Execute(cmd *cobra.Command, args []string) error {
	acceptors := config.GetConfig().Acceptors
	initiators := config.GetConfig().Initiators

	// Acceptors
	for _, acceptor := range acceptors {
		if len(acceptor.SQLStoreDriver) == 0 {
			continue
		}

		var file string
		switch acceptor.SQLStoreDriver {
		case "sqlite3":
			file = acceptor.SQLStoreDataSourceName
			if len(file) == 0 {
				file = os.ExpandEnv(strings.Join([]string{"$HOME", ".fix", "acceptor.db"}, string(os.PathSeparator)))
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
			file = acceptor.SQLStoreDataSourceName
		}

		dbConn, err := sql.Open(acceptor.SQLStoreDriver, file)

		if err != nil {
			return err
		}

		defer dbConn.Close()

		dir, _ := qsql.FS.ReadDir(acceptor.SQLStoreDriver)

		for _, f := range dir {
			if !f.Type().Type().IsRegular() {
				continue
			}
			fd, _ := qsql.FS.Open(strings.Join([]string{acceptor.SQLStoreDriver, f.Name()}, string(os.PathSeparator)))
			sqlBytes, _ := io.ReadAll(fd)

			if _, err := dbConn.Exec(string(sqlBytes)); err != nil {
				return err
			}
		}
	}

	// Initiators
	for _, initiator := range initiators {
		if len(initiator.SQLStoreDriver) == 0 {
			continue
		}

		var file string
		switch initiator.SQLStoreDriver {
		case "sqlite3":
			file = initiator.SQLStoreDataSourceName
			if len(file) == 0 {
				file = os.ExpandEnv(strings.Join([]string{"$HOME", ".fix", "initiator.db"}, string(os.PathSeparator)))
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
			file = initiator.SQLStoreDataSourceName
		}

		dbConn, err := sql.Open(initiator.SQLStoreDriver, file)

		if err != nil {
			return err
		}

		defer dbConn.Close()

		dir, _ := qsql.FS.ReadDir(initiator.SQLStoreDriver)

		for _, f := range dir {
			if !f.Type().Type().IsRegular() {
				continue
			}
			fd, _ := qsql.FS.Open(strings.Join([]string{initiator.SQLStoreDriver, f.Name()}, string(os.PathSeparator)))
			sqlBytes, _ := io.ReadAll(fd)

			if _, err := dbConn.Exec(string(sqlBytes)); err != nil {
				return err
			}
		}
	}

	return nil
}
