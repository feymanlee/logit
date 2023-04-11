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
	Name string
	// 日志级别
	LogLevel zapcore.Level
	// CallerSkip，默认值 3
	CallerSkip int
	// 慢请求时间阈值 请求处理时间超过该值则使用 Warn 级别打印日志
	SlowThreshold time.Duration
	// 日志输出路径，默认 []string{"console"}
	// Optional.
	OutputPaths []string
	// 日志初始字段
	// Optional.
	InitialFields map[string]interface{}
	// 是否关闭打印 caller，默认 false
	// Optional.
	DisableCaller bool
	// 是否关闭打印 stack strace，默认 false
	// Optional.
	DisableStacktrace bool
	// 配置日志字段 key 的名称
	// Optional.
	EncoderConfig *zapcore.EncoderConfig
	// lumberjack sink 支持日志文件 rotate
	// Optional.
	LumberjackSink *LumberjackSink
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
	_logger       *zap.Logger
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

// CtxLogger 创建打印日志的 ctx logger
func (g GormLogger) CtxLogger(ctx context.Context) *zap.Logger {
	_, ctxLogger := NewCtxLogger(ctx, g._logger, "")
	return ctxLogger
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
	l := g.CtxLogger(ctx).Named("sql")
	switch {
	case err != nil:
		l.Error("sql trace", zap.String("sql", sql), zap.Float64("latency", latency), zap.Int64("rows", rows), zap.String("error", err.Error()))
	case g.slowThreshold != 0 && latency > g.slowThreshold.Seconds():
		l.Warn("sql trace[slow]", zap.String("sql", sql), zap.Float64("latency", latency), zap.Int64("rows", rows), zap.Float64("threshold", g.slowThreshold.Seconds()))
	default:
		l.Info("sql trace", zap.String("sql", sql), zap.Float64("latency", latency), zap.Int64("rows", rows))
	}
}

//
// NewGormLogger
//  @Description: 创建实现了 gorm logger interface 的 logger
//  @param opt
//  @return GormLogger
//  @return error
//
func NewGormLogger(opt GormLoggerOptions) (GormLogger, error) {
	l := GormLogger{
		name:          GormLoggerName,
		callerSkip:    GormLoggerCallerSkip,
		logLevel:      opt.LogLevel,
		slowThreshold: opt.SlowThreshold,
	}
	if opt.Name != "" {
		l.name = opt.Name
	}
	if opt.CallerSkip != 0 {
		l.callerSkip = opt.CallerSkip
	}
	var err error
	l._logger, err = NewLogger(Options{
		Level:             "debug",
		Format:            "json",
		OutputPaths:       opt.OutputPaths,
		InitialFields:     opt.InitialFields,
		DisableCaller:     opt.DisableCaller,
		DisableStacktrace: opt.DisableStacktrace,
		EncoderConfig:     opt.EncoderConfig,
		LumberjackSink:    opt.LumberjackSink,
	})
	l._logger = l._logger.Named(l.name)
	return l, err
}
