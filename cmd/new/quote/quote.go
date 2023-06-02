package newquote

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fixt11"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
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
	optionSymbol, optionQuoteID               string
	optionBuyQuantities, optionSellQuantities []int64
	optionBuyPrices, optionSellPrices         []float64
	optionAutoPriceUpdate                     bool
	optionOrderOrigination                    string
	partyIdOptions                            *options.PartyIdOptions
	attributeOptions                          *options.AttributeOptions
	optionExecReports                         int
	optionExecReportsTimeout                  time.Duration
	optionExecReportsTimeoutReset             bool
	optionStopOnFinalState                    bool
	optionUpdateQuotePeriod                   time.Duration
	optionCancelAfterUpdates                  bool
	optionCancelAfterXUpdates                 int

	priceIteration int = 0
)

var NewQuoteCmd = &cobra.Command{
	Use:               "quote",
	Short:             "New quote",
	Long:              "Send a new quote after initiating a session with a FIX acceptor.",
	Args:              cobra.ExactArgs(0),
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: utils.MakePersistentPreRunE(Validate),
	RunE:              Execute,
}

func init() {
	NewQuoteCmd.Flags().StringVar(&optionQuoteID, "id", "", "Quote id (uuid autogenerated if not given)")
	NewQuoteCmd.Flags().StringVar(&optionSymbol, "symbol", "", "Quote symbol")
	NewQuoteCmd.Flags().Int64SliceVar(&optionBuyQuantities, "buy-quantities", []int64{}, "Quote buy quantities")
	NewQuoteCmd.Flags().Int64SliceVar(&optionSellQuantities, "sell-quantities", []int64{}, "Quote sell quantities")
	NewQuoteCmd.Flags().Float64SliceVar(&optionBuyPrices, "buy-prices", []float64{}, "Quote buy prices")
	NewQuoteCmd.Flags().Float64SliceVar(&optionSellPrices, "sell-prices", []float64{}, "Quote sell prices")
	NewQuoteCmd.Flags().BoolVar(&optionAutoPriceUpdate, "auto-price-update", false, "Generate price oscilation to send an infinity of quote updates")
	NewQuoteCmd.Flags().StringVar(&optionOrderOrigination, "origination", "", "Order origination")

	partyIdOptions = options.NewPartyIdOptions(NewQuoteCmd)
	attributeOptions = options.NewAttributeOptions(NewQuoteCmd)

	NewQuoteCmd.Flags().IntVar(&optionExecReports, "exec-reports", 1, "Expect given number of execution reports before logging out (0 wait indefinitely)")
	NewQuoteCmd.Flags().DurationVar(&optionExecReportsTimeout, "exec-reports-timeout", 5*time.Second, "Log out if execution reports not received within timeout (0s wait indefinitely)")
	NewQuoteCmd.Flags().BoolVar(&optionExecReportsTimeoutReset, "exec-reports-timeout-reset", false, "Reset execution reports timeout each time an execution report is received")

	NewQuoteCmd.Flags().BoolVar(&optionStopOnFinalState, "stop-on-final-state", false, "Stop application when receiving an order with a final state")

	NewQuoteCmd.Flags().DurationVar(&optionUpdateQuotePeriod, "update-period", 10*time.Second, "Period for recurring update quote")
	NewQuoteCmd.Flags().BoolVar(&optionCancelAfterUpdates, "cancel-quote-after-updates", false, "Cancel quote when all updates are done")
	NewQuoteCmd.Flags().IntVar(&optionCancelAfterXUpdates, "cancel-quote-after-x-updates", -1, "Cancel quote after X updates")

	NewQuoteCmd.MarkFlagRequired("symbol")

	NewQuoteCmd.RegisterFlagCompletionFunc("symbol", cobra.NoFileCompletions)
	NewQuoteCmd.RegisterFlagCompletionFunc("origination", complete.OrderOriginationRole)
}

func Validate(cmd *cobra.Command, args []string) error {
	if len(optionBuyPrices) != len(optionBuyQuantities) {
		return fmt.Errorf("number of buy prices (%d) must be equal to buy quantities (%d)", len(optionBuyPrices), len(optionBuyQuantities))
	}

	if len(optionSellPrices) != len(optionSellQuantities) {
		return fmt.Errorf("number of sell prices (%d) must be equal to sell quantities (%d)", len(optionSellPrices), len(optionSellQuantities))
	}

	if optionAutoPriceUpdate && (len(optionBuyPrices) > 1 || len(optionSellPrices) > 1) {
		return fmt.Errorf("auto-price-update option works with only 1 price per side")
	}

	if len(optionOrderOrigination) > 0 {
		originations := utils.PrettyOptionValues(dict.OrderOriginations)
		search := utils.Search(originations, strings.ToLower(optionOrderOrigination))
		if search < 0 {
			return errors.OptionOrderOriginationUnknown
		}
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

	// Prepare quote
	quote, err := buildMessage(*session)
	if err != nil {
		return err
	}

	// Send the quote
	err = quickfix.Send(quote)
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

	var updatePeriod <-chan time.Time
	if optionUpdateQuotePeriod > 0 {
		updatePeriod = make(<-chan time.Time)
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

		case <-updatePeriod:
			if !optionAutoPriceUpdate && priceIteration >= len(optionBuyPrices) && priceIteration >= len(optionSellPrices) {
				optionUpdateQuotePeriod = 0

				if optionCancelAfterUpdates {
					cancelMsg, err := buildCancelMessage(*session)
					if err != nil {
						return err
					}
					err = quickfix.Send(cancelMsg)
					if err != nil {
						return err
					}
				}
				continue LOOP
			}

			if optionCancelAfterXUpdates > 0 && priceIteration%optionCancelAfterXUpdates == 0 {
				// Prepare quote cancel
				quoteMsg, err := buildCancelMessage(*session)
				if err != nil {
					return err
				}

				// Send the cancellation
				err = quickfix.Send(quoteMsg)
				if err != nil {
					return err
				}
			}

			// Prepare quote
			quoteMsg, err := buildMessage(*session)
			if err != nil {
				return err
			}

			// Send the quote
			err = quickfix.Send(quoteMsg)
			if err != nil {
				return err
			}

		case msg, ok := <-app.FromAppMessages:
			if !ok {
				break LOOP
			}

			if err := processReponse(app, msg); err != nil {
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

		if optionUpdateQuotePeriod > 0 {
			updatePeriod = time.After(optionUpdateQuotePeriod)
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
		quoteId = uuid.New().String()
	}

	switch session.BeginString {
	case quickfix.BeginStringFIXT11:
		switch session.DefaultApplVerID {
		case "FIX.5.0SP2":
			header.Set(field.NewMsgType(enum.MsgType_QUOTE))
			message.Body.Set(field.NewQuoteID(quoteId))
			message.Body.Set(field.NewTransactTime(time.Now()))
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

	message.Body.Set(field.NewSymbol(optionSymbol))
	if optionAutoPriceUpdate {
		priceAdjustment := 0.01 * float64(priceIteration%100)
		if len(optionBuyPrices) == 1 {
			message.Body.Set(field.NewBidSize(decimal.NewFromInt(optionBuyQuantities[0]), 2))
			message.Body.Set(field.NewBidPx(decimal.NewFromFloat(optionBuyPrices[0]-priceAdjustment), 2))
		}
		if len(optionSellPrices) == 1 {
			message.Body.Set(field.NewOfferSize(decimal.NewFromInt(optionSellQuantities[0]), 2))
			message.Body.Set(field.NewOfferPx(decimal.NewFromFloat(optionSellPrices[0]+priceAdjustment), 2))
		}
	} else {
		if priceIteration < len(optionBuyPrices) {
			message.Body.Set(field.NewBidSize(decimal.NewFromInt(optionBuyQuantities[priceIteration]), 2))
			message.Body.Set(field.NewBidPx(decimal.NewFromFloat(optionBuyPrices[priceIteration]), 2))
		}
		if priceIteration < len(optionSellPrices) {
			message.Body.Set(field.NewOfferSize(decimal.NewFromInt(optionSellQuantities[priceIteration]), 2))
			message.Body.Set(field.NewOfferPx(decimal.NewFromFloat(optionSellPrices[priceIteration]), 2))
		}
	}

	priceIteration++

	return message, nil
}

func buildCancelMessage(session config.Session) (quickfix.Messagable, error) {
	// Message
	message := quickfix.NewMessage()
	header := fixt11.NewHeader(&message.Header)

	quoteId := optionQuoteID
	if len(quoteId) == 0 {
		quoteId = uuid.New().String()
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
	inst := instruments.Add()
	inst.Set(field.NewSymbol(optionSymbol))
	message.Body.SetGroup(instruments)

	return message, nil
}

func processReponse(app *application.NewOrder, msg *quickfix.Message) error {
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
	} else if msgType.Value() == enum.MsgType_REJECT || msgType.Value() == enum.MsgType_BUSINESS_MESSAGE_REJECT {
		return makeError(errors.FixOrderRejected)
	} else if msgType.Value() == enum.MsgType_QUOTE_STATUS_REPORT {
		app.WriteMessageBodyAsTable(os.Stdout, msg)
		quoteStatus := field.QuoteStatusField{}
		err = msg.Body.GetField(tag.QuoteStatus, &quoteStatus)
		if err != nil {
			return err
		}
		if quoteStatus.Value() == enum.QuoteStatus_REJECTED {
			return makeError(errors.FixOrderRejected)
		}
		return nil
	} else if msgType.Value() != enum.MsgType_EXECUTION_REPORT {
		return quickfix.InvalidMessageType()
	}

	// OrdStatus
	ordStatus := field.OrdStatusField{}
	err = msg.Body.GetField(tag.OrdStatus, &ordStatus)
	if err != nil {
		return err
	}

	app.WriteMessageBodyAsTable(os.Stdout, msg)
	return nil
}

func isFinalStatus(msg *quickfix.Message) bool {
	ordStatus := field.OrdStatusField{}
	err := msg.Body.GetField(tag.OrdStatus, &ordStatus)
	if err != nil {
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
		return true // because there is no modification in our use case
	default:
		return false
	}
}
