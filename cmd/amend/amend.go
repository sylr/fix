package amend

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/cmd/amend/order"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/utils"
)

// AmendCmd represents the buy command
var AmendCmd = &cobra.Command{
	Use:   "amend",
	Short: "Send a cancel/replace FIX message",
	Long:  "Send a cancel/replace FIX message after initiating a session with a FIX acceptor.",
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
	initiator.AddPersistentFlags(AmendCmd)
	initiator.AddPersistentFlagCompletions(AmendCmd)
	initiator.AddPersistentFlagCompletions(amendorder.AmendOrderCmd)

	AmendCmd.AddCommand(amendorder.AmendOrderCmd)
}
