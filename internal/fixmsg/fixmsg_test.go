package fixmsg

import (
	"errors"
	"strings"
	"testing"

	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/session"

	fixerrors "sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix43"
	"sylr.dev/fix/v2/fix44"
	"sylr.dev/fix/v2/fix50sp2"
)

func marketDataRequest(beginString string, entryTypes ...fix50sp2.MDEntryType) *fix50sp2.MarketDataRequestBuilder {
	b := fix50sp2.BuildMarketDataRequest(beginString).
		SetMDReqID("REQ-1").
		SetSubscriptionRequestType(fix50sp2.SubscriptionRequestTypeSnapshot).
		SetMarketDepth(0).
		AddNoRelatedSymEntry(fix50sp2.NewNoRelatedSymEntry().SetSymbol("EURUSD"))

	for _, t := range entryTypes {
		e := fix50sp2.NewNoMDEntryTypesEntry()
		if t != "" {
			e.SetMDEntryType(t)
		}
		b.AddNoMDEntryTypesEntry(e)
	}

	return b
}

func check(t *testing.T, m session.Outgoing, d *dictionary.Dict) error {
	t.Helper()

	msg, err := Message(m, "FIX.4.4", d)
	if err != nil {
		t.Fatal(err)
	}

	return Check(msg, d, dictionary.CheckAll)
}

func TestCheck(t *testing.T) {
	d := fix44.Dictionary()

	if err := check(t, marketDataRequest("FIX.4.4", fix50sp2.MDEntryTypeBid), d); err != nil {
		t.Errorf("valid MarketDataRequest: %v", err)
	}

	// Values beyond the enumeration are the counterparty's to judge.
	if err := check(t, marketDataRequest("FIX.4.4", "101"), d); err != nil {
		t.Errorf("MarketDataRequest with MDEntryType 101: %v", err)
	}

	// A group entry without its required field: its count says one entry,
	// but none follows.
	if err := check(t, marketDataRequest("FIX.4.4", ""), d); err == nil || !strings.Contains(err.Error(), "0 entries follow") {
		t.Errorf("MarketDataRequest with an empty NoMDEntryTypes entry: %v, want its entry missing", err)
	}

	// A required group missing.
	if err := check(t, marketDataRequest("FIX.4.4"), d); err == nil || !strings.Contains(err.Error(), "NoMDEntryTypes") {
		t.Errorf("MarketDataRequest without NoMDEntryTypes: %v, want NoMDEntryTypes missing", err)
	}
}

// TestCheckPartySubIDs checks typed party sub IDs, which FIX 4.3 cannot
// carry, are refused with the field named.
func TestCheckPartySubIDs(t *testing.T) {
	order := fix50sp2.BuildNewOrderSingle("FIX.4.3").
		SetClOrdID("C1").
		SetHandlInst(fix50sp2.HandlInstAutomatedExecutionOrderPrivateNoBrokerIntervention).
		SetSide(fix50sp2.SideBuy).
		SetOrdType(fix50sp2.OrdTypeMarket).
		SetSymbol("EURUSD").
		AddNoPartyIDsEntry(fix50sp2.NewNoPartyIDsEntry().
			SetPartyID("P1").
			SetPartyIDSource(fix50sp2.PartyIDSourceProprietary).
			SetPartyRole(fix50sp2.PartyRoleClientId).
			AddNoPartySubIDsEntry(fix50sp2.NewNoPartySubIDsEntry().SetPartySubID("S1")))
	order.SetString(fix50sp2.TagTransactTime, "20261023-12:00:00")

	msg, err := Message(order, "FIX.4.3", fix43.Dictionary())
	if err != nil {
		t.Fatal(err)
	}

	err = Check(msg, fix43.Dictionary(), dictionary.CheckAll)
	if err == nil || !strings.Contains(err.Error(), "NoPartySubIDs(802)") {
		t.Errorf("Check() = %v, want NoPartySubIDs(802) refused", err)
	}
}

func newOrderSingle(beginString string) *fix50sp2.NewOrderSingleBuilder {
	order := fix50sp2.BuildNewOrderSingle(beginString).
		SetClOrdID("C1").
		SetSide(fix50sp2.SideBuy).
		SetOrdType(fix50sp2.OrdTypeMarket).
		SetSymbol("EURUSD")
	order.SetString(fix50sp2.TagTransactTime, "20261023-12:00:00")

	return order
}

// TestCheckValues checks values outside their enumeration are left to the
// counterparty, but empty ones, which a multiple-value field may hold
// between two spaces, are not.
func TestCheckValues(t *testing.T) {
	d := fix44.Dictionary()

	for _, tc := range []struct {
		name     string
		execInst string
		checks   dictionary.Checks
		want     error
	}{
		{"valid", "1 2", dictionary.CheckAll, nil},
		{"outside the enumeration", "1 ZZ", dictionary.CheckAll, nil},
		{"empty", "", dictionary.CheckAll, dictionary.ErrEmptyValue},
		{"empty, not checked", "", dictionary.CheckAll &^ dictionary.CheckEmptyValues, dictionary.ErrInvalidEnumValue},
		{"two spaces", "1  2", dictionary.CheckAll, dictionary.ErrInvalidEnumValue},
		{"trailing space", "1 ", dictionary.CheckAll, dictionary.ErrInvalidEnumValue},
	} {
		t.Run(tc.name, func(t *testing.T) {
			order := newOrderSingle("FIX.4.4")
			order.SetString(fix50sp2.TagExecInst, tc.execInst)

			msg, err := Message(order, "FIX.4.4", d)
			if err != nil {
				t.Fatal(err)
			}

			err = Check(msg, d, tc.checks)
			switch {
			case tc.want == nil && err != nil:
				t.Errorf("Check() = %v, want nil", err)
			case tc.want != nil && (!errors.Is(err, tc.want) || !errors.Is(err, ErrNotSupported)):
				t.Errorf("Check() = %v, want %v", err, tc.want)
			}
		})
	}
}

// TestCheckUnknownTags checks tags the message does not define are
// refused as dictionary.ValidateWith refuses them: user-defined ones, 5000
// and up, with CheckUserDefinedTags, the others with CheckUnknownTags.
func TestCheckUnknownTags(t *testing.T) {
	d := fix44.Dictionary()

	for _, tc := range []struct {
		tag    int
		checks dictionary.Checks
		ok     bool
	}{
		{9999, dictionary.CheckAll, false},
		{9999, dictionary.CheckAll &^ dictionary.CheckUserDefinedTags, true},
		{9999, dictionary.CheckAll &^ dictionary.CheckUnknownTags, false},
		{4999, dictionary.CheckAll &^ dictionary.CheckUnknownTags, true},
		{4999, dictionary.CheckAll &^ dictionary.CheckUserDefinedTags, false},
	} {
		order := newOrderSingle("FIX.4.4")
		order.SetString(tc.tag, "X")

		msg, err := Message(order, "FIX.4.4", d)
		if err != nil {
			t.Fatal(err)
		}

		err = Check(msg, d, tc.checks)
		if tc.ok && err != nil {
			t.Errorf("tag %d, checks %b: Check() = %v, want nil", tc.tag, tc.checks, err)
		} else if !tc.ok && !errors.Is(err, ErrNotSupported) {
			t.Errorf("tag %d, checks %b: Check() = %v, want %v", tc.tag, tc.checks, err, ErrNotSupported)
		}
	}
}

// TestCheckStructure checks MsgType and NUMINGROUP fields, outside their
// enumeration, are reported rather than replaced: replacing them would
// change the message's type or its groups.
func TestCheckStructure(t *testing.T) {
	ref := func(tags ...int) []dictionary.FieldRef {
		var refs []dictionary.FieldRef
		for _, tag := range tags {
			refs = append(refs, dictionary.FieldRef{Tag: tag, Required: true})
		}
		return refs
	}

	d := dictionary.New(&dictionary.Dict{
		Type: "FIX", Major: "4", Minor: "4",
		Fields: map[int]*dictionary.FieldDef{
			8:   {Tag: 8, Name: "BeginString", Type: "STRING"},
			9:   {Tag: 9, Name: "BodyLength", Type: "LENGTH"},
			10:  {Tag: 10, Name: "CheckSum", Type: "STRING"},
			11:  {Tag: 11, Name: "ClOrdID", Type: "STRING"},
			34:  {Tag: 34, Name: "MsgSeqNum", Type: "SEQNUM"},
			35:  {Tag: 35, Name: "MsgType", Type: "STRING", Enums: map[string]string{"A": "LOGON"}},
			49:  {Tag: 49, Name: "SenderCompID", Type: "STRING"},
			52:  {Tag: 52, Name: "SendingTime", Type: "UTCTIMESTAMP"},
			54:  {Tag: 54, Name: "Side", Type: "CHAR", Enums: map[string]string{"1": "BUY", "2": "SELL"}},
			56:  {Tag: 56, Name: "TargetCompID", Type: "STRING"},
			552: {Tag: 552, Name: "NoSides", Type: "NUMINGROUP", Enums: map[string]string{"1": "ONE_SIDE"}},
		},
		Header:  ref(8, 9, 35, 49, 56, 34, 52),
		Trailer: ref(10),
		Messages: map[string]*dictionary.MessageDef{
			"D": {
				Name: "NewOrderSingle", MsgType: "D", MsgCat: "app",
				Fields: []dictionary.FieldRef{{Tag: 11, Required: true}, {Tag: 552}, {Tag: 54}},
				Groups: []dictionary.GroupDef{{
					CounterTag: 552, Name: "NoSides", DelimTag: 54,
					Fields: ref(54), Layout: []int{54},
				}},
			},
		},
	})

	body := fix.AppendTagValue(nil, tagMsgType, "D")
	body = fix.AppendTagValue(body, tagSenderCompID, "SENDER")
	body = fix.AppendTagValue(body, tagTargetCompID, "TARGET")
	body = fix.AppendTagInt(body, tagMsgSeqNum, 1)
	body = fix.AppendTagValue(body, tagSendingTime, "20261023-12:00:00.000")
	body = fix.AppendTagValue(body, 11, "C1")
	body = fix.AppendTagValue(body, 552, "2")
	body = fix.AppendTagValue(body, 54, "1")
	body = fix.AppendTagValue(body, 54, "2")

	msg, err := parse("FIX.4.4", body, d)
	if err != nil {
		t.Fatal(err)
	}

	err = Check(msg, d, dictionary.CheckAll)
	var verr *dictionary.ValidationError
	if !errors.As(err, &verr) || !errors.Is(err, dictionary.ErrInvalidEnumValue) || verr.Tag != tagMsgType {
		t.Errorf("Check() = %v, want MsgType D's value refused", err)
	}

	d.Fields[35].Enums["D"] = "ORDER_SINGLE"
	d = dictionary.New(d)

	err = Check(msg, d, dictionary.CheckAll)
	if !errors.As(err, &verr) || !errors.Is(err, dictionary.ErrInvalidEnumValue) || verr.Tag != 552 {
		t.Errorf("Check() = %v, want NoSides 2's value refused", err)
	}
}

func TestCheckErrors(t *testing.T) {
	if err := Check(&fix.Message{}, nil, dictionary.CheckAll); !errors.Is(err, fixerrors.ErrFix) {
		t.Errorf("Check() without a dictionary = %v, want %v", err, fixerrors.ErrFix)
	}

	// fix's FIX 5.0 SP2 dictionary lacks IOI, which FIX 5.0 SP2 has.
	d := fix50sp2.Dictionary()
	msg, err := parse("FIXT.1.1", fix.AppendTagValue(nil, tagMsgType, "6"), d)
	if err != nil {
		t.Fatal(err)
	}

	err = Check(msg, d, dictionary.CheckAll)
	if !errors.Is(err, ErrNotSupported) || !errors.Is(err, fixerrors.ErrFix) || !strings.Contains(err.Error(), "the session's FIX.5.0SP2 dictionary does not define MsgType 6") {
		t.Errorf("Check() = %v, want MsgType 6 not defined", err)
	}
}

func TestValidEnum(t *testing.T) {
	single := &dictionary.FieldDef{Type: "CHAR", Enums: map[string]string{"2": "", "1": ""}}
	multiple := &dictionary.FieldDef{Type: "MULTIPLECHARVALUE", Enums: single.Enums}

	for _, tc := range []struct {
		def     *dictionary.FieldDef
		value   string
		want    string
		changed bool
	}{
		{single, "2", "2", false},
		{single, "", "", false},
		{single, "Z", "1", true},
		{single, "1 2", "1", true},
		{multiple, "2 1", "2 1", false},
		{multiple, "2 Z", "2 1", true},
		{multiple, "2  Z", "2  1", true},
		{multiple, "", "", false},
	} {
		if got, changed := validEnum(tc.def, tc.value); got != tc.want || changed != tc.changed {
			t.Errorf("validEnum(%s, %q) = %q, %t, want %q, %t", tc.def.Type, tc.value, got, changed, tc.want, tc.changed)
		}
	}
}
