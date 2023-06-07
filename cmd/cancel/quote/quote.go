package cancelorder

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
	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/initiator/application"
	"sylr.dev/fix/pkg/utils"
)

var (
	optionOrderSymbols       []string
	optionQuoteID            string
	optionExecReportsTimeout time.Duration
	partyIdOptions           *options.PartyIdOptions
)

var CancelQuoteCmd = &cobra.Command{
	Use:               "quote",
	Short:             "Cancel quote",
	Long:              "Send a cancel quote request after initiating a session with a FIX acceptor.",
	Args:              cobra.ExactArgs(0),
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: utils.MakePersistentPreRunE(Validate),
	RunE:              Execute,
}

func init() {
	CancelQuoteCmd.Flags().StringVar(&optionQuoteID, "id", "", "Quote id")

	CancelQuoteCmd.Flags().StringSliceVar(&optionOrderSymbols, "symbols", []string{}, "Order symbol")
	CancelQuoteCmd.Flags().DurationVar(&optionExecReportsTimeout, "exec-reports-timeout", 5*time.Second, "Log out if execution reports not received within timeout (0s wait indefinitely)")

	partyIdOptions = options.NewPartyIdOptions(CancelQuoteCmd)

	CancelQuoteCmd.RegisterFlagCompletionFunc("side", complete.OrderSide)
}

func Validate(cmd *cobra.Command, args []string) error {
	if len(optionQuoteID) == 0 {
		optionQuoteID = uuid.NewString()
	}
	if len(optionOrderSymbols) == 0 {
		return errors.OptionsNoSymbolGiven
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

	// Prepare cancel message
	cancelMsg, err := buildMessage(*session)
	if err != nil {
		return err
	}

	// Send the cancel message
	err = quickfix.SendToTarget(cancelMsg, sessionId)
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
				break LOOP
			}

			if err := processResponse(app, msg); err != nil {
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

func buildMessage(session config.Session) (quickfix.Messagable, error) {
	// Message
	message := quickfix.NewMessage()
	header := fixt11.NewHeader(&message.Header)

	quoteId := optionQuoteID
	if len(quoteId) == 0 {
		quoteId = uuid.NewString()
	}

	switch session.BeginString {
	case quickfix.BeginStringFIXT11:
		switch session.DefaultApplVerID {
		case "FIX.5.0SP2":
			header.Set(field.NewMsgType(enum.MsgType_QUOTE_CANCEL))
			message.Body.Set(field.NewQuoteID(quoteId))
			message.Body.Set(field.NewQuoteMsgID("msg_" + quoteId))
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

	message.Body.Set(field.NewQuoteCancelType(enum.QuoteCancelType_CANCEL_FOR_ONE_OR_MORE_SECURITIES))
	instruments := quickfix.NewRepeatingGroup(tag.NoQuoteEntries, nil)
	for _, symbol := range optionOrderSymbols {
		inst := instruments.Add()
		inst.Set(field.NewSymbol(symbol))
	}
	message.Body.SetGroup(instruments)

	return message, nil
}

func processResponse(app *application.CancelOrder, msg *quickfix.Message) error {
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

	if msgType.Value() == enum.MsgType_ORDER_CANCEL_REJECT {
		app.WriteMessageBodyAsTable(os.Stdout, msg)
		return makeError(errors.FixOrderRejected)
	}

	if msgType.Value() == enum.MsgType_EXECUTION_REPORT {
		app.WriteMessageBodyAsTable(os.Stdout, msg)
		ordStatus := field.OrdStatusField{}
		if err = msg.Body.GetField(tag.OrdStatus, &ordStatus); err != nil {
			return err
		}
		if ordStatus.Value() == enum.OrdStatus_REJECTED {
			return makeError(errors.FixOrderRejected)
		}
	} else {
		return quickfix.InvalidMessageType()
	}

	return nil
}
