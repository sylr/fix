package status_order

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fixt11"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/cli/complete"
	"sylr.dev/fix/pkg/cli/options"
	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/initiator/application"
	"sylr.dev/fix/pkg/utils"
)

var (
	optionOrderID            string
	optionOrigClOrdID        string
	optionSymbol             string
	optionSide               string
	partyIdOptions           *options.PartyIdOptions
	optionExecReportsTimeout time.Duration
)

var StatusOrderCmd = &cobra.Command{
	Use:               "order",
	Short:             "order status",
	Long:              "Send an Order Status Request after initiating a session with a FIX acceptor.",
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
	StatusOrderCmd.Flags().StringVar(&optionOrderID, "id", "", "Order id (required if origclordid empty)")
	StatusOrderCmd.Flags().StringVar(&optionOrigClOrdID, "origclordid", "", "Client order id (required if id empty)")
	StatusOrderCmd.Flags().StringVar(&optionSide, "side", "", "Order side (buy, sell ... etc)")
	StatusOrderCmd.Flags().StringVar(&optionSymbol, "symbol", "", "Order symbol")

	partyIdOptions = options.NewPartyIdOptions(StatusOrderCmd)

	StatusOrderCmd.Flags().DurationVar(&optionExecReportsTimeout, "exec-reports-timeout", 5*time.Second, "Log out if execution reports not received within timeout (0s wait indefinitely)")

	StatusOrderCmd.MarkFlagRequired("side")
	StatusOrderCmd.MarkFlagRequired("symbol")

	StatusOrderCmd.RegisterFlagCompletionFunc("side", complete.OrderSide)
	StatusOrderCmd.RegisterFlagCompletionFunc("type", complete.OrderType)
	StatusOrderCmd.RegisterFlagCompletionFunc("expiry", complete.OrderTimeInForce)
	StatusOrderCmd.RegisterFlagCompletionFunc("symbol", cobra.NoFileCompletions)
}

func Validate(cmd *cobra.Command, args []string) error {
	if len(optionSymbol) == 0 {
		return fmt.Errorf("%w: --symbol can not be empty", errors.Options)
	}

	return partyIdOptions.Validate()
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
	initiatorConfig, err := context.GetInitiator()
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

	app := application.NewNewOrder()
	app.Logger = logger
	app.Settings = settings
	app.TransportDataDictionary = transportDict
	app.AppDataDictionary = appDict

	var quickfixLogger *zerolog.Logger
	if options.QuickFixLogging {
		quickfixLogger = logger
	}

	// Choose right timeout cli option > config > default value (5s)
	var timeout time.Duration
	if options.Timeout != time.Duration(0) {
		timeout = options.Timeout
	} else if initiatorConfig.SocketTimeout != time.Duration(0) {
		timeout = initiatorConfig.SocketTimeout
	} else {
		timeout = 5 * time.Second
	}

	init, err := initiator.Initiate(app, settings, quickfixLogger)
	if err != nil {
		return err
	}

	// Start session
	if err = init.Start(); err != nil {
		return err
	}

	defer func() {
		app.Stop()
		init.Stop()
	}()

	// Wait for session connection
	select {
	case <-time.After(timeout):
		return errors.ConnectionTimeout
	case _, ok := <-app.Connected:
		if !ok {
			return errors.FixLogout
		}
	}

	// Prepare order
	order, err := buildMessage(*session)
	if err != nil {
		return err
	}

	// Send the order
	err = quickfix.Send(order)
	if err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var waitTimeout <-chan time.Time
	if optionExecReportsTimeout > 0 {
		waitTimeout = time.After(optionExecReportsTimeout)
	} else {
		waitTimeout = make(<-chan time.Time)
	}

LOOP:
	for {
		select {
		case signal := <-interrupt:
			logger.Debug().Msgf("Received signal: %s", signal)
			break LOOP

		case <-waitTimeout:
			logger.Warn().Msgf("Timeout while expecting execution report")
			break LOOP

		case msg, ok := <-app.FromAppMessages:
			if !ok {
				break LOOP
			}

			if err := processResponse(app, msg); err != nil {
				if errors.Is(err, quickfix.InvalidMessageType()) {
					continue LOOP
				}

				return err
			}
		}

		logger.Debug().Msgf("Exiting response loop")
		break LOOP
	}

	return nil
}

func buildMessage(session config.Session) (quickfix.Messagable, error) {
	eSide, err := dict.OrderSideStringToEnum(optionSide)
	if err != nil {
		return nil, err
	}

	// Message
	message := quickfix.NewMessage()
	header := fixt11.NewHeader(&message.Header)

	switch session.BeginString {
	case quickfix.BeginStringFIXT11:
		switch session.DefaultApplVerID {
		case "FIX.5.0SP2":
			header.Set(field.NewMsgType(enum.MsgType_ORDER_STATUS_REQUEST))
			message.Body.Set(field.NewClOrdID(uuid.NewString()))
			if len(optionOrderID) > 0 {
				message.Body.Set(field.NewOrderID(optionOrderID))
			}
			if len(optionOrigClOrdID) > 0 {
				message.Body.Set(field.NewOrigClOrdID(optionOrigClOrdID))
			}
			message.Body.Set(field.NewSide(eSide))
			message.Body.Set(field.NewSymbol(optionSymbol))
			partyIdOptions.EnrichMessageBody(&message.Body, session)

		default:
			return nil, errors.FixVersionNotImplemented
		}
	default:
		return nil, errors.FixVersionNotImplemented
	}

	utils.QuickFixMessagePartSetString(&message.Header, session.TargetCompID, field.NewTargetCompID)
	utils.QuickFixMessagePartSetString(&message.Header, session.TargetSubID, field.NewTargetSubID)
	utils.QuickFixMessagePartSetString(&message.Header, session.SenderCompID, field.NewSenderCompID)
	utils.QuickFixMessagePartSetString(&message.Header, session.SenderSubID, field.NewSenderSubID)

	return message, nil
}

func processResponse(app *application.NewOrder, msg *quickfix.Message) error {
	msgType := field.MsgTypeField{}
	ordStatus := field.OrdStatusField{}
	text := field.TextField{}

	// Text
	if msg.Body.Has(tag.Text) {
		if err := msg.Body.GetField(tag.Text, &text); err != nil {
			return err
		}
	}

	makeError := func(errType error) error {
		if len(text.String()) > 0 {
			return fmt.Errorf("%w: %s", errType, text.String())
		} else {
			return errType
		}
	}

	// MsgType
	err := msg.Header.GetField(tag.MsgType, &msgType)
	if err != nil {
		return err
	} else if msgType.Value() == enum.MsgType_REJECT ||
		msgType.Value() == enum.MsgType_BUSINESS_MESSAGE_REJECT ||
		msgType.Value() == enum.MsgType_ORDER_CANCEL_REJECT {
		return makeError(errors.FixOrderRejected)
	} else if msgType.Value() != enum.MsgType_EXECUTION_REPORT {
		return quickfix.InvalidMessageType()
	}

	// OrdStatus
	err = msg.Body.GetField(tag.OrdStatus, &ordStatus)
	if err != nil {
		return err
	}

	app.WriteMessageBodyAsTable(os.Stdout, msg)
	return nil
}
