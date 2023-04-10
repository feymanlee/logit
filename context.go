// context 中不能使用 global 中的方法打印日志， global 会调用 context 的方法，会陷入循环

package logit

import (
	"context"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

// CtxKey context key 类型
type CtxKey string

const (
	// CtxLoggerName define the ctx baseLogger name
	CtxLoggerName CtxKey = "ctx_logger"
	// TraceIDKeyName define the trace id key name
	TraceIDKeyName CtxKey = "trace_id"
)

//
// CtxLogger
//  @Description: get the ctxLogger in context
//  @param c
//  @param fields
//  @return *zap.Logger
//
func CtxLogger(c context.Context, fields ...zap.Field) *zap.Logger {
	if c == nil {
		c = context.Background()
	}
	var ctxLoggerItf interface{}
	if gc, ok := c.(*gin.Context); ok {
		ctxLoggerItf, _ = gc.Get(string(CtxLoggerName))
	} else {
		ctxLoggerItf = c.Value(CtxLoggerName)
	}

	var ctxLogger *zap.Logger
	if ctxLoggerItf != nil {
		ctxLogger = ctxLoggerItf.(*zap.Logger)
	} else {
		_, ctxLogger = NewCtxLogger(c, CloneLogger(string(CtxLoggerName)), CtxTraceID(c))
	}

	if len(fields) > 0 {
		ctxLogger = ctxLogger.With(fields...)
	}
	return ctxLogger
}

// CtxTraceID get trace id from context
// Modify TraceIDPrefix change the prefix
func CtxTraceID(c context.Context) string {
	if c == nil {
		c = context.Background()
	}
	// first get from gin context
	if gc, ok := c.(*gin.Context); ok {
		if traceID := gc.GetString(string(TraceIDKeyName)); traceID != "" {
			return traceID
		}
		if traceID := gc.Query(string(TraceIDKeyName)); traceID != "" {
			return traceID
		}
		if traceID := jsoniter.Get(GetGinRequestBody(gc), string(TraceIDKeyName)).ToString(); traceID != "" {
			return traceID
		}

	} else {
		// get from go context
		traceIDItf := c.Value(TraceIDKeyName)
		if traceIDItf != nil {
			return traceIDItf.(string)
		}
	}
	// return default value
	return xid.New().String()
}

//
// NewCtxLogger
//  @Description: return a context with baseLogger and trace id and a baseLogger with trace id
//  @param c
//  @param logger
//  @param traceID
//  @return context.Context
//  @return *zap.Logger
//
func NewCtxLogger(c context.Context, logger *zap.Logger, traceID string) (context.Context, *zap.Logger) {
	if c == nil {
		c = context.Background()
	}
	if traceID == "" {
		traceID = CtxTraceID(c)
	}
	ctxLogger := logger.With(zap.String(string(TraceIDKeyName), traceID))
	if gc, ok := c.(*gin.Context); ok {
		// set traceID in gin.Context
		gc.Set(string(TraceIDKeyName), traceID)
	}
	// set traceID in context.Context
	c = context.WithValue(c, TraceIDKeyName, traceID)
	return c, ctxLogger
}
