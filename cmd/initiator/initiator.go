package initiator

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/utils"
)

var (
	optionNoLogout bool
)

var InitiatorCmd = &cobra.Command{
	Use:   "initiator",
	Short: "Launch a FIX initiator",
	Long:  "Launch a FIX initiator and wait for messages.",
	RunE:  Execute,
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
	initiator.AddPersistentFlags(InitiatorCmd)
	initiator.AddPersistentFlagCompletions(InitiatorCmd)

	InitiatorCmd.Flags().BoolVar(&optionNoLogout, "no-logout", false, "Do not send FIX logout")
}

func Execute(cmd *cobra.Command, args []string) error {
	logger := config.GetLogger()

	c, err := initiator.Connect(cmd.Context(), initiator.Options{})
	if err != nil {
		return err
	}
	defer c.Close(!optionNoLogout)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

LOOP:
	for {
		select {
		case signal := <-interrupt:
			logger.Debug().Msgf("Received signal: %s", signal)
			break LOOP
		case <-c.Done():
			return errors.FixLogout
		case msg := <-c.Messages():
			utils.WriteMessageBodyAsTable(os.Stdout, msg, c.Dictionary)
		}
	}

	return nil
}
