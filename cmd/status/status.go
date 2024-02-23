package status

import (
	"github.com/spf13/cobra"

	status_order "sylr.dev/fix/cmd/status/order"
	status_security "sylr.dev/fix/cmd/status/security"
	status_tradingsession "sylr.dev/fix/cmd/status/tradingsession"
	"sylr.dev/fix/pkg/initiator"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Send a TradingSessionStatus FIX message",
	Long:  "Send a TradingSessionStatus FIX message after initiating a session with a FIX acceptor.",
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
	initiator.AddPersistentFlags(StatusCmd)
	initiator.AddPersistentFlagCompletions(StatusCmd)
	initiator.AddPersistentFlagCompletions(status_order.StatusOrderCmd)
	initiator.AddPersistentFlagCompletions(status_security.StatusSecurityCmd)
	initiator.AddPersistentFlagCompletions(status_tradingsession.StatusTradingSessionCmd)

	StatusCmd.AddCommand(status_order.StatusOrderCmd)
	StatusCmd.AddCommand(status_security.StatusSecurityCmd)
	StatusCmd.AddCommand(status_tradingsession.StatusTradingSessionCmd)
}
