package logger

import (
	"context"
	"github.com/rs/zerolog"
	"os"
)

func FromContext(ctx context.Context) zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
