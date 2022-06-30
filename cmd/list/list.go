package list

import (
	"github.com/spf13/cobra"

	listsecurity "sylr.dev/fix/cmd/list/security"
	"sylr.dev/fix/pkg/initiator"
)

// ListCmd represents the buy command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Send a list FIX message",
	Long:  "Send a list FIX message after initiating a sesion with a FIX acceptor.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := initiator.ValidateOptions(cmd, args)
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
	initiator.AddPersistentFlags(ListCmd)
	initiator.AddPersistentFlagCompletions(ListCmd)
	initiator.AddPersistentFlagCompletions(listsecurity.ListSecurityCmd)

	ListCmd.AddCommand(listsecurity.ListSecurityCmd)
}
