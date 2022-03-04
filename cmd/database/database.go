package database

import (
	"github.com/spf13/cobra"

	dbinit "sylr.dev/fix/cmd/database/init"
	fdb "sylr.dev/fix/pkg/database"
)

// ConfigCmd represents the buy command
var DatabaseCmd = &cobra.Command{
	Use:     "database",
	Aliases: []string{"db"},
	Short:   "Manage quickfix database",
	Long:    "Manage quickfix database.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := fdb.ValidateOptions(cmd, args)
		if err != nil {
			return err
		}

		if cmd.HasParent() {
			parent := cmd.Parent()
			if parent.PersistentPreRunE != nil {
				err = parent.PersistentPreRunE(parent, args)
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func init() {
	DatabaseCmd.AddCommand(dbinit.DatabaseInitCmd)
}
