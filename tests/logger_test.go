package test

import (
	"fmt"
	"os"
	"strings"
	te "testing"
	"time"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
)

func TestLoggerConfig(t *te.T) {
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

	err := errors.New("error message")
	logger.Error().Stack().Err(err).Msg("")
	logger.Info().Msg("ok")
	t.Fail()

	// ////-----------------------------------

}

func printedMarshalStack(err error) any {
	fmt.Printf("%+v\n", err)

	return "up"
}
