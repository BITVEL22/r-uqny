package logging

import (
	"log/slog"
	"testing"
)

func TestInit(t *testing.T) {
	Init(slog.LevelInfo)

	// Verify the logger is usable after Init.
	slog.Info("test log message")
}

func TestInitDebug(t *testing.T) {
	Init(slog.LevelDebug)

	slog.Debug("test debug message")
}

