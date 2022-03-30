package acceptor

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/acceptor"
	"sylr.dev/fix/pkg/acceptor/application"
	"sylr.dev/fix/pkg/utils"
)

var (
	optionNatsEmbeded      bool
	optionNatsURL          string
	optionNatsOrderSubject string
)

var AcceptorCmd = &cobra.Command{
	Use:   "acceptor",
	Short: "Launch a FIX acceptor",
	Long:  "Launch a FIX acceptor.",
	RunE:  Execute,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := acceptor.ValidateOptions(cmd, args)
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
}

func init() {
	AcceptorCmd.Flags().StringVar(&optionNatsURL, "nats-url", "nats://127.0.0.1:4222", "NATS URL used to forward FIX messages")
	AcceptorCmd.Flags().StringVar(&optionNatsOrderSubject, "nats-order-subject", "orders.{{.Symbol}}.{{.Side}}.{{.Type}}", "NATS order subject")
	utils.AddBothBoolFlags(AcceptorCmd.Flags(), &optionNatsEmbeded, "nats-embeded", "", true, "Launch embeded NATS server")

	acceptor.AddPersistentFlags(AcceptorCmd)
	acceptor.AddPersistentFlagCompletions(AcceptorCmd)
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

	settings, err := context.ToQuickFixAcceptorSettings()
	if err != nil {
		return err
	}

	transportDict, appDict, err := sessions[0].GetFIXDictionaries()
	if err != nil {
		return err
	}

	acceptorOptions := application.AcceptorOptions{
		NATSEmbeded:      optionNatsEmbeded,
		NATSURL:          optionNatsURL,
		NATSOrderSubject: optionNatsOrderSubject,
	}

	app, err := application.NewAcceptor(&acceptorOptions)
	if err != nil {
		return err
	}

	app.TransportDataDictionary = transportDict
	app.AppDataDictionary = appDict
	app.Logger = logger

	var quickfixLogger *zerolog.Logger
	if options.QuickFixLogging {
		quickfixLogger = logger
	}

	acceptor, err := acceptor.NewAcceptor(app, settings, quickfixLogger)
	if err != nil {
		return err
	}

	err = acceptor.Start()
	if err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	acceptor.Stop()
	os.Exit(0)

	return nil
}
