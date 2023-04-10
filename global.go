// 开箱即用的方法
// 第一个参数为 context ，会尝试从其中获取带 trace id 的 baseLogger 进行打印

package logit

import (
	"context"

	"go.uber.org/zap"
)

// Debugs 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来记录 debug 级别的日志
// logit.Debugs(nil, "abc", 123)
func Debugs(c context.Context, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Debug(args...)
}

// Infos 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来记录 info 级别的日志
func Infos(c context.Context, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Info(args...)
}

// Warns 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来记录 warn 级别的日志
func Warns(c context.Context, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Warn(args...)
}

// Errors 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来记录 Error 级别的日志
func Errors(c context.Context, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Error(args...)
}

// Panics 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来记录 Panic 级别的日志
func Panics(c context.Context, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Panic(args...)
}

// Fatals 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来记录 Fatal 级别的日志
func Fatals(c context.Context, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Fatal(args...)
}

// Debugf 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来模板字符串记录 debug 级别的日志
// logit.Debugf(nil, "str:%s", "abd")
func Debugf(c context.Context, template string, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Debugf(template, args...)
}

// Infof 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来模板字符串记录 info 级别的日志
func Infof(c context.Context, template string, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Infof(template, args...)
}

// Warnf 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来模板字符串记录 warn 级别的日志
func Warnf(c context.Context, template string, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Warnf(template, args...)
}

// Errorf 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来模板字符串记录 error 级别的日志
func Errorf(c context.Context, template string, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Errorf(template, args...)
}

// Panicf 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来模板字符串记录 panic 级别的日志
func Panicf(c context.Context, template string, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Panicf(template, args...)
}

// Fatalf 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来模板字符串记录 fatal 级别的日志
func Fatalf(c context.Context, template string, args ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Fatalf(template, args...)
}

// Debugw 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来 kv 记录 debug 级别的日志
// logit.Debugw(nil, "msg", "k1", "v1", "k2", "v2")
func Debugw(c context.Context, msg string, keysAndValues ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Debugw(msg, keysAndValues...)
}

// Infow 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来 kv 记录 info 级别的日志
func Infow(c context.Context, msg string, keysAndValues ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Infow(msg, keysAndValues...)
}

// Warnw 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来 kv 记录 warn 级别的日志
func Warnw(c context.Context, msg string, keysAndValues ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Warnw(msg, keysAndValues...)
}

// Errorw 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来 kv 记录 error 级别的日志
func Errorw(c context.Context, msg string, keysAndValues ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Errorw(msg, keysAndValues...)
}

// Panicw 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来 kv 记录 panic 级别的日志
func Panicw(c context.Context, msg string, keysAndValues ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Panicw(msg, keysAndValues...)
}

// Fatalw 尝试从 Context 中获取带 trace id 的 sugared baseLogger 来 kv 记录 fatal 级别的日志
func Fatalw(c context.Context, msg string, keysAndValues ...interface{}) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().Fatalw(msg, keysAndValues...)
}

// Debug 尝试从 Context 中获取带 trace id 的 baseLogger 记录 debug 级别的日志
func Debug(c context.Context, msg string, fields ...zap.Field) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

// Info 尝试从 Context 中获取带 trace id 的 baseLogger 记录 info 级别的日志
func Info(c context.Context, msg string, fields ...zap.Field) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

// Warn 尝试从 Context 中获取带 trace id 的 baseLogger 记录 warn 级别的日志
func Warn(c context.Context, msg string, fields ...zap.Field) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

// Error 尝试从 Context 中获取带 trace id 的 baseLogger 记录 error 级别的日志
func Error(c context.Context, msg string, fields ...zap.Field) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// Panic 尝试从 Context 中获取带 trace id 的 baseLogger 记录 panic 级别的日志
func Panic(c context.Context, msg string, fields ...zap.Field) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Panic(msg, fields...)
}

// Fatal 尝试从 Context 中获取带 trace id 的 baseLogger 记录 fatal 级别的日志
func Fatal(c context.Context, msg string, fields ...zap.Field) {
	CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

// 获取一个 SugaredLogger
func For(c context.Context, fields ...interface{}) *zap.SugaredLogger {
	return CtxLogger(c).WithOptions(zap.AddCallerSkip(1)).Sugar().With(fields...)
}

func Sync(c context.Context) error {
	return CtxLogger(c).Sync()
}

// ExtraField 顺序传入 kv 对，返回以 extra 为 key ，传入的 kv 对组成的 map 为值的 zap Reflect Field
// 在需要固定日志外层 json 字段有需要添加新字段时可以使用
func ExtraField(keysAndValues ...interface{}) zap.Field {
	fieldMap := map[string]interface{}{}
	for i := 0; i < len(keysAndValues); {
		k, v := keysAndValues[i], keysAndValues[i+1]
		if kstr, ok := k.(string); ok {
			fieldMap[kstr] = v
		}
		i += 2
	}
	return zap.Reflect("extra", fieldMap)
}
