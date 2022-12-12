//go:build validator
// +build validator

package marketdatavalidator

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fixt11"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/initiator/application"
	"sylr.dev/fix/pkg/utils"
)

var (
	optionSymbol string
)

var MarketDataValidatorCmd = &cobra.Command{
	Use:               "validator",
	Short:             "Validates market data retrieved",
	Long:              "Validates market data retrieved.",
	Args:              cobra.ExactArgs(0),
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: utils.MakePersistentPreRunE(Validate),
	RunE:              Execute,
}

func init() {
	MarketDataValidatorCmd.Flags().StringVar(&optionSymbol, "symbol", "", "Symbol")

	MarketDataValidatorCmd.RegisterFlagCompletionFunc("symbol", cobra.NoFileCompletions)
}

func Validate(cmd *cobra.Command, args []string) error {
	err := utils.ReconcileBoolFlags(cmd.Flags())
	if err != nil {
		return err
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

	ctxInitiator, err := context.GetInitiator()
	if err != nil {
		return err
	}

	session := sessions[0]
	transportDict, appDict, err := session.GetFIXDictionaries()
	if err != nil {
		return err
	}

	settings, err := context.ToQuickFixInitiatorSettings()
	if err != nil {
		return err
	}

	app := application.NewMarketDataValidator(logger, optionSymbol)
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

	defer func() {
		app.Stop()
		init.Stop()
	}()

	// Choose right timeout cli option > config > default value (5s)
	var timeout time.Duration
	if options.Timeout != time.Duration(0) {
		timeout = options.Timeout
	} else if ctxInitiator.SocketTimeout != time.Duration(0) {
		timeout = ctxInitiator.SocketTimeout
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

	// Prepare market data request
	marketDataRequest, err := buildMessage(*session, optionSymbol)
	if err != nil {
		return err
	}

	// Send the order
	err = quickfix.Send(marketDataRequest)
	if err != nil {
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

LOOP:
	for {
		select {
		case signal := <-interrupt:
			switch signal {
			case syscall.SIGINT, syscall.SIGTERM:
				logger.Debug().Msgf("Received signal: %v", signal)
				break LOOP
			default:
				logger.Info().Msgf("Received unhandled signal: %v", signal)
			}

		case _, ok := <-app.Connected:
			if !ok {
				logger.Info().Msgf("Fix application not connected anymore")
				break LOOP
			} else {
				logger.Info().Msgf("Fix application wrote on Connected chan, wants to exit")
				break LOOP
			}
		}
	}

	logger.Debug().Msgf("exiting")
	return nil
}

func buildMessage(session config.Session, symbol string) (quickfix.Messagable, error) {
	uid := uuid.New()
	mdReqID := field.NewMDReqID(uid.String())
	subReqType := field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES)
	marketDepth := field.NewMarketDepth(0)

	// Message
	message := quickfix.NewMessage()
	header := fixt11.NewHeader(&message.Header)

	header.Set(field.NewMsgType(enum.MsgType_MARKET_DATA_REQUEST))
	message.Body.Set(mdReqID)
	message.Body.Set(subReqType)
	message.Body.Set(marketDepth)
	message.Body.Set(field.NewMDUpdateType(enum.MDUpdateType_INCREMENTAL_REFRESH))

	entryTypes := quickfix.NewRepeatingGroup(
		tag.NoMDEntryTypes,
		quickfix.GroupTemplate{quickfix.GroupElement(tag.MDEntryType)},
	)

	entryTypes.Add().Set(field.NewMDEntryType(enum.MDEntryType_BID))
	entryTypes.Add().Set(field.NewMDEntryType(enum.MDEntryType_OFFER))
	entryTypes.Add().Set(field.NewMDEntryType(enum.MDEntryType_TRADE))
	entryTypes.Add().Set(field.NewMDEntryType(enum.MDEntryType("101"))) // Trade History

	message.Body.SetGroup(entryTypes)

	relatedSym := quickfix.NewRepeatingGroup(
		tag.NoRelatedSym,
		quickfix.GroupTemplate{quickfix.GroupElement(tag.Symbol)},
	)

	relatedSym.Add().Set(field.NewSymbol(symbol))
	message.Body.SetGroup(relatedSym)

	utils.QuickFixMessagePartSetString(&message.Header, session.TargetCompID, field.NewTargetCompID)
	utils.QuickFixMessagePartSetString(&message.Header, session.TargetSubID, field.NewTargetSubID)
	utils.QuickFixMessagePartSetString(&message.Header, session.SenderCompID, field.NewSenderCompID)
	utils.QuickFixMessagePartSetString(&message.Header, session.SenderSubID, field.NewSenderSubID)

	return message, nil
}
