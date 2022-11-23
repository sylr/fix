package cancelmassorder

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/cli/complete"
	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/initiator/application"
	"sylr.dev/fix/pkg/utils"
)

var (
	optionOrderSide          string
	optionOrderSymbol        string
	optionExecReportsTimeout time.Duration
)

var MassCancelOrderCmd = &cobra.Command{
	Use:               "mass",
	Short:             "Cancel mass order",
	Long:              "Send mass cancel order request after initiating a session with a FIX acceptor.",
	Args:              cobra.ExactArgs(0),
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := utils.ValidateRequiredFlags(cmd); err != nil {
			return err
		}

		if err := Validate(cmd, args); err != nil {
			return err
		}

		if cmd.HasParent() {
			parent := cmd.Parent()
			if parent.PersistentPreRunE != nil {
				return parent.PersistentPreRunE(cmd, args)
			}
		}
		return nil
	},
	RunE: Execute,
}

func init() {
	MassCancelOrderCmd.Flags().StringVar(&optionOrderSide, "side", "", "Order side (buy, sell ... etc)")
	MassCancelOrderCmd.Flags().StringVar(&optionOrderSymbol, "symbol", "", "Order symbol")

	MassCancelOrderCmd.Flags().DurationVar(&optionExecReportsTimeout, "exec-reports-timeout", 5*time.Second, "Log out if execution reports not received within timeout (0s wait indefinitely)")

	MassCancelOrderCmd.MarkFlagRequired("side")
	MassCancelOrderCmd.MarkFlagRequired("symbol")

	MassCancelOrderCmd.RegisterFlagCompletionFunc("side", complete.OrderSide)
}

func Validate(cmd *cobra.Command, args []string) error {
	sides := utils.PrettyOptionValues(dict.OrderSides)
	search := utils.Search(sides, strings.ToLower(optionOrderSide))
	if search < 0 {
		return errors.OptionOrderSideUnknown
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

	app := application.NewCancelOrder()
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
	var sessionId quickfix.SessionID
	var ok bool
	select {
	case <-time.After(timeout):
		return errors.ConnectionTimeout
	case sessionId, ok = <-app.Connected:
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
	err = quickfix.SendToTarget(order, sessionId)
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
			logger.Warn().Msgf("Timeout while expecting execution report or cancel request reject")
			break LOOP

		case msg, ok := <-app.FromAppMessages:
			if !ok {
				logger.Error().Msg("Bad FIX application message received")
				break LOOP
			}

			if err := processReponse(app, msg); err != nil {
				if errors.Is(err, quickfix.InvalidMessageType()) {
					continue LOOP
				}
				return err
			}

			break LOOP
		}
	}

	return nil
}

func processReponse(app *application.CancelOrder, msg *quickfix.Message) error {
	msgType := field.MsgTypeField{}
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
	}

	if msgType.Value() == enum.MsgType_REJECT || msgType.Value() == enum.MsgType_BUSINESS_MESSAGE_REJECT {
		return makeError(errors.FixOrderRejected)
	}
	if msgType.Value() == enum.MsgType_ORDER_MASS_CANCEL_REPORT {
		app.WriteMessageBodyAsTable(os.Stdout, msg)
		resp := field.MassCancelResponseField{}
		if err = msg.Body.GetField(tag.MassCancelResponse, &resp); err != nil {
			return err
		}
		if resp.Value() == enum.MassCancelResponse_CANCEL_REQUEST_REJECTED {
			return makeError(errors.FixOrderRejected)
		}
	} else {
		return quickfix.InvalidMessageType()
	}
	return nil
}

func buildMessage(session config.Session) (quickfix.Messagable, error) {
	eside, err := dict.OrderSideStringToEnum(optionOrderSide)
	if err != nil {
		return nil, err
	}

	switch session.BeginString {
	case quickfix.BeginStringFIXT11:
		switch session.DefaultApplVerID {
		case "FIX.5.0SP2":
			message := quickfix.NewMessage()
			message.Header.Set(field.NewMsgType(enum.MsgType_ORDER_MASS_CANCEL_REQUEST))
			message.Body.Set(field.NewClOrdID(optionOrderSymbol + "_CANCELREQ"))
			message.Body.Set(field.NewMassCancelRequestType(enum.MassCancelRequestType_CANCEL_ORDERS_FOR_A_SECURITY))
			message.Body.Set(field.NewTransactTime(time.Now()))
			message.Body.Set(field.NewSide(eside))
			message.Body.Set(field.NewSymbol(optionOrderSymbol))
			return message, nil

		default:
			return nil, errors.FixVersionNotImplemented
		}

	default:
		return nil, errors.FixVersionNotImplemented
	}
}
