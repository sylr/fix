package list

import (
	"github.com/spf13/cobra"

	listsecurity "sylr.dev/fix/v2/cmd/list/security"
	"sylr.dev/fix/v2/initiator"
)

// ListCmd groups the commands requesting lists.
var ListCmd = &cobra.Command{
	Use:               "list",
	Short:             "Send a list FIX message",
	Long:              "Send a list FIX message after initiating a session with a FIX acceptor.",
	PersistentPreRunE: initiator.ValidateOptions,
}

func init() {
	initiator.AddPersistentFlags(ListCmd)
	if err := initiator.AddPersistentFlagCompletions(ListCmd); err != nil {
		panic(err)
	}

	ListCmd.AddCommand(listsecurity.ListSecurityCmd)
}
