package cancelorder

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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
	optionOrderSide          string
	optionOrderSymbol        string
	optionOrderID            string
	optionClientOrderID      string
	optionOrigClientOrderID  string
	optionExecReportsTimeout time.Duration
	partyIdOptions           *options.PartyIdOptions
)

var CancelOrderCmd = &cobra.Command{
	Use:               "order",
	Short:             "Cancel single order",
	Long:              "Send a cancel order request after initiating a session with a FIX acceptor.",
	Args:              cobra.ExactArgs(0),
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: utils.MakePersistentPreRunE(Validate),
	RunE:              Execute,
}

func init() {
	CancelOrderCmd.Flags().StringVar(&optionOrderID, "id", "", "Order id")
	CancelOrderCmd.Flags().StringVar(&optionClientOrderID, "clordid", "", "Client order id")
	CancelOrderCmd.Flags().StringVar(&optionOrigClientOrderID, "origclordid", "", "Orig client order id")

	CancelOrderCmd.Flags().StringVar(&optionOrderSide, "side", "", "Order side (buy, sell ... etc)")
	CancelOrderCmd.Flags().StringVar(&optionOrderSymbol, "symbol", "", "Order symbol")
	CancelOrderCmd.Flags().DurationVar(&optionExecReportsTimeout, "exec-reports-timeout", 5*time.Second, "Log out if execution reports not received within timeout (0s wait indefinitely)")

	partyIdOptions = options.NewPartyIdOptions(CancelOrderCmd)

	CancelOrderCmd.MarkFlagRequired("clordid")
	CancelOrderCmd.MarkFlagRequired("side")

	CancelOrderCmd.RegisterFlagCompletionFunc("side", complete.OrderSide)
}

func Validate(cmd *cobra.Command, args []string) error {
	sides := utils.PrettyOptionValues(dict.OrderSides)
	search := utils.Search(sides, strings.ToLower(optionOrderSide))
	if search < 0 {
		return errors.OptionOrderSideUnknown
	}

	return partyIdOptions.Validate()
}

func Execute(cmd *cobra.Command, args []string) error {
	logger := config.GetLogger()

	c, err := initiator.Connect(cmd.Context(), initiator.Options{
		MsgTypes: []string{
			string(fix50sp2.MsgTypeExecutionreport),
			string(fix50sp2.MsgTypeOrdercancelreject),
			string(fix50sp2.MsgTypeOrdermasscancelreport),
		},
	})
	if err != nil {
		return err
	}
	defer c.Close(true)

	// Prepare cancel message
	cancelMsg, err := buildMessage(c.Session)
	if err != nil {
		return err
	}

	// Send the cancel message
	if _, err = c.Send(cmd.Context(), cancelMsg); err != nil {
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

		case <-c.Done():
			break LOOP

		case msg := <-c.Messages():
			if err := processResponse(c, msg); err != nil {
				if errors.Is(err, errInvalidMessageType) {
					continue LOOP
				}
				return err
			}
			break LOOP
		}
	}

	return nil
}

func buildMessage(session *config.Session) (*fix50sp2.OrderCancelRequestBuilder, error) {
	eside, err := dict.OrderSideStringToEnum(optionOrderSide)
	if err != nil {
		return nil, err
	}

	if err := session.CheckSupported(); err != nil {
		return nil, err
	}

	message := fix50sp2.BuildOrderCancelRequest(session.BeginString).
		SetClOrdID(optionClientOrderID).
		SetSide(eside).
		SetTransactTime(time.Now())
	if len(optionOrderID) > 0 {
		message.SetOrderID(optionOrderID)
	}
	if len(optionOrigClientOrderID) > 0 {
		message.SetOrigClOrdID(optionOrigClientOrderID)
	}
	if len(optionOrderSymbol) > 0 {
		message.SetSymbol(optionOrderSymbol)
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
		utils.WriteMessageBodyAsTable(os.Stdout, msg, c.Dictionary)
		return makeError(fix50sp2.WrapOrderCancelReject(msg).Text())

	case string(fix50sp2.MsgTypeExecutionreport):
		utils.WriteMessageBodyAsTable(os.Stdout, msg, c.Dictionary)
		report := fix50sp2.WrapExecutionReport(msg)
		ordStatus, ok := report.OrdStatus()
		if !ok {
			return fmt.Errorf("%w: missing OrdStatus", errors.Fix)
		}
		if ordStatus == fix50sp2.OrdStatusRejected {
			return makeError(report.Text())
		}

	default:
		return errInvalidMessageType
	}

	return nil
}
