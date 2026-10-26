package new

import (
	"github.com/spf13/cobra"

	neworder "sylr.dev/fix/v2/cmd/new/order"
	newquote "sylr.dev/fix/v2/cmd/new/quote"
	"sylr.dev/fix/v2/initiator"
)

// NewCmd groups the commands sending new orders and quotes.
var NewCmd = &cobra.Command{
	Use:               "new",
	Short:             "Send a new FIX message",
	Long:              "Send a new FIX message after initiating a session with a FIX acceptor.",
	PersistentPreRunE: initiator.ValidateOptions,
}

func init() {
	initiator.AddPersistentFlags(NewCmd)
	if err := initiator.AddPersistentFlagCompletions(NewCmd); err != nil {
		panic(err)
	}

	NewCmd.AddCommand(neworder.NewOrderCmd)
	NewCmd.AddCommand(newquote.NewQuoteCmd)
}
