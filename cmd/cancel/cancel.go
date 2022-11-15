package cancel

import (
	"github.com/spf13/cobra"

	masscancelorder "sylr.dev/fix/cmd/cancel/mass"
	cancelorder "sylr.dev/fix/cmd/cancel/order"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/utils"
)

// CancelCmd represents the buy command
var CancelCmd = &cobra.Command{
	Use:               "cancel",
	Short:             "Send a cancel FIX message",
	Long:              "Send a cancel FIX message after initiating a sesion with a FIX acceptor.",
	RunE:              Execute,
	PersistentPreRunE: utils.MakePersistentPreRunE(initiator.ValidateOptions),
}

func init() {
	initiator.AddPersistentFlags(CancelCmd)
	initiator.AddPersistentFlagCompletions(CancelCmd)
	initiator.AddPersistentFlagCompletions(cancelorder.CancelOrderCmd)

	CancelCmd.AddCommand(cancelorder.CancelOrderCmd)
	CancelCmd.AddCommand(masscancelorder.MassCancelOrderCmd)
}

func Execute(cmd *cobra.Command, args []string) error {
	return nil
}
