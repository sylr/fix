package cancel

import (
	"github.com/spf13/cobra"

	cancelmass "sylr.dev/fix/v2/cmd/cancel/mass"
	cancelorder "sylr.dev/fix/v2/cmd/cancel/order"
	cancelquote "sylr.dev/fix/v2/cmd/cancel/quote"
	"sylr.dev/fix/v2/initiator"
)

// CancelCmd groups the commands cancelling orders and quotes.
var CancelCmd = &cobra.Command{
	Use:               "cancel",
	Short:             "Send a cancel FIX message",
	Long:              "Send a cancel FIX message after initiating a session with a FIX acceptor.",
	PersistentPreRunE: initiator.ValidateOptions,
}

func init() {
	initiator.AddPersistentFlags(CancelCmd)
	if err := initiator.AddPersistentFlagCompletions(CancelCmd); err != nil {
		panic(err)
	}

	CancelCmd.AddCommand(cancelorder.CancelOrderCmd)
	CancelCmd.AddCommand(cancelmass.MassCancelOrderCmd)
	CancelCmd.AddCommand(cancelquote.CancelQuoteCmd)
}
