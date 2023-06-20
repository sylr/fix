package amendorder

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fixt11"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"

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
	optionOrderSide, optionOrderType string
	optionOrderSymbol                string
	optionClOrdID                    string
	optionOrderID, optionOrigClOrdID string
	optionOrderExpiry                string
	optionOrderQuantity              int64
	optionOrderPrice                 float64
	partyIdOptions                   *options.PartyIdOptions
	attributeOptions                 *options.AttributeOptions
	optionExecReports                int
	optionExecReportsTimeout         time.Duration
	optionExecReportsTimeoutReset    bool
	optionStopOnFinalState           bool
)

var AmendOrderCmd = &cobra.Command{
	Use:               "order",
	Short:             "New single order",
	Long:              "Send a new single order after initiating a session with a FIX acceptor.",
	Args:              cobra.ExactArgs(0),
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: utils.MakePersistentPreRunE(Validate),
	RunE:              Execute,
}

func init() {
	AmendOrderCmd.Flags().StringVar(&optionOrderID, "id", "", "Order id (required if origclordid empty)")
	AmendOrderCmd.Flags().StringVar(&optionOrigClOrdID, "origclordid", "", "Client order id (required if id empty)")
	AmendOrderCmd.Flags().StringVar(&optionClOrdID, "clordid", "", "Client order id of the amendment")
	AmendOrderCmd.Flags().StringVar(&optionOrderSide, "side", "", "Order side (buy, sell ... etc)")
	AmendOrderCmd.Flags().StringVar(&optionOrderType, "type", "", "Order type (market, limit, stop ... etc)")
	AmendOrderCmd.Flags().StringVar(&optionOrderSymbol, "symbol", "", "Order symbol")
	AmendOrderCmd.Flags().Int64Var(&optionOrderQuantity, "quantity", 1, "Order quantity")
	AmendOrderCmd.Flags().StringVar(&optionOrderExpiry, "expiry", "day", "Order expiry (day, good_till_cancel ... etc)")
	AmendOrderCmd.Flags().Float64Var(&optionOrderPrice, "price", 0.0, "Order price")

	partyIdOptions = options.NewPartyIdOptions(AmendOrderCmd)
	attributeOptions = options.NewAttributeOptions(AmendOrderCmd)

	AmendOrderCmd.Flags().IntVar(&optionExecReports, "exec-reports", 1, "Expect given number of execution reports before logging out (0 wait indefinitely)")
	AmendOrderCmd.Flags().DurationVar(&optionExecReportsTimeout, "exec-reports-timeout", 5*time.Second, "Log out if execution reports not received within timeout (0s wait indefinitely)")
	AmendOrderCmd.Flags().BoolVar(&optionExecReportsTimeoutReset, "exec-reports-timeout-reset", false, "Reset execution reports timeout each time an execution report is received")

	AmendOrderCmd.Flags().BoolVar(&optionStopOnFinalState, "stop-on-final-state", false, "Stop application when receiving an order with a final state")

	AmendOrderCmd.MarkFlagRequired("side")
	AmendOrderCmd.MarkFlagRequired("type")
	AmendOrderCmd.MarkFlagRequired("symbol")
	AmendOrderCmd.MarkFlagRequired("quantity")

	AmendOrderCmd.RegisterFlagCompletionFunc("side", complete.OrderSide)
	AmendOrderCmd.RegisterFlagCompletionFunc("type", complete.OrderType)
	AmendOrderCmd.RegisterFlagCompletionFunc("expiry", complete.OrderTimeInForce)
	AmendOrderCmd.RegisterFlagCompletionFunc("symbol", cobra.NoFileCompletions)
}

func Validate(cmd *cobra.Command, args []string) error {
	sides := utils.PrettyOptionValues(dict.OrderSides)
	search := utils.Search(sides, strings.ToLower(optionOrderSide))
	if search < 0 {
		return errors.OptionOrderSideUnknown
	}

	types := utils.PrettyOptionValues(dict.OrderTypes)
	search = utils.Search(types, strings.ToLower(optionOrderType))
	if search < 0 {
		return errors.OptionOrderTypeUnknown
	}

	if len(optionOrderID) == 0 && len(optionOrigClOrdID) == 0 {
		return fmt.Errorf("order id or original client order id must be filled")
	}

	if len(optionClOrdID) == 0 {
		optionClOrdID = uuid.NewString()
	}

	if strings.ToLower(optionOrderType) == "market" && optionOrderPrice > 0 {
		return errors.OptionsInvalidMarketPrice
	} else if strings.ToLower(optionOrderType) != "market" && optionOrderPrice == 0 {
		return errors.OptionsNoPriceGiven
	}

	if err := attributeOptions.Validate(); err != nil {
		return err
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

	execReports := 0
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
			logger.Warn().Msgf("Timeout while expecting execution reports (%d/%d)", execReports, optionExecReports)
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

			if optionStopOnFinalState && isFinalStatus(msg) {
				break LOOP
			}

			// Reset timeout
			if optionExecReportsTimeoutReset && optionExecReportsTimeout > 0 {
				waitTimeout = time.After(optionExecReportsTimeout)
			}

			execReports = execReports + 1
		}

		if optionExecReports != 0 && execReports >= optionExecReports {
			logger.Debug().Msgf("Exiting response loop, execution reports: %d/%d", execReports, optionExecReports)
			break LOOP
		}
	}

	return nil
}

func buildMessage(session config.Session) (quickfix.Messagable, error) {
	eSide, err := dict.OrderSideStringToEnum(optionOrderSide)
	if err != nil {
		return nil, err
	}

	eType, err := dict.OrderTypeStringToEnum(optionOrderType)
	if err != nil {
		return nil, err
	}

	eExpiry, err := dict.OrderTimeInForceStringToEnum(optionOrderExpiry)
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
			header.Set(field.NewMsgType(enum.MsgType_ORDER_CANCEL_REPLACE_REQUEST))
			if len(optionOrderID) > 0 {
				message.Body.Set(field.NewOrderID(optionOrderID))
			}
			if len(optionOrigClOrdID) > 0 {
				message.Body.Set(field.NewOrigClOrdID(optionOrigClOrdID))
			}
			message.Body.Set(field.NewClOrdID(optionClOrdID))
			message.Body.Set(field.NewSide(eSide))
			message.Body.Set(field.NewTransactTime(time.Now()))
			message.Body.Set(field.NewOrdType(eType))
			message.Body.Set(field.NewTimeInForce(eExpiry))
			message.Body.Set(field.NewSymbol(optionOrderSymbol))
			message.Body.Set(field.NewOrderQty(decimal.NewFromInt(optionOrderQuantity), 2))
			message.Body.Set(field.NewPrice(decimal.NewFromFloat(optionOrderPrice), 2))
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
	switch ordStatus.Value() {
	case enum.OrdStatus_NEW:
		break
	case enum.OrdStatus_PARTIALLY_FILLED:
		break
	case enum.OrdStatus_FILLED:
		break
	case enum.OrdStatus_DONE_FOR_DAY:
		break
	case enum.OrdStatus_REPLACED:
		break
	case enum.OrdStatus_PENDING_CANCEL:
		break
	case enum.OrdStatus_STOPPED:
		break
	case enum.OrdStatus_SUSPENDED:
		break
	case enum.OrdStatus_PENDING_NEW:
		break
	case enum.OrdStatus_CALCULATED:
		break
	case enum.OrdStatus_EXPIRED:
		break
	case enum.OrdStatus_ACCEPTED_FOR_BIDDING:
		break
	case enum.OrdStatus_PENDING_REPLACE:
		break
	case enum.OrdStatus_CANCELED:
		return makeError(errors.FixOrderCanceled)
	case enum.OrdStatus_REJECTED:
		return makeError(errors.FixOrderRejected)
	default:
		return makeError(errors.FixOrderStatusUnknown)
	}

	return nil
}

func isFinalStatus(msg *quickfix.Message) bool {
	ordStatus := field.OrdStatusField{}

	if err := msg.Body.GetField(tag.OrdStatus, &ordStatus); err != nil {
		return false
	}

	switch ordStatus.Value() {
	case enum.OrdStatus_FILLED:
		return true
	case enum.OrdStatus_DONE_FOR_DAY:
		return true
	case enum.OrdStatus_STOPPED:
		return true
	case enum.OrdStatus_EXPIRED:
		return true
	case enum.OrdStatus_CANCELED:
		return true
	case enum.OrdStatus_REJECTED:
		return true
	default:
		return false
	}
}
