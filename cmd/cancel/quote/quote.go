package cancelorder

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

func buildMessage(session *config.Session) (*fix50sp2.QuoteCancelBuilder, error) {
	if err := session.CheckSupported(); err != nil {
		return nil, err
	}

	quoteId := optionQuoteID
	if len(quoteId) == 0 {
		quoteId = uuid.NewString()
	}

	message := fix50sp2.BuildQuoteCancel(session.BeginString).
		SetQuoteID(quoteId)
	// FIX 4.3 and 4.4 do not define QuoteMsgID in a QuoteCancel.
	if session.Defines(string(fix50sp2.MsgTypeQuotecancel), fix50sp2.TagQuoteMsgID) {
		message.SetQuoteMsgID("msg_" + quoteId)
	}
	message.SetQuoteCancelType(fix50sp2.QuoteCancelTypeCancelForOneOrMoreSecurities)
	for _, entry := range partyIdOptions.PartyEntries(session) {
		message.AddNoPartyIDsEntry(entry)
	}
	for _, symbol := range optionOrderSymbols {
		message.AddNoQuoteEntriesEntry(fix50sp2.NewNoQuoteEntriesEntry().SetSymbol(symbol))
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
