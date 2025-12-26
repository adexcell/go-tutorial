package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

func New(level string, jsonFormat bool) zerolog.Logger {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		l = zerolog.InfoLevel
	}

	var output io.Writer

	if jsonFormat {
		output = os.Stdout
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: zerolog.TimeFormatUnix,
		}
	}

	logger := zerolog.New(output).
		Level(l).
		With().
		Timestamp().
		Caller(). // добавляет файл и строку, где вызван лог
		Logger()

	return logger
}
