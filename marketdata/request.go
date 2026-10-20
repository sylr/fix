package marketdata

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/olekukonko/tablewriter"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"

	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/fixt11"
	"sylr.dev/fix/v2/utils"
)

const nilstr = "<nil>"

var (
	ErrRejected        = fmt.Errorf("%w: message rejected", errors.Fix)
	ErrRequestRejected = fmt.Errorf("%w: market data request rejected", errors.Fix)
)

var (
	type2sym = map[string]string{
		"Bid":   "+",
		"Trade": "=",
		"Offer": "-",
	}
)

// RequestHandler handles the answers to a MarketDataRequest.
type RequestHandler struct {
	fix50sp2.UnimplementedHandler

	// Writer receives the entries of the refreshes as tables, nothing is
	// printed if nil.
	Writer io.Writer
	// Dictionary describes the enums printed, may be nil.
	Dictionary *dictionary.Dict
}

// Handle handles m, a message received from initiator.Client.Messages:
// Rejects and MarketDataRequestRejects are returned as errors.
func (h *RequestHandler) Handle(m *fix.Message) error {
	if string(m.MsgType()) == string(fixt11.MsgTypeReject) {
		return fmt.Errorf("%w: %s", ErrRejected, m.Get(fixt11.TagText))
	}

	return fix50sp2.Dispatch(m, h)
}

func (h *RequestHandler) OnMarketDataSnapshotFullRefresh(msg *fix50sp2.MarketDataSnapshotFullRefresh) error {
	if h.Writer != nil {
		printFIX50NoMDEntriesFull(h.Writer, msg, h.Dictionary)
	}

	return nil
}

func (h *RequestHandler) OnMarketDataIncrementalRefresh(msg *fix50sp2.MarketDataIncrementalRefresh) error {
	if h.Writer != nil {
		printFIX50NoMDEntriesInc(h.Writer, msg, h.Dictionary)
	}

	return nil
}

func (h *RequestHandler) OnMarketDataRequestReject(msg *fix50sp2.MarketDataRequestReject) error {
	reason, _ := msg.MDReqRejReason()
	text, _ := msg.Text()

	return fmt.Errorf("%w: reason=%s text=%s", ErrRequestRejected, enumName(h.Dictionary, fix50sp2.TagMDReqRejReason, string(reason)), text)
}

// entry holds the columns printed for an MD entry.
type entry struct {
	typSign, typ, orderTradeID, action, symbol, price, size, tim string
}

func newEntry(e *fix50sp2.NoMDEntriesEntry, d *dictionary.Dict) entry {
	var en entry

	if updateAction, ok := e.MDUpdateAction(); ok {
		en.action = enumName(d, fix50sp2.TagMDUpdateAction, string(updateAction))
	} else {
		en.action = nilstr
	}

	if entryType, ok := e.MDEntryType(); ok {
		en.typ = enumName(d, fix50sp2.TagMDEntryType, string(entryType))
		en.typSign = type2sym[en.typ]
	} else {
		en.typ = nilstr
	}

	if orderType, ok := e.OrdType(); ok {
		en.typ = fmt.Sprintf("%s (%s)", en.typ, enumName(d, fix50sp2.TagOrdType, string(orderType)))
	}

	if orderID, ok := e.OrderID(); ok && len(orderID) > 0 {
		en.orderTradeID = orderID
	} else if tradeID, ok := e.TradeID(); ok {
		en.orderTradeID = tradeID
	} else {
		en.orderTradeID = nilstr
	}

	en.symbol = valueOrNil(e.Symbol())
	en.price = valueOrNil(e.MDEntryPxBytes())
	en.size = valueOrNil(e.MDEntrySizeBytes())

	stringDate, okDate := e.MDEntryDate()
	stringTime, okTime := e.MDEntryTimeBytes()

	if okDate && okTime {
		timeDate, _ := time.Parse("20060102", stringDate)
		timeTime, _ := time.Parse("15:04:05.999999999", string(stringTime))
		en.tim = utils.CombineDateAndTime(timeDate, timeTime).Format(time.RFC3339Nano)
	} else if okDate {
		timeDate, _ := time.Parse("20060102", stringDate)
		en.tim = timeDate.Format("2006-01-02")
	} else if okTime {
		timeTime, _ := time.Parse("15:04:05.999999999", string(stringTime))
		en.tim = timeTime.Format("15:04:05.999")
	}

	return en
}

func printFIX50NoMDEntriesFull(w io.Writer, msg *fix50sp2.MarketDataSnapshotFullRefresh, d *dictionary.Dict) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"SYMBOL", "ORDER ID", "TYPE", "PRICE", "SIZE", "TIME"})
	table.SetBorders(tablewriter.Border{Left: false, Right: false, Top: false, Bottom: true})
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("-")

	symbol := valueOrNil(msg.Symbol())

	for _, e := range msg.NoMDEntriesEntries() {
		en := newEntry(&e, d)
		table.Append([]string{fmt.Sprintf("%s %s", en.typSign, symbol), en.orderTradeID, en.typ, en.price, en.size, en.tim})
	}

	if last, ok, err := msg.LastUpdateTime(); ok && err == nil {
		table.SetFooter([]string{"", "", "", "", "Last Time", last.Format(time.RFC3339Nano)})
	}

	table.Render()
}

func printFIX50NoMDEntriesInc(w io.Writer, msg *fix50sp2.MarketDataIncrementalRefresh, d *dictionary.Dict) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"SYMBOL", "ID", "ACTION", "TYPE", "PRICE", "SIZE", "TIME"})
	table.SetBorders(tablewriter.Border{Left: false, Right: false, Top: false, Bottom: true})
	table.SetColMinWidth(2, 8)
	table.SetColMinWidth(3, 14)
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("-")

	for _, e := range msg.NoMDEntriesEntries() {
		en := newEntry(&e, d)
		table.Append([]string{fmt.Sprintf("%s %s", en.typSign, en.symbol), en.orderTradeID, en.action, en.typ, en.price, en.size, en.tim})
	}

	table.Render()
}

// enumName returns the description d gives to value of tag in CamelCase,
// value itself if none.
func enumName(d *dictionary.Dict, tag int, value string) string {
	if d != nil {
		if def := d.Fields[tag]; def != nil {
			if desc, ok := def.Enums[value]; ok {
				return strcase.ToCamel(strings.ToLower(desc))
			}
		}
	}

	return value
}

func valueOrNil[T string | []byte](v T, ok bool) string {
	if !ok {
		return nilstr
	}

	return string(v)
}
