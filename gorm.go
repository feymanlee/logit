// gorm v2

package logit

import (
	"context"
	"time"

	"github.com/axiaoxin-com/goutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

const (
	// GormLoggerName gorm baseLogger 名称
	GormLoggerName = "gorm"
	// GormLoggerCallerSkip caller skip
	GormLoggerCallerSkip = 3
)

type GormLoggerOptions struct {
	Name             string
	LoggerCallerSkip int
	// 日志级别
	LogLevel zapcore.Level
	// 指定慢查询时间
	SlowThreshold time.Duration
	// Trace 方法打印日志是使用的日志 level
	TraceWithLevel zapcore.Level
}

// GormLogger 使用 zap 来打印 gorm 的日志
// 初始化时在内部的 baseLogger 中添加 trace id 可以追踪 sql 执行记录
type GormLogger struct {
	name       string
	callerSkip int
	// 日志级别
	logLevel zapcore.Level
	// 指定慢查询时间
	slowThreshold time.Duration
	// Trace 方法打印日志是使用的日志 level
	traceWithLevel zapcore.Level
	_logger        *zap.Logger
}

var gormLogLevelMap = map[gormlogger.LogLevel]zapcore.Level{
	gormlogger.Info:  zap.InfoLevel,
	gormlogger.Warn:  zap.WarnLevel,
	gormlogger.Error: zap.ErrorLevel,
}

// LogMode 实现 gorm baseLogger 接口方法
func (g GormLogger) LogMode(gormLogLevel gormlogger.LogLevel) gormlogger.Interface {
	level, exists := gormLogLevelMap[gormLogLevel]
	if !exists {
		level = zap.DebugLevel
	}
	newLogger := g
	newLogger.logLevel = level
	return &newLogger
}

// CtxLogger 创建打印日志的 ctxlogger
func (g GormLogger) CtxLogger(ctx context.Context) *zap.Logger {
	_, ctxLogger := NewCtxLogger(ctx, g._logger, "")
	return ctxLogger.WithOptions(zap.AddCallerSkip(g.callerSkip))
}

// Info 实现 gorm baseLogger 接口方法
func (g GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.InfoLevel {
		g.CtxLogger(ctx).Sugar().Infof(msg, data...)
	}
}

// Warn 实现 gorm baseLogger 接口方法
func (g GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.WarnLevel {
		g.CtxLogger(ctx).Sugar().Warnf(msg, data...)
	}
}

// Error 实现 gorm baseLogger 接口方法
func (g GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.ErrorLevel {
		g.CtxLogger(ctx).Sugar().Errorf(msg, data...)
	}
}

// Trace 实现 gorm baseLogger 接口方法
func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	now := time.Now()
	latency := now.Sub(begin).Seconds()
	sql, rows := fc()
	sql = goutils.RemoveDuplicateWhitespace(sql, true)
	l := g.CtxLogger(ctx)
	switch {
	case err != nil:
		l.Named("sql").Error("sql trace", zap.String("sql", sql), zap.Float64("latency", latency), zap.Int64("rows", rows), zap.String("error", err.Error()))
	case g.slowThreshold != 0 && latency > g.slowThreshold.Seconds():
		l.Named("sql").Warn("sql trace[slow]", zap.String("sql", sql), zap.Float64("latency", latency), zap.Int64("rows", rows), zap.Float64("threshold", g.slowThreshold.Seconds()))
	default:
		log := l.Debug
		if g.traceWithLevel == zap.InfoLevel {
			log = l.Info
		} else if g.traceWithLevel == zap.WarnLevel {
			log = l.Warn
		} else if g.traceWithLevel == zap.ErrorLevel {
			log = l.Error
		}
		log("sql trace", zap.String("sql", sql), zap.Float64("latency", latency), zap.Int64("rows", rows))
	}
}

// NewGormLogger 返回带 zap baseLogger 的 GormLogger
func NewGormLogger(opt GormLoggerOptions) GormLogger {
	l := GormLogger{
		name:           GormLoggerName,
		callerSkip:     GormLoggerCallerSkip,
		logLevel:       opt.LogLevel,
		slowThreshold:  opt.SlowThreshold,
		traceWithLevel: opt.TraceWithLevel,
	}
	if opt.Name != "" {
		l.name = opt.Name
	}
	if opt.LoggerCallerSkip != 0 {
		l.callerSkip = opt.LoggerCallerSkip
	}
	l._logger = CloneLogger(l.name)
	return l
}
