// 默认的 logit 全局开箱即用的方法（如： logit.Debug , logit.Debugf 等）都是使用默认 logger 执行的，
// 再使用 ReplaceLogger 方法替换默认 logger 为新的 logger 来解决。

package main

import (
	"github.com/feymanlee/logit"
)

func main() {
	logit.Info(nil, "aaaa")
	// 默认 logger 输出到 stderr，不会输出日志到文件
	// Output:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.661457","logger":"logit.ctx_logger","msg":"aaaa","pid":73847}
	// 创建一个支持 lumberjack 的 logger
	options := logit.Options{
		Name:        "replacedLogger",
		OutputPaths: []string{"stderr", "lumberjack:"},
	}
	logger, _ := logit.NewLogger(options)
	// 替换默认 logger
	resetLogger := logit.ReplaceLogger(logger)

	logit.Error(nil, "ReplaceLogger")
	// Output并保存到文件:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.661927","logger":"replacedLogger.ctx_logger","caller":"logit/global.go:Error:166","msg":"ReplaceLogger","pid":73847,"stacktrace":"github.com/axiaoxin-com/logit.Error\n\t/Users/ashin/go/src/logit/global.go:166\nmain.main\n\t/Users/ashin/go/src/logit/example/replace.go:30\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}

	// 恢复为默认 logger
	resetLogger()

	// 全局方法将恢复使用原始的 logger
	logit.Error(nil, "ResetLogger")
	// Output:
	// {"level":"ERROR","time":"2020-04-15 20:09:23.742995","logger":"logit.ctx_logger","msg":"ResetLogger","pid":73847}
}
