package log

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger_Error(t *testing.T) {
	Error("test", zap.String("info", "test"))
}

func TestLogger_Info(t *testing.T) {
	Info("test", zap.String("info", "test"))
}

func TestOverrideLoggerWithOption(t *testing.T) {
	OverrideLoggerWithOption(map[string]interface{}{
		"service-name": "test_service",
	}, WithJSONEncoding(), WrapLevelEncoder(zapcore.CapitalLevelEncoder))
	Info("test", zap.String("info", "test"))
	Error("test", zap.String("info", "test"))
}
