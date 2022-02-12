package initiator

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/initiator/application"
)

var InitiatorCmd = &cobra.Command{
	Use:   "initiator",
	Short: "Launch a FIX initiator",
	Long:  "Launch a FIX initiator and wait for messages.",
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
	initiator.AddPersistentFlags(InitiatorCmd)
	initiator.AddPersistentFlagCompletions(InitiatorCmd)
}

func Execute(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()
	logger := config.GetLogger()

	context, err := config.GetCurrentContext()
	if err != nil {
		return err
	}

	session, err := context.GetSession()
	if err != nil {
		return err
	}

	acceptor, err := context.GetAcceptor()
	if err != nil {
		return err
	}

	settings, err := context.ToQuickFixSettings()
	if err != nil {
		return err
	}

	transportDict, appDict, err := session.GetFIXDictionaries()
	if err != nil {
		return err
	}

	app := application.NewInitiator()
	app.Settings = settings
	app.TransportDataDictionary = transportDict
	app.AppDataDictionary = appDict
	app.Logger = logger

	init, err := initiator.Initiate(app, settings)
	if err != nil {
		return err
	}

	err = init.Start()
	if err != nil {
		return err
	}

	defer init.Stop()

	// Choose right timeout cli option > config > default value (5s)
	var timeout time.Duration
	if options.Timeout != time.Duration(0) {
		timeout = options.Timeout
	} else if acceptor.SocketTimeout != time.Duration(0) {
		timeout = acceptor.SocketTimeout
	} else {
		timeout = 5 * time.Second
	}

	// Wait for session connection
	select {
	case <-time.After(timeout):
		return errors.ConnectionTimeout
	case _, ok := <-app.Connected:
		if !ok {
			return errors.FixLogout
		}
	}

	message := quickfix.NewMessage()
	message.Header.SetField(tag.MsgType, field.NewMsgType(enum.MsgType_HEARTBEAT))
	err = quickfix.SendToTarget(message, app.SessionID)
	if err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

LOOP:
	for {
		select {
		case <-interrupt:
			break LOOP
		case _, ok := <-app.Messages:
			if !ok {
				return errors.FixLogout
			}
		}
	}

	return nil
}
