package utils

import (
	"github.com/spf13/cobra"
)

// MakePersistentPreRunE returns a command's PersistentPreRunE, which checks
// its required flags, then runs validator. With cobra.EnableTraverseRunHooks,
// which the root command sets, cobra runs the parents' PersistentPreRunE
// before it.
func MakePersistentPreRunE(validator func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := cmd.ValidateRequiredFlags(); err != nil {
			return err
		}

		return validator(cmd, args)
	}
}
