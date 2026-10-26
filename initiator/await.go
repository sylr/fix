package initiator

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/fixt11"
)

// TimeoutMode tells what AwaitOptions.Timeout bounds.
type TimeoutMode int

const (
	// TimeoutTotal bounds the whole wait, until the handler reports done.
	TimeoutTotal TimeoutMode = iota
	// TimeoutIdle bounds each wait for a response: the timeout restarts
	// with every message handled.
	TimeoutIdle
	// TimeoutFirst bounds the wait for the first response only: the
	// messages that follow are handled until the handler reports done or
	// ctx is done, as for a subscription.
	TimeoutFirst
)

// AwaitOptions configure Await.
type AwaitOptions struct {
	// What names what is awaited in the timeout error, as in "execution
	// report".
	What string
	// Timeout bounds the wait as Mode tells, none if not positive.
	Timeout time.Duration
	Mode    TimeoutMode
	// Tick, if not nil, makes Await call OnTick at each of its values,
	// between two messages: a time.Ticker's C, say, for periodic updates.
	// Ticks and messages are served alike, neither starving the other.
	// OnTick may Send: the Rejects of what it sends fail Await too. Its
	// error ends Await.
	Tick   <-chan time.Time
	OnTick func() error
}

// ErrSkip is what an Await handler returns for a message that does not
// answer what it awaits, as a response to another request: Await goes on
// waiting, its timeout unaffected.
var ErrSkip = errors.New("message skipped")

// RejectError is a Reject (35=3) or a BusinessMessageReject (35=j) of a
// message sent with Send. It wraps errors.ErrFixRejected.
type RejectError struct {
	// RejectInfo holds the reject's MsgType, its SessionRejectReason (373)
	// or BusinessRejectReason (380), RefTagID (371) and Text (58).
	session.RejectInfo
	// RefSeqNum and RefMsgType are the MsgSeqNum and MsgType of the
	// message rejected; RefSeqNum is 0 if the reject did not tell.
	RefSeqNum  uint64
	RefMsgType string

	description string
}

func (e *RejectError) Error() string {
	return fmt.Sprintf("%s: %s", errors.ErrFixRejected, e.description)
}

func (e *RejectError) Unwrap() error {
	return errors.ErrFixRejected
}

// Await hands handle the messages received, in order, until it reports
// done, then returns nil. Use it after sending a request:
//
//	if _, err := c.Send(ctx, request); err != nil {
//		return err
//	}
//	return c.Await(ctx, initiator.AwaitOptions{What: "execution report", Timeout: timeout},
//		func(m *fix.Message) (bool, error) {
//			utils.WriteMessageBodyAsTable(os.Stdout, m, c.Dictionary)
//			return isFinal(m), nil
//		})
//
// with ctx the command's, done on SIGINT and SIGTERM:
//
//	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
//	defer stop()
//	context.AfterFunc(ctx, stop) // a second signal kills
//
// The messages are those Messages delivers, the ones Options.MsgTypes
// selects and the rejects. A Reject (35=3) or BusinessMessageReject (35=j)
// of a message sent with Send, matched by its RefSeqNum (45) in the sequence
// number generation it was received in or, for a BusinessMessageReject
// without one, by its RefMsgType (372), fails Await
// with a *RejectError; the others are logged and skipped, never handed to
// handle. An error from handle or OnTick ends Await with it.
//
// Await returns an error wrapping errors.ErrResponseTimeout once
// opts.Timeout elapsed before done, as opts.Mode tells; one wrapping
// errors.ErrFixLogout and Err once the session ended, after every message
// received before was handled; and nil once ctx is done, the caller
// interrupting it: ctx.Err tells. Messages and Await must not be read from
// elsewhere meanwhile; Await may be called again afterwards.
func (c *Client) Await(ctx context.Context, opts AwaitOptions, handle func(*fix.Message) (done bool, err error)) error {
	var timeout <-chan time.Time
	var timer *time.Timer
	if opts.Timeout > 0 {
		timer = time.NewTimer(opts.Timeout)
		defer timer.Stop()
		timeout = timer.C
	}

	handled := 0
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-c.done:
			if err := c.Err(); err != nil {
				return fmt.Errorf("%w: %w", errors.ErrFixLogout, err)
			}
			return fmt.Errorf("%w: client closed", errors.ErrFix)

		case <-timeout:
			return fmt.Errorf("%w: %s (%d received)", errors.ErrResponseTimeout, opts.What, handled)

		case <-opts.Tick:
			if err := opts.OnTick(); err != nil {
				return err
			}

		case m := <-c.Messages():
			if isReject(m) {
				if err := c.rejectError(m); err != nil {
					return err
				}
				c.logger.Warn().Msgf("%s of a message not sent ignored: %s", c.msgTypeName(string(m.MsgType())), m.Get(fixt11.TagText))
				continue
			}

			done, err := handle(m)
			if errors.Is(err, ErrSkip) {
				continue
			}
			if err != nil || done {
				return err
			}

			handled++
			switch {
			case timer == nil:
			case opts.Mode == TimeoutIdle:
				timer.Reset(opts.Timeout)
			case opts.Mode == TimeoutFirst:
				timer.Stop()
				timeout = nil
			}
		}
	}
}

// isReject reports whether m is a Reject or a BusinessMessageReject.
func isReject(m *fix.Message) bool {
	switch string(m.MsgType()) {
	case string(fixt11.MsgTypeReject), string(fix50sp2.MsgTypeBusinessmessagereject):
		return true
	}

	return false
}

// rejectError returns the *RejectError of reject m if it rejects a message
// sent with Send, nil otherwise.
func (c *Client) rejectError(m *fix.Message) error {
	e := &RejectError{RejectInfo: session.RejectInfo{MsgType: string(m.MsgType()), Text: string(m.Get(fixt11.TagText))}}
	e.RefMsgType = string(m.Get(fixt11.TagRefMsgType))

	if v, err := fix.ParseInt(m.Get(fixt11.TagRefSeqNum)); err == nil && v > 0 {
		e.RefSeqNum = uint64(v)
	}

	var msgType string
	var sent bool
	switch {
	case e.RefSeqNum > 0:
		msgType, sent = c.rejectedSent(e.RefSeqNum)
	case e.MsgType == string(fix50sp2.MsgTypeBusinessmessagereject):
		c.mu.Lock()
		for _, messages := range c.sent {
			if sent = slices.ContainsFunc(messages, func(s sentMessage) bool { return s.msgType == e.RefMsgType }); sent {
				break
			}
		}
		c.mu.Unlock()
	}

	if !sent {
		return nil
	}
	if e.RefMsgType == "" {
		e.RefMsgType = msgType
	}

	reasonTag := fixt11.TagSessionRejectReason
	if e.MsgType == string(fix50sp2.MsgTypeBusinessmessagereject) {
		reasonTag = fix50sp2.TagBusinessRejectReason
	}
	if v, err := fix.ParseInt(m.Get(reasonTag)); err == nil {
		e.Reason = v
	}
	if v, err := fix.ParseInt(m.Get(fixt11.TagRefTagID)); err == nil {
		e.RefTagID = v
	}

	e.description = c.msgTypeName(e.MsgType) + " of " + c.msgTypeName(e.RefMsgType)
	if e.RefSeqNum > 0 {
		e.description += " " + strconv.FormatUint(e.RefSeqNum, 10)
	}
	if v := m.Get(reasonTag); v != nil {
		e.description += ": reason " + dict.EnumName(c.Dictionary, reasonTag, string(v))
	}
	if e.RefTagID > 0 {
		e.description += ", tag " + c.tagName(e.RefTagID)
	}
	if e.Text != "" {
		e.description += ": " + e.Text
	}

	return e
}

// rejectedSent returns the MsgType of the message sent with Send of MsgSeqNum
// seq the session saw rejected, and forgets it: its receipt tells whether
// the Reject is of it, or of a message of another generation.
func (c *Client) rejectedSent(seq uint64) (string, bool) {
	c.mu.Lock()
	messages := slices.Clone(c.sent[seq])
	c.mu.Unlock()

	for _, m := range messages {
		// Rejection takes the session's lock: not under c.mu.
		if _, ok := m.receipt.Rejection(); !ok {
			continue
		}

		c.mu.Lock()
		c.sent[seq] = slices.DeleteFunc(c.sent[seq], func(s sentMessage) bool { return s == m })
		if len(c.sent[seq]) == 0 {
			delete(c.sent, seq)
		}
		c.mu.Unlock()

		return m.msgType, true
	}

	return "", false
}

// tagName names tag after the session's dictionary.
func (c *Client) tagName(tag int) string {
	if c.Dictionary != nil {
		if def := c.Dictionary.Fields[tag]; def != nil {
			return fmt.Sprintf("%s(%d)", def.Name, tag)
		}
	}

	return strconv.Itoa(tag)
}
