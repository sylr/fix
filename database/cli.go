package database

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/config"
)

func ValidateOptions(cmd *cobra.Command, args []string) error {
	return config.Load()
}
