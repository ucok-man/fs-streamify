package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func New(level string, env string) (*zerolog.Logger, error) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	var writer io.Writer = os.Stdout
	if env == "development" {
		writer = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = os.Stdout
			w.TimeFormat = "02 Jan, 2006 15:04:05"
			w.FormatCaller = func(i interface{}) string {
				return filepath.Base(fmt.Sprintf("%s", i))
			}
		})
	}

	logger := zerolog.
		New(writer).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Logger()

	return &logger, nil
}
