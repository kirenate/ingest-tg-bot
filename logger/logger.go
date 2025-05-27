package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func MakeLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-4s |", i))
	}

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}
	zerolog.ErrorStackMarshaler = printedMarshalStack
	logger := zerolog.New(output).With().Timestamp().Logger()
	logger = logger.With().Caller().Logger()
	logger = logger.With().Stack().Logger()
	log.Logger = logger

	return logger
}

func printedMarshalStack(err error) any {
	fmt.Printf("%+v\n", err)

	return "up"
}

func init() {
	log.Logger = MakeLogger()
}
