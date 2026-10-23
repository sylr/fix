package status_security

import (
	"testing"

	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

func TestBuildMessageVersions(t *testing.T) {
	optionSecurityStatReqID = "REQ-1"
	optionSymbol = "EURUSD"
	optionSubType = "snapshot"

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
