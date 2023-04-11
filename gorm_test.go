package logit

import (
	"os"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Product test model
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestGormLogger(t *testing.T) {
	// defer clear
	defer os.Remove("./sqlite3.db")
	logger, err := NewGormLogger(GormLoggerOptions{
		Name:          "gorm",
		CallerSkip:    3,
		LogLevel:      zapcore.InfoLevel,
		SlowThreshold: 5 * time.Second,
		OutputPaths:   []string{"stdout"},
		InitialFields: map[string]interface{}{
			"key1": "value1",
		},
		DisableCaller:     false,
		DisableStacktrace: false,
		EncoderConfig:     &defaultEncoderConfig,
		LumberjackSink:    nil,
	})
	if err != nil {
		t.Errorf("new gorm logger failed: %v", err)
	}
	if logger == (GormLogger{}) {
		t.Error("CtxGormLogger return empty GormLogger")
	}
	// Create gorm db instance
	db, err := gorm.Open(sqlite.Open("./sqlite3.db"), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		t.Error(err)
	}
	// Migrate the schema
	db.AutoMigrate(&Product{})
	var ret Product
	db.Model(Product{}).Find(&ret)
	_, err = db.Table("aa").Rows()
	if err != nil {
		t.Log(err)
	}
}
