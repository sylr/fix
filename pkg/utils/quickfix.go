package utils

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/quickfixgo/quickfix"
	"github.com/rs/zerolog"

	"sylr.dev/fix/pkg/dump"
	"sylr.dev/fix/pkg/dump/dict"
)

var (
	humanPrinters map[string]dump.HumanLogPrinter
	dumpers       map[string]dump.Dumper
)

func LogMessage(logger zerolog.Logger, level zerolog.Level, message *quickfix.Message, fixVersion string) {
	var writer bytes.Buffer

	var ok bool
	var humanPrinter dump.HumanLogPrinter
	var dumper dump.Dumper

	if humanPrinter, ok = humanPrinters[fixVersion]; !ok {
		humanPrinter = dump.HumanLogPrinter{
			CommonPrinter:  dump.CommonPrinter{Indent: "  ", Writer: &writer},
			FixDictDefault: dict.Get(fixVersion),
		}
		dumper = dump.NewDumper(false, &humanPrinter)
	} else {
		dumper = dumpers[fixVersion]
	}

	reader := bytes.NewBufferString(message.String())
	dumper.Dump(reader)

	for {
		line, err := writer.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				break
			} else if len(line) == 0 {
				break
			}
		}
		logger.WithLevel(level).Msg(strings.TrimRight(line, "\n "))
	}
}

func MessageTable(w io.Writer, message *quickfix.Message, fixVersion string) {
	bodyTags := message.Body.Tags()
	tw := tabwriter.NewWriter(w, 10, 0, 2, ' ', 0)
	dic := dict.Get(fixVersion)

	tw.Write([]byte(fmt.Sprintf("FIX TAG\tDESCRIPTION\tVALUE\t\n")))

	var line []string
	for i := 0; i < len(bodyTags); i++ {
		tagString, tok := dic.TagName(int(bodyTags[i]))
		valueString, _ := message.Body.GetString(bodyTags[i])
		valueName, vok := dic.ValueName(int(bodyTags[i]), valueString)

		line = []string{strconv.Itoa(int(bodyTags[i]))}

		if tok {
			line = append(line, tagString)
		} else {
			line = append(line, "<unknown>")
		}

		if vok {
			line = append(line, fmt.Sprintf("%s(%s)", valueString, valueName))
		} else {
			line = append(line, valueString)
		}

		tw.Write([]byte(fmt.Sprintf("%s\t\n", strings.Join(line, "\t"))))
	}

	tw.Flush()
}
