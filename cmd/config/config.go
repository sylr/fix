package config

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/cmd/config/init"
)

// ConfigCmd represents the buy command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage fix configuration",
	Long:  "Manage fix configuration.",
}

func init() {
	ConfigCmd.AddCommand(configinit.ConfigInitCmd)
}
