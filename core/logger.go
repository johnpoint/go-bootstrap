package core

import "go.uber.org/zap"

type Logger interface {
	Info(msg string, field ...zap.Field)
	Error(msg string, field ...zap.Field)
	Debug(msg string, field ...zap.Field)
}

type DefaultLogger struct {
	l *zap.Logger
}

func (d *DefaultLogger) Debug(msg string, field ...zap.Field) {
	d.l.Debug(msg, field...)
}

func (d *DefaultLogger) Info(msg string, field ...zap.Field) {
	d.l.Info(msg, field...)
}

func (d *DefaultLogger) Error(msg string, field ...zap.Field) {
	d.l.Error(msg, field...)
}

func NewDefaultLogger() Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return &DefaultLogger{
		logger,
	}
}
