package status_order

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
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

// errInvalidMessageType is returned by processResponse for messages that are
// not an answer to the request.
var errInvalidMessageType = fmt.Errorf("invalid message type")

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
	logger := config.GetLogger()

	c, err := initiator.Connect(cmd.Context(), initiator.Options{
		MsgTypes: []string{
			string(fix50sp2.MsgTypeExecutionreport),
			string(fix50sp2.MsgTypeQuotestatusreport),
			string(fix50sp2.MsgTypeOrdercancelreject),
		},
	})
	if err != nil {
		return err
	}
	defer c.Close(true)

	// Prepare order
	order, err := buildMessage(c.Session)
	if err != nil {
		return err
	}

	// Send the order
	if _, err = c.Send(cmd.Context(), order); err != nil {
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

		case <-c.Done():
			break LOOP

		case msg := <-c.Messages():
			if err := processResponse(c, msg); err != nil {
				if errors.Is(err, errInvalidMessageType) {
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

func buildMessage(session *config.Session) (*fix50sp2.OrderStatusRequestBuilder, error) {
	eSide, err := dict.OrderSideStringToEnum(optionSide)
	if err != nil {
		return nil, err
	}

	if err := session.CheckSupported(); err != nil {
		return nil, err
	}

	message := fix50sp2.BuildOrderStatusRequest(session.BeginString).
		SetClOrdID(uuid.NewString()).
		SetSide(eSide).
		SetSymbol(optionSymbol)
	if len(optionOrderID) > 0 {
		message.SetOrderID(optionOrderID)
	}
	if len(optionOrigClOrdID) > 0 {
		if session.Defines(string(fix50sp2.MsgTypeOrderstatusrequest), fix50sp2.TagOrigClOrdID) {
			// OrigClOrdID has no setter on the trimmed OrderStatusRequest
			message.SetString(fix50sp2.TagOrigClOrdID, optionOrigClOrdID)
		} else {
			// FIX 4.3 and 4.4 do not define OrigClOrdID in an
			// OrderStatusRequest: its ClOrdID names the order.
			message.SetClOrdID(optionOrigClOrdID)
		}
	}
	for _, entry := range partyIdOptions.PartyEntries(session) {
		message.AddNoPartyIDsEntry(entry)
	}

	return message, nil
}

func processResponse(c *initiator.Client, msg *fix.Message) error {
	makeError := func(text string, _ bool) error {
		if len(text) > 0 {
			return fmt.Errorf("%w: %s", errors.FixOrderRejected, text)
		}
		return errors.FixOrderRejected
	}

	switch string(msg.MsgType()) {
	case string(fixt11.MsgTypeReject):
		return makeError(fixt11.WrapReject(msg).Text())

	case string(fix50sp2.MsgTypeBusinessmessagereject):
		return makeError(fix50sp2.WrapBusinessMessageReject(msg).Text())

	case string(fix50sp2.MsgTypeOrdercancelreject):
		return makeError(fix50sp2.WrapOrderCancelReject(msg).Text())

	case string(fix50sp2.MsgTypeExecutionreport):
		if _, ok := fix50sp2.WrapExecutionReport(msg).OrdStatus(); !ok {
			return fmt.Errorf("%w: missing OrdStatus", errors.Fix)
		}

		utils.WriteMessageBodyAsTable(os.Stdout, msg, c.Dictionary)
		return nil

	default:
		return errInvalidMessageType
	}
}
