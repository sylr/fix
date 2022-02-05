package listsecurity

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	slr50sp1 "github.com/quickfixgo/fix50sp1/securitylistrequest"
	slr50sp2 "github.com/quickfixgo/fix50sp2/securitylistrequest"
	"github.com/quickfixgo/quickfix"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/application"
	"sylr.dev/fix/pkg/initiator"
)

var (
	optionSide, optionType string
	optionSymbol, optionID string
	optionExpiry           string
	optionQuantity         int64
	optionPrice            float64
)

var ListSecurityCmd = &cobra.Command{
	Use:               "security",
	Aliases:           []string{"securities"},
	Short:             "List securities",
	Long:              "Send a securitylist FIX Message after initiating a session with a FIX acceptor.",
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
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC822,
	}
	multi := zerolog.MultiLevelWriter(consoleWriter)
	logger := zerolog.New(multi).With().Timestamp().Logger().Level(config.IntToZerologLevel(options.Verbose))

	if options.LogCaller {
		logger = logger.With().Caller().Logger()
	}

	context, err := config.GetCurrentContext()
	if err != nil {
		return err
	}

	session, err := context.GetSession()
	if err != nil {
		return err
	}

	acceptor, err := context.GetAcceptor()
	if err != nil {
		return err
	}

	transportDict, appDict, err := session.GetFIXDictionaries()
	if err != nil {
		return err
	}

	settings, err := context.ToQuickFixSettings()
	if err != nil {
		return err
	}

	app := application.NewSecurityList()
	app.Logger = &logger
	app.Settings = settings
	app.TransportDataDictionary = transportDict
	app.AppDataDictionary = appDict

	init, err := initiator.Initiate(app, settings)
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
	} else if acceptor.SocketTimeout != time.Duration(0) {
		timeout = acceptor.SocketTimeout
	} else {
		timeout = 5 * time.Second
	}

	// Wait for session connection
	select {
	case <-time.After(timeout):
		return fmt.Errorf("connection timeout")
	case <-app.Connected:
	}

	// Prepare securitylist
	securitylist, err := new(*session)
	if err != nil {
		return err
	}

	// Send the order
	err = quickfix.Send(securitylist)
	if err != nil {
		return err
	}

	// Wait for the order response
	var responseMessage *quickfix.Message
	select {
	case <-time.After(timeout):
		return fmt.Errorf("timeout while waiting for order response")
	case responseMessage = <-app.FromAppChan:
	}

	app.WriteMessageBodyAsTable(os.Stdout, responseMessage)

	return nil
}

func new(session config.Session) (quickfix.Messagable, error) {
	var order quickfix.Messagable

	target := field.NewTargetCompID(session.TargetCompID)
	sender := field.NewSenderCompID(session.SenderCompID)
	stype := field.NewSecurityListRequestType(enum.SecurityListRequestType_SYMBOL)
	reqid := field.NewSecurityReqID("1")

	switch session.BeginString {
	case quickfix.BeginStringFIXT11:
		switch session.DefaultApplVerID {
		case "FIX.5.0SP1":
			securitylist50sp1 := slr50sp1.New(reqid, stype)
			securitylist50sp1.Header.Set(target)
			securitylist50sp1.Header.Set(sender)

			order = securitylist50sp1
		case "FIX.5.0SP2":
			securitylist50sp2 := slr50sp2.New(reqid, stype)
			securitylist50sp2.Header.Set(target)
			securitylist50sp2.Header.Set(sender)

			order = securitylist50sp2
		default:
			return nil, errors.New("FIX version not implemented")
		}
	default:
		return nil, errors.New("FIX version not implemented")
	}

	return order, nil
}
