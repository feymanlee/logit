// 默认的 logit 全局开箱即用的方法（如： logit.Debug , logit.Debugf 等）都是使用默认 logger 执行的，
// 默认 logger 不支持 Sentry 和输出日志到文件，可以通过创建一个新的 logger，
// 再使用 ReplaceLogger 方法替换默认 logger 为新的 logger 来解决。

package main

import (
	"os"

	"github.com/feymanlee/logit"
)

func main() {
	// 默认使用全局方法不会保存到文件和上报 Sentry
	logit.Error(nil, "default logger no sentry and file")
	// Output:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.661457","logger":"logit.ctx_logger","msg":"default logger no sentry and file","pid":73847}

	// 创建一个支持 sentry 和 lumberjack 的 logger
	sentryClient, _ := logit.NewSentryClient(os.Getenv("dsn"), true)
	options := logit.Options{
		Name:           "replacedLogger",
		OutputPaths:    []string{"stderr", "lumberjack:"},
		LumberjackSink: logit.NewLumberjackSink("lumberjack", "/tmp/replace.log", 1, 1, 10, true, true),
		SentryClient:   sentryClient,
	}
	logger, _ := logit.NewLogger(options)
	// 替换默认 logger
	resetLogger := logit.ReplaceLogger(logger)

	// 全局方法将使用新的 logger，上报 sentry 并输出到文件
	logit.Error(nil, "ReplaceLogger")
	// Output并保存到文件:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.661927","logger":"replacedLogger.ctx_logger","caller":"logit/global.go:Error:166","msg":"ReplaceLogger","pid":73847,"stacktrace":"github.com/axiaoxin-com/logit.Error\n\t/Users/ashin/go/src/logit/global.go:166\nmain.main\n\t/Users/ashin/go/src/logit/example/replace.go:30\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}

	// 重置默认 logger
	resetLogger()

	// 全局方法将恢复使用原始的 logger，不再上报 sentry 和输出到文件
	logit.Error(nil, "ResetLogger")
	// Output:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.742995","logger":"logit.ctx_logger","msg":"ResetLogger","pid":73847}
}
