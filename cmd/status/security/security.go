package status_security

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fixt11"
	"github.com/quickfixgo/quickfix"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/cli/complete"
	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/initiator/application"
	"sylr.dev/fix/pkg/utils"
)

var (
	optionSecurityStatReqID string
	optionSubType           string
	optionSymbol            string
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

	if _, ok := dict.SubscriptionRequestTypes[strings.ToUpper(optionSubType)]; !ok {
		return fmt.Errorf("%w: unkonwn subscription type `%s`", errors.Options, optionSubType)
	}

	return nil
}

func Execute(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()
	logger := config.GetLogger()

	context, err := config.GetCurrentContext()
	if err != nil {
		return err
	}

	sessions, err := context.GetSessions()
	if err != nil {
		return err
	}

	session := sessions[0]
	initiatior, err := context.GetInitiator()
	if err != nil {
		return err
	}

	transportDict, appDict, err := session.GetFIXDictionaries()
	if err != nil {
		return err
	}

	settings, err := context.ToQuickFixInitiatorSettings()
	if err != nil {
		return err
	}

	app := application.NewSecurityStatusRequest()
	app.Logger = logger
	app.Settings = settings
	app.TransportDataDictionary = transportDict
	app.AppDataDictionary = appDict

	var quickfixLogger *zerolog.Logger
	if options.QuickFixLogging {
		quickfixLogger = logger
	}

	init, err := initiator.Initiate(app, settings, quickfixLogger)
	if err != nil {
		return err
	}

	// Start session
	err = init.Start()
	if err != nil {
		return err
	}

	defer init.Stop()

	// Choose right timeout cli option > config > default value (5s)
	var timeout time.Duration
	if options.Timeout != time.Duration(0) {
		timeout = options.Timeout
	} else if initiatior.SocketTimeout != time.Duration(0) {
		timeout = initiatior.SocketTimeout
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

	// Prepare security Status Request
	tssr, err := buildMessage(*session)
	if err != nil {
		return err
	}

	// Send the security status request
	err = quickfix.Send(tssr)
	if err != nil {
		return err
	}

	// Wait for the order response
	var ok bool
	var responseMessage *quickfix.Message

	select {
	case <-time.After(timeout):
		return errors.ResponseTimeout
	case responseMessage, ok = <-app.FromAdminChan:
		if !ok {
			return errors.FixLogout
		}
	case responseMessage, ok = <-app.FromAppChan:
		if !ok {
			return errors.FixLogout
		}
	}

	app.WriteMessageBodyAsTable(os.Stdout, responseMessage)

	return nil
}

func buildMessage(session config.Session) (quickfix.Messagable, error) {
	// Message
	message := quickfix.NewMessage()
	header := fixt11.NewHeader(&message.Header)
	header.Set(field.NewMsgType(enum.MsgType_SECURITY_STATUS_REQUEST))

	utils.QuickFixMessagePartSetString(&message.Header, session.TargetCompID, field.NewTargetCompID)
	utils.QuickFixMessagePartSetString(&message.Header, session.TargetSubID, field.NewTargetSubID)
	utils.QuickFixMessagePartSetString(&message.Header, session.SenderCompID, field.NewSenderCompID)
	utils.QuickFixMessagePartSetString(&message.Header, session.SenderSubID, field.NewSenderSubID)

	utils.QuickFixMessagePartSetString(&message.Body, optionSecurityStatReqID, field.NewSecurityStatusReqID)
	utils.QuickFixMessagePartSetString(&message.Body, optionSymbol, field.NewSymbol)
	utils.QuickFixMessagePartSetString(&message.Body, dict.SubscriptionRequestTypes[strings.ToUpper(optionSubType)], field.NewSubscriptionRequestType)

	return message, nil
}
