package cancel

import (
	"github.com/spf13/cobra"

	cancelorder "sylr.dev/fix/cmd/cancel/order"
	"sylr.dev/fix/pkg/initiator"
)

// CancelCmd represents the buy command
var CancelCmd = &cobra.Command{
	Use:               "cancel",
	Short:             "Send a cancel FIX message",
	Long:              "Send a cancel FIX message after initiating a sesion with a FIX acceptor.",
	RunE:              Execute,
	PersistentPreRunE: initiator.ValidateOptions,
}

func init() {
	initiator.AddPersistentFlags(CancelCmd)
	initiator.AddPersistentFlagCompletions(CancelCmd)
	initiator.AddPersistentFlagCompletions(cancelorder.CancelOrderCmd)

	CancelCmd.AddCommand(cancelorder.CancelOrderCmd)
}

func Execute(cmd *cobra.Command, args []string) error {
	return nil
}
