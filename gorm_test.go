package logit

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCtxGormLogger(t *testing.T) {
	logger := NewGormLogger(GormLoggerOptions{
		Name:             "gorm",
		LoggerCallerSkip: 3,
		LogLevel:         zapcore.InfoLevel,
		SlowThreshold:    5 * time.Second,
		TraceWithLevel:   zap.DebugLevel,
	})
	if logger == (GormLogger{}) {
		t.Error("CtxGormLogger return empty GormLogger")
	}
}
