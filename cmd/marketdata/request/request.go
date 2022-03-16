package marketdatarequest

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	mdr50sp1 "github.com/quickfixgo/fix50sp1/marketdatarequest"
	mdr50sp2 "github.com/quickfixgo/fix50sp2/marketdatarequest"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/initiator"
	"sylr.dev/fix/pkg/initiator/application"
	"sylr.dev/fix/pkg/utils"
)

var (
	optionType string
)

var MarketDataRequestCmd = &cobra.Command{
	Use:               "request",
	Short:             "Send MarketDataRequest FIX message",
	Long:              "Send a MarketDataRequest FIX Message after initiating a session with a FIX ctxInitiator.",
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
				return parent.PersistentPreRunE(cmd, args)
			}
		}
		return nil
	},
	RunE: Execute,
}

func init() {

}

func Validate(cmd *cobra.Command, args []string) error {
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

	app := application.NewMarketDataRequest()
	app.Logger = logger
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

	defer init.Stop()

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

	// Prepare securitylist
	securitylist, err := buildMessage(*session)
	if err != nil {
		return err
	}

	// Send the order
	err = quickfix.Send(securitylist)
	if err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

LOOP:
	for {
		select {
		case signal := <-interrupt:
			logger.Debug().Msgf("Received signal: %s", signal)
			break LOOP
		case _, ok := <-app.FromAppChan:
			if !ok {
				break LOOP
			}
			//app.WriteMessageBodyAsTable(os.Stdout, responseMessage.ToMessage())
		}
	}

	return nil
}

func buildMessage(session config.Session) (quickfix.Messagable, error) {
	var messagable quickfix.Messagable

	mdReqID := field.NewMDReqID("2")
	subReqType := field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT)
	marketDepth := field.NewMarketDepth(0)

	switch session.BeginString {
	case quickfix.BeginStringFIXT11:
		switch session.DefaultApplVerID {
		case "FIX.5.0SP1":
			request := mdr50sp1.New(mdReqID, subReqType, marketDepth)
			request.SetMDUpdateType(enum.MDUpdateType_FULL_REFRESH)

			entryTypes := mdr50sp1.NewNoMDEntryTypesRepeatingGroup()
			entryTypes.Add().SetMDEntryType(enum.MDEntryType_BID)
			entryTypes.Add().SetMDEntryType(enum.MDEntryType_OFFER)
			entryTypes.Add().SetMDEntryType(enum.MDEntryType_TRADE)
			request.SetNoMDEntryTypes(entryTypes)

			relatedSym := mdr50sp1.NewNoRelatedSymRepeatingGroup()
			relatedSym.Add().SetSymbol("LADYW_EUR")
			relatedSym.Add().SetSymbol("LADYW_EUR")
			request.SetNoRelatedSym(relatedSym)

			messagable = request.Message
		case "FIX.5.0SP2":
			request := mdr50sp2.New(mdReqID, subReqType, marketDepth)
			request.SetMDUpdateType(enum.MDUpdateType_FULL_REFRESH)

			entryTypes := mdr50sp2.NewNoMDEntryTypesRepeatingGroup()
			entryTypes.Add().SetMDEntryType(enum.MDEntryType_BID)
			entryTypes.Add().SetMDEntryType(enum.MDEntryType_OFFER)
			entryTypes.Add().SetMDEntryType(enum.MDEntryType_TRADE)
			request.SetNoMDEntryTypes(entryTypes)

			relatedSym := mdr50sp2.NewNoRelatedSymRepeatingGroup()
			relatedSym.Add().SetSymbol("LADYW_EUR")
			relatedSym.Add().SetSymbol("LADYW_EUR")
			request.SetNoRelatedSym(relatedSym)

			messagable = request.Message
		default:
			return nil, errors.FixVersionNotImplemented
		}
	default:
		return nil, errors.FixVersionNotImplemented
	}

	message := messagable.ToMessage()
	utils.QuickFixMessagePartSet(&message.Header, session.TargetCompID, field.NewTargetCompID)
	utils.QuickFixMessagePartSet(&message.Header, session.TargetSubID, field.NewTargetSubID)
	utils.QuickFixMessagePartSet(&message.Header, session.SenderCompID, field.NewSenderCompID)
	utils.QuickFixMessagePartSet(&message.Header, session.SenderSubID, field.NewSenderSubID)

	return message, nil
}
