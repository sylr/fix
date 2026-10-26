package options_test

import (
	"testing"
	"time"

	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/cli/options"
	"sylr.dev/fix/v2/cli/options/optionstest"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

func TestDecimal(t *testing.T) {
	for _, tc := range []struct {
		in, want string
		err      bool
	}{
		{in: "1.08345", want: "1.08345"},
		{in: "10", want: "10"},
		{in: "1.50", want: "1.5"},
		{in: "-2.5", want: "-2.5"},
		{in: "1e3", want: "1000"},
		{in: "", want: ""},
		{in: "NaN", err: true},
		{in: "Inf", err: true},
		{in: "-Inf", err: true},
		{in: "1,5", err: true},
		{in: "abc", err: true},
	} {
		var d options.Decimal
		err := d.Set(tc.in)
		if tc.err {
			if !errors.Is(err, errors.ErrOptions) {
				t.Errorf("Set(%q) = %v, want an options error", tc.in, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("Set(%q) = %v", tc.in, err)
			continue
		}
		if got := d.String(); got != tc.want {
			t.Errorf("Set(%q): String() = %q, want %q", tc.in, got, tc.want)
		}
		if d.IsSet() != (tc.in != "") {
			t.Errorf("Set(%q): IsSet() = %v", tc.in, d.IsSet())
		}
	}
}

func TestDecimals(t *testing.T) {
	var d options.Decimals
	if err := d.Set("1.1,2.25"); err != nil {
		t.Fatal(err)
	}
	if err := d.Set("3"); err != nil {
		t.Fatal(err)
	}
	if got := d.String(); got != "[1.1,2.25,3]" {
		t.Errorf("String() = %q", got)
	}
	if err := d.Set("4,NaN"); !errors.Is(err, errors.ErrOptions) {
		t.Errorf("Set(4,NaN) = %v, want an options error", err)
	}
	if err := d.Replace(nil); err != nil || len(d) != 0 {
		t.Errorf("Replace(nil) = %v, %v", err, d)
	}
}

func TestPositive(t *testing.T) {
	for in, ok := range map[string]bool{"1": true, "0.0001": true, "0": false, "-1": false} {
		var d options.Decimal
		if err := d.Set(in); err != nil {
			t.Fatal(err)
		}
		if err := options.Positive("price", d.Value()); (err == nil) != ok {
			t.Errorf("Positive(%s) = %v", in, err)
		}
	}
}

func TestCheckPrice(t *testing.T) {
	price := func(s string) options.Decimal {
		var d options.Decimal
		if err := d.Set(s); err != nil {
			t.Fatal(err)
		}
		return d
	}

	for _, tc := range []struct {
		ordType fix50sp2.OrdType
		price   string
		want    error
	}{
		{fix50sp2.OrdTypeMarket, "", nil},
		{fix50sp2.OrdTypeMarket, "1", errors.ErrOptionsInvalidMarketPrice},
		{fix50sp2.OrdTypeLimit, "", errors.ErrOptionsNoPriceGiven},
		{fix50sp2.OrdTypeLimit, "0", nil},
		{fix50sp2.OrdTypeLimit, "-0.5", nil},
		{fix50sp2.OrdTypeLimit, "1.08345", nil},
	} {
		err := options.CheckPrice(tc.ordType, price(tc.price))
		if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
			t.Errorf("CheckPrice(%s, %q) = %v, want %v", tc.ordType, tc.price, err, tc.want)
		}
	}
}

func TestPartyIDOptionsValidate(t *testing.T) {
	cmd := &cobra.Command{}
	o := options.NewPartyIDOptions(cmd)

	party := []string{"party-id", "P1", "party-id-source", "proprietary", "party-role", "executing_firm"}
	for _, tc := range []struct {
		name  string
		flags []string
		want  error
	}{
		{"none", nil, nil},
		{"party", party, nil},
		{"sub ids", append(party[:6:6], "party-sub-ids", "A B", "party-sub-id-types", "firm person"), nil},
		{"sub id types without sub ids", append(party[:6:6], "party-sub-id-types", "firm"), errors.ErrOptions},
		{"sub id types count", append(party[:6:6], "party-sub-ids", "A B", "party-sub-id-types", "firm"), errors.ErrOptionsInconsistentValues},
		{"missing role", party[:4], errors.ErrOptionsInconsistentValues},
		{"unknown role", []string{"party-id", "P1", "party-id-source", "proprietary", "party-role", "nope"}, errors.ErrOptionOrderRoleUnknown},
		{"unknown sub id type", append(party[:6:6], "party-sub-ids", "A", "party-sub-id-types", "nope"), errors.ErrOptionPartySubIDTypeUnknown},
	} {
		t.Run(tc.name, func(t *testing.T) {
			optionstest.Set(t, cmd, tc.flags...)
			err := o.Validate()
			if (tc.want == nil) != (err == nil) || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("Validate() = %v, want %v", err, tc.want)
			}
		})
	}
}

// TestPartyEntries checks the sub ids and their types are sent, and the
// deprecated --copy-credentials-fom-config still works.
func TestPartyEntries(t *testing.T) {
	cmd := &cobra.Command{}
	o := options.NewPartyIDOptions(cmd)
	optionstest.Set(t, cmd,
		"party-id", "P1", "party-id-source", "proprietary", "party-role", "executing_firm",
		"party-sub-ids", "A B", "party-sub-id-types", "firm person",
		"copy-credentials-fom-config", "true")

	s := fixmsgtest.Session("FIX.4.4")
	order := fix50sp2.BuildNewOrderSingle(s.BeginString).
		SetClOrdID("ORD-1").
		SetSide(fix50sp2.SideBuy).
		SetTransactTime(time.Now()).
		SetOrdType(fix50sp2.OrdTypeMarket).
		SetSymbol("EURUSD").
		SetOrderQty(1)
	for _, entry := range o.PartyEntries(s) {
		order.AddNoPartyIDsEntry(entry)
	}

	msg := fixmsgtest.Validate(t, s, order)
	var got []string
	for _, f := range msg.Fields() {
		switch f.Tag {
		case fix50sp2.TagPartyID, fix50sp2.TagPartySubID, fix50sp2.TagPartySubIDType:
			got = append(got, string(f.Value))
		}
	}
	want := []string{"P1", "A", string(fix50sp2.PartySubIDTypeFirm), "B", string(fix50sp2.PartySubIDTypePerson), "user"}
	if len(got) != len(want) {
		t.Fatalf("parties = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("parties = %q, want %q", got, want)
		}
	}

	if err := s.Check(order); err != nil {
		t.Errorf("Check() = %v", err)
	}
}

// TestCheck checks Check refuses what the session's FIX version can't
// carry.
func TestCheck(t *testing.T) {
	s := fixmsgtest.Session("FIX.4.4")
	order := fix50sp2.BuildNewOrderSingle(s.BeginString).
		SetClOrdID("ORD-1").
		SetSide(fix50sp2.SideBuy).
		SetTransactTime(time.Now()).
		SetOrdType(fix50sp2.OrdTypeMarket).
		SetSymbol("EURUSD").
		SetOrderQty(1).
		SetOrderOrigination(fix50sp2.OrderOriginationOrderReceivedFromACustomer)

	if err := s.Check(order); err == nil {
		t.Error("Check() = nil, want an error")
	}
}
