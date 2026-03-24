package core

import (
	"log/slog"
	"os"
)

// NewDefaultLogger creates a new default logger with debug level and text handler output to stderr.
func NewDefaultLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
