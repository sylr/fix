package marketdata

import (
	"github.com/spf13/cobra"

	marketdatarequest "sylr.dev/fix/v2/cmd/marketdata/request"
	"sylr.dev/fix/v2/initiator"
)

// MarketDataCmd groups the market data commands.
var MarketDataCmd = &cobra.Command{
	Use:               "marketdata",
	Short:             "Send a MarketData FIX message",
	Long:              "Send a MarketData FIX message after initiating a session with a FIX acceptor.",
	PersistentPreRunE: initiator.ValidateOptions,
}

func init() {
	initiator.AddPersistentFlags(MarketDataCmd)
	if err := initiator.AddPersistentFlagCompletions(MarketDataCmd); err != nil {
		panic(err)
	}

	MarketDataCmd.AddCommand(marketdatarequest.MarketDataRequestCmd)
}
