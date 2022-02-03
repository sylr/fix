package new

import (
	"github.com/spf13/cobra"

	neworder "sylr.dev/fix/cmd/new/order"
	"sylr.dev/fix/pkg/initiator"
)

// NewCmd represents the buy command
var NewCmd = &cobra.Command{
	Use:               "new",
	Short:             "Send a new FIX message",
	Long:              "Send a new FIX message after initiating a sesion with a FIX acceptor.",
	RunE:              Execute,
	PersistentPreRunE: initiator.ValidateOptions,
}

func init() {
	initiator.AddPersistentFlags(NewCmd)
	initiator.AddPersistentFlagCompletions(NewCmd)
	initiator.AddPersistentFlagCompletions(neworder.NewOrderCmd)

	NewCmd.AddCommand(neworder.NewOrderCmd)
}

func Execute(cmd *cobra.Command, args []string) error {
	return nil
}
