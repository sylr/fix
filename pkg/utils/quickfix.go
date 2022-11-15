package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/quickfixgo/enum"
	qtag "github.com/quickfixgo/tag"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"sylr.dev/fix/pkg/dict"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/datadictionary"
)

type QuickFixMessagePartSetter interface {
	Set(field quickfix.FieldWriter) *quickfix.FieldMap
}

func QuickFixMessagePartSetString[T quickfix.FieldWriter, T2 ~string](setter QuickFixMessagePartSetter, value T2, f func(T2) T) {
	if len(value) > 0 {
		setter.Set(f(value))
	}
}

func QuickFixMessagePartSetDecimal[T quickfix.FieldWriter, T2 ~string](setter QuickFixMessagePartSetter, value T2, f func(decimal.Decimal, int32) T, scale int32) {
	decimal.NewFromString(string(value))
	if len(value) > 0 {
		setter.Set(f(MustNot(decimal.NewFromString(string(value))), scale))
	}
}

type QuickFixAppMessageLogger struct {
	Logger                  *zerolog.Logger
	TransportDataDictionary *datadictionary.DataDictionary
	AppDataDictionary       *datadictionary.DataDictionary
}

func (app *QuickFixAppMessageLogger) LogMessageType(message *quickfix.Message, sessionID quickfix.SessionID, log string) {
	msgType, err := message.MsgType()
	if err != nil {
		app.Logger.Debug().CallerSkipFrame(1).Msgf("%s %s", sessionID.String(), log)
	} else {
		desc := MapSearch(dict.MessageTypes, enum.MsgType(msgType))
		if desc != nil {
			app.Logger.Debug().CallerSkipFrame(1).Msgf("%s %s %s(%s)", sessionID.String(), log, msgType, *desc)
		} else {
			app.Logger.Debug().CallerSkipFrame(1).Msgf("%s %s %s", sessionID.String(), log, msgType)
		}
	}
}

func (app *QuickFixAppMessageLogger) LogMessage(level zerolog.Level, message *quickfix.Message, sessionID quickfix.SessionID, sending bool) {
	if app.Logger.GetLevel() > level {
		return
	}

	message.Cook()
	sort.Sort(message.Header)

	loggedMessage := quickfix.NewMessage()
	message.CopyInto(loggedMessage)

	if sending {
		// When we are sending messages the quickfix.Message.fields is empty,
		// we need to call ParseMessageWithDataDictionary to get it populated.
		br := bytes.NewBufferString(message.String())
		err := quickfix.ParseMessageWithDataDictionary(loggedMessage, br, app.TransportDataDictionary, app.AppDataDictionary)
		if err != nil {
			app.Logger.Error().Err(err).Msg("")
		}
	}

	fields := loggedMessage.GetFields()

	parts := []struct {
		prefix string
		fm     quickfix.FieldMap
		dict   *datadictionary.DataDictionary
	}{
		{"Headers: ", loggedMessage.Header.FieldMap, app.TransportDataDictionary},
		{"Body:    ", loggedMessage.Body.FieldMap, app.AppDataDictionary},
		{"Trailers:", loggedMessage.Trailer.FieldMap, app.TransportDataDictionary},
	}
	for _, part := range parts {
		if len(part.fm.Tags()) == 0 {
			continue
		}

		w := bytes.NewBuffer([]byte{})
		bw := bufio.NewWriter(w)
		br := bufio.NewReader(w)

		skipped := 0
		for i, field := range fields {
			if !part.fm.Has(field.Tag()) {
				skipped++
				continue
			}

			if i-skipped > 0 {
				bw.Write([]byte(","))
			}

			app.WriteField(bw, field)

			bw.Flush()
		}

		str, _ := io.ReadAll(br)

		formatStr := "<- %s %s"
		if sending {
			formatStr = "-> %s %s"
		}

		app.Logger.WithLevel(level).Msgf(formatStr, part.prefix, str)
	}
}

func (app *QuickFixAppMessageLogger) WriteField(w io.Writer, field quickfix.TagValue) {
	fieldTag := field.Tag()
	fieldTagDescription := "<unknown>"
	fieldValue := field.Value()

	if fieldTag == qtag.Password {
		fieldValue = "<redacted>"
	}

	if app.AppDataDictionary != nil {
		tagField, tok := app.AppDataDictionary.FieldTypeByTag[int(fieldTag)]
		if tok {
			fieldTagDescription = tagField.Name()
		}

		if tagField != nil && len(tagField.Enums) > 0 {
			if en, ok := tagField.Enums[fieldValue]; ok {
				fieldValue += fmt.Sprintf("(%s)", en.Description)
			}
		}

		w.Write([]byte(fmt.Sprintf("%d(%s)=%s", fieldTag, fieldTagDescription, fieldValue)))
	} else {
		w.Write([]byte(fmt.Sprintf("%d=%s", fieldTag, fieldValue)))
	}
}

func (app *QuickFixAppMessageLogger) WriteMessageBodyAsTable(w io.Writer, message *quickfix.Message) {
	bodyTags := message.Body.Tags()

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"TAG", "DESCRIPTION", "VALUES"})
	table.SetBorders(tablewriter.Border{false, false, false, true})
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("-")
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})

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

		table.Append(line)
	}

	table.Render()
}

func MapSearch[K comparable, V comparable](m map[K]V, search V) *K {
	for k, v := range m {
		if search == v {
			return &k
		}
	}

	return nil
}
