package status

import (
	"github.com/spf13/cobra"

	status_security "sylr.dev/fix/cmd/status/security"
	status_tradingsession "sylr.dev/fix/cmd/status/tradingsession"
	"sylr.dev/fix/pkg/initiator"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
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
	initiator.AddPersistentFlags(StatusCmd)

	StatusCmd.AddCommand(status_tradingsession.StatusTradingSessionCmd)
	StatusCmd.AddCommand(status_security.StatusSecurityCmd)
}
