package utils_test

import (
	"bytes"
	"testing"

	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
	"sylr.dev/fix/v2/utils"
)

func TestWriteMessageBodyAsTableGroups(t *testing.T) {
	b := fix50sp2.BuildMarketDataSnapshotFullRefresh(session.FIXT11).
		SetMDReqID("REQ-1").
		SetSymbol("EURUSD").
		AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().
			SetMDEntryType(fix50sp2.MDEntryTypeBid).
			SetMDEntryPx(1.1).
			AddNoPartyIDsEntry(fix50sp2.NewNoPartyIDsEntry().SetPartyID("A").SetPartyRole(fix50sp2.PartyRoleExecutingFirm)).
			AddNoPartyIDsEntry(fix50sp2.NewNoPartyIDsEntry().SetPartyID("B"))).
		AddNoMDEntriesEntry(fix50sp2.NewNoMDEntriesEntry().
			SetMDEntryType(fix50sp2.MDEntryTypeOffer).
			SetMDEntryPx(1.2))

	m, err := fix.Parse(b.Build())
	if err != nil {
		t.Fatal(err)
	}

	_, d, err := fixmsgtest.Session("FIX.5.0SP2").FIXDictionaries()
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	utils.WriteMessageBodyAsTable(&out, &m, d)

	want := `
  TAG    DESCRIPTION          VALUE         
--------------------------------------------
  262   MDReqID         REQ-1               
   55   Symbol          EURUSD              
  268   NoMDEntries     2                   
  269   - MDEntryType   0 (BID)             
  270     MDEntryPx     1.1                 
  453     NoPartyIDs    2                   
  448     - PartyID     A                   
  452       PartyRole   1 (EXECUTING_FIRM)  
  448     - PartyID     B                   
  269   - MDEntryType   1 (OFFER)           
  270     MDEntryPx     1.2                 
--------------------------------------------
`
	if got := "\n" + out.String(); got != want {
		t.Errorf("table:%s\nwant:%s", got, want)
	}
}
