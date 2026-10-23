package amendorder

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/cli/complete"
	"sylr.dev/fix/v2/cli/options"
	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/fixt11"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/utils"
)

var errInvalidMessageType = fmt.Errorf("invalid message type")

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
	logger := config.GetLogger()

	c, err := initiator.Connect(cmd.Context(), initiator.Options{
		MsgTypes: []string{
			string(fix50sp2.MsgTypeExecutionreport),
			string(fix50sp2.MsgTypeOrdercancelreject),
			string(fix50sp2.MsgTypeBusinessmessagereject),
		},
	})
	if err != nil {
		return err
	}

	defer c.Close(true)

	session := c.Session

	// Prepare order
	order, err := buildMessage(c.BeginString(), session)
	if err != nil {
		return err
	}

	// Send the order
	if _, err = c.Send(cmd.Context(), order); err != nil {
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

		case <-c.Done():
			break LOOP

		case msg := <-c.Messages():
			if err := processResponse(c, msg); err != nil {
				if errors.Is(err, errInvalidMessageType) {
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

func buildMessage(beginString string, session *config.Session) (*fix50sp2.OrderCancelReplaceRequestBuilder, error) {
	if err := session.CheckSupported(); err != nil {
		return nil, err
	}

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
	message := fix50sp2.BuildOrderCancelReplaceRequest(beginString)
	if len(optionOrderID) > 0 {
		message.SetOrderID(optionOrderID)
	}
	if len(optionOrigClOrdID) > 0 {
		message.SetOrigClOrdID(optionOrigClOrdID)
	}
	message.SetClOrdID(optionClOrdID).
		SetSide(eSide).
		SetTransactTime(time.Now()).
		SetOrdType(eType).
		SetTimeInForce(eExpiry).
		SetSymbol(optionOrderSymbol)

	// FIX 4.3 requires HandlInst.
	if session.Requires(string(fix50sp2.MsgTypeOrdercancelreplacerequest), fix50sp2.TagHandlInst) {
		message.SetHandlInst(fix50sp2.HandlInstAutomatedExecutionOrderPrivateNoBrokerIntervention)
	}

	for _, entry := range partyIdOptions.PartyEntries(session) {
		message.AddNoPartyIDsEntry(entry)
	}

	message.SetString(fix50sp2.TagOrderQty, decimal.NewFromInt(optionOrderQuantity).StringFixed(2))
	message.SetString(fix50sp2.TagPrice, decimal.NewFromFloat(optionOrderPrice).StringFixed(2))

	return message, nil
}

func processResponse(c *initiator.Client, msg *fix.Message) error {
	text := string(msg.Get(fix50sp2.TagText))

	makeError := func(errType error) error {
		if len(text) > 0 {
			return fmt.Errorf("%w: %s", errType, text)
		} else {
			return errType
		}
	}

	// MsgType
	switch msgType := string(msg.MsgType()); msgType {
	case string(fixt11.MsgTypeReject),
		string(fix50sp2.MsgTypeBusinessmessagereject),
		string(fix50sp2.MsgTypeOrdercancelreject):
		return makeError(errors.FixOrderRejected)
	case string(fix50sp2.MsgTypeExecutionreport):
	default:
		return errInvalidMessageType
	}

	// OrdStatus
	ordStatus, ok := fix50sp2.WrapExecutionReport(msg).OrdStatus()
	if !ok {
		return fmt.Errorf("missing OrdStatus")
	}

	utils.WriteMessageBodyAsTable(os.Stdout, msg, c.Dictionary)
	switch ordStatus {
	case fix50sp2.OrdStatusNew,
		fix50sp2.OrdStatusPartiallyFilled,
		fix50sp2.OrdStatusFilled,
		fix50sp2.OrdStatusDoneForDay,
		fix50sp2.OrdStatusReplaced,
		fix50sp2.OrdStatusPendingCancel,
		fix50sp2.OrdStatusStopped,
		fix50sp2.OrdStatusSuspended,
		fix50sp2.OrdStatusPendingNew,
		fix50sp2.OrdStatusCalculated,
		fix50sp2.OrdStatusExpired,
		fix50sp2.OrdStatusAcceptedForBidding,
		fix50sp2.OrdStatusPendingReplace:
		return nil
	case fix50sp2.OrdStatusCanceled:
		return makeError(errors.FixOrderCanceled)
	case fix50sp2.OrdStatusRejected:
		return makeError(errors.FixOrderRejected)
	default:
		return makeError(errors.FixOrderStatusUnknown)
	}
}

func isFinalStatus(msg *fix.Message) bool {
	ordStatus, ok := fix50sp2.WrapExecutionReport(msg).OrdStatus()
	if !ok {
		return false
	}

	switch ordStatus {
	case fix50sp2.OrdStatusFilled,
		fix50sp2.OrdStatusDoneForDay,
		fix50sp2.OrdStatusStopped,
		fix50sp2.OrdStatusExpired,
		fix50sp2.OrdStatusCanceled,
		fix50sp2.OrdStatusRejected:
		return true
	default:
		return false
	}
}
