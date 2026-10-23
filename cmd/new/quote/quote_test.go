package newquote

import (
	"testing"

	"github.com/sylr/fix2/fix"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

func testSession() *config.Session {
	return &config.Session{BeginString: "FIXT.1.1", DefaultApplVerID: "FIX.5.0SP2"}
}

func TestBuildMessage(t *testing.T) {
	optionQuoteID = "quote-1"
	optionSymbol = "EURUSD"
	optionAutoPriceUpdate = false
	optionBuyQuantities, optionBuyPrices = []int64{10}, []float64{1.1}
	optionSellQuantities, optionSellPrices = []int64{20}, []float64{1.2}
	priceIteration = 0

	builder, err := buildMessage("FIXT.1.1", testSession())
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fix.Parse(builder.Build())
	if err != nil {
		t.Fatal(err)
	}

	for tag, want := range map[int]string{
		fix50sp2.TagQuoteID:   "quote-1",
		fix50sp2.TagSymbol:    "EURUSD",
		fix50sp2.TagBidSize:   "10.00",
		fix50sp2.TagBidPx:     "1.10",
		fix50sp2.TagOfferSize: "20.00",
		fix50sp2.TagOfferPx:   "1.20",
	} {
		if got := string(msg.Get(tag)); got != want {
			t.Errorf("tag %d: got %q, want %q", tag, got, want)
		}
	}

	if priceIteration != 1 {
		t.Errorf("priceIteration: got %d, want 1", priceIteration)
	}
}

func TestBuildCancelMessage(t *testing.T) {
	optionQuoteID = "quote-1"
	optionSymbol = "EURUSD"

	builder, err := buildCancelMessage("FIXT.1.1", testSession())
	if err != nil {
		t.Fatal(err)
	}

	msg, err := fix.Parse(builder.Build())
	if err != nil {
		t.Fatal(err)
	}

	for tag, want := range map[int]string{
		fix50sp2.TagQuoteID:         "quote-1",
		fix50sp2.TagQuoteMsgID:      "msg_quote-1",
		fix50sp2.TagQuoteCancelType: string(fix50sp2.QuoteCancelTypeCancelForOneOrMoreSecurities),
		fix50sp2.TagNoQuoteEntries:  "1",
		fix50sp2.TagSymbol:          "EURUSD",
	} {
		if got := string(msg.Get(tag)); got != want {
			t.Errorf("tag %d: got %q, want %q", tag, got, want)
		}
	}
}

func TestBuildMessageVersions(t *testing.T) {
	optionQuoteID = "quote-1"
	optionSymbol = "EURUSD"
	optionAutoPriceUpdate = false
	optionBuyQuantities, optionBuyPrices = []int64{10}, []float64{1.1}
	optionSellQuantities, optionSellPrices = []int64{20}, []float64{1.2}

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			priceIteration = 0
			session := fixmsgtest.Session(version)

			builder, err := buildMessage(session.BeginString, session)
			if err != nil {
				t.Fatal(err)
			}

			fixmsgtest.Validate(t, session, builder)
		})
	}
}

func TestBuildCancelMessageVersions(t *testing.T) {
	optionQuoteID = "quote-1"
	optionSymbol = "EURUSD"

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			builder, err := buildCancelMessage(session.BeginString, session)
			if err != nil {
				t.Fatal(err)
			}

			fixmsgtest.Validate(t, session, builder)
		})
	}
}
