package send

import (
	"fmt"

	"github.com/spf13/cobra"

	sendorder "sylr.dev/fix/cmd/send/order"
	"sylr.dev/fix/config"
)

// SendCmd represents the buy command
var SendCmd = &cobra.Command{
	Use:               "send",
	Short:             "Send a FIX message",
	Long:              "Send a FIX message after initiating a sesion with a FIX acceptor.",
	RunE:              Execute,
	PersistentPreRunE: Validate,
}

func init() {
	options := config.GetOptions()

	SendCmd.PersistentFlags().StringVar(&options.Context, "context", "", "Context to use")
	SendCmd.PersistentFlags().StringVar(&options.Acceptor, "acceptor", "", "Acceptor to use (can't be used with --context)")
	SendCmd.PersistentFlags().StringVar(&options.Session, "session", "", "Session to use (can't be used with --context)")
	SendCmd.PersistentFlags().DurationVar(&options.Timeout, "timeout", 0, "Duration for timeouts")

	AddPersistentFlagCompletion(SendCmd)
	AddPersistentFlagCompletion(sendorder.SendOrderCmd)

	SendCmd.AddCommand(sendorder.SendOrderCmd)
}

func Validate(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()
	fixConfig := config.GetConfig()

	if len(options.Context) > 0 {
		if len(options.Acceptor) > 0 || len(options.Session) > 0 {
			return fmt.Errorf("can't use --acceptor/--session with --context")
		}
	} else {
		if len(options.Acceptor) == 0 || len(options.Session) == 0 {
			return fmt.Errorf("you need to specify either --context or --acceptor/--session")
		}
	}

	conf, err := config.ReadYAML(options.Config, options.Interactive)

	if err != nil {
		return fmt.Errorf("unable to read configuration: %w", err)
	}

	if len(options.Context) == 0 {
		const contextName = "from-cli-options"
		conf.Contexts = append(conf.Contexts, config.Context{
			Name:     contextName,
			Acceptor: options.Acceptor,
			Session:  options.Session,
		})
		(*options).Context = contextName
	}

	*fixConfig = *conf

	context, err := config.GetCurrentContext()
	if err != nil {
		return fmt.Errorf("%s: %w", config.ErrBadConfig, err)
	}

	_, err = context.GetAcceptor()
	if err != nil {
		return fmt.Errorf("%s: %w", config.ErrBadConfig, err)
	}

	_, err = context.GetSession()
	if err != nil {
		return fmt.Errorf("%s: %w", config.ErrBadConfig, err)
	}

	return nil
}

func Execute(cmd *cobra.Command, args []string) error {
	return nil
}

func AddPersistentFlagCompletion(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc("context", completeContext)
	cmd.RegisterFlagCompletionFunc("acceptor", completeAcceptor)
	cmd.RegisterFlagCompletionFunc("session", completeSession)
}
