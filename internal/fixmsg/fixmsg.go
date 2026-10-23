// Package fixmsg checks the messages fix builds against a session's
// dictionary before they are sent.
//
// Commands build every FIX version's messages with package fix50sp2's
// builders, whose tags and values FIX 4.3 and 4.4 share: what a version lacks
// or requires on top is caught here, rather than by the counterparty.
package fixmsg

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
)

// standard header and trailer tags.
const (
	tagBeginString  = 8
	tagBodyLength   = 9
	tagCheckSum     = 10
	tagMsgSeqNum    = 34
	tagMsgType      = 35
	tagSenderCompID = 49
	tagSendingTime  = 52
	tagTargetCompID = 56
)

// ErrNotSupported is the error Check wraps.
var ErrNotSupported = errors.Fix

// Message returns m as a whole message of BeginString beginString, parsed
// with d's profile if d is not nil. Its header holds the fields a session
// writes, CompIDs, MsgSeqNum and SendingTime, with placeholder values, so
// that d can validate it.
func Message(m session.Outgoing, beginString string, d *dictionary.Dict) (*fix.Message, error) {
	body := fix.AppendTagValue(nil, tagMsgType, m.MsgType())
	body = fix.AppendTagValue(body, tagSenderCompID, "SENDER")
	body = fix.AppendTagValue(body, tagTargetCompID, "TARGET")
	body = fix.AppendTagInt(body, tagMsgSeqNum, 1)
	body = fix.AppendTagValue(body, tagSendingTime, time.Now().UTC().Format("20060102-15:04:05.000"))

	body, err := m.AppendBody(body, profile(d), func(int) bool { return false })
	if err != nil {
		return nil, err
	}

	return parse(beginString, body, d)
}

// Check reports what a counterparty validating msg with d, running checks,
// would reject it for, but its enumerated values: a counterparty may accept
// values beyond its dictionary's.
func Check(msg *fix.Message, d *dictionary.Dict, checks dictionary.Checks) error {
	msgType := string(msg.MsgType())
	version := versionOf(d)

	// What a FIX version lacks, told plainly.
	mdef := d.Messages[msgType]
	if mdef == nil {
		return fmt.Errorf("%w: MsgType %s is not part of %s", ErrNotSupported, msgType, version)
	}

	if checks&dictionary.CheckUnknownTags != 0 {
		for _, f := range msg.Fields() {
			if !d.Defines(msgType, f.Tag) {
				return fmt.Errorf("%w: %s(%d) is not part of %s's %s", ErrNotSupported, fieldName(d, f.Tag), f.Tag, version, mdef.Name)
			}
		}
	}

	for _, f := range mdef.Fields {
		if f.Required && msg.Get(f.Tag) == nil {
			return fmt.Errorf("%w: %s's %s requires %s(%d)", ErrNotSupported, version, mdef.Name, fieldName(d, f.Tag), f.Tag)
		}
	}

	// The rest as the counterparty validates it: group entries, required
	// fields at every level, formats...
	valid, err := withValidEnums(msg, d)
	if err != nil {
		return err
	}

	if err := d.ValidateWith(valid, checks); err != nil {
		return fmt.Errorf("%w: %s's %s: %w", ErrNotSupported, version, mdef.Name, err)
	}

	return nil
}

// withValidEnums returns msg, its values outside their field's enumeration
// replaced with one inside, so that validating it judges everything else.
func withValidEnums(msg *fix.Message, d *dictionary.Dict) (*fix.Message, error) {
	var body []byte
	var beginString string
	changed := false

	for _, f := range msg.Fields() {
		switch f.Tag {
		case tagBeginString:
			beginString = string(f.Value)
			continue
		case tagBodyLength, tagCheckSum:
			continue
		}

		value := f.Value
		if def := d.Fields[f.Tag]; def != nil && len(def.Enums) > 0 {
			if v, ok := validEnum(def, string(value)); !ok {
				value, changed = []byte(v), true
			}
		}
		body = fix.AppendTagValue(body, f.Tag, value)
	}

	if !changed {
		return msg, nil
	}

	return parse(beginString, body, d)
}

// validEnum returns value if def's enumeration holds it, every value of it
// for fields holding several, else one it holds, and false.
func validEnum(def *dictionary.FieldDef, value string) (string, bool) {
	first := slices.Sorted(maps.Keys(def.Enums))[0]

	if !strings.HasPrefix(strings.ToUpper(def.Type), "MULTIPLE") {
		if _, ok := def.Enums[value]; ok {
			return value, true
		}
		return first, false
	}

	tokens := strings.Fields(value)
	ok := len(tokens) > 0
	for i, t := range tokens {
		if _, in := def.Enums[t]; !in {
			tokens[i], ok = first, false
		}
	}
	if len(tokens) == 0 {
		return first, false
	}

	return strings.Join(tokens, " "), ok
}

// parse returns the message of BeginString beginString and body, which
// starts with its MsgType, parsed with d's profile.
func parse(beginString string, body []byte, d *dictionary.Dict) (*fix.Message, error) {
	frame := fix.AppendTagValue(nil, tagBeginString, beginString)
	frame = fix.AppendTagValue(frame, tagBodyLength, strconv.Itoa(len(body)))
	frame = fix.AppendTrailer(append(frame, body...), 0)

	var msg fix.Message
	if err := fix.ParseInto(&msg, frame, profile(d)); err != nil {
		return nil, err
	}

	return &msg, nil
}

func profile(d *dictionary.Dict) *fix.Profile {
	if d == nil {
		return nil
	}

	return d.Profile()
}

// versionOf names d's FIX version, as FIX.4.4 or FIX.5.0SP2.
func versionOf(d *dictionary.Dict) string {
	version := d.Type + "." + d.Major + "." + d.Minor
	if d.ServicePack != "" && d.ServicePack != "0" {
		version += "SP" + d.ServicePack
	}

	return version
}

// fieldName names tag as d does, or as FIX 5.0 SP2 does a tag d lacks.
func fieldName(d *dictionary.Dict, tag int) string {
	for _, dict := range []*dictionary.Dict{d, fix50sp2.Dictionary()} {
		if f := dict.Fields[tag]; f != nil {
			return f.Name
		}
	}

	return "tag"
}
