package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"sylr.dev/fix/cmd/list"
	"sylr.dev/fix/cmd/new"
	"sylr.dev/fix/config"
)

var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var FixCmd = &cobra.Command{
	Use:          "fix",
	SilenceUsage: true,
	Version:      Version,
}

func init() {
	options := config.GetOptions()

	FixCmd.AddCommand(new.NewCmd)
	FixCmd.AddCommand(list.ListCmd)

	configPath := strings.Join([]string{"$HOME", ".fix", "config"}, string(os.PathSeparator))

	FixCmd.PersistentFlags().StringVar(&options.Config, "config", os.ExpandEnv(configPath), "Config file")
	FixCmd.PersistentFlags().CountVarP(&options.Verbose, "verbose", "v", "Increase verbosity")
	FixCmd.PersistentFlags().BoolVar(&options.LogCaller, "log-caller", false, "Add caller info to log lines")
	FixCmd.PersistentFlags().BoolVar(&options.Interactive, "interactive", true, "Enable interactive mode")
	FixCmd.PersistentFlags().BoolP("help", "h", false, "Help for fix")
	FixCmd.PersistentFlags().Bool("version", false, "Version for fix")
}
