package utils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/quickfixgo/quickfix"
	"github.com/rs/zerolog"
)

type quickFixLog struct {
	prefix string
	logger *zerolog.Logger
	limit  int
}

func (l quickFixLog) OnIncoming(s []byte) {
	if l.logger != nil {
		if l.limit == 0 || len(s) <= l.limit {
			l.logger.Debug().Msgf("quickfix(%s, incoming): %s", l.prefix, bytes.ReplaceAll(s, []byte("\001"), []byte("|")))
		} else {
			l.logger.Debug().Msgf(
				"quickfix(%s, incoming[0:%d:%s]): %s...",
				l.prefix,
				l.limit,
				strings.ReplaceAll(humanize.Bytes(uint64(len(s))), " ", ""),
				bytes.ReplaceAll(s[0:l.limit], []byte("\001"), []byte("|")),
			)
		}
	}
}

func (l quickFixLog) OnOutgoing(s []byte) {
	if l.logger != nil {
		if l.limit == 0 || len(s) <= l.limit {
			l.logger.Debug().Msgf("quickfix(%s, outgoing): %s", l.prefix, bytes.ReplaceAll(s, []byte("\001"), []byte("|")))
		} else {
			l.logger.Debug().Msgf(
				"quickfix(%s, outgoing[0:%d:%s]): %s...",
				l.prefix,
				l.limit,
				strings.ReplaceAll(humanize.Bytes(uint64(len(s))), " ", ""),
				bytes.ReplaceAll(s[0:l.limit], []byte("\001"), []byte("|")),
			)
		}
	}
}

func (l quickFixLog) OnEvent(s string) {
	if l.logger != nil {
		l.logger.Debug().Msgf("quickfix(%s, event): %s", l.prefix, s)
	}
}

func (l quickFixLog) OnEventf(format string, a ...interface{}) {
	l.OnEvent(fmt.Sprintf(format, a...))
}

type quickfixLogFactory struct {
	logger *zerolog.Logger
	limit  int
}

func (q quickfixLogFactory) Create() (quickfix.Log, error) {
	log := quickFixLog{prefix: "global", logger: q.logger, limit: q.limit}
	return log, nil
}

func (q quickfixLogFactory) CreateSessionLog(sessionID quickfix.SessionID) (quickfix.Log, error) {
	log := quickFixLog{prefix: sessionID.String(), logger: q.logger, limit: q.limit}
	return log, nil
}

// NewQuickFixLogFactory creates an instance of LogFactory that writes messages and events to stdout.
func NewQuickFixLogFactory(logger *zerolog.Logger, limit int) quickfix.LogFactory {
	return quickfixLogFactory{logger: logger, limit: limit}
}
