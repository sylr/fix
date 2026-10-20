package acceptor

import (
	"context"

	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
)

// Serve serves the sessions of the current context on its acceptor's address
// until ctx is done, each session's messages handled by the handler handler
// returns for it.
func Serve(ctx context.Context, handler func(*config.Session) (session.Handler, error)) error {
	context, err := config.GetCurrentContext()
	if err != nil {
		return err
	}

	acceptorConfig, err := context.GetAcceptor()
	if err != nil {
		return err
	}

	sessions, err := context.GetSessions()
	if err != nil {
		return err
	}

	acceptor := session.NewAcceptor(session.AcceptorOptions{
		Logger:       config.GetSlogLogger(),
		LogonTimeout: config.Timeout(0),
	})

	for _, s := range sessions {
		h, err := handler(s)
		if err != nil {
			return err
		}

		_, display, err := s.FIXDictionaries()
		if err != nil {
			return err
		}

		observer := session.SlogObserver{Logger: config.GetMessageLogger(), Dictionary: display}
		cfg, err := s.FIXConfig(h, config.GetSlogLogger(), observer)
		if err != nil {
			return err
		}

		if cfg.Store, err = acceptorConfig.Store.Open("acceptor.db", cfg.ID); err != nil {
			return err
		}

		if _, err := acceptor.Add(cfg); err != nil {
			return err
		}
	}

	listener, err := acceptorConfig.Listen()
	if err != nil {
		return err
	}

	config.GetLogger().Info().Msgf("Listening on %s", listener.Addr())
	acceptor.Listen(listener)

	return acceptor.Run(ctx)
}
