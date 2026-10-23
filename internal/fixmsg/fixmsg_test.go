package fixmsg

import (
	"strings"
	"testing"

	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/session"

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
