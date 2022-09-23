package cancelorder

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/utils"
)

var (
	optionSide     string
	optionSymbol   string
	optionID       string
	optionQuantity int64
)

var CancelOrderCmd = &cobra.Command{
	Use:               "order",
	Short:             "Cancel single order",
	Long:              "Send a cancel order request after initiating a session with a FIX acceptor.",
	Args:              cobra.ExactArgs(0),
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := utils.ValidateRequiredFlags(cmd); err != nil {
			return err
		}

		if err := Validate(cmd, args); err != nil {
			return err
		}

		if cmd.HasParent() {
			parent := cmd.Parent()
			if parent.PersistentPreRunE != nil {
				return parent.PersistentPreRunE(cmd, args)
			}
		}
		return nil
	},
	RunE: Execute,
}

func init() {
	CancelOrderCmd.Flags().StringVar(&optionID, "id", "", "Order id")
}

func Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func Execute(cmd *cobra.Command, args []string) error {
	return errors.NotImplemented
}
