//go:build validator
// +build validator

package marketdatavalidator

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/marketdata"
	"sylr.dev/fix/v2/utils"
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
	logger := config.GetLogger()

	c, err := initiator.Connect(cmd.Context(), initiator.Options{
		MsgTypes: []string{
			string(fix50sp2.MsgTypeMarketdatasnapshotfullrefresh),
			string(fix50sp2.MsgTypeMarketdataincrementalrefresh),
			string(fix50sp2.MsgTypeBusinessmessagereject),
		},
	})
	if err != nil {
		return err
	}
	defer c.Close(true)

	validator := marketdata.NewValidator(logger, optionSymbol)

	// Prepare market data request
	marketDataRequest, err := buildMessage(c, optionSymbol)
	if err != nil {
		return err
	}

	// Send the market data request
	if _, err = c.Send(cmd.Context(), marketDataRequest); err != nil {
		return err
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

		case <-c.Done():
			logger.Info().Msgf("Fix application not connected anymore")
			break LOOP

		case msg := <-c.Messages():
			err := validator.Handle(msg)
			if errors.Is(err, marketdata.ErrBusinessMessageReject) {
				logger.Info().Msgf("Fix application received a BusinessMessageReject, wants to exit: %s", err)
				break LOOP
			} else if err != nil {
				return err
			}
		}
	}

	logger.Debug().Msgf("exiting")
	return nil
}

func buildMessage(c *initiator.Client, symbol string) (session.Outgoing, error) {
	message := fix50sp2.BuildMarketDataRequest(c.BeginString()).
		SetMDReqID(uuid.NewString()).
		SetSubscriptionRequestType(fix50sp2.SubscriptionRequestTypeSnapshotPlusUpdates).
		SetMarketDepth(0).
		SetMDUpdateType(fix50sp2.MDUpdateTypeIncrementalRefresh).
		AddNoMDEntryTypesEntry(fix50sp2.NewNoMDEntryTypesEntry().SetMDEntryType(fix50sp2.MDEntryTypeBid)).
		AddNoMDEntryTypesEntry(fix50sp2.NewNoMDEntryTypesEntry().SetMDEntryType(fix50sp2.MDEntryTypeOffer)).
		AddNoMDEntryTypesEntry(fix50sp2.NewNoMDEntryTypesEntry().SetMDEntryType(fix50sp2.MDEntryTypeTrade)).
		AddNoMDEntryTypesEntry(fix50sp2.NewNoMDEntryTypesEntry().SetMDEntryType(fix50sp2.MDEntryType("101"))). // Trade History
		AddNoRelatedSymEntry(fix50sp2.NewNoRelatedSymEntry().SetSymbol(symbol))

	return message, nil
}
