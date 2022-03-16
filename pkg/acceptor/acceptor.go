package acceptor

import (
	"github.com/quickfixgo/quickfix"
	"github.com/rs/zerolog"

	"sylr.dev/fix/pkg/utils"
)

func NewAcceptor(app quickfix.Application, settings *quickfix.Settings, logger *zerolog.Logger) (*quickfix.Acceptor, error) {
	var msgStoreFactory quickfix.MessageStoreFactory

	if settings.GlobalSettings().HasSetting("SQLStoreDriver") {
		if driver, err := settings.GlobalSettings().Setting("SQLStoreDriver"); err == nil && driver == "sqlite3" {
			msgStoreFactory = quickfix.NewSQLStoreFactory(settings)
		}
	}

	if msgStoreFactory == nil {
		msgStoreFactory = quickfix.NewMemoryStoreFactory()
	}

	return quickfix.NewAcceptor(app, msgStoreFactory, settings, utils.NewQuickFixLogFactory(logger))
}
