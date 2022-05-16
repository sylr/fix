package tradingsessionstatus

import (
	"github.com/spf13/cobra"

	tradingsessionstatusrequest "sylr.dev/fix/cmd/tradingsessionstatus/request"
	"sylr.dev/fix/pkg/initiator"
)
// TradingSessionStatusCmd represents the buy command
var TradingSessionStatusCmd = &cobra.Command{
	Use:   "tradingsessionstatus",
	Short: "Send a TradingSessionStatus FIX message",
	Long:  "Send a TradingSessionStatus FIX message after initiating a sesion with a FIX acceptor.",
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
	initiator.AddPersistentFlags(TradingSessionStatusCmd)

	TradingSessionStatusCmd.AddCommand(tradingsessionstatusrequest.TradingSessionStatusRequestCmd)
}


