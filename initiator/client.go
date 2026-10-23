package initiator

import (
	"context"
	stderrors "errors"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/session"
	"github.com/sylr/fix2/transport"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/fixt11"
	"sylr.dev/fix/v2/internal/fixmsg"
)

// Options configure a Client.
type Options struct {
	// MsgTypes are the application messages delivered on Messages, every
	// one if empty: the others are logged and dropped. Session-level
	// Rejects (35=3) are always delivered.
	MsgTypes []string
	// Buffer is how many messages the client holds for the reader of
	// Messages, 1024 if 0. Once full, the session waits for the reader.
	Buffer int
}

// event is what the session hands the reader of Messages, in order: a
// message, or the logout that closes Done.
type event struct {
	message *fix.Message
	logout  bool
}

// Client is a session logged on to an acceptor.
type Client struct {
	// Session is the session's configuration.
	Session *config.Session
	// Dictionary names and describes the fields of the session's messages,
	// header and trailer included; nil without dictionaries.
	Dictionary *dictionary.Dict

	session   *session.Session
	initiator *session.Initiator
	logger    *zerolog.Logger
	opts      Options

	// events holds what the session received; pump hands its messages one
	// by one on messages, which is unbuffered: once pump closes done, the
	// reader has received every message before the logout.
	events   chan event
	messages chan *fix.Message
	loggedOn chan struct{}
	done     chan struct{}
	stopping chan struct{}

	// logonAccepted is set on the reader's goroutine when the counterparty's
	// Logon established the session, before OnLogon runs on the delivery's:
	// a Logout after it is not a refusal.
	logonAccepted atomic.Bool

	mu         sync.Mutex
	reason     error
	logonOnce  sync.Once
	doneOnce   sync.Once
	closeOnce  sync.Once
	logoutOnce sync.Once

	cancel  context.CancelFunc
	stopped chan struct{}
}

// Connect connects to the acceptor of the current context, which must hold
// one session, and waits for the session to log on.
func Connect(ctx context.Context, opts Options) (*Client, error) {
	context, err := config.GetCurrentContext()
	if err != nil {
		return nil, err
	}

	initiatorConfig, err := context.GetInitiator()
	if err != nil {
		return nil, err
	}

	sessions, err := context.GetSessions()
	if err != nil {
		return nil, err
	} else if len(sessions) != 1 {
		return nil, errors.ConfigContextMultipleSessions
	}

	return ConnectSession(ctx, initiatorConfig, sessions[0], opts)
}

// ConnectSession connects s to the acceptor of initiatorConfig and waits for
// it to log on, until --timeout, the session's LogonTimeout or 5s.
func ConnectSession(ctx context.Context, initiatorConfig *config.Initiator, s *config.Session, opts Options) (*Client, error) {
	if opts.Buffer == 0 {
		opts.Buffer = 1024
	}

	c := &Client{
		Session:  s,
		logger:   config.GetLogger(),
		opts:     opts,
		events:   make(chan event, opts.Buffer),
		messages: make(chan *fix.Message),
		loggedOn: make(chan struct{}),
		done:     make(chan struct{}),
		stopping: make(chan struct{}),
		stopped:  make(chan struct{}),
	}
	_, display, err := s.FIXDictionaries()
	if err != nil {
		return nil, err
	}
	c.Dictionary = display

	observer := session.SlogObserver{Logger: config.GetMessageLogger(), Dictionary: display}
	cfg, err := s.FIXConfig(c, config.GetSlogLogger(), observer)
	if err != nil {
		return nil, err
	}

	dialer, err := initiatorConfig.Dialer()
	if err != nil {
		return nil, err
	}

	if cfg.Store, err = initiatorConfig.Store.Open("initiator.db", cfg.ID); err != nil {
		return nil, err
	}

	c.initiator = session.NewInitiator(session.InitiatorOptions{
		Logger:       config.GetSlogLogger(),
		DialTimeout:  initiatorConfig.DialTimeoutOrDefault(),
		ReconnectMin: initiatorConfig.ReconnectMin,
		ReconnectMax: initiatorConfig.ReconnectMax,
	})

	if c.session, err = c.initiator.Add(cfg, &logonDialer{Dialer: dialer, client: c}); err != nil {
		cfg.Store.Close()
		return nil, err
	}

	runCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	c.cancel = cancel
	go c.pump()
	go func() {
		defer close(c.stopped)
		c.initiator.Run(runCtx)
	}()

	select {
	case <-c.loggedOn:
		return c, nil
	case <-c.done:
		select {
		case <-c.loggedOn:
			// Logged on, then out at once: the reader gets what the
			// session received in between, then Done.
			return c, nil
		default:
		}

		err := c.Err()
		c.abort()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errors.FixLogout, err)
		}
		return nil, errors.FixLogout
	case <-time.After(config.Timeout(s.LogonTimeout)):
		c.abort()
		return nil, errors.ConnectionTimeout
	case <-ctx.Done():
		c.abort()
		return nil, ctx.Err()
	}
}

// abort stops the session without waiting for it, so that the initiator
// stops reconnecting.
func (c *Client) abort() {
	c.Close(false)
	c.cancel()
}

// Send sends m, once checked against the session's dictionary (see
// fixmsg.Check): a message its FIX version does not define, or that lacks
// what it requires, is not sent.
func (c *Client) Send(ctx context.Context, m session.Outgoing) (session.Receipt, error) {
	if c.Dictionary != nil {
		msg, err := fixmsg.Message(m, c.Session.BeginString, c.Dictionary)
		if err != nil {
			return session.Receipt{}, err
		}
		if err := fixmsg.Check(msg, c.Dictionary, c.Session.Checks()); err != nil {
			return session.Receipt{}, err
		}
	}

	return c.session.Send(ctx, m)
}

// BeginString returns the BeginString (8) of the session's messages.
func (c *Client) BeginString() string {
	return c.Session.BeginString
}

// Messages returns the application messages received, and the Rejects. It
// is never closed: Done tells when the session logged out.
func (c *Client) Messages() <-chan *fix.Message {
	return c.messages
}

// Done returns a channel closed once the session logged out or was
// disconnected, or its Logon was refused; Err then tells why. By then, the
// reader of Messages received every message received before.
func (c *Client) Done() <-chan struct{} {
	return c.done
}

// pump hands the events received to the reader of Messages, in order,
// closing done at the first logout.
func (c *Client) pump() {
	for {
		select {
		case ev := <-c.events:
			if ev.logout {
				c.doneOnce.Do(func() { close(c.done) })
				continue
			}
			select {
			case c.messages <- ev.message:
			case <-c.stopping:
				return
			}
		case <-c.stopping:
			return
		}
	}
}

// push queues ev for pump, waiting for room unless the client is closing.
func (c *Client) push(ev event) {
	select {
	case c.events <- ev:
	case <-c.stopping:
	}
}

// loggedOut records why the session logged out, unless asked to, and ends
// the session: a logout is final for a Client. It queues the logout once the
// session stopped, behind every message received before, those fix2 holds
// past OnLogout while it cannot record them included.
func (c *Client) loggedOut(reason error) {
	c.mu.Lock()
	select {
	case <-c.stopping:
	default:
		if c.reason == nil {
			c.reason = reason
		}
	}
	c.mu.Unlock()

	c.logoutOnce.Do(func() {
		go func() {
			c.cancel()
			<-c.session.Done()
			c.push(event{logout: true})
		}()
	})
}

// Err returns why the session logged out, nil if it has not or was asked to.
func (c *Client) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.reason
}

// Close stops the session, logging it out if logout is true; otherwise the
// connection is left to close with the process.
func (c *Client) Close(logout bool) {
	c.closeOnce.Do(func() {
		close(c.stopping)

		if logout {
			c.cancel()
			<-c.stopped
		}
	})
}

func (c *Client) msgTypeName(msgType string) string {
	if name, err := dict.SearchValue(dict.MessageTypes, fix50sp2.MsgType(msgType)); err == nil {
		return fmt.Sprintf("%s(%s)", msgType, name)
	}

	return msgType
}

// Authenticate implements session.Handler: an initiator trusts its acceptor.
func (c *Client) Authenticate(context.Context, *session.LogonRequest) error {
	return nil
}

// OnLogon implements session.Handler.
func (c *Client) OnLogon(_ context.Context, s *session.Session) {
	c.logger.Debug().Msgf("Logon: %s", s.ID())
	c.logonOnce.Do(func() { close(c.loggedOn) })
}

// OnLogout implements session.Handler.
func (c *Client) OnLogout(_ context.Context, s *session.Session, reason error) {
	c.logger.Debug().Msgf("Logout: %s: %v", s.ID(), reason)
	c.loggedOut(reason)
}

// ToAdmin implements session.Handler: it sets the session's Username (553)
// and Password (554) on Logons.
func (c *Client) ToAdmin(_ context.Context, _ *session.Session, msgType string, b *fix.Builder) error {
	if msgType != string(fixt11.MsgTypeLogon) {
		return nil
	}

	if len(c.Session.Username) > 0 {
		c.logger.Debug().Msg("Username injected in logon message")
		b.SetString(fixt11.TagUsername, c.Session.Username)
	}

	if len(c.Session.Password) > 0 {
		c.logger.Debug().Msg("Password injected in logon message")
		b.SetString(fixt11.TagPassword, c.Session.Password)
	}

	return nil
}

// FromAdmin implements session.Handler: it delivers Rejects, and ends a
// Connect whose Logon the counterparty refused, which OnLogout is not told
// of. It must return quickly: a Reject is dropped if the reader lags far
// behind.
func (c *Client) FromAdmin(_ context.Context, s *session.Session, m *fix.Message, _ session.Inbound) {
	switch string(m.MsgType()) {
	case string(fixt11.MsgTypeLogon):
		// FromAdmin is told of a Logon the session refused too, logging
		// out (an invalid NextExpectedMsgSeqNum, say): only one that
		// established the session is accepted.
		switch s.State() {
		case session.StateLoggedOn, session.StateStalled, session.StateResettingSeqs:
			c.logonAccepted.Store(true)
		}

	case string(fixt11.MsgTypeReject):
		select {
		case c.events <- event{message: m.Clone()}:
		default:
			c.logger.Warn().Msgf("Reject dropped, messages buffer full: %s", m.Get(fixt11.TagText))
		}

	case string(fixt11.MsgTypeLogout):
		if c.logonAccepted.Load() {
			return
		}

		reason := "logon refused"
		if text := m.Get(fixt11.TagText); len(text) > 0 {
			reason += ": " + string(text)
		}

		c.refused(reason)
	}
}

// refused ends a Connect whose Logon failed.
func (c *Client) refused(reason string) {
	c.mu.Lock()
	if c.reason == nil {
		c.reason = stderrors.New(reason)
	}
	c.mu.Unlock()

	c.doneOnce.Do(func() { close(c.done) })
}

// logonDialer dials the counterparty for a Client. A connection that ended
// before the session logged on had its Logon refused, the counterparty
// closing it without a Logout: rather than send the credentials again,
// logonDialer ends Connect. Dial failures are retried until its timeout.
type logonDialer struct {
	transport.Dialer
	client    *Client
	connected atomic.Bool
}

func (d *logonDialer) Dial(ctx context.Context) (transport.Conn, error) {
	if d.connected.Load() && !d.client.logonAccepted.Load() {
		d.client.refused("logon refused: the counterparty closed the connection")
		return nil, stderrors.New("logon refused")
	}

	conn, err := d.Dialer.Dial(ctx)
	if err == nil {
		d.connected.Store(true)
	}

	return conn, err
}

// FromApp implements session.Handler.
func (c *Client) FromApp(_ context.Context, _ *session.Session, m *fix.Message, _ session.Inbound) error {
	msgType := string(m.MsgType())

	if len(c.opts.MsgTypes) > 0 && !slices.Contains(c.opts.MsgTypes, msgType) {
		c.logger.Info().Msgf("Received unexpected message type: %s", c.msgTypeName(msgType))
		return nil
	}

	c.push(event{message: m.Clone()})

	return nil
}

// OnResend implements session.Handler.
func (c *Client) OnResend(context.Context, *session.Session, uint64, *fix.Message) bool {
	return true
}

// OnRecoveryLoss implements session.Handler.
func (c *Client) OnRecoveryLoss(_ context.Context, _ *session.Session, begin, end uint64, cause error) {
	c.logger.Warn().Msgf("Messages %d to %d lost: %v", begin, end, cause)
}

var _ session.Handler = (*Client)(nil)
