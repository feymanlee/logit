# logit

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/feymanlee/logit?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/feymanlee/logit)](https://goreportcard.com/report/github.com/feymanlee/logit)
[![Unit-Tests](https://github.com/feymanlee/logit/workflows/Unit-Tests/badge.svg)](https://github.com/feymanlee/logit/actions)
[![Coverage Status](https://coveralls.io/repos/github/feymanlee/logit/badge.svg?branch=main)](https://coveralls.io/github/feymanlee/logit?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/feymanlee/logit.svg)](https://pkg.go.dev/github.com/feymanlee/logit)

logit 简单封装了在日常使用 [zap](https://github.com/uber-go/zap) 打日志时的常用方法。

- 提供快速使用 zap 打印日志的方法，除 zap 的 DPanic 、 DPanicf 方法外所有日志打印方法开箱即用
- 提供多种快速创建 `logger` 的方法
- 支持从 Context 中创建、获取带有 **Trace ID** 的 logger
- 提供 `gin` 的日志中间件，支持通过配置自定义记录 `TraceId` `context keys` `Request Header` `Request Form` `Request Body` `Response Body` 以及其他的 HTTP 请求信息
- 支持 `Gorm`，记录 `TraceId` `请求时间` `SQL` `慢 SQL` `ERR`
- 支持 `go-redis` 记录 `TraceId` `redis 命令` `请求结果` `耗时` `慢请求` `pipline`，目前只支持 `go-redis/v8`, 后续会增加对 `go-redis/v9` 的支持
- 支持将日志保存到文件并自动 rotate
- 支持自定义 logger Encoder 配置

logit 只提供 zap 使用时的常用方法汇总，不是对 zap 进行二次开发，拒绝过度封装。

## 开箱即用

```shell
go get github.com/feymanlee/logit
```

在 `logit` 被 import 时，会生成内部使用的默认 logger 。
默认 logger 使用 JSON 格式打印日志内容到 stderr 。
默认带有初始字段 pid 打印进程 ID 。

开箱即用的方法第一个参数为 context.Context, 可以传入 gin.Context ，会尝试从其中获取 Trace ID 进行日志打印，无需 Trace ID 可以直接传 nil

```go
ctx := context.Background()
/* zap Debug */
logit.Debug(ctx, "Debug message", zap.Int("intType", 123), zap.Bool("boolType", false), zap.Ints("sliceInt", []int{1, 2, 3}), zap.Reflect("map", map[string]interface{}{"i": 1, "s": "s"}))
// Output:
// {"level":"DEBUG","time":"2020-04-15 18:12:11.991006","logger":"logit.ctx_logger","msg":"Debug message","pid":45713,"intType":123,"boolType":false,"sliceInt":[1,2,3],"map":{"i":1,"s":"s"}}

/* zap sugared logger Debug */
logit.Debugs(ctx, "Debugs message", 123, false, []int{1, 2, 3}, map[string]interface{}{"i": 1, "s": "s"})
// Output:
// {"level":"DEBUG","time":"2020-04-15 18:12:11.991239","logger":"logit.ctx_logger","msg":"Debugs message123 false [1 2 3] map[i:1 s:s]","pid":45713}

/* zap sugared logger Debugf */
logit.Debugf(ctx, "Debugf message, %s", "ok")
// Output:
// {"level":"DEBUG","time":"2020-04-15 18:12:11.991268","logger":"logit.ctx_logger","msg":"Debugf message, ok","pid":45713}

/* zap sugared logger Debugw */
logit.Debugw(ctx, "Debugw message", "name", "axiaoxin", "age", 18)
// Output:
// {"level":"DEBUG","time":"2020-04-15 18:12:11.991277","logger":"logit.ctx_logger","msg":"Debugw message","pid":45713,"name":"axiaoxin","age":18}

/* with context */
c, _ := logit.NewCtxLogger(context.Background(), logit.CloneLogger("myname"), "trace-id-123")
logit.Debug(c, "Debug with trace id")
// Output:
// {"level":"DEBUG","time":"2020-04-15 18:12:11.991314","logger":"logit.myname","msg":"Debug with trace id","pid":45713,"traceID":"trace-id-123"}

/* extra fields */
logit.Debug(c, "extra fields demo", logit.ExtraField("k1", "v1", "k2", 2, "k3", true))
// Output:
// {"level":"DEBUG","time":"2020-04-15 18:12:11.991348","logger":"logit.myname","msg":"extra fields demo","pid":45713,"traceID":"trace-id-123","extra":{"k1":"v1","k2":2,"k3":true}}
```

**详细示例 [example/logit.go](_example/logit.go)**

## 替换默认 logger

```go
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

```

**示例 [example/replace.go](_example/replace.go)**

## 快速获取、创建你的 Logger

`logit` 提供多种方式快速获取一个 logger 来打印日志

**示例 [example/logger.go](_example/logging.go)**

## 带 Trace ID 的 CtxLogger

每一次函数或者 gin 的 http 接口调用，在最顶层入口处都将一个带有唯一 trace id 的 logger 放入 context.Context 或 gin.Context ，
后续函数在内部打印日志时从 Context 中获取带有本次调用 trace id 的 logger 来打印日志几个进行调用链路跟踪。

**示例 1 普通函数中打印打印带 Trace ID 的日志 [example/context.go](_example/context.go)**

**示例 2 gin 中打印带 Trace ID 的日志 [example/gin.go](_example/gintraceid.go)**

## 日志保存到文件并自动 rotate

使用 lumberjack 将日志保存到文件并 rotate.

```go
package main

import "github.com/feymanlee/logit"

// Options 传入 LumberjacSink ，并在 OutputPaths 中添加对应 scheme 就能将日志保存到文件并自动 rotate
func main() {
	// scheme 为 lumberjack ，日志文件为 /tmp/x.log , 保存 7 天，保留 10 份文件，文件大小超过 100M ，使用压缩备份，压缩文件名使用 localtime
	sink := logit.NewLumberjackSink("/tmp/x.log", 7, 10, 100, true, true)
	err := logit.RegisterSink("lumberjack", sink)
	if err != nil {
		panic(err)
	}
	options := logit.Options{
		// 使用 sink 中设置的 scheme 即 lumberjack: 或 lumberjack:// 并指定保存日志到指定文件，日志文件将自动按 LumberjackSink 的配置做 rotate
		OutputPaths: []string{"lumberjack:"},
	}
	logger, _ := logit.NewLogger(options)
	logger.Debug("xxx")

	sink2 := logit.NewLumberjackSink("/tmp/x2.log", 7, 10, 100, true, true)
	err = logit.RegisterSink("lumberjack2", sink2)
	if err != nil {
		panic(err)
	}
	options2 := logit.Options{
		// 使用 sink 中设置的 scheme 即 lumberjack: 或 lumberjack:// 并指定保存日志到指定文件，日志文件将自动按 LumberjackSink 的配置做 rotate
		OutputPaths: []string{"lumberjack2:"},
	}
	logger2, _ := logit.NewLogger(options2)
	logger2.Debug("yyy")
}

```

**示例 [example/lumberjack.go](_example/lumberjack.go)**

## 支持 Gorm 日志打印

使用 gorm v2 支持 context logger 打印 trace id

```go
package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 模拟一个 ctx ，并将 logger 和 traceID 设置到 ctx 中
	gormLogger, err := logit.NewGormLogger(logit.GormLoggerOptions{
		Name:              "gorm",
		CallerSkip:        3,
		LogLevel:          zapcore.InfoLevel,
		SlowThreshold:     5 * time.Second,
		OutputPaths:       []string{"stdout", "lumberjack:", "/tem/a-xx.log"},
		InitialFields:     nil,
		DisableCaller:     false,
		DisableStacktrace: false,
	})
	if err != nil {
		panic(err)
	}
	// 新建会话模式设置 logger，也可以在 Open 时 使用 Config 设置
	db = db.Session(&gorm.Session{
		Logger: gormLogger,
	})
}

```

**示例 [example/gorm.go](_example/gorm.go)**

## 支持 Go-redis 日志打印

使用 go-redis v8 并支持打印 trace id

```go
package main

import (
	"time"

	"github.com/feymanlee/logit"
	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	// 这里可以添加一写自定义的配置
	logHook, err := logit.NewRedisLogger(logit.RedisLoggerOptions{
		Name:          "redis",
		CallerSkip:    3,
		SlowThreshold: time.Millisecond * 10, // 慢查询阈值，会使用 Warn 打印日志
		InitialFields: map[string]interface{}{
			"key1": "value1",
		},
		OutputPaths:       []string{"stdout", "lumberjack:", "/tem/a-xx.log"},
		DisableCaller:     false, // 禁用 caller 打印
		DisableStacktrace: false, // 禁用 Stacktrace
		EncoderConfig:     nil,
	})
	if err != nil {
		panic(err)
	}
	client.AddHook(logHook)
}

```

**示例 [example/gorm.go](_example/redis.go)**

## gin middleware: GinLogger

支持打印 gin 日志

```go
package main

import (
	"fmt"
	"time"

	"github.com/feymanlee/logit"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	// you can custom the config or use logit.GinLogger() by default config
	conf := logit.GinLoggerConfig{
		Name: "access",
		Formatter: func(c *gin.Context, ext logit.GinLogExtends) string {
			return fmt.Sprintf("%s use %s request %s at %v, handler %s use %f seconds to respond it with %d",
				c.ClientIP(),
				c.Request.Method,
				c.Request.Host,
				c.Request.RequestURI,
				ext.HandleName,
				ext.Latency,
				c.Writer.Status())
		},
		SkipPaths:           []string{"/user/list"},
		EnableDetails:       false,
		TraceIDFunc:         func(c *gin.Context) string { return "my-trace-id" },
		SkipPathRegexps:     []string{"/user/.*?"},
		EnableContextKeys:   false,       // 记录 context 里面的 key
		EnableRequestHeader: false,       // 记录 header
		EnableRequestForm:   false,       // 记录 request form
		EnableRequestBody:   false,       // 记录 request body
		EnableResponseBody:  false,       // 记录 response body
		SlowThreshold:       time.Second, // 慢查询阈值，超时这个时间会答应 Warn 日志
		OutputPaths:         []string{"stdout", "lumberjack:", "/tem/a-xx.log"},
		InitialFields:       map[string]interface{}{"key1": "value1"}, // 一些初始化的打印字段
		DisableCaller:       false,                                    // 禁用 caller 打印
		DisableStacktrace:   false,                                    // 禁用 Stacktrace
		EncoderConfig:       nil,
	}
	app.Use(logit.NewGinLogger(conf))
	app.POST("/ping", func(c *gin.Context) {
		// panic("xx")
		// time.Sleep(300 * time.Millisecond)
		c.JSON(200, string(logit.GetGinRequestBody(c)))
	})
	app.Run(":8888")
}

```

示例： [example/ginlogger.go](_example/ginlogger.go)

## 自定义 logger Encoder 配置

**示例 [example/encoder.go](_example/encoder.go)**

## 感谢

* 从 [axiaoxin-com/logging](https://github.com/axiaoxin-com/logging) 获得灵感并参考了很多的代码
