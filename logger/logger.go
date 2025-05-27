package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func loggerSettings() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-4s |", i))
	}

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}
	zerolog.ErrorStackMarshaler = printedMarshalStack
	log := zerolog.New(output).With().Timestamp().Logger()
	log = log.With().Caller().Logger()
	log = log.With().Stack().Logger()
}

func printedMarshalStack(err error) any {
	fmt.Printf("%+v\n", err)

	return "up"
}
