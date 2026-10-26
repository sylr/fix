package initiator

import (
	"bytes"
	"cmp"
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/session"
	"github.com/sylr/fix2/transport"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/fixt11"
)

// Options configure a Client.
type Options struct {
	// MsgTypes are the application messages delivered on Messages, every
	// one if empty: the others are logged and dropped. Rejects (35=3) and
	// BusinessMessageRejects (35=j) are always delivered.
	MsgTypes []string
	// Buffer is how many messages the client holds for the reader of
	// Messages, 1024 if not positive. Once full, the session waits for the
	// reader.
	Buffer int
}

// Client is a session logged on to an acceptor.
//
// It delivers the messages received on Messages in the order they were
// received: a Reject (35=3), which the session handles apart from
// application messages, is held until every application message received
// before it was delivered, or skipped by a GapFill, or the session ended.
type Client struct {
	// Session is the session's configuration.
	Session *config.Session
	// Dictionary names and describes the fields of the session's messages,
	// header and trailer included; nil without dictionaries.
	Dictionary *dictionary.Dict

	session *session.Session
	logger  *zerolog.Logger
	opts    Options

	// stop stops the initiator, which closes stopped once it did.
	stop     context.CancelFunc
	stopped  chan struct{}
	stopOnce sync.Once

	loggedOn chan struct{}
	// messages is unbuffered: once pump closes done, the reader received
	// every message queued.
	messages chan *fix.Message
	// closing is closed by Close, ended once the session stopped, done by
	// pump once it handed every message, or Close discarded them.
	closing chan struct{}
	ended   chan struct{}
	done    chan struct{}
	wake    chan struct{}

	mu sync.Mutex
	// queue holds the messages for pump; room is closed when it has room
	// again for an application message.
	queue  []*fix.Message
	room   chan struct{}
	reason error
	ending bool
	closed bool
	// accepted is set when the counterparty's Logon established the
	// session: a Logout or a disconnection after it is not a refusal.
	accepted bool
	conn     transport.Conn
	dialed   bool
	dialErr  error
	// sent holds the messages sent with Send by MsgSeqNum, for Await to
	// match the Rejects: their receipts tell which, of the generations a
	// MsgSeqNum was sent in, the counterparty rejected.
	sent map[uint64][]sentMessage
	// inbound holds where the application messages received, not
	// delivered yet, stand, in order; held are the Rejects waiting for
	// them. gapFills are the GapFills received not processed yet, whose
	// NewSeqNo may skip messages received ahead of a gap: those are never
	// delivered.
	inbound  []pos
	held     []heldReject
	gapFills []gapFill
}

// sentMessage is a message sent with Send.
type sentMessage struct {
	receipt session.Receipt
	msgType string
}

// pos is where a message received stands: its sequence number generation,
// then its MsgSeqNum. A reset starts a later generation.
type pos struct {
	gen, seq uint64
}

func (p pos) compare(q pos) int {
	return cmp.Or(cmp.Compare(p.gen, q.gen), cmp.Compare(p.seq, q.seq))
}

// heldReject is a Reject waiting for the application messages received
// before it to be delivered.
type heldReject struct {
	pos
	message *fix.Message
}

// gapFill is a SequenceReset-GapFill received, at pos, to NewSeqNo end.
type gapFill struct {
	pos
	end uint64
}

// Connect connects to the acceptor of the current context, which must hold
// one session, and waits for the session to log on.
func Connect(ctx context.Context, opts Options) (*Client, error) {
	cfgContext, err := config.GetCurrentContext()
	if err != nil {
		return nil, err
	}

	initiatorConfig, err := cfgContext.GetInitiator()
	if err != nil {
		return nil, err
	}

	sessions, err := cfgContext.GetSessions()
	if err != nil {
		return nil, err
	}

	switch len(sessions) {
	case 0:
		return nil, errors.ErrConfigContextNoSession
	case 1:
	default:
		return nil, errors.ErrConfigContextMultipleSessions
	}

	return ConnectSession(ctx, initiatorConfig, sessions[0], opts)
}

// ConnectSession connects s to the acceptor of initiatorConfig and waits for
// it to log on, until --timeout, the session's LogonTimeout or 5s. It fails
// with an error wrapping errors.ErrFixLogonRefused if the counterparty
// refuses the Logon, errors.ErrConnectionTimeout and the last dial error if
// the session did not log on in time, or ctx's error.
func ConnectSession(ctx context.Context, initiatorConfig *config.Initiator, s *config.Session, opts Options) (*Client, error) {
	if opts.Buffer <= 0 {
		opts.Buffer = 1024
	}

	c := &Client{
		Session:  s,
		logger:   config.GetLogger(),
		opts:     opts,
		stopped:  make(chan struct{}),
		loggedOn: make(chan struct{}),
		messages: make(chan *fix.Message),
		closing:  make(chan struct{}),
		ended:    make(chan struct{}),
		done:     make(chan struct{}),
		wake:     make(chan struct{}, 1),
		sent:     make(map[uint64][]sentMessage),
	}
	_, display, err := s.FIXDictionaries()
	if err != nil {
		return nil, err
	}
	c.Dictionary = display

	obs := observer{SlogObserver: session.SlogObserver{Logger: config.GetMessageLogger(), Dictionary: display}, client: c}
	cfg, err := s.FIXConfig(c, config.GetSlogLogger(), obs)
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

	initiator := session.NewInitiator(session.InitiatorOptions{
		Logger:       config.GetSlogLogger(),
		DialTimeout:  initiatorConfig.DialTimeoutOrDefault(),
		ReconnectMin: initiatorConfig.ReconnectMin,
		ReconnectMax: initiatorConfig.ReconnectMax,
	})

	if c.session, err = initiator.Add(cfg, &logonDialer{Dialer: dialer, client: c}); err != nil {
		cfg.Store.Close()
		return nil, err
	}

	runCtx, stop := context.WithCancel(context.WithoutCancel(ctx))
	c.stop = stop
	go c.pump()
	go func() {
		defer close(c.stopped)
		// Run only fails if the initiator already runs.
		_ = initiator.Run(runCtx)
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
		c.Close(false)
		switch {
		case errors.Is(err, errors.ErrFixLogonRefused):
			return nil, err
		case err != nil:
			return nil, fmt.Errorf("%w: %w", errors.ErrFixLogout, err)
		}
		return nil, errors.ErrFixLogout
	case <-time.After(config.Timeout(s.LogonTimeout)):
		err := c.timeoutError()
		c.Close(false)
		return nil, err
	case <-ctx.Done():
		c.Close(false)
		return nil, ctx.Err()
	}
}

// timeoutError is the error of a Connect that timed out: why the last dial
// failed, and why the session halted, if they did.
func (c *Client) timeoutError() error {
	c.mu.Lock()
	dialErr := c.dialErr
	c.mu.Unlock()

	err := errors.ErrConnectionTimeout
	if dialErr != nil {
		err = fmt.Errorf("%w: %w", err, dialErr)
	}
	if serr := c.session.Err(); serr != nil {
		err = fmt.Errorf("%w: %w", err, serr)
	}

	return err
}

// Send sends m, once checked against the session's dictionary (see
// fixmsg.Check): a message its FIX version does not define, or that lacks
// what it requires, is not sent. Its Receipt's Seq is the RefSeqNum (45) of
// the Rejects of it, which fail Await.
func (c *Client) Send(ctx context.Context, m session.Outgoing) (session.Receipt, error) {
	if err := c.Session.Check(m); err != nil {
		return session.Receipt{}, err
	}

	receipt, err := c.session.Send(ctx, m)
	if err != nil {
		return receipt, err
	}

	c.mu.Lock()
	c.sent[receipt.Seq()] = append(c.sent[receipt.Seq()], sentMessage{receipt: receipt, msgType: m.MsgType()})
	c.mu.Unlock()

	return receipt, nil
}

// BeginString returns the BeginString (8) of the session's messages.
func (c *Client) BeginString() string {
	return c.Session.BeginString
}

// Messages returns the application messages received, and the Rejects, in
// the order they were received. It is never closed: Done tells when the
// session ended.
func (c *Client) Messages() <-chan *fix.Message {
	return c.messages
}

// Done returns a channel closed once the session ended and stopped, its
// store closed: the counterparty logged it out or disconnected, or refused
// its Logon, or Close was called; Err then tells why. Unless Close ended
// it, the reader of Messages received every message received before.
func (c *Client) Done() <-chan struct{} {
	return c.done
}

// Err returns why the session ended, nil if it has not or Close ended it.
func (c *Client) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.reason
}

// Close ends the session, logging it out if logout is true, otherwise
// dropping its connection without a Logout, and waits for it to stop and
// its store to be closed. The messages the reader of Messages did not
// receive yet are discarded. Closing a closed Client only waits for Done.
func (c *Client) Close(logout bool) {
	c.mu.Lock()
	var conn transport.Conn
	if !c.closed {
		c.closed = true
		c.queue = nil
		close(c.closing)
		if !logout {
			conn = c.conn
		}
	}
	c.mu.Unlock()

	if conn != nil {
		// Our Logout, if the session still tries, fails to be written.
		conn.Close()
	}

	c.end(nil)
	<-c.done
}

// end records why the session ended, unless Close ended it, and stops it:
// an end is final for a Client. Every end goes through it.
func (c *Client) end(reason error) {
	c.mu.Lock()
	if c.reason == nil && !c.closed {
		c.reason = reason
	}
	c.ending = true
	c.mu.Unlock()

	c.stopOnce.Do(func() { go c.shutdown() })
}

// shutdown stops the initiator, then, once no callback of the session runs
// anymore, queues the Rejects still held, and closes ended.
func (c *Client) shutdown() {
	c.stop()
	<-c.stopped
	<-c.session.Done()

	c.mu.Lock()
	for _, r := range c.held {
		c.queueLocked(r.message)
	}
	c.held = nil
	c.mu.Unlock()

	close(c.ended)
}

// pump hands the messages queued to the reader of Messages, in order, and
// closes done once the session ended and they were all handed, or Close
// discarded them.
func (c *Client) pump() {
	defer close(c.done)

	for {
		m, ok := c.next()
		if !ok {
			return
		}

		select {
		case c.messages <- m:
		case <-c.closing:
			<-c.ended
			return
		}
	}
}

// next returns the next message queued, waiting for one, or false once the
// session ended and none is left.
func (c *Client) next() (*fix.Message, bool) {
	for {
		c.mu.Lock()
		if len(c.queue) > 0 {
			m := c.queue[0]
			c.queue[0] = nil
			c.queue = c.queue[1:]
			if c.room != nil && len(c.queue) < c.opts.Buffer {
				close(c.room)
				c.room = nil
			}
			c.mu.Unlock()
			return m, true
		}
		c.mu.Unlock()

		select {
		case <-c.wake:
		case <-c.ended:
			c.mu.Lock()
			empty := len(c.queue) == 0
			c.mu.Unlock()
			if empty {
				return nil, false
			}
		}
	}
}

// queueLocked queues m for pump, unless Close discarded the messages.
func (c *Client) queueLocked(m *fix.Message) {
	if c.closed {
		return
	}

	c.queue = append(c.queue, m)
	select {
	case c.wake <- struct{}{}:
	default:
	}
}

// push queues application message m for pump, waiting for room unless the
// client is closed.
func (c *Client) push(m *fix.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for len(c.queue) >= c.opts.Buffer && !c.closed {
		if c.room == nil {
			c.room = make(chan struct{})
		}
		room := c.room
		c.mu.Unlock()
		select {
		case <-room:
		case <-c.closing:
		}
		c.mu.Lock()
	}

	c.queueLocked(m)
}

// received records frame f, received while the session expected next, before
// the session handles it, on the reader goroutine: the frames before it were
// handled. The application messages and GapFills from next on are recorded:
// those before are duplicates the session ignores. Of those received ahead
// of a gap with the same MsgSeqNum, the session keeps the first.
func (c *Client) received(f session.Frame, next uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// A reset (a Logon with ResetSeqNumFlag) waits for the messages
	// received before it to be delivered, and drops those held ahead of a
	// gap: what is left of earlier generations never will be.
	n := 0
	for n < len(c.inbound) && c.inbound[n].gen < f.Generation {
		n++
	}
	c.inbound = slices.Delete(c.inbound, 0, n)

	// The GapFills handled since, unless rejected, skipped the messages up
	// to their NewSeqNo.
	c.gapFills = slices.DeleteFunc(c.gapFills, func(g gapFill) bool {
		if g.gen == f.Generation && g.seq >= next {
			return false
		}
		c.inbound = slices.DeleteFunc(c.inbound, func(p pos) bool {
			return p.gen == g.gen && p.seq >= g.seq && p.seq < g.end
		})
		return true
	})
	c.releaseLocked()

	msgType := frameField(f.Raw, fixt11.TagMsgType)
	if f.Seq < next || msgType == nil {
		return
	}
	p := pos{gen: f.Generation, seq: f.Seq}
	i, found := slices.BinarySearchFunc(c.inbound, p, pos.compare)
	if found || slices.ContainsFunc(c.gapFills, func(g gapFill) bool { return g.pos == p }) {
		return
	}
	switch {
	case !isAdmin(msgType):
		c.inbound = slices.Insert(c.inbound, i, p)
	case string(msgType) == string(fixt11.MsgTypeSequenceReset) && string(frameField(f.Raw, fixt11.TagGapFillFlag)) == "Y":
		if end, err := fix.ParseInt(frameField(f.Raw, fixt11.TagNewSeqNo)); err == nil && end > 0 && uint64(end) > f.Seq+1 {
			c.gapFills = append(c.gapFills, gapFill{pos: p, end: uint64(end)})
		}
	}
}

// rejected records that the session rejected message p: a GapFill's
// NewSeqNo is not applied; a message counted as received is never
// delivered.
func (c *Client) rejected(p pos, counted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.gapFills = slices.DeleteFunc(c.gapFills, func(g gapFill) bool { return g.pos == p })
	if i, found := slices.BinarySearchFunc(c.inbound, p, pos.compare); found && counted {
		c.inbound = slices.Delete(c.inbound, i, i+1)
	}
	c.releaseLocked()
}

// passLocked records that the application messages before p were
// delivered, or never will be, and releases the Rejects they held.
func (c *Client) passLocked(p pos) {
	i, _ := slices.BinarySearchFunc(c.inbound, p, pos.compare)
	c.inbound = slices.Delete(c.inbound, 0, i)
	c.releaseLocked()
}

// releaseLocked queues the Rejects held that no application message
// received before them waits to be delivered.
func (c *Client) releaseLocked() {
	n := 0
	for _, r := range c.held {
		if len(c.inbound) > 0 && c.inbound[0].compare(r.pos) < 0 {
			break
		}
		c.queueLocked(r.message)
		n++
	}
	c.held = slices.Delete(c.held, 0, n)
}

// refused ends a Connect whose Logon failed, unless the session logged on.
func (c *Client) refused(reason string) {
	c.mu.Lock()
	accepted := c.accepted
	c.mu.Unlock()

	if !accepted {
		c.end(fmt.Errorf("%w: %s", errors.ErrFixLogonRefused, reason))
	}
}

// msgTypeName names msgType after the session's dictionary.
func (c *Client) msgTypeName(msgType string) string {
	if c.Dictionary != nil {
		if def, ok := c.Dictionary.Messages[msgType]; ok {
			return fmt.Sprintf("%s(%s)", def.Name, msgType)
		}
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

	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.loggedOn:
	default:
		close(c.loggedOn)
	}
}

// OnLogout implements session.Handler.
func (c *Client) OnLogout(_ context.Context, s *session.Session, reason error) {
	c.logger.Debug().Msgf("Logout: %s: %v", s.ID(), reason)
	c.end(reason)
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

// FromAdmin implements session.Handler: it holds the Rejects until the
// application messages received before them are delivered, and ends a
// Connect whose Logon was refused, which OnLogout is not told of. It runs
// on the reader goroutine.
func (c *Client) FromAdmin(_ context.Context, s *session.Session, m *fix.Message, in session.Inbound) {
	c.mu.Lock()
	accepted := c.accepted
	c.mu.Unlock()

	switch string(m.MsgType()) {
	case string(fixt11.MsgTypeLogon):
		// FromAdmin is told of a Logon the session refused too, logging
		// out (an invalid NextExpectedMsgSeqNum, say): only one that
		// established the session is accepted.
		switch s.State() {
		case session.StateLoggedOn, session.StateStalled, session.StateResettingSeqs:
			c.mu.Lock()
			c.accepted = true
			c.mu.Unlock()
		default:
			c.refused("the counterparty's Logon was rejected")
		}

	case string(fixt11.MsgTypeReject):
		if !accepted {
			c.refused(withText("the counterparty rejected it", m))
			return
		}

		c.mu.Lock()
		c.held = append(c.held, heldReject{pos: pos{gen: in.Generation, seq: in.Seq}, message: m.Clone()})
		c.releaseLocked()
		c.mu.Unlock()

	case string(fixt11.MsgTypeLogout):
		if !accepted {
			c.refused(withText("the counterparty logged out", m))
		}
	}
}

// withText appends the Text (58) of m to reason, if any.
func withText(reason string, m *fix.Message) string {
	if text := m.Get(fixt11.TagText); len(text) > 0 {
		return reason + ": " + string(text)
	}

	return reason
}

// FromApp implements session.Handler.
func (c *Client) FromApp(_ context.Context, _ *session.Session, m *fix.Message, in session.Inbound) error {
	// Delivered in sequence: the messages before it were.
	c.mu.Lock()
	c.passLocked(pos{gen: in.Generation, seq: in.Seq})
	c.mu.Unlock()

	msgType := string(m.MsgType())
	if len(c.opts.MsgTypes) > 0 && !slices.Contains(c.opts.MsgTypes, msgType) && msgType != string(fix50sp2.MsgTypeBusinessmessagereject) {
		c.logger.Info().Msgf("Received unexpected message type: %s", c.msgTypeName(msgType))
	} else {
		c.push(m.Clone())
	}

	c.mu.Lock()
	c.passLocked(pos{gen: in.Generation, seq: in.Seq + 1})
	c.mu.Unlock()

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

// observer logs the frames, and tells the Client of the application
// messages and GapFills received, in order, for the Rejects to wait for the
// messages received before them.
type observer struct {
	session.SlogObserver
	client *Client
}

// OnInbound implements session.MessageObserver.
func (o observer) OnInbound(s *session.Session, f session.Frame) {
	o.SlogObserver.OnInbound(s, f)
	o.client.received(f, s.NextTargetSeqNum())
}

// OnRejected implements session.RejectObserver.
func (o observer) OnRejected(_ *session.Session, r session.Rejected) {
	if r.Seq > 0 {
		o.client.rejected(pos{gen: r.Generation, seq: r.Seq}, r.Counted)
	}
}

var _ session.RejectObserver = observer{}

// frameField returns the value of the first field tag of frame, nil if it
// has none.
func frameField(frame []byte, tag int) []byte {
	_, rest, ok := bytes.Cut(frame, append(strconv.AppendInt([]byte{0x01}, int64(tag), 10), '='))
	if !ok {
		return nil
	}
	value, _, _ := bytes.Cut(rest, []byte{0x01})

	return value
}

// isAdmin reports whether msgType is that of a session-level message.
func isAdmin(msgType []byte) bool {
	return len(msgType) == 1 && bytes.IndexByte([]byte("012345A"), msgType[0]) >= 0
}

// logonDialer dials the counterparty for a Client, remembering why it
// failed. A connection the counterparty closes before the session logged on
// had its Logon refused: the Client ends. So does one that ended before
// for a reason of ours, rather than send the credentials again.
type logonDialer struct {
	transport.Dialer
	client *Client
}

func (d *logonDialer) Dial(ctx context.Context) (transport.Conn, error) {
	c := d.client

	c.mu.Lock()
	dialed, accepted, ending := c.dialed, c.accepted, c.ending
	c.mu.Unlock()

	switch {
	case ending:
		return nil, stderrors.New("session ended")
	case dialed && !accepted:
		c.refused("the connection closed before the logon")
		return nil, stderrors.New("logon refused")
	}

	conn, err := d.Dialer.Dial(ctx)

	c.mu.Lock()
	defer c.mu.Unlock()
	if err != nil {
		c.dialErr = err
		return nil, err
	}
	c.dialed, c.dialErr = true, nil
	c.conn = &clientConn{Conn: conn, client: c}

	return c.conn, nil
}

// clientConn is a connection of a Client, which tells it of the
// counterparty closing it.
type clientConn struct {
	transport.Conn
	client *Client
}

func (cc *clientConn) Read(p []byte) (int, error) {
	n, err := cc.Conn.Read(p)
	if err != nil && (stderrors.Is(err, io.EOF) || stderrors.Is(err, syscall.ECONNRESET)) {
		cc.client.refused("the counterparty closed the connection")
	}

	return n, err
}
