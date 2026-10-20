package new

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/cmd/new/order"
	"sylr.dev/fix/v2/cmd/new/quote"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/utils"
)

// NewCmd represents the buy command
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Send a new FIX message",
	Long:  "Send a new FIX message after initiating a session with a FIX acceptor.",
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
	initiator.AddPersistentFlags(NewCmd)
	initiator.AddPersistentFlagCompletions(NewCmd)
	initiator.AddPersistentFlagCompletions(neworder.NewOrderCmd)
	initiator.AddPersistentFlagCompletions(newquote.NewQuoteCmd)

	NewCmd.AddCommand(neworder.NewOrderCmd)
	NewCmd.AddCommand(newquote.NewQuoteCmd)
}
