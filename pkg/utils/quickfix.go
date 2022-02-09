package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/iancoleman/strcase"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/datadictionary"
	"github.com/rs/zerolog"
)

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

		app.WriteTags(bw, i.fm)
		bw.Flush()

		str, _ := io.ReadAll(br)

		formatStr := "%s <- %s"
		if sending {
			formatStr = "%s -> %s"
		}

		app.Logger.Trace().Msgf(formatStr, i.prefix, str)
	}
}

func (app *AppMessageLogger) WriteTags(w io.Writer, fieldMap quickfix.FieldMap) {
	tags := fieldMap.Tags()

	for i, tag := range tags {
		tagString := strconv.Itoa(int(tag))
		tagDescription := "<unknown>"
		values, _ := fieldMap.GetStrings(tag)
		valueDescription := ""

		for j := range values {
			if app.AppDataDictionary != nil {
				tagField, tok := app.AppDataDictionary.FieldTypeByTag[int(tag)]
				if tok {
					tagDescription = tagField.Name()
				}
				if len(tagField.Enums) > 0 {
					if en, ok := tagField.Enums[values[j]]; ok {
						valueDescription = strcase.ToCamel(strings.ToLower(en.Description))
						values[j] += fmt.Sprintf("(%s)", valueDescription)
					}
				}

				formatStr := "%s(%s)=%s"
				if i < len(tags)-1 || j < len(values)-1 {
					formatStr += ","
				}

				w.Write([]byte(fmt.Sprintf(formatStr, tagString, tagDescription, values[j])))
			} else {
				formatStr := "%s=%s"
				if i < len(tags)-1 || j < len(values)-1 {
					formatStr += ","
				}

				w.Write([]byte(fmt.Sprintf(formatStr, tagString, values[j])))
			}
		}
	}
}

func (app *AppMessageLogger) WriteMessageBodyAsTable(w io.Writer, message *quickfix.Message) {
	bodyTags := message.Body.Tags()

	tw := tabwriter.NewWriter(w, 10, 0, 2, ' ', 0)
	tw.Write([]byte(fmt.Sprintf("TAG\tTYPE\tVALUE\n")))

	var line []string
	for _, tag := range bodyTags {
		var tagDescription = "<unknown>"
		var valueString = ""
		var valueDescription = ""

		values, _ := message.Body.GetStrings(tag)

		for i := range values {
			if app.AppDataDictionary != nil {
				tagField, tok := app.AppDataDictionary.FieldTypeByTag[int(tag)]
				if tok {
					tagDescription = tagField.Name()
				}
				if len(tagField.Enums) > 0 {
					if en, ok := tagField.Enums[values[i]]; ok {
						valueDescription = strcase.ToCamel(strings.ToLower(en.Description))
						values[i] += fmt.Sprintf(" (%s)", valueDescription)
					}
				}
			}
		}

		valueString = strings.Join(values, ", ")
		line = []string{
			strconv.Itoa(int(tag)),
			tagDescription,
			valueString,
		}

		tw.Write([]byte(fmt.Sprintf("%s\t\n", strings.Join(line, "\t"))))
	}

	tw.Flush()
}
