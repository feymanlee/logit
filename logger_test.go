package logit

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	if baseLogger == nil {
		t.Error("baseLogger is nil")
	}
}

func TestNewLoggerNoParam(t *testing.T) {
	logger, err := NewLogger(Options{})
	if err != nil {
		t.Error(err)
	}
	if logger == nil {
		t.Error("return a nil baseLogger")
	}
	logger.Debug("TestNewLoggerNoParam Debug")
}

func TestNewLogger(t *testing.T) {
	options := Options{
		Name:              "tlogger",
		Level:             "debug",
		Format:            "json",
		OutputPaths:       []string{"stderr"},
		InitialFields:     map[string]interface{}{"service_name": "testing"},
		DisableCaller:     false,
		DisableStacktrace: false,
	}
	logger, err := NewLogger(options)
	if err != nil {
		t.Error(err)
	}
	logger.Debug("TestNewLogger Debug")
	logger.Error("TestNewLogger Error")
}

func TestCloneLogger(t *testing.T) {
	nlogger := CloneLogger("cloned")
	if reflect.DeepEqual(nlogger, baseLogger) {
		t.Error("CloneLogger should not be default baseLogger")
	}
	if &nlogger == &baseLogger {
		t.Error("CloneLogger should not be default baseLogger")
	}
}

func TestSetLevel(t *testing.T) {
	baseLogger.Debug("TestChangeLevel raw debug level")
	t.Log("current level:", atomicLevel.Level())
	atomicLevel.SetLevel(zap.InfoLevel)
	t.Log("new level:", atomicLevel.Level())
	baseLogger.Debug("TestChangeLevel raw debug level should not be logged")
	// reset
	atomicLevel.SetLevel(zap.DebugLevel)
}

func TestTextLevel(t *testing.T) {
	level := TextLevel()
	if level != "debug" {
		t.Error(level)
	}
}
