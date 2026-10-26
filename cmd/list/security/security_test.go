package listsecurity

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/config"
	fixerrors "sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator/initiatortest"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

// setOptions sets the options of the command as its flags would, restoring
// their defaults at the end of t.
func setOptions(t *testing.T, typ, symbol string) {
	t.Helper()

	optionType, optionSymbol, optionSecurityReqID = typ, symbol, ""
	t.Cleanup(func() {
		optionType, optionSymbol, optionSecurityReqID = "all_securities", "", ""
	})
}

func TestValidate(t *testing.T) {
	for _, tc := range []struct {
		typ, symbol string
		wantErr     bool
	}{
		{"all_securities", "", false},
		{"SYMBOL", "EURUSD", false},
		{"symbol", "", true},
		{"product", "EURUSD", true},
		{"never", "", true},
	} {
		t.Run(tc.typ+"/"+tc.symbol, func(t *testing.T) {
			setOptions(t, tc.typ, tc.symbol)

			err := Validate(ListSecurityCmd, nil)
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Fatalf("Validate() = %v, want error %v", err, tc.wantErr)
			}
			if err != nil && !errors.Is(err, fixerrors.ErrOptions) {
				t.Errorf("Validate() = %v, want an options error", err)
			}
			if err == nil && optionSecurityReqID == "" {
				t.Error("Validate() did not generate the SecurityReqID")
			}
		})
	}
}

func TestBuildMessageVersions(t *testing.T) {
	setOptions(t, "symbol", "EURUSD")
	if err := Validate(ListSecurityCmd, nil); err != nil {
		t.Fatal(err)
	}

	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			s := fixmsgtest.Session(version)

			b, err := buildMessage(s)
			if err != nil {
				t.Fatal(err)
			}

			m := fixmsgtest.Validate(t, s, b)
			if id := string(m.Get(fix50sp2.TagSecurityReqID)); id != optionSecurityReqID {
				t.Errorf("SecurityReqID = %q, want %q", id, optionSecurityReqID)
			}
			if symbol := string(m.Get(fix50sp2.TagSymbol)); symbol != "EURUSD" {
				t.Errorf("Symbol = %q, want EURUSD", symbol)
			}
		})
	}
}

// venue answers each SecurityListRequest with a SecurityList of another
// request, then one fragment per entry of fragments, holding that many
// securities. It tells the last with LastFragment, or the total with
// TotNoRelatedSym if total is set, but leaves out the last lost ones.
// result, if set, is the SecurityRequestResult; err, if set, rejects the
// request instead. chatter, if set, makes it send a SecurityList of another
// request every 100ms afterwards.
type venue struct {
	session.BaseHandler
	fragments []int
	total     bool
	lost      int
	result    fix50sp2.SecurityRequestResult
	err       error
	chatter   bool
}

func (v venue) FromApp(ctx context.Context, s *session.Session, m *fix.Message, _ session.Inbound) error {
	if v.err != nil {
		return v.err
	}

	reqID := string(m.Get(fix50sp2.TagSecurityReqID))
	if _, err := s.Send(ctx, fix50sp2.BuildSecurityList(s.ID().BeginString).SetSecurityReqID("another")); err != nil {
		return err
	}

	total := 0
	for _, n := range v.fragments {
		total += n
	}

	listed := 0
	for i, n := range v.fragments[:len(v.fragments)-v.lost] {
		list := fix50sp2.BuildSecurityList(s.ID().BeginString).SetSecurityReqID(reqID)
		if v.result != "" {
			list.SetSecurityRequestResult(v.result)
		}
		if v.total {
			list.SetTotNoRelatedSym(total)
		} else if len(v.fragments) > 1 {
			last := fix50sp2.LastFragmentNo
			if i == len(v.fragments)-1 {
				last = fix50sp2.LastFragmentYes
			}
			list.SetLastFragment(last)
		}
		for range n {
			listed++
			list.AddNoRelatedSymEntry(fix50sp2.NewNoRelatedSymEntry().SetSymbol(fmt.Sprintf("SYM-%d", listed)))
		}
		if _, err := s.Send(ctx, list); err != nil {
			return err
		}
	}

	if v.chatter {
		go func() {
			for range 50 {
				time.Sleep(100 * time.Millisecond)
				if _, err := s.Send(context.Background(), fix50sp2.BuildSecurityList(s.ID().BeginString).SetSecurityReqID("another")); err != nil {
					return
				}
			}
		}()
	}

	return nil
}

func TestExecute(t *testing.T) {
	for _, tc := range []struct {
		name    string
		handler venue
		want    error
		listed  int
	}{
		{"one list", venue{fragments: []int{3}}, nil, 3},
		{"last fragment", venue{fragments: []int{2, 2, 1}}, nil, 5},
		{"total", venue{fragments: []int{2, 2, 1}, total: true}, nil, 5},
		{"lost fragment", venue{fragments: []int{2, 2, 1}, total: true, lost: 1}, fixerrors.ErrResponseTimeout, 4},
		{"lost fragment, other lists", venue{fragments: []int{2, 2, 1}, total: true, lost: 1, chatter: true}, fixerrors.ErrResponseTimeout, 4},
		{"invalid request", venue{fragments: []int{0}, result: fix50sp2.SecurityRequestResultInvalidOrUnsupportedRequest}, fixerrors.ErrFixRejected, 0},
		{"timeout", venue{}, fixerrors.ErrResponseTimeout, 0},
		{"reject", venue{err: &session.BusinessRejectError{Reason: session.BusinessRejectUnsupportedMessageType, Text: "no lists"}}, fixerrors.ErrFixRejected, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// The acceptor reads the timeout too: it is restored once it
			// stopped.
			config.GetOptions().Timeout = 500 * time.Millisecond
			t.Cleanup(func() { config.GetOptions().Timeout = 0 })
			initiatortest.Serve(t, tc.handler, "FIX.5.0SP2")
			setOptions(t, "all_securities", "")
			if err := Validate(ListSecurityCmd, nil); err != nil {
				t.Fatal(err)
			}

			// A timeout defeated would wait until ctx is done.
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			var out bytes.Buffer
			ListSecurityCmd.SetContext(ctx)
			ListSecurityCmd.SetOut(&out)
			t.Cleanup(func() { ListSecurityCmd.SetOut(nil) })

			err := Execute(ListSecurityCmd, nil)
			if !errors.Is(err, tc.want) || (tc.want == nil && err != nil) {
				t.Errorf("Execute() = %v, want %v", err, tc.want)
			}
			if got := strings.Count(out.String(), "SYM-"); got != tc.listed {
				t.Errorf("listed %d securities, want %d:\n%s", got, tc.listed, out.String())
			}
		})
	}
}
