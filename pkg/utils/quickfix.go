package utils

import (
	"bytes"
	"io"
	"strings"

	"github.com/quickfixgo/quickfix"
	"github.com/rs/zerolog"
	"sylr.dev/fix/pkg/dump"
	"sylr.dev/fix/pkg/dump/dict"
)

func LogMessage(logger zerolog.Logger, level zerolog.Level, message *quickfix.Message, fixVersion string) {
	var writer bytes.Buffer

	humanPrinter := dump.HumanLogPrinter{
		CommonPrinter:  dump.CommonPrinter{Indent: "  ", Writer: &writer},
		FixDictDefault: dict.Get(fixVersion),
	}

	dumper := dump.NewDumper(false, &humanPrinter)
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
