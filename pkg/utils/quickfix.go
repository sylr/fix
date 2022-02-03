package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/iancoleman/strcase"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/datadictionary"
	"github.com/rs/zerolog"
)

type SortableTag []quickfix.Tag

func (t SortableTag) Len() int           { return len(t) }
func (t SortableTag) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t SortableTag) Less(i, j int) bool { return t[i] < t[j] }

type AppMessageLogger struct {
	Logger                  *zerolog.Logger
	TransportDataDictionary *datadictionary.DataDictionary
	AppDataDictionary       *datadictionary.DataDictionary
}

func (app *AppMessageLogger) LogMessage(message *quickfix.Message, sessionID quickfix.SessionID, sending bool) {
	if app.Logger.GetLevel() > zerolog.TraceLevel {
		return
	}

	loop := []struct {
		prefix string
		fm     quickfix.FieldMap
		dict   *datadictionary.DataDictionary
	}{
		{"Headers ", message.Header.FieldMap, app.TransportDataDictionary},
		{"Body    ", message.Body.FieldMap, app.AppDataDictionary},
		{"Trailers", message.Trailer.FieldMap, app.TransportDataDictionary},
	}
	for _, i := range loop {
		if len(i.fm.Tags()) == 0 {
			continue
		}
		w := bytes.NewBuffer([]byte{})
		bw := bufio.NewWriter(w)
		br := bufio.NewReader(w)

		app.WriteTags(bw, i.fm, i.dict)
		bw.Flush()

		str, _ := io.ReadAll(br)

		formatStr := "%s <- %s"
		if sending {
			formatStr = "%s -> %s"
		}

		app.Logger.Trace().Msgf(formatStr, i.prefix, str)
	}
}

func (app *AppMessageLogger) WriteTags(w io.Writer, fieldMap quickfix.FieldMap, dict *datadictionary.DataDictionary) {
	tags := fieldMap.Tags()
	stags := SortableTag(tags)
	sort.Sort(stags)
	for i, tag := range stags {
		tagString := strconv.Itoa(int(tag))
		tagDescription := "<unknown>"
		valueString, _ := fieldMap.GetString(tag)
		valueDescription := ""

		if dict != nil {
			tagField, tok := dict.FieldTypeByTag[int(tag)]
			if tok {
				tagDescription = tagField.Name()
			}
			if len(tagField.Enums) > 0 {
				if en, ok := tagField.Enums[valueString]; ok {
					valueString = en.Value
					valueDescription = strcase.ToCamel(strings.ToLower(en.Description))
					valueString += fmt.Sprintf("(%s)", valueDescription)
				}
			}
		}

		formatStr := "%s(%s)=%s"
		if i < len(tags)-1 {
			formatStr += ","
		}

		w.Write([]byte(fmt.Sprintf(formatStr, tagString, tagDescription, valueString)))
	}
}

func (app *AppMessageLogger) WriteMessageBodyAsTable(w io.Writer, message *quickfix.Message) {
	bodyTags := message.Body.Tags()
	sBodyTags := SortableTag(bodyTags)
	sort.Sort(sBodyTags)

	tw := tabwriter.NewWriter(w, 10, 0, 2, ' ', 0)

	tw.Write([]byte(fmt.Sprintf("TAG\tTYPE\tVALUE\n")))

	var line []string
	for _, tag := range sBodyTags {
		var tagDescription = "<unknown>"
		var valueString = ""
		var valueDescription = ""
		if app.AppDataDictionary != nil {
			tagField, tok := app.AppDataDictionary.FieldTypeByTag[int(tag)]
			if tok {
				tagDescription = tagField.Name()
			}
			valueString, _ = message.Body.GetString(tag)
			if len(tagField.Enums) > 0 {
				if en, ok := tagField.Enums[valueString]; ok {
					valueString = en.Value
					valueDescription = strcase.ToCamel(strings.ToLower(en.Description))
					valueString += fmt.Sprintf(" (%s)", valueDescription)
				}
			}
		}

		line = []string{
			strconv.Itoa(int(tag)),
			tagDescription,
			valueString,
		}

		tw.Write([]byte(fmt.Sprintf("%s\t\n", strings.Join(line, "\t"))))
	}

	tw.Flush()
}
