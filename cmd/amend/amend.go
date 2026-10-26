package amend

import (
	"github.com/spf13/cobra"

	amendorder "sylr.dev/fix/v2/cmd/amend/order"
	"sylr.dev/fix/v2/initiator"
)

// AmendCmd groups the commands amending orders.
var AmendCmd = &cobra.Command{
	Use:               "amend",
	Short:             "Send a cancel/replace FIX message",
	Long:              "Send a cancel/replace FIX message after initiating a session with a FIX acceptor.",
	PersistentPreRunE: initiator.ValidateOptions,
}

func init() {
	initiator.AddPersistentFlags(AmendCmd)
	if err := initiator.AddPersistentFlagCompletions(AmendCmd); err != nil {
		panic(err)
	}

	AmendCmd.AddCommand(amendorder.AmendOrderCmd)
}
