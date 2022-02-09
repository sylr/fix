package new

import (
	"github.com/quickfixgo/quickfix"
	"github.com/spf13/cobra"
	"sylr.dev/fix/pkg/application"
)

// NewCmd represents the buy command
var NewCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch a FIX server",
	Long:  "Launch a FIX server accepting orders.",
	RunE:  Execute,
}

func init() {

}

func Execute(cmd *cobra.Command, args []string) error {
	app, err := application.NewServer()
	if err != nil {
		return err
	}

	acceptor, err := quickfix.NewInitiator(app, nil, &quickfix.Settings{}, nil)
	if err != nil {
		return err
	}

	acceptor.Start()
	if err != nil {
		return err
	} else {
		defer acceptor.Stop()
	}

	return nil
}
