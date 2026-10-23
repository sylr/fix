package marketdatarequest

import (
	"testing"

	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

func TestBuildMessage(t *testing.T) {
	optionMDReqID = "REQ-1"
	optionMarketDepth = 5
	optionTypes = []string{"bid", "offer"}
	optionSymbols = []string{"EURUSD", "GBPUSD"}
	SubType = fix50sp2.SubscriptionRequestTypeSnapshotPlusUpdates
	MDUpdateType = fix50sp2.MDUpdateTypeIncrementalRefresh

	c := &initiator.Client{Session: &config.Session{BeginString: session.FIXT11, DefaultApplVerID: "FIX.5.0SP2"}}
	b, err := buildMessage(c)
	if err != nil {
		t.Fatal(err)
	}

	m, err := fix.Parse(b.(*fix50sp2.MarketDataRequestBuilder).Build())
	if err != nil {
		t.Fatal(err)
	}
	request := fix50sp2.WrapMarketDataRequest(&m)

	if id, _ := request.MDReqID(); id != "REQ-1" {
		t.Errorf("MDReqID = %q, want REQ-1", id)
	}
	if typ, _ := request.SubscriptionRequestType(); typ != SubType {
		t.Errorf("SubscriptionRequestType = %q, want %q", typ, SubType)
	}
	if depth, _, _ := request.MarketDepth(); depth != 5 {
		t.Errorf("MarketDepth = %d, want 5", depth)
	}
	if typ, _ := request.MDUpdateType(); typ != MDUpdateType {
		t.Errorf("MDUpdateType = %q, want %q", typ, MDUpdateType)
	}

	var types []fix50sp2.MDEntryType
	for _, e := range request.NoMDEntryTypesEntries() {
		typ, _ := e.MDEntryType()
		types = append(types, typ)
	}
	if len(types) != 2 || types[0] != fix50sp2.MDEntryTypeBid || types[1] != fix50sp2.MDEntryTypeOffer {
		t.Errorf("MDEntryTypes = %v, want [0 1]", types)
	}

	var symbols []string
	for _, e := range request.NoRelatedSymEntries() {
		symbol, _ := e.Symbol()
		symbols = append(symbols, symbol)
	}
	if len(symbols) != 2 || symbols[0] != "EURUSD" || symbols[1] != "GBPUSD" {
		t.Errorf("Symbols = %v, want [EURUSD GBPUSD]", symbols)
	}
}

func TestBuildMessageVersions(t *testing.T) {
	optionMDReqID = "REQ-1"
	optionMarketDepth = 5
	optionTypes = []string{"bid", "offer"}
	optionSymbols = []string{"EURUSD", "GBPUSD"}
	SubType = fix50sp2.SubscriptionRequestTypeSnapshotPlusUpdates
	MDUpdateType = fix50sp2.MDUpdateTypeIncrementalRefresh

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			b, err := buildMessage(&initiator.Client{Session: session})
			if err != nil {
				t.Fatal(err)
			}

			fixmsgtest.Validate(t, session, b)
		})
	}
}
