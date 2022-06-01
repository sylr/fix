package acceptor

import (
	"github.com/rs/zerolog"

	"github.com/quickfixgo/quickfix"

	"sylr.dev/fix/pkg/utils"
)

func NewAcceptor(app quickfix.Application, settings *quickfix.Settings, logger *zerolog.Logger) (*quickfix.Acceptor, error) {
	var msgStoreFactory quickfix.MessageStoreFactory

	if settings.GlobalSettings().HasSetting("SQLStoreDriver") {
		driver, err := settings.GlobalSettings().Setting("SQLStoreDriver")
		if err != nil {
			return nil, err
		}
		switch driver {
		case "sqlite3":
		case "postgres":
			msgStoreFactory = quickfix.NewSQLStoreFactory(settings)
		}
	}

	if msgStoreFactory == nil {
		msgStoreFactory = quickfix.NewMemoryStoreFactory()
	}

	return quickfix.NewAcceptor(app, msgStoreFactory, settings, utils.NewQuickFixLogFactory(logger))
}
