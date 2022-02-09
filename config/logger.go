package config

import "github.com/rs/zerolog"

var (
	logger *zerolog.Logger
)

func GetLogger() *zerolog.Logger {
	return logger
}

func SetLogger(l *zerolog.Logger) {
	logger = l
}
