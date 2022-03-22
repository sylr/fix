package marketdata

import (
	"github.com/spf13/cobra"

	markedatarequest "sylr.dev/fix/cmd/marketdata/request"
	"sylr.dev/fix/pkg/initiator"
)

// MarketDataCmd represents the buy command
var MarketDataCmd = &cobra.Command{
	Use:   "marketdata",
	Short: "Send a MarketData FIX message",
	Long:  "Send a MarketData FIX message after initiating a sesion with a FIX acceptor.",
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
	initiator.AddPersistentFlags(MarketDataCmd)

	MarketDataCmd.AddCommand(markedatarequest.MarketDataRequestCmd)
}
