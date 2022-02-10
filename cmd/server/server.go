package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/quickfixgo/quickfix"
	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/application"
	"sylr.dev/fix/pkg/initiator"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch a FIX server",
	Long:  "Launch a FIX server accepting orders.",
	RunE:  Execute,
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
	initiator.AddPersistentFlags(ServerCmd)
	initiator.AddPersistentFlagCompletions(ServerCmd)
}

func Execute(cmd *cobra.Command, args []string) error {
	logger := config.GetLogger()

	context, err := config.GetCurrentContext()
	if err != nil {
		return err
	}

	settings, err := context.ToQuickFixSettings()
	if err != nil {
		return err
	}

	app, err := application.NewServer()
	if err != nil {
		return err
	}

	app.Logger = logger

	acceptor, err := quickfix.NewAcceptor(app, quickfix.NewMemoryStoreFactory(), settings, quickfix.NewScreenLogFactory())
	if err != nil {
		return err
	}

	err = acceptor.Start()
	if err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	acceptor.Stop()
	os.Exit(0)

	return nil
}
