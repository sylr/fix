package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"sylr.dev/fix/cmd/cancel"
	configcmd "sylr.dev/fix/cmd/config"
	"sylr.dev/fix/cmd/database"
	"sylr.dev/fix/cmd/initiator"
	"sylr.dev/fix/cmd/list"
	"sylr.dev/fix/cmd/marketdata"
	"sylr.dev/fix/cmd/new"
	"sylr.dev/fix/cmd/status"
	"sylr.dev/fix/config"
)

var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var FixCmd = &cobra.Command{
	Use:          "fix",
	SilenceUsage: true,
	Version:      Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return InitLogger(cmd, args)
	},
}

func init() {
	options := config.GetOptions()

	FixCmd.AddCommand(configcmd.ConfigCmd)
	FixCmd.AddCommand(database.DatabaseCmd)
	FixCmd.AddCommand(new.NewCmd)
	FixCmd.AddCommand(marketdata.MarketDataCmd)
	FixCmd.AddCommand(cancel.CancelCmd)
	FixCmd.AddCommand(list.ListCmd)
	FixCmd.AddCommand(initiator.InitiatorCmd)
	FixCmd.AddCommand(status.StatusCmd)

	configPath := strings.Join([]string{"$HOME", ".fix", "config"}, string(os.PathSeparator))

	FixCmd.PersistentFlags().StringVar(&options.Config, "config", os.ExpandEnv(configPath), "Config file")
	FixCmd.PersistentFlags().CountVarP(&options.Verbose, "verbose", "v", "Increase verbosity")
	FixCmd.PersistentFlags().BoolVar(&options.LogCaller, "log-caller", false, "Add caller info to log lines")
	FixCmd.PersistentFlags().BoolVar(&options.Interactive, "interactive", true, "Enable interactive mode")
	FixCmd.PersistentFlags().BoolP("help", "h", false, "Help for fix")
	FixCmd.PersistentFlags().Bool("version", false, "Version for fix")
}

func InitLogger(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()
	zerolog.TimeFieldFormat = time.RFC3339Nano
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "Jan 2 15:04:05.000-0700",
	}
	multi := zerolog.MultiLevelWriter(consoleWriter)
	logger := zerolog.New(multi).With().Timestamp().Logger().Level(config.IntToZerologLevel(options.Verbose))

	if options.LogCaller {
		logger = logger.With().Caller().Logger()
	}
	config.SetLogger(&logger)
	return nil
}
