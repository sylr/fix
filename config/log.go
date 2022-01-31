package config

import "github.com/rs/zerolog"

func IntToZerologLevel(level int) zerolog.Level {
	switch level {
	case 0:
		return zerolog.InfoLevel
	case 1:
		return zerolog.DebugLevel
	default:
		return zerolog.TraceLevel
	}
}
