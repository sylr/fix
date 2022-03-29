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

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/datadictionary"
	qtag "github.com/quickfixgo/tag"
	"github.com/rs/zerolog"
)

type QuickFixMessagePartSetter interface {
	Set(field quickfix.FieldWriter) *quickfix.FieldMap
}

func QuickFixMessagePartSet[T quickfix.FieldWriter](setter QuickFixMessagePartSetter, value string, f func(string) T) {
	if len(value) > 0 {
		setter.Set(f(value))
	}
}

type QuickFixAppMessageLogger struct {
	Logger                  *zerolog.Logger
	TransportDataDictionary *datadictionary.DataDictionary
	AppDataDictionary       *datadictionary.DataDictionary
}

func (app *QuickFixAppMessageLogger) WriteMessage(w io.Writer, message *quickfix.Message, sessionID quickfix.SessionID, sending bool) {
	if app.Logger.GetLevel() > zerolog.TraceLevel {
		return
	}

	message.Cook()
	sort.Sort(message.Header)

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

		b := bytes.NewBuffer([]byte{})
		bw := bufio.NewWriter(b)
		br := bufio.NewReader(b)

		app.WriteTags(bw, i.fm)
		bw.Flush()

		str, _ := io.ReadAll(br)

		formatStr := "%s <- %s\n"
		if sending {
			formatStr = "%s -> %s\n"
		}

		fmt.Fprintf(w, formatStr, i.prefix, str)
	}
}

func (app *QuickFixAppMessageLogger) LogMessage(level zerolog.Level, message *quickfix.Message, sessionID quickfix.SessionID, sending bool) {
	if app.Logger.GetLevel() > level {
		return
	}

	message.Cook()
	sort.Sort(message.Header)

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

		app.Logger.WithLevel(level).Msgf(formatStr, i.prefix, str)
	}
}

func (app *QuickFixAppMessageLogger) WriteTags(w io.Writer, fieldMap quickfix.FieldMap) {
	tags := fieldMap.Tags()

	for i, tag := range tags {
		values := fieldMap.Values(tag)
		for j, t := range values {
			tagString := strconv.Itoa(int(tag))
			tagDescription := "<unknown>"
			stringValues, _ := fieldMap.GetStrings(tag)
			valueDescription := ""

			if tag == qtag.Password {
				stringValues[j] = "<redacted>"
			}

			if app.AppDataDictionary != nil {
				tagField, tok := app.AppDataDictionary.FieldTypeByTag[int(t.Tag())]
				if tok {
					tagDescription = tagField.Name()
				}

				if tagField != nil && len(tagField.Enums) > 0 {
					if en, ok := tagField.Enums[stringValues[j]]; ok {
						valueDescription = en.Description
						stringValues[j] += fmt.Sprintf("(%s)", valueDescription)
					}
				}
				formatStr := "%s(%s)=%s"
				if i < len(tags)-1 || j < len(values)-1 {
					formatStr += ","
				}

				w.Write([]byte(fmt.Sprintf(formatStr, tagString, tagDescription, stringValues[j])))
			} else {
				formatStr := "%s=%s"
				if i < len(tags)-1 || j < len(values)-1 {
					formatStr += ","
				}

				w.Write([]byte(fmt.Sprintf(formatStr, tagString, stringValues[j])))
			}
		}
	}
}

func (app *QuickFixAppMessageLogger) WriteMessageBodyAsTable(w io.Writer, message *quickfix.Message) {
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
						valueDescription = en.Description
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
