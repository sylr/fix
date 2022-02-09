package cancelorder

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/pkg/errors"
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
		err := Validate(cmd, args)
		if err != nil {
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
	CancelOrderCmd.Flags().StringVar(&optionSide, "side", "", "Order side (buy, sell ... etc)")
	CancelOrderCmd.Flags().StringVar(&optionSymbol, "symbol", "", "Order symbol")
	CancelOrderCmd.Flags().Int64Var(&optionQuantity, "quantity", 1, "Order quantity")
	CancelOrderCmd.Flags().StringVar(&optionID, "id", "", "Order id")

	//CancelOrderCmd.MarkFlagRequired("side")
	//CancelOrderCmd.MarkFlagRequired("symbol")
	//CancelOrderCmd.MarkFlagRequired("quantity")
	//CancelOrderCmd.MarkFlagRequired("id")

	//CancelOrderCmd.RegisterFlagCompletionFunc("side", complete.OrderSide)
}

func Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func Execute(cmd *cobra.Command, args []string) error {
	return errors.NotImplemented
}
