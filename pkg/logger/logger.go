package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	// JSON-логи + timestamp для ELK/Grafana
	return zerolog.New(os.Stderr).With().Timestamp().Logger()
}
