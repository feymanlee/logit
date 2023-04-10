/**
 * @Author: feymanlee@gmail.com
 * @Description:
 * @File:  redis
 * @Date: 2023/4/6 18:06
 */

package logit

import (
	"context"
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
	ctxRedisStartKey = "_log_redis_start_"
)

type RedisLoggerOptions struct {
	Name       string
	CallerSkip int
}

type RedisLogger struct {
	name string
	// 日志级别
	logLevel zapcore.Level
	// 指定慢查询时间
	slowThreshold time.Duration
	// Trace 方法打印日志是使用的日志 level
	traceWithLevel zapcore.Level
	callerSkip     int
	_logger        *zap.Logger
}

func NewRedisLogger(opt RedisLoggerOptions) RedisLogger {
	l := RedisLogger{
		name:       defaultRedisLoggerName,
		callerSkip: defaultRedisLoggerCallerSkip,
	}
	if opt.CallerSkip != 0 {
		l.callerSkip = opt.CallerSkip
	}
	if opt.Name != "" {
		l.name = opt.Name
	}
	l._logger = CloneLogger(l.name)
	return l
}

//
// CtxLogger
//  @Description: 创建打印日志的 ctx logger
//  @receiver l
//  @param ctx
//  @return *zap.Logger
//
func (l RedisLogger) CtxLogger(ctx context.Context) *zap.Logger {
	_, ctxLogger := NewCtxLogger(ctx, l._logger, "")
	return ctxLogger.WithOptions(zap.AddCallerSkip(l.callerSkip))
}

//
// BeforeProcess
//  @Description: 实现 go-redis HOOK BeforeProcess 方法
//  @receiver l
//  @param ctx
//  @param cmd
//  @return context.Context
//  @return error
//
func (l RedisLogger) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if gc, ok := ctx.(*gin.Context); ok {
		// set start time in gin.Context
		gc.Set(ctxRedisStartKey, time.Now())
		return ctx, nil
	}
	// set start time in context
	return context.WithValue(ctx, ctxRedisStartKey, time.Now()), nil
}

//
// AfterProcess
//  @Description: 实现 go-redis HOOK AfterProcess 方法
//  @receiver l
//  @param ctx
//  @param cmd
//  @return error
//
func (l RedisLogger) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	logger := l.CtxLogger(ctx)
	cost := l.getCost(ctx)
	if err := cmd.Err(); err != nil {
		logger.Warn("redis trace", zap.String("command", cmd.FullName()), zap.Float64("latency_ms", cost), zap.Error(err))
	} else {
		logger.Info("redis trace", zap.String("command", cmd.FullName()), zap.String("args", cmd.String()), zap.Float64("latency_ms", cost))
	}
	return nil
}

//
// BeforeProcessPipeline
//  @Description:
//  @receiver l
//  @param ctx
//  @param cmds
//  @return context.Context
//  @return error
//
func (l RedisLogger) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

//
// AfterProcessPipeline
//  @Description: 实现 go-redis HOOK AfterProcessPipeline 方法
//  @receiver l
//  @param ctx
//  @param cmds
//  @return error
//
func (l RedisLogger) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	logger := l.CtxLogger(ctx)
	cost := l.getCost(ctx)
	pipelineCmds := make([]string, 0, len(cmds))
	pipelineArgs := make([]string, 0, len(cmds))
	pipelineErrs := make([]error, 0, len(cmds))
	for _, cmd := range cmds {
		pipelineCmds = append(pipelineCmds, cmd.FullName())
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

//
// getCost
//  @Description: 获取命令执行耗时
//  @receiver l
//  @param ctx
//  @return cost
//
func (l RedisLogger) getCost(ctx context.Context) (cost float64) {
	var startTime time.Time
	if gc, ok := ctx.(*gin.Context); ok {
		// set start time in gin.Context
		startTime = gc.GetTime(ctxRedisStartKey)
	} else {
		startTime = ctx.Value(ctxRedisStartKey).(time.Time)
	}
	if !startTime.IsZero() {
		cost = time.Since(startTime).Seconds() * 1e3
	}
	return
}
