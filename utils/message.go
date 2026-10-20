package utils

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"
)

// WriteMessageBodyAsTable writes the body of message as a table, a row per
// tag, its values joined, named and described by d if not nil.
func WriteMessageBodyAsTable(w io.Writer, message *fix.Message, d *dictionary.Dict) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"TAG", "DESCRIPTION", "VALUES"})
	table.SetBorders(tablewriter.Border{Bottom: true})
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("-")
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})

	var sections []dictionary.Section
	if d != nil {
		sections = d.Sections(message)
	}

	var tags []int
	values := map[int][]string{}

	for i, field := range message.Fields() {
		if sections != nil && sections[i] != dictionary.SectionBody {
			continue
		} else if sections == nil && isEnvelopeTag(field.Tag) {
			continue
		}

		value := string(field.Value)
		if def := fieldDef(d, field.Tag); def != nil {
			if desc, ok := def.Enums[value]; ok {
				value += fmt.Sprintf(" (%s)", desc)
			}
		}

		if !slices.Contains(tags, field.Tag) {
			tags = append(tags, field.Tag)
		}
		values[field.Tag] = append(values[field.Tag], value)
	}

	for _, tag := range tags {
		description := "<unknown>"
		if def := fieldDef(d, tag); def != nil {
			description = def.Name
		}

		table.Append([]string{
			strconv.Itoa(tag),
			description,
			strings.Join(values[tag], ", "),
		})
	}

	table.Render()
}

// FieldName returns the name d gives tag, "<unknown>" if none.
func FieldName(d *dictionary.Dict, tag int) string {
	if def := fieldDef(d, tag); def != nil {
		return def.Name
	}

	return "<unknown>"
}

func fieldDef(d *dictionary.Dict, tag int) *dictionary.FieldDef {
	if d == nil {
		return nil
	}

	return d.Fields[tag]
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

func MapSearch[K comparable, V comparable](m map[K]V, search V) *K {
	for k, v := range m {
		if search == v {
			return &k
		}
	}

	return nil
}
