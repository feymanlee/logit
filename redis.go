/**
 * @Author: feymanlee@gmail.com
 * @Description:
 * @File:  redis
 * @Date: 2023/4/6 18:06
 */

package logit

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	//  默认 logger 的名称
	defaultRedisLoggerName = "redis"
	// 默认 caller skip
	defaultRedisLoggerCallerSkip = 4
	// 上下文中保存开始时间的 key
	ctxRedisStartKey CtxKey = "_log_redis_start_"
	// 默认的 redis 慢查询时间，30ms
	defaultSlowThreshold = time.Millisecond * 30
)

type RedisLoggerOptions struct {
	Name string
	// CallerSkip，默认值 4
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
	// nil err level
	NilErrLevel string
}

type RedisLogger struct {
	name string
	// 指定慢查询时间
	slowThreshold time.Duration
	callerSkip    int
	_logger       *zap.Logger
	nilErrLevel   string
}

func NewRedisLogger(opt RedisLoggerOptions) (RedisLogger, error) {
	l := RedisLogger{
		name:          defaultRedisLoggerName,
		callerSkip:    defaultRedisLoggerCallerSkip,
		slowThreshold: defaultSlowThreshold,
		nilErrLevel:   opt.NilErrLevel,
	}
	if opt.CallerSkip != 0 {
		l.callerSkip = opt.CallerSkip
	}
	if opt.Name != "" {
		l.name = opt.Name
	}
	if opt.SlowThreshold > 0 {
		l.slowThreshold = opt.SlowThreshold
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
	})
	l._logger = l._logger.Named(l.name)
	return l, err
}

// CtxLogger
//
//	@Description: 创建打印日志的 ctx logger
//	@receiver l
//	@param ctx
//	@return *zap.Logger
func (l RedisLogger) CtxLogger(ctx context.Context) *zap.Logger {
	_, ctxLogger := NewCtxLogger(ctx, l._logger, "")
	return ctxLogger.WithOptions(zap.AddCallerSkip(l.callerSkip))
}

// BeforeProcess
//
//	@Description: 实现 go-redis HOOK BeforeProcess 方法
//	@receiver l
//	@param ctx
//	@param cmd
//	@return context.Context
//	@return error
func (l RedisLogger) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if gc, ok := ctx.(*gin.Context); ok {
		// set start time in gin.Context
		gc.Set(string(ctxRedisStartKey), time.Now())
		return ctx, nil
	}
	// set start time in context
	return context.WithValue(ctx, ctxRedisStartKey, time.Now()), nil
}

// AfterProcess
//
//	@Description: 实现 go-redis HOOK AfterProcess 方法
//	@receiver l
//	@param ctx
//	@param cmd
//	@return error
func (l RedisLogger) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	logger := l.CtxLogger(ctx)
	cost := l.getCost(ctx)
	if err := cmd.Err(); err != nil {
		level := zap.ErrorLevel
		if errors.Is(err, redis.Nil) {
			var err1 error
			if level, err1 = zapcore.ParseLevel(l.nilErrLevel); err1 != nil {
				level = zap.ErrorLevel
			}
		}
		logger.Log(level, "redis trace", zap.String("command", cmd.FullName()), zap.String("args", cmd.String()), zap.Float64("latency_ms", cost), zap.Error(err))
	} else {
		log := logger.Info
		if cost > float64(l.slowThreshold) {
			log = logger.Warn
		}
		log("redis trace", zap.String("command", cmd.FullName()), zap.String("args", cmd.String()), zap.Float64("latency_ms", cost))
	}
	return nil
}

// BeforeProcessPipeline
//
//	@Description:
//	@receiver l
//	@param ctx
//	@param cmds
//	@return context.Context
//	@return error
func (l RedisLogger) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if gc, ok := ctx.(*gin.Context); ok {
		// set start time in gin.Context
		gc.Set(string(ctxRedisStartKey), time.Now())
		return ctx, nil
	}
	// set start time in context
	return context.WithValue(ctx, ctxRedisStartKey, time.Now()), nil
}

// AfterProcessPipeline
//
//	@Description: 实现 go-redis HOOK AfterProcessPipeline 方法
//	@receiver l
//	@param ctx
//	@param cmds
//	@return error
func (l RedisLogger) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	logger := l.CtxLogger(ctx)
	cost := l.getCost(ctx)
	pipelineArgs := make([]string, 0, len(cmds))
	pipelineErrs := make([]error, 0, len(cmds))
	for _, cmd := range cmds {
		pipelineArgs = append(pipelineArgs, cmd.String())
		if err := cmd.Err(); err != nil {
			pipelineErrs = append(pipelineErrs, err)
		}
	}
	if len(pipelineErrs) > 0 {
		logger.Warn("redis trace", zap.Any("args", pipelineArgs), zap.Bool("pipeline", true), zap.Float64("latency_ms", cost), zap.Errors("errors", pipelineErrs))
	} else {
		logger.Info("redis trace", zap.Any("args", pipelineArgs), zap.Bool("pipeline", true), zap.Float64("latency_ms", cost))
	}
	return nil
}

// getCost
//
//	@Description: 获取命令执行耗时
//	@receiver l
//	@param ctx
//	@return cost
func (l RedisLogger) getCost(ctx context.Context) (cost float64) {
	var startTime time.Time
	if gc, ok := ctx.(*gin.Context); ok {
		// set start time in gin.Context
		startTime = gc.GetTime(string(ctxRedisStartKey))
	} else {
		startTime = ctx.Value(ctxRedisStartKey).(time.Time)
	}
	if !startTime.IsZero() {
		cost = time.Since(startTime).Seconds() * 1e3
	}
	return
}
