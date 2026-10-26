package status

import (
	"github.com/spf13/cobra"

	statusorder "sylr.dev/fix/v2/cmd/status/order"
	statussecurity "sylr.dev/fix/v2/cmd/status/security"
	statustradingsession "sylr.dev/fix/v2/cmd/status/tradingsession"
	"sylr.dev/fix/v2/initiator"
)

// StatusCmd groups the commands requesting the status of orders,
// securities and trading sessions.
var StatusCmd = &cobra.Command{
	Use:               "status",
	Short:             "Send a status request FIX message",
	Long:              "Send a status request FIX message after initiating a session with a FIX acceptor.",
	PersistentPreRunE: initiator.ValidateOptions,
}

func init() {
	initiator.AddPersistentFlags(StatusCmd)
	if err := initiator.AddPersistentFlagCompletions(StatusCmd); err != nil {
		panic(err)
	}

	StatusCmd.AddCommand(statusorder.StatusOrderCmd)
	StatusCmd.AddCommand(statussecurity.StatusSecurityCmd)
	StatusCmd.AddCommand(statustradingsession.StatusTradingSessionCmd)
}
