package list

import (
	"github.com/spf13/cobra"

	listsecurity "sylr.dev/fix/v2/cmd/list/security"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/utils"
)

// ListCmd represents the buy command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Send a list FIX message",
	Long:  "Send a list FIX message after initiating a session with a FIX acceptor.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := utils.ValidateRequiredFlags(cmd); err != nil {
			return err
		}

		if err := initiator.ValidateOptions(cmd, args); err != nil {
			return err
		}

		if cmd.HasParent() {
			parent := cmd.Parent()
			if parent.PersistentPreRunE != nil {
				return parent.PersistentPreRunE(parent, args)
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
