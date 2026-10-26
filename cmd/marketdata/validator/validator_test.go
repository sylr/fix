//go:build validator

package marketdatavalidator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
	fixerrors "sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator/initiatortest"
	"sylr.dev/fix/v2/internal/fixmsg"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
	"sylr.dev/fix/v2/marketdata"
)

// setOptions sets the options of the command as its flags would, restoring
// their defaults at the end of t.
func setOptions(t *testing.T, symbol string) {
	t.Helper()

	optionSymbol, optionMDReqID, optionTradeHistory = symbol, "", true
	t.Cleanup(func() { optionSymbol, optionMDReqID, optionTradeHistory = "", "", true })
}

func TestValidate(t *testing.T) {
	setOptions(t, "")
	if err := Validate(MarketDataValidatorCmd, nil); !errors.Is(err, fixerrors.ErrOptions) {
		t.Errorf("Validate(no --symbol) = %v, want an options error", err)
	}

	setOptions(t, "EURUSD")
	if err := Validate(MarketDataValidatorCmd, nil); err != nil {
		t.Fatal(err)
	}
	if optionMDReqID == "" {
		t.Error("Validate() did not generate the MDReqID")
	}
}

// The validator deliberately subscribes to MDEntryType "101", a non-standard
// value (Trade History): the request is checked for its structure only, as
// full validation would refuse the value.
func TestBuildMessageVersions(t *testing.T) {
	setOptions(t, "EURUSD")
	optionMDReqID = "REQ-1"

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			b, err := buildMessage(session)
			if err != nil {
				t.Fatal(err)
			}

			_, d, err := session.FIXDictionaries()
			if err != nil {
				t.Fatal(err)
			}

			msg, err := fixmsg.Message(b, session.BeginString, d)
			if err != nil {
				t.Fatal(err)
			}

			if err := fixmsg.Check(msg, d, session.Checks()); err != nil {
				t.Fatalf("%s: %v", session.Version(), err)
			}
			if id := string(msg.Get(fix50sp2.TagMDReqID)); id != "REQ-1" {
				t.Errorf("MDReqID = %q, want REQ-1", id)
			}
		})
	}
}

// venue answers each MarketDataRequest with a snapshot, of another request
// if other is set, then logs out if logout is set; or rejects it with a
// MarketDataRequestReject if reject is set, or with err.
type venue struct {
	session.BaseHandler
	reject bool
	other  bool
	logout bool
	err    error
}

func (v venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	if v.err != nil {
		return v.err
	}

	reqID := string(m.Get(fix50sp2.TagMDReqID))
	beginString := s.ID().BeginString

	if v.reject {
		_, err := s.Send(ctx, fix50sp2.BuildMarketDataRequestReject(beginString).
			SetMDReqID(reqID).
			SetMDReqRejReason(fix50sp2.MDReqRejReasonUnknownSymbol).
			SetText("unknown symbol"))
		return err
	}

	if v.other {
		reqID = "another"
	}
	snapshot := fix50sp2.BuildMarketDataSnapshotFullRefresh(beginString).
		SetMDReqID(reqID).
		SetSymbol("EURUSD").
		SetLastUpdateTime(time.Now()).
		AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().
			SetMDEntryType(fix50sp2.MDEntryTypeBid).
			SetOrderID("1").
			SetOrdType(fix50sp2.OrdTypeLimit).
			SetMDEntryPx(1).
			SetMDEntrySize(1))
	if _, err := s.Send(ctx, snapshot); err != nil {
		return err
	}

	if v.logout {
		go func() { _ = s.Logout(context.Background(), "bye") }()
	}

	return nil
}

func TestExecute(t *testing.T) {
	for _, tc := range []struct {
		name    string
		handler venue
		want    error
	}{
		{"subscription", venue{}, nil},
		{"another request", venue{other: true}, fixerrors.ErrResponseTimeout},
		{"request reject", venue{reject: true}, marketdata.ErrRequestRejected},
		{"business reject", venue{err: &session.BusinessRejectError{Reason: session.BusinessRejectUnsupportedMessageType, Text: "no market data"}}, fixerrors.ErrFixRejected},
		{"logout", venue{logout: true}, fixerrors.ErrFixLogout},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// The acceptor reads the timeout too: it is restored once it
			// stopped.
			config.GetOptions().Timeout = 500 * time.Millisecond
			t.Cleanup(func() { config.GetOptions().Timeout = 0 })
			initiatortest.Serve(t, tc.handler, "FIX.5.0SP2")
			// The acceptor refuses MDEntryType 101, not a CHAR.
			setOptions(t, "EURUSD")
			optionTradeHistory = false
			if err := Validate(MarketDataValidatorCmd, nil); err != nil {
				t.Fatal(err)
			}

			// A subscription runs until interrupted, past the timeout.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			MarketDataValidatorCmd.SetContext(ctx)

			err := Execute(MarketDataValidatorCmd, nil)
			if !errors.Is(err, tc.want) || (tc.want == nil && err != nil) {
				t.Errorf("Execute() = %v, want %v", err, tc.want)
			}
		})
	}
}
