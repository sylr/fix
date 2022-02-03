package list

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/pkg/initiator"
)

// ListCmd represents the buy command
var ListCmd = &cobra.Command{
	Use:               "list",
	Short:             "Send a list FIX Message",
	Long:              "Send a list FIX message after initiating a sesion with a FIX acceptor.",
	RunE:              Execute,
	PersistentPreRunE: initiator.ValidateOptions,
}

func init() {
	initiator.AddPersistentFlags(ListCmd)
	initiator.AddPersistentFlagCompletions(ListCmd)
}

func Execute(cmd *cobra.Command, args []string) error {
	return nil
}
