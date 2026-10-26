// Package options holds the flags and helpers the order and quote commands
// share.
package options

import (
	"fmt"
	"strings"

	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/fix50sp2"
)

// Must panics if err is not nil: for the errors of setting up flags at
// init, which are programming errors.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Lookup returns the value of m named v, regardless of case.
func Lookup[T ~string](m map[string]T, v string) (T, bool) {
	value, ok := m[strings.ToUpper(v)]
	return value, ok
}

// WithText wraps err with the Text (58) of m, if any.
func WithText(err error, m *fix.Message) error {
	if text := m.Get(fix50sp2.TagText); len(text) > 0 {
		return fmt.Errorf("%w: %s", err, text)
	}

	return err
}
