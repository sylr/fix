package application

import (
	"testing"
	"time"

	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/internal/fixmsg"
	"sylr.dev/fix/v2/internal/fixmsg/fixmsgtest"
)

func TestBuildExecutionReportVersions(t *testing.T) {
	for _, version := range fixmsgtest.Versions {
		t.Run(version, func(t *testing.T) {
			session := fixmsgtest.Session(version)

			order := fix50sp2.BuildNewOrderSingle(session.BeginString).
				SetClOrdID("order-1").
				SetSide(fix50sp2.SideBuy).
				SetTransactTime(time.Now()).
				SetOrdType(fix50sp2.OrdTypeLimit).
				SetSymbol("EURUSD")
			order.SetString(fix50sp2.TagOrderQty, "10.00")

			msg, err := fixmsg.Message(order, session.BeginString, nil)
			if err != nil {
				t.Fatal(err)
			}

			report, err := buildExecutionReport(session.BeginString, fix50sp2.WrapNewOrderSingle(msg), fix50sp2.SideBuy, fix50sp2.OrdStatusNew)
			if err != nil {
				t.Fatal(err)
			}

			fixmsgtest.Validate(t, session, report)
		})
	}
}
