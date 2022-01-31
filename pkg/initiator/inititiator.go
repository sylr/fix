package initiator

import (
	"github.com/quickfixgo/quickfix"
)

func Initiate(app quickfix.Application, settings *quickfix.Settings) (*quickfix.Initiator, error) {
	return quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), settings, quickfix.NewNullLogFactory())
}
