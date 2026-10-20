package acceptor

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/acceptor"
	"sylr.dev/fix/v2/acceptor/application"
	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/utils"
)

var (
	optionNatsEmbeded      bool
	optionNatsURL          string
	optionNatsOrderSubject string
)

var AcceptorCmd = &cobra.Command{
	Use:               "acceptor",
	Short:             "Launch a FIX acceptor",
	Long:              "Launch a FIX acceptor.",
	RunE:              Execute,
	PersistentPreRunE: utils.MakePersistentPreRunE(acceptor.ValidateOptions),
}

func init() {
	AcceptorCmd.Flags().StringVar(&optionNatsURL, "nats-url", "nats://127.0.0.1:4222", "NATS URL used to forward FIX messages")
	AcceptorCmd.Flags().StringVar(&optionNatsOrderSubject, "nats-order-subject", "orders.{{.Symbol}}.{{.Side}}.{{.Type}}", "NATS order subject")
	utils.AddBothBoolFlags(AcceptorCmd.Flags(), &optionNatsEmbeded, "nats-embeded", "", true, "Launch embeded NATS server")

	acceptor.AddPersistentFlags(AcceptorCmd)
	acceptor.AddPersistentFlagCompletions(AcceptorCmd)
}

func Execute(cmd *cobra.Command, args []string) error {
	logger := config.GetLogger()

	acceptorOptions := application.AcceptorOptions{
		NATSEmbeded:      optionNatsEmbeded,
		NATSURL:          optionNatsURL,
		NATSOrderSubject: optionNatsOrderSubject,
	}

	app, err := application.NewAcceptor(&acceptorOptions)
	if err != nil {
		return err
	}
	defer app.Close()

	app.Logger = logger

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return acceptor.Serve(ctx, func(*config.Session) (session.Handler, error) {
		return app, nil
	})
}
