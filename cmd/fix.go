package cmd

import (
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/cmd/amend"
	"sylr.dev/fix/v2/cmd/cancel"
	initcmd "sylr.dev/fix/v2/cmd/init"
	"sylr.dev/fix/v2/cmd/initiator"
	"sylr.dev/fix/v2/cmd/list"
	"sylr.dev/fix/v2/cmd/marketdata"
	"sylr.dev/fix/v2/cmd/new"
	"sylr.dev/fix/v2/cmd/probe"
	"sylr.dev/fix/v2/cmd/status"
	"sylr.dev/fix/v2/config"
)

var Version = "dev"

// FixCmd is the root command.
var FixCmd = &cobra.Command{
	Use:          "fix",
	SilenceUsage: true,
	Version:      Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := InitLogger(cmd, args); err != nil {
			return err
		}
		return InitHTTP(cmd, args)
	},
}

func init() {
	// Every command's PersistentPreRunE runs, the root's first: the logger
	// is set up, then the parents load the configuration, before the
	// commands validate their options.
	cobra.EnableTraverseRunHooks = true

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

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".fix", "config")

	FixCmd.PersistentFlags().StringVar(&options.Config, "config", configPath, "Config file")
	FixCmd.PersistentFlags().CountVarP(&options.Verbose, "verbose", "v", "Increase verbosity")
	FixCmd.PersistentFlags().BoolVar(&options.LogCaller, "log-caller", false, "Add caller info to log lines")
	FixCmd.PersistentFlags().BoolVar(&options.Interactive, "interactive", true, "Enable interactive mode")
	FixCmd.PersistentFlags().BoolP("help", "h", false, "Help for fix")
	FixCmd.PersistentFlags().Bool("version", false, "Version for fix")
	FixCmd.PersistentFlags().BoolVar(&options.Metrics, "metrics", false, "Enable metrics")
	FixCmd.PersistentFlags().BoolVar(&options.PProf, "pprof", false, "Enable pprof")
	FixCmd.PersistentFlags().IntVar(&options.HTTPPort, "port", 8080, "HTTP port")
}

// ANSI colours of the log levels.
const (
	colorRed     = 31
	colorGreen   = 32
	colorYellow  = 33
	colorMagenta = 35
	colorCyan    = 36
)

func consoleDefaultFormatMessage() zerolog.Formatter {
	return func(i any) string {
		if i == nil || i == "" {
			return ""
		}

		return fmt.Sprintf("%s", i)
	}
}

// InitLogger sets up the logger, which writes to stderr, apart from the
// commands' output.
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
		Out:           os.Stderr,
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

// InitHTTP serves the metrics and pprof endpoints, if enabled, on --port.
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

	listener, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(options.HTTPPort)))
	if err != nil {
		return fmt.Errorf("metrics and pprof endpoints: %w", err)
	}

	server := &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	go func() {
		if err := server.Serve(listener); err != nil {
			config.GetLogger().Error().Err(err).Msg("Metrics and pprof endpoints stopped")
		}
	}()

	return nil
}
