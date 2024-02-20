package initiator

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/store/sql"

	"sylr.dev/fix/pkg/utils"
)

func Initiate(app quickfix.Application, settings *quickfix.Settings, logger *zerolog.Logger) (*quickfix.Initiator, error) {
	var msgStoreFactory quickfix.MessageStoreFactory

	if settings.GlobalSettings().HasSetting("SQLStoreDriver") {
		driver, err := settings.GlobalSettings().Setting("SQLStoreDriver")
		if err != nil {
			return nil, err
		}
		switch driver {
		case "sqlite3", "postgres":
			msgStoreFactory = sql.NewStoreFactory(settings)
		default:
			return nil, fmt.Errorf("Unsupported SQLStoreDriver: %s", driver)
		}
	}

	if msgStoreFactory == nil {
		msgStoreFactory = quickfix.NewMemoryStoreFactory()
	}

	return quickfix.NewInitiator(app, msgStoreFactory, settings, utils.NewQuickFixLogFactory(logger))
}
