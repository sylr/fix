package utils

import (
	"fmt"

	"github.com/quickfixgo/quickfix"
	"github.com/rs/zerolog"
)

type quickFixLog struct {
	prefix string
	logger *zerolog.Logger
}

func (l quickFixLog) OnIncoming(s []byte) {
	l.logger.Debug().Msgf("quickfix(%s, incoming): %s", l.prefix, s)
}

func (l quickFixLog) OnOutgoing(s []byte) {
	l.logger.Debug().Msgf("quickfix(%s, outgoing): %s", l.prefix, s)
}

func (l quickFixLog) OnEvent(s string) {
	l.logger.Debug().Msgf("quickfix(%s, event): %s", l.prefix, s)
}

func (l quickFixLog) OnEventf(format string, a ...interface{}) {
	l.OnEvent(fmt.Sprintf(format, a...))
}

type quickfixLogFactory struct {
	logger *zerolog.Logger
}

func (q quickfixLogFactory) Create() (quickfix.Log, error) {
	log := quickFixLog{prefix: "global", logger: q.logger}
	return log, nil
}

func (q quickfixLogFactory) CreateSessionLog(sessionID quickfix.SessionID) (quickfix.Log, error) {
	log := quickFixLog{prefix: sessionID.String(), logger: q.logger}
	return log, nil
}

//NewQuickFixLogFactory creates an instance of LogFactory that writes messages and events to stdout.
func NewQuickFixLogFactory(logger *zerolog.Logger) quickfix.LogFactory {
	return quickfixLogFactory{logger: logger}
}
