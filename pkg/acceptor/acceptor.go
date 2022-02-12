package acceptor

import (
	"github.com/quickfixgo/quickfix"
)

func NewAcceptor(app quickfix.Application, settings *quickfix.Settings) (*quickfix.Acceptor, error) {
	return quickfix.NewAcceptor(app, quickfix.NewMemoryStoreFactory(), settings, quickfix.NewNullLogFactory())
}
