package logging

import (
	"log/slog"
	"os"
)

// Init configures the default structured logger with the given level.
func Init(level slog.Level) {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})

	slog.SetDefault(slog.New(handler))
}

