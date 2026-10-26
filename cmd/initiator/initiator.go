package initiator

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/utils"
)

var (
	optionNoLogout bool
)

var InitiatorCmd = &cobra.Command{
	Use:               "initiator",
	Short:             "Launch a FIX initiator",
	Long:              "Launch a FIX initiator and wait for messages.",
	RunE:              Execute,
	PersistentPreRunE: utils.MakePersistentPreRunE(initiator.ValidateOptions),
}

func init() {
	initiator.AddPersistentFlags(InitiatorCmd)
	if err := initiator.AddPersistentFlagCompletions(InitiatorCmd); err != nil {
		panic(err)
	}

	InitiatorCmd.Flags().BoolVar(&optionNoLogout, "no-logout", false, "Do not send FIX logout")
}

func Execute(cmd *cobra.Command, args []string) error {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// A second signal kills.
	context.AfterFunc(ctx, stop)

	c, err := initiator.Connect(ctx, initiator.Options{})
	if err != nil {
		return err
	}
	defer c.Close(!optionNoLogout)

	return c.Await(ctx, initiator.AwaitOptions{What: "messages"}, func(msg *fix.Message) (bool, error) {
		utils.WriteMessageBodyAsTable(os.Stdout, msg, c.Dictionary)
		return false, nil
	})
}
