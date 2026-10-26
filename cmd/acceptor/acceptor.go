package acceptor

import (
	"github.com/spf13/cobra"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/acceptor"
	"sylr.dev/fix/v2/acceptor/application"
	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/utils"
)

var (
	optionNatsEmbedded     bool
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
	AcceptorCmd.Flags().BoolVar(&optionNatsEmbedded, "nats-embedded", true, "Launch an embedded NATS server")

	acceptor.AddPersistentFlags(AcceptorCmd)
	if err := acceptor.AddPersistentFlagCompletions(AcceptorCmd); err != nil {
		panic(err)
	}
}

func Execute(cmd *cobra.Command, args []string) error {
	logger := config.GetLogger()

	acceptorOptions := application.AcceptorOptions{
		NATSEmbedded:     optionNatsEmbedded,
		NATSURL:          optionNatsURL,
		NATSOrderSubject: optionNatsOrderSubject,
	}

	app, err := application.NewAcceptor(&acceptorOptions)
	if err != nil {
		return err
	}
	defer app.Close()

	app.Logger = logger

	return acceptor.Serve(cmd.Context(), func(*config.Session) (session.Handler, error) {
		return app, nil
	})
}
