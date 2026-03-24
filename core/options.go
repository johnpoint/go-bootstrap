package core

import (
	"context"
	"io"
	"log/slog"
)

type BootOption func(*Helper)

// LogOutput returns a BootOption that sets the output writer for logs.
func LogOutput(writer io.Writer) BootOption {
	return func(helper *Helper) {
		helper.logWriter = writer
	}
}

// SetLoggerType returns a BootOption that sets the logger type (text or JSON).
func SetLoggerType(loggerType LoggerType) BootOption {
	return func(helper *Helper) {
		helper.loggerType = loggerType
	}
}

// WithComponents returns a BootOption that adds components to be initialized.
func WithComponents(components ...Component) BootOption {
	return func(helper *Helper) {
		helper.components = append(helper.components, components...)
	}
}

// Level returns a BootOption that sets the log level.
func Level(level slog.Level) BootOption {
	return func(helper *Helper) {
		helper.level = level
	}
}

// WithContext returns a BootOption that sets the context for the boot helper.
func WithContext(ctx context.Context) BootOption {
	return func(helper *Helper) {
		helper.ctx = ctx
	}
}
