package new

import (
	"github.com/spf13/cobra"

	neworder "sylr.dev/fix/cmd/new/order"
	"sylr.dev/fix/pkg/initiator"
)

// NewCmd represents the buy command
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Send a new FIX message",
	Long:  "Send a new FIX message after initiating a sesion with a FIX acceptor.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := initiator.ValidateOptions(cmd, args)
		if err != nil {
			return err
		}

		if cmd.HasParent() {
			parent := cmd.Parent()
			if parent.PersistentPreRunE != nil {
				err = parent.PersistentPreRunE(parent, args)
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func init() {
	initiator.AddPersistentFlags(NewCmd)
	initiator.AddPersistentFlagCompletions(NewCmd)
	initiator.AddPersistentFlagCompletions(neworder.NewOrderCmd)

	NewCmd.AddCommand(neworder.NewOrderCmd)
}
