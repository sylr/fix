package cmd

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"sylr.dev/fix/cmd/amend"
	"sylr.dev/fix/cmd/cancel"
	initcmd "sylr.dev/fix/cmd/init"
	"sylr.dev/fix/cmd/initiator"
	"sylr.dev/fix/cmd/list"
	"sylr.dev/fix/cmd/marketdata"
	"sylr.dev/fix/cmd/new"
	"sylr.dev/fix/cmd/probe"
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
		InitHTTP(cmd, args)
		return InitLogger(cmd, args)
	},
}

func init() {
	options := config.GetOptions()

	FixCmd.AddCommand(amend.AmendCmd)
	FixCmd.AddCommand(cancel.CancelCmd)
	FixCmd.AddCommand(initcmd.InitCmd)
	FixCmd.AddCommand(initiator.InitiatorCmd)
	FixCmd.AddCommand(list.ListCmd)
	FixCmd.AddCommand(marketdata.MarketDataCmd)
	FixCmd.AddCommand(new.NewCmd)
	FixCmd.AddCommand(probe.ProbeCmd)
	FixCmd.AddCommand(status.StatusCmd)

	configPath := filepath.Join("$HOME", ".fix", "config")

	FixCmd.PersistentFlags().StringVar(&options.Config, "config", os.ExpandEnv(configPath), "Config file")
	FixCmd.PersistentFlags().CountVarP(&options.Verbose, "verbose", "v", "Increase verbosity")
	FixCmd.PersistentFlags().BoolVar(&options.LogCaller, "log-caller", false, "Add caller info to log lines")
	FixCmd.PersistentFlags().BoolVar(&options.Interactive, "interactive", true, "Enable interactive mode")
	FixCmd.PersistentFlags().BoolP("help", "h", false, "Help for fix")
	FixCmd.PersistentFlags().Bool("version", false, "Version for fix")
	FixCmd.PersistentFlags().BoolVar(&options.Metrics, "metrics", false, "Enable metrics")
	FixCmd.PersistentFlags().BoolVar(&options.PProf, "pprof", false, "Enable pprof")
	FixCmd.PersistentFlags().IntVar(&options.HTTPPort, "port", 8080, "HTTP port")
}

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

func consoleDefaultFormatMessage() zerolog.Formatter {
	return func(i interface{}) string {
		if i == nil || i == "" {
			return ""
		}

		return fmt.Sprintf("%s", i)
	}
}

func InitLogger(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.LevelColors = map[zerolog.Level]int{
		zerolog.TraceLevel: colorMagenta,
		zerolog.DebugLevel: colorYellow,
		zerolog.InfoLevel:  colorGreen,
		zerolog.WarnLevel:  colorCyan,
		zerolog.ErrorLevel: colorRed,
		zerolog.FatalLevel: colorRed,
		zerolog.PanicLevel: colorRed,
	}
	consoleWriter := zerolog.ConsoleWriter{
		Out:           os.Stdout,
		TimeFormat:    "Jan 2 15:04:05.000-0700",
		FormatMessage: consoleDefaultFormatMessage(),
	}
	multi := zerolog.MultiLevelWriter(consoleWriter)
	logger := zerolog.New(multi).With().Timestamp().Logger().Level(config.IntToZerologLevel(options.Verbose))

	if options.LogCaller {
		logger = logger.With().Caller().Logger()
	}
	config.SetLogger(&logger)
	return nil
}

func InitHTTP(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()

	if !options.Metrics && !options.PProf {
		return nil
	}

	mux := http.NewServeMux()

	if options.Metrics {
		mux.Handle("/metrics", promhttp.Handler())
	}
	if options.PProf {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	go http.ListenAndServe(fmt.Sprintf(":%d", options.HTTPPort), mux)

	return nil
}
