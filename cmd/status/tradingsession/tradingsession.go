package status_tradingsession

import (
	"fmt"
	"os"
	"strings"
	"time"

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
	optionTradingSessionID string
	optionSubType          string
)

var StatusTradingSessionCmd = &cobra.Command{
	Use:               "tradingsession",
	Short:             "trading session status",
	Long:              "Send a Trading Session Status Request after initiating a session with a FIX acceptor.",
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
	StatusTradingSessionCmd.Flags().StringVar(&optionTradingSessionID, "trading-session-id", uuid.NewString(), "Trading session ID")
	StatusTradingSessionCmd.Flags().StringVar(&optionSubType, "subscription-type", "snapshot", "Subscription type")

	StatusTradingSessionCmd.RegisterFlagCompletionFunc("subscription-type", complete.SubscriptionRequestTypes)
}

func Validate(cmd *cobra.Command, args []string) error {
	if len(optionTradingSessionID) == 0 {
		return fmt.Errorf("%w: --trading-session-id can not be empty", errors.Options)
	}

	if _, ok := dict.SubscriptionRequestTypes[strings.ToUpper(optionSubType)]; !ok {
		return fmt.Errorf("%w: unkonwn subscription type `%s`", errors.Options, optionSubType)
	}

	return nil
}

func Execute(cmd *cobra.Command, args []string) error {
	c, err := initiator.Connect(cmd.Context(), initiator.Options{
		MsgTypes: []string{string(fix50sp2.MsgTypeTradingsessionstatus)},
	})
	if err != nil {
		return err
	}
	defer c.Close(true)

	// Prepare Trading Session Status Request
	tssr, err := buildMessage(c)
	if err != nil {
		return err
	}

	// Send the trading session status request
	if _, err = c.Send(cmd.Context(), tssr); err != nil {
		return err
	}

	// Wait for the response
	select {
	case <-time.After(config.Timeout(0)):
		return errors.ResponseTimeout
	case <-c.Done():
		return errors.FixLogout
	case responseMessage := <-c.Messages():
		utils.WriteMessageBodyAsTable(os.Stdout, responseMessage, c.Dictionary)
	}

	return nil
}

func buildMessage(c *initiator.Client) (session.Outgoing, error) {
	if err := c.Session.CheckSupported(); err != nil {
		return nil, err
	}

	return fix50sp2.BuildTradingSessionStatusRequest(c.BeginString()).
		SetTradSesReqID(optionTradingSessionID).
		SetSubscriptionRequestType(dict.SubscriptionRequestTypes[strings.ToUpper(optionSubType)]), nil
}
