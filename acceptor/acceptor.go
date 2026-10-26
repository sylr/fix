package acceptor

import (
	"context"

	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
)

// Serve serves the sessions of the current context (see ServeContext).
func Serve(ctx context.Context, handler func(*config.Session) (session.Handler, error)) error {
	cfgContext, err := config.GetCurrentContext()
	if err != nil {
		return err
	}

	return ServeContext(ctx, cfgContext, handler)
}

// ServeContext serves the sessions of cfgContext on its acceptor's address
// until ctx is done, each session's messages handled by the handler handler
// returns for it. If it fails to set them up, it stops those it set up.
func ServeContext(ctx context.Context, cfgContext *config.Context, handler func(*config.Session) (session.Handler, error)) error {
	acceptorConfig, err := cfgContext.GetAcceptor()
	if err != nil {
		return err
	}

	sessions, err := cfgContext.GetSessions()
	if err != nil {
		return err
	} else if len(sessions) == 0 {
		return errors.ErrConfigContextNoSession
	}

	for _, s := range sessions {
		if err := s.CheckSupported(); err != nil {
			return err
		}
	}

	listener, err := acceptorConfig.Listen()
	if err != nil {
		return err
	}

	acceptor := session.NewAcceptor(session.AcceptorOptions{
		Logger:       config.GetSlogLogger(),
		LogonTimeout: config.Timeout(0),
	})

	if err := add(acceptor, acceptorConfig, sessions, handler); err != nil {
		_ = listener.Close()
		for _, s := range acceptor.Sessions() {
			_ = acceptor.Remove(context.Background(), s.ID())
		}
		return err
	}

	config.GetLogger().Info().Msgf("Listening on %s", listener.Addr())
	acceptor.Listen(listener)

	return acceptor.Run(ctx)
}

// add adds sessions to acceptor, each with the handler handler returns for
// it and a store of acceptorConfig.
func add(acceptor *session.Acceptor, acceptorConfig *config.Acceptor, sessions []*config.Session, handler func(*config.Session) (session.Handler, error)) error {
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
			cfg.Store.Close()
			return err
		}
	}

	return nil
}
