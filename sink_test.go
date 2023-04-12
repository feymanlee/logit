/**
 * @Author: lifameng@changba.com
 * @Description:
 * @File:  sink_test
 * @Date: 2023/4/12 09:57
 */

package logit

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLumberjackSink(t *testing.T) {
	scheme := "lumberjack"
	filename := "test.log"
	maxAge := 1
	maxBackups := 2
	maxSize := 1

	lumberjackSink := NewLumberjackSink(filename, maxAge, maxBackups, maxSize, true, true)

	err := RegisterSink(scheme, lumberjackSink)
	if err != nil {
		t.Fatalf("failed to register lumberjack sink: %v", err)
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.OutputPaths = []string{scheme + "://" + filename}

	logger, err := config.Build()
	if err != nil {
		t.Fatalf("failed to build logger: %v", err)
	}

	logger.Info("Hello, lumberjack!")

	// Clean up the test log file
	err = os.Remove(filename)
	if err != nil {
		t.Fatalf("failed to remove test log file: %v", err)
	}
}
