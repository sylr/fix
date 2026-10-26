package acceptor

import (
	"fmt"

	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/cli/complete"
	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
)

// ValidateOptions loads the configuration and resolves the context the
// command runs in (see resolveContext), which must have an acceptor and at
// least one session.
func ValidateOptions(cmd *cobra.Command, args []string) error {
	if err := config.Load(); err != nil {
		return err
	}

	cfgContext, err := resolveContext()
	if err != nil {
		return err
	}

	if len(cfgContext.Acceptor) == 0 {
		return fmt.Errorf("%w: context %s has no acceptor", errors.ErrConfig, cfgContext.Name)
	}

	if _, err := cfgContext.GetAcceptor(); err != nil {
		return err
	}

	sessions, err := cfgContext.GetSessions()
	if err != nil {
		return err
	} else if len(sessions) == 0 {
		return errors.ErrConfigContextNoSession
	}

	return nil
}

// resolveContext returns the context --context names, or the one
// --acceptor and --session make up, which becomes the current-context, or
// else the configuration's current-context.
func resolveContext() (*config.Context, error) {
	options := config.GetOptions()

	switch {
	case len(options.Context) > 0:
		if len(options.Acceptor) > 0 || len(options.Session) > 0 {
			return nil, fmt.Errorf("%w: can't use --acceptor/--session with --context", errors.ErrOptions)
		}
	case len(options.Acceptor) > 0 && len(options.Session) > 0:
		config.SetCurrentContext(&config.Context{
			Name:     fmt.Sprintf("(%s/%s)", options.Acceptor, options.Session),
			Acceptor: options.Acceptor,
			Sessions: []string{options.Session},
		})
	case len(options.Acceptor) > 0 || len(options.Session) > 0:
		return nil, fmt.Errorf("%w: --acceptor and --session go together", errors.ErrOptions)
	case len(config.GetConfig().CurrentContext) == 0:
		return nil, fmt.Errorf("%w: you need to specify either --context or --acceptor/--session", errors.ErrOptions)
	}

	return config.GetCurrentContext()
}

func AddPersistentFlags(cmd *cobra.Command) {
	options := config.GetOptions()

	cmd.PersistentFlags().StringVar(&options.Context, "context", "", "Context to use")
	cmd.PersistentFlags().StringVar(&options.Acceptor, "acceptor", "", "Acceptor to use, with --session (can't be used with --context)")
	cmd.PersistentFlags().StringVar(&options.Session, "session", "", "Session to use, with --acceptor (can't be used with --context)")
	cmd.PersistentFlags().DurationVar(&options.Timeout, "timeout", 0, "Duration for timeouts")
	cmd.PersistentFlags().BoolVar(&options.FixLogging, "fix-logging", false, "Log FIX session events")
}

func AddPersistentFlagCompletions(cmd *cobra.Command) error {
	if err := cmd.RegisterFlagCompletionFunc("context", complete.Context); err != nil {
		return err
	}

	if err := cmd.RegisterFlagCompletionFunc("acceptor", complete.Acceptor); err != nil {
		return err
	}

	if err := cmd.RegisterFlagCompletionFunc("session", complete.Session); err != nil {
		return err
	}

	return nil
}
