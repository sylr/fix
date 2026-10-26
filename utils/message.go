package utils

import (
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"
)

// WriteMessageBodyAsTable writes the body of message as a table, a row per
// field in the message's order, named and described by d if not nil. The
// entries of the repeating groups d defines follow their counter, indented,
// the first field of each marked with "- ".
func WriteMessageBodyAsTable(w io.Writer, message *fix.Message, d *dictionary.Dict) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"TAG", "DESCRIPTION", "VALUE"})
	table.SetBorders(tablewriter.Border{Bottom: true})
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("-")
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.SetAutoWrapText(false)

	var sections []dictionary.Section
	var groups []dictionary.GroupDef
	if d != nil {
		sections = d.Sections(message)
		if def := d.Messages[string(message.MsgType())]; def != nil {
			groups = def.Groups
		}
	}

	// open are the groups the field read belongs to, the innermost last.
	var open []*openGroup

	for i, field := range message.Fields() {
		if sections != nil && sections[i] != dictionary.SectionBody {
			continue
		} else if sections == nil && isEnvelopeTag(field.Tag) {
			continue
		}

		for len(open) > 0 && !open[len(open)-1].holds(field.Tag) {
			open = open[:len(open)-1]
		}

		indent := ""
		if depth := len(open); depth > 0 {
			indent = strings.Repeat("  ", depth-1)
			if open[depth-1].next(field.Tag) {
				indent += "- "
			} else {
				indent += "  "
			}
		}

		value := string(field.Value)
		description := "<unknown>"
		if d != nil {
			if def := d.Fields[field.Tag]; def != nil {
				description = def.Name
				if desc, ok := def.Enums[value]; ok {
					value += " (" + desc + ")"
				}
			}
		}

		table.Append([]string{strconv.Itoa(field.Tag), indent + description, value})

		// A counter opens its group, among the message's or the entry's.
		defs := groups
		if len(open) > 0 {
			defs = open[len(open)-1].def.Groups
		}
		for k := range defs {
			if defs[k].CounterTag == field.Tag {
				count, _ := strconv.Atoi(string(field.Value))
				open = append(open, &openGroup{def: &defs[k], count: count})
				break
			}
		}
	}

	table.Render()
}

// openGroup is a repeating group being read.
type openGroup struct {
	def *dictionary.GroupDef
	// count is the number of entries the counter tells, entries the number
	// read; at is the position in def.Layout of the last field read.
	count, entries, at int
}

// holds reports whether a field tag read next belongs to the group: it
// starts an entry the counter tells, or follows the last field of the entry
// in the group's layout.
func (g *openGroup) holds(tag int) bool {
	switch at := slices.Index(g.def.Layout, tag); {
	case at < 0:
		return false
	case at == 0:
		return g.entries < g.count
	default:
		return g.entries > 0 && at > g.at
	}
}

// next records the field tag, which the group holds, and reports whether it
// starts an entry.
func (g *openGroup) next(tag int) bool {
	g.at = slices.Index(g.def.Layout, tag)
	if g.at == 0 {
		g.entries++
		return true
	}

	return false
}

// isEnvelopeTag reports whether tag belongs to the standard header or
// trailer, for messages printed without a dictionary.
func isEnvelopeTag(tag int) bool {
	switch tag {
	case 8, 9, 10, 34, 35, 43, 49, 50, 52, 56, 57, 97, 115, 116, 122, 128, 129, 142, 143, 144, 145, 369, 1128, 1129:
		return true
	}

	return false
}
