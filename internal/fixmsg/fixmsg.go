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

// ErrNotSupported is the error Check wraps for a message its counterparty
// would reject.
var ErrNotSupported = fmt.Errorf("%w: not supported", errors.ErrFix)

// userDefinedTagMin is the first user-defined tag, which
// dictionary.CheckUserDefinedTags checks rather than CheckUnknownTags.
const userDefinedTagMin = 5000

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
	if d == nil {
		return fmt.Errorf("%w: no dictionary to check messages with", errors.ErrFix)
	}

	msgType := string(msg.MsgType())
	dict := dictName(d)

	// What the dictionary lacks, told plainly.
	mdef := d.Messages[msgType]
	if mdef == nil {
		return fmt.Errorf("%w: %s does not define MsgType %s", ErrNotSupported, dict, msgType)
	}

	for _, f := range msg.Fields() {
		if checks&unknownCheck(f.Tag) != 0 && !d.Defines(msgType, f.Tag) {
			return fmt.Errorf("%w: %s does not define %s(%d) in %s", ErrNotSupported, dict, fieldName(d, f.Tag), f.Tag, mdef.Name)
		}
	}

	for _, f := range mdef.Fields {
		if f.Required && msg.Get(f.Tag) == nil {
			return fmt.Errorf("%w: %s requires %s(%d) in %s", ErrNotSupported, dict, fieldName(d, f.Tag), f.Tag, mdef.Name)
		}
	}

	// The rest as the counterparty validates it: group entries, required
	// fields at every level, formats...
	valid, err := withValidEnums(msg, d)
	if err != nil {
		return fmt.Errorf("%w: %s: replacing enumerated values: %w", errors.ErrFix, mdef.Name, err)
	}

	if err := d.ValidateWith(valid, checks); err != nil {
		return fmt.Errorf("%w: %s rejects %s: %w", ErrNotSupported, dict, mdef.Name, err)
	}

	return nil
}

// unknownCheck returns the check rejecting tag where the message does not
// define it, as dictionary.Dict.ValidateWith chooses it.
func unknownCheck(tag int) dictionary.Checks {
	if tag >= userDefinedTagMin {
		return dictionary.CheckUserDefinedTags
	}

	return dictionary.CheckUnknownTags
}

// withValidEnums returns msg, its non-empty values outside their field's
// enumeration replaced with one inside, so that validating it judges
// everything else. Empty values are left for validation to report. MsgType
// and NUMINGROUP fields are left too: they give the message its structure.
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
		if def := d.Fields[f.Tag]; f.Tag != tagMsgType && def != nil && len(def.Enums) > 0 && def.Type != "NUMINGROUP" {
			if v, ok := validEnum(def, string(value)); ok {
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

// validEnum returns value, its non-empty values outside def's enumeration
// replaced with one inside, and whether it replaced any. A multiple-value
// field's values are separated by single spaces, as fix2 splits them.
func validEnum(def *dictionary.FieldDef, value string) (string, bool) {
	valid := func(v string) bool {
		_, ok := def.Enums[v]
		return ok || v == ""
	}

	if !isMultipleValue(def.Type) {
		if valid(value) {
			return value, false
		}
		return slices.Min(slices.Collect(maps.Keys(def.Enums))), true
	}

	values := strings.Split(value, " ")
	changed := false
	for i, v := range values {
		if !valid(v) {
			values[i], changed = slices.Min(slices.Collect(maps.Keys(def.Enums))), true
		}
	}

	return strings.Join(values, " "), changed
}

// isMultipleValue reports whether a field of type fixType holds values
// separated by spaces, as fix2's dictionary validation tells them.
func isMultipleValue(fixType string) bool {
	switch fixType {
	case "MULTIPLEVALUESTRING", "MULTIPLECHARVALUE", "MULTIPLESTRINGVALUE":
		return true
	}

	return false
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

// dictName names d in errors: the session's dictionary, rather than its
// FIX version, may lack what the version has, as fix's built-in FIX 5.0 SP2
// one does what fix does not use.
func dictName(d *dictionary.Dict) string {
	return "the session's " + versionOf(d) + " dictionary"
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
