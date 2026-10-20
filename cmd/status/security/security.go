package status_security

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/cli/complete"
	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/utils"
)

var (
	optionSecurityStatReqID string
	optionSubType           string
	optionSymbol            string

	SubType fix50sp2.SubscriptionRequestType
)

var StatusSecurityCmd = &cobra.Command{
	Use:               "security",
	Short:             "security status",
	Long:              "Send a security Session Status Request after initiating a session with a FIX acceptor.",
	Args:              cobra.ExactArgs(0),
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := Validate(cmd, args)
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
	RunE: Execute,
}

func init() {
	StatusSecurityCmd.Flags().StringVar(&optionSymbol, "symbol", "", "Symbol")
	StatusSecurityCmd.Flags().StringVar(&optionSubType, "subscription-type", "snapshot", "Subscription type")
	StatusSecurityCmd.Flags().StringVar(&optionSecurityStatReqID, "security-status-request-id", uuid.NewString(), "Security Status Request id")
	StatusSecurityCmd.RegisterFlagCompletionFunc("subscription-type", complete.SubscriptionRequestTypes)
}

func Validate(cmd *cobra.Command, args []string) error {
	if len(optionSymbol) == 0 {
		return fmt.Errorf("%w: --symbol can not be empty", errors.Options)
	}

	if len(optionSecurityStatReqID) == 0 {
		return fmt.Errorf("%w: --security-status-request-id can not be empty", errors.Options)
	}

	var ok bool
	if SubType, ok = dict.SubscriptionRequestTypes[strings.ToUpper(optionSubType)]; !ok {
		return fmt.Errorf("%w: unknown subscription type `%s`", errors.Options, optionSubType)
	}

	return nil
}

func Execute(cmd *cobra.Command, args []string) error {
	logger := config.GetLogger()

	c, err := initiator.Connect(cmd.Context(), initiator.Options{
		MsgTypes: []string{string(fix50sp2.MsgTypeSecuritystatus)},
	})
	if err != nil {
		return err
	}
	defer c.Close(true)

	// Prepare security Status Request
	ssr, err := buildMessage(c)
	if err != nil {
		return err
	}

	// Send the security status request
	if _, err = c.Send(cmd.Context(), ssr); err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

LOOP:
	for {
		select {
		case signal := <-interrupt:
			logger.Debug().Msgf("Received signal: %s", signal)

			break LOOP

		case <-c.Done():
			break LOOP

		case responseMessage := <-c.Messages():
			utils.WriteMessageBodyAsTable(os.Stdout, responseMessage, c.Dictionary)

			if SubType != fix50sp2.SubscriptionRequestTypeSnapshotPlusUpdates {
				break LOOP
			}
		}
	}

	return nil
}

func buildMessage(c *initiator.Client) (session.Outgoing, error) {
	return fix50sp2.BuildSecurityStatusRequest(c.BeginString()).
		SetSecurityStatusReqID(optionSecurityStatReqID).
		SetSymbol(optionSymbol).
		SetSubscriptionRequestType(dict.SubscriptionRequestTypes[strings.ToUpper(optionSubType)]), nil
}
