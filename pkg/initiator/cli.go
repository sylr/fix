package initiator

import (
	"fmt"

	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/cli/complete"
	"sylr.dev/fix/pkg/errors"
)

func ValidateOptions(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()
	conf, err := config.ReadYAML(options.Config, options.Interactive)

	if err != nil {
		return fmt.Errorf("unable to read configuration: %w", err)
	}

	err = conf.Validate()
	if err != nil {
		return err
	}

	// Retrieve the global config pointer
	fixConfig := config.GetConfig()

	// Set the config retrieved in the config file into the global config pointer
	*fixConfig = *conf

	// Initialize the context name with the config current-context value
	contextName := fixConfig.CurrentContext

	if len(options.Context) > 0 {
		if len(options.Initiator) > 0 || len(options.Session) > 0 {
			return fmt.Errorf("%w: can't use --initiator/--session with --context", errors.Options)
		}
		contextName = options.Context
	} else if len(contextName) == 0 {
		if len(options.Initiator) == 0 || len(options.Session) == 0 {
			return fmt.Errorf("%w: you need to specify either --context or --initiator/--session", errors.Options)
		}
	}

	if len(contextName) == 0 {
		contextName = "no-context"
		fixConfig.Contexts = append(fixConfig.Contexts, &config.Context{
			Name:      contextName,
			Initiator: options.Initiator,
			Sessions:  []string{options.Session},
		})
		(*options).Context = contextName
	}

	context, err := config.GetCurrentContext()
	if err != nil {
		return err
	}

	_, err = context.GetInitiator()
	if err != nil {
		return err
	}

	sessions, err := context.GetSessions()
	if err != nil {
		return err
	} else {
		switch len(sessions) {
		case 0:
			return errors.ConfigContextNoSession
		case 1:
			// OK
		default:
			return errors.ConfigContextMultipleSessions
		}
	}

	return nil
}

func AddPersistentFlags(cmd *cobra.Command) {
	options := config.GetOptions()

	cmd.PersistentFlags().StringVar(&options.Context, "context", "", "Context to use")
	cmd.PersistentFlags().StringVar(&options.Initiator, "initiator", "", "Initiator to use (can't be used with --context)")
	cmd.PersistentFlags().StringVar(&options.Session, "session", "", "Session to use (can't be used with --context)")
	cmd.PersistentFlags().DurationVar(&options.Timeout, "timeout", 0, "Duration for timeouts")
	cmd.PersistentFlags().BoolVar(&options.QuickFixLogging, "quickfix-logging", false, "Enable quickfix logging")
}

func AddPersistentFlagCompletions(cmd *cobra.Command) error {
	if err := cmd.RegisterFlagCompletionFunc("context", complete.Context); err != nil {
		return err
	}

	if err := cmd.RegisterFlagCompletionFunc("initiator", complete.Initiator); err != nil {
		return err
	}

	if err := cmd.RegisterFlagCompletionFunc("session", complete.Session); err != nil {
		return err
	}

	return nil
}
