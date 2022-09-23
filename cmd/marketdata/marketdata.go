package marketdata

import (
	"github.com/spf13/cobra"

	markedatarequest "sylr.dev/fix/cmd/marketdata/request"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/utils"
)

// MarketDataCmd represents the buy command
var MarketDataCmd = &cobra.Command{
	Use:   "marketdata",
	Short: "Send a MarketData FIX message",
	Long:  "Send a MarketData FIX message after initiating a sesion with a FIX acceptor.",
	Args:  cobra.ExactValidArgs(1),
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
	initiator.AddPersistentFlags(MarketDataCmd)
	initiator.AddPersistentFlags(markedatarequest.MarketDataRequestCmd)

	if err := initiator.AddPersistentFlagCompletions(MarketDataCmd); err != nil {
		panic(err)
	}
	if err := initiator.AddPersistentFlagCompletions(markedatarequest.MarketDataRequestCmd); err != nil {
		panic(err)
	}

	MarketDataCmd.AddCommand(markedatarequest.MarketDataRequestCmd)
}
