package config

import (
	"log/slog"

	"github.com/rs/zerolog"
)

var (
	logger *zerolog.Logger
)

func GetLogger() *zerolog.Logger {
	return logger
}

func SetLogger(l *zerolog.Logger) {
	logger = l
}

// GetSlogLogger returns the logger fix2 logs session events to: the zerolog
// logger with --fix-logging, else one discarding them.
func GetSlogLogger() *slog.Logger {
	if !options.FixLogging || logger == nil {
		return slog.New(slog.DiscardHandler)
	}

	return slog.New(zerolog.NewSlogHandler(*logger))
}

// GetMessageLogger returns the logger fix2 logs the messages sent and
// received to, field by field, at debug level: the zerolog logger, whose
// level is set by --verbose.
func GetMessageLogger() *slog.Logger {
	if logger == nil {
		return slog.New(slog.DiscardHandler)
	}

	return slog.New(zerolog.NewSlogHandler(*logger))
}
