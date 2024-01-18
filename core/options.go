package core

import (
	"context"
	"io"
	"log/slog"
)

type BootOption func(*Helper)

func LogOutput(writer io.Writer) BootOption {
	return func(helper *Helper) {
		helper.logWriter = writer
	}
}

func SetLoggerType(loggerType LoggerType) BootOption {
	return func(helper *Helper) {
		helper.loggerType = loggerType
	}
}

func WithComponents(components ...Component) BootOption {
	return func(helper *Helper) {
		helper.components = append(helper.components, components...)
	}
}

func Level(level slog.Level) BootOption {
	return func(helper *Helper) {
		helper.level = level
	}
}

func WithContext(ctx context.Context) BootOption {
	return func(helper *Helper) {
		helper.ctx = ctx
	}
}
