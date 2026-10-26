package options

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/spf13/pflag"

	"sylr.dev/fix/v2/errors"
)

// Decimal is a pflag.Value holding an exact decimal, as a price or a
// quantity: 1.08345 is sent as 1.08345, never rounded through a float. NaN
// and Inf are refused. Its zero value is unset; setting it to "" unsets it.
type Decimal struct {
	value decimal.Decimal
	set   bool
}

var _ pflag.Value = (*Decimal)(nil)

// Set implements pflag.Value.
func (d *Decimal) Set(s string) error {
	if s == "" {
		*d = Decimal{}
		return nil
	}

	v, err := parseDecimal(s)
	if err != nil {
		return err
	}
	*d = Decimal{value: v, set: true}

	return nil
}

// String implements pflag.Value: the decimal as sent, "" if unset.
func (d *Decimal) String() string {
	if !d.set {
		return ""
	}

	return d.value.String()
}

// Type implements pflag.Value.
func (d *Decimal) Type() string {
	return "decimal"
}

// IsSet reports whether d was given a value.
func (d Decimal) IsSet() bool {
	return d.set
}

// Value returns the decimal, 0 if unset.
func (d Decimal) Value() decimal.Decimal {
	return d.value
}

// Decimals is a pflag.SliceValue of exact decimals, as Decimal.
type Decimals []decimal.Decimal

var _ pflag.SliceValue = (*Decimals)(nil)

// Set implements pflag.Value: it appends the comma separated decimals of s.
func (d *Decimals) Set(s string) error {
	return d.Append(s)
}

// String implements pflag.Value.
func (d *Decimals) String() string {
	return "[" + strings.Join(d.GetSlice(), ",") + "]"
}

// Type implements pflag.Value.
func (d *Decimals) Type() string {
	return "decimals"
}

// Append implements pflag.SliceValue: it appends the comma separated
// decimals of s.
func (d *Decimals) Append(s string) error {
	for v := range strings.SplitSeq(s, ",") {
		parsed, err := parseDecimal(v)
		if err != nil {
			return err
		}
		*d = append(*d, parsed)
	}

	return nil
}

// Replace implements pflag.SliceValue.
func (d *Decimals) Replace(values []string) error {
	parsed := make(Decimals, 0, len(values))
	for _, v := range values {
		if err := parsed.Append(v); err != nil {
			return err
		}
	}
	*d = parsed

	return nil
}

// GetSlice implements pflag.SliceValue.
func (d *Decimals) GetSlice() []string {
	values := make([]string, len(*d))
	for i, v := range *d {
		values[i] = v.String()
	}

	return values
}

// parseDecimal parses s as a decimal: digits, with an optional sign, decimal
// point and exponent.
func parseDecimal(s string) (decimal.Decimal, error) {
	v, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("%w: invalid decimal %q", errors.ErrOptions, s)
	}

	return v, nil
}

// Positive returns an error unless v, the value of flag name, is above 0.
func Positive(name string, v decimal.Decimal) error {
	if !v.IsPositive() {
		return fmt.Errorf("%w: --%s must be positive, not %s", errors.ErrOptions, name, v)
	}

	return nil
}
