package main

import (
	"context"

	"github.com/feymanlee/logit"
	"go.uber.org/zap"
)

func main() {
	/* 克隆一个带有初始字段的默认 logger */
	// 初始字段可以不传，克隆的 logger 名称会是 logit.subname ，该 logger 打印的日志都会带上传入的字段
	cloneDefaultLogger := logit.CloneLogger("subname", zap.String("str_field", "field_value"))
	cloneDefaultLogger.Debug("CloneDefaultLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548271","logger":"logit.subname","msg":"CloneDefaultLogger","pid":68701,"str_field":"field_value"}

	/* 使用 Options 创建 logger */
	// 可以直接使用空 Options 创建默认配置项的 logger
	// 不支持 sentry 和 http 动态修改日志级别，日志输出到 stderr
	emptyOptionsLogger, _ := logit.NewLogger(logit.Options{})
	emptyOptionsLogger.Debug("emptyOptionsLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548323","logger":"logit","caller":"example/logging.go:main:48","msg":"emptyOptionsLogger","pid":68701}

	// 配置 Options 创建 logger
	options := logit.Options{
		Name:              "logit",            // logger 名称
		Level:             "debug",            // zap 的 AtomicLevel ， logger 日志级别
		Format:            "json",             // 日志输出格式为 json
		OutputPaths:       []string{"stderr"}, // 日志输出位置为 stderr
		InitialFields:     nil,                // DefaultInitialFields 初始 logger 带有 pid 字段
		DisableCaller:     false,              // 是否打印调用的代码行位置
		DisableStacktrace: false,              // 错误日志是否打印调用栈信息
	}
	optionsLogger, _ := logit.NewLogger(options)
	optionsLogger.Debug("optionsLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548363","logger":"logit","caller":"example/logging.go:main:67","msg":"optionsLogger","pid":68701}

	/* 从 context.Context 或*gin.Context 中获取或创建 logger */
	ctx := context.Background()
	ctxLogger := logit.CtxLogger(ctx, zap.String("field1", "xxx"))
	ctxLogger.Debug("ctxLogger")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:39:37.548414","logger":"logit.ctx_logger","msg":"ctxLogger","pid":68701,"field1":"xxx"}
}
