# logit

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/feymanlee/logit?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/feymanlee/logit)](https://goreportcard.com/report/github.com/feymanlee/logit)
[![Unit-Tests](https://github.com/feymanlee/logit/workflows/Unit-Tests/badge.svg)](https://github.com/feymanlee/logit/actions)
[![Coverage Status](https://coveralls.io/repos/github/feymanlee/logit/badge.svg?branch=main)](https://coveralls.io/github/feymanlee/logit?branch=main)

logit 简单封装了在日常使用 [zap](https://github.com/uber-go/zap) 打日志时的常用方法。

- 提供快速使用 zap 打印日志的方法，除 zap 的 DPanic 、 DPanicf 方法外所有日志打印方法开箱即用
- 提供多种快速创建 `logger` 的方法
- 支持从 Context 中创建、获取带有 **Trace ID** 的 logger
- 提供 `gin` 的日志中间件，支持 Trace ID，可以记录更加详细的请求和响应信息，支持通过配置自定义
- 支持 `Gorm` 日志并打印 Trace ID
- 支持 `go-redis` 日志并打印 Trace ID，目前只支持 `go-redis/v8`, 后续会增加对 `go-redis/v9` 的支持
- 支持服务内部函数方式和外部 HTTP 方式 **动态调整日志级别**，无需修改配置、重启服务
- 支持自定义 logger Encoder 配置
- 支持将日志保存到文件并自动 rotate

logit 只提供 zap 使用时的常用方法汇总，不是对 zap 进行二次开发，拒绝过度封装。

## 开箱即用

`logit` 提供的开箱即用方法都是使用自身默认 logger 克隆出的 CtxLogger 实际执行的。
在 `logit` 被 import 时，会生成内部使用的默认 logger 。
默认 logger 使用 JSON 格式打印日志内容到 stderr 。
默认不带 Sentry 上报功能，可以通过设置环境变量或者替换 logger 方法支持。
默认 logger 可通过代码内部动态修改日志级别， 默认不支持 HTTP 方式动态修改日志级别，需要指定端口创建新的 logger 来支持。
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

**示例 [example/logit.go](_example/logit.go)**

## 替换默认 log

**示例 [example/replace.go](_example/replace.go)**

## 快速获取、创建你的 Logger

`logit` 提供多种方式快速获取一个 logger 来打印日志

**示例 [example/logger.go](_example/logger.go)**

## 带 Trace ID 的 CtxLogger

每一次函数或者 gin 的 http 接口调用，在最顶层入口处都将一个带有唯一 trace id 的 logger 放入 context.Context 或 gin.Context ，
后续函数在内部打印日志时从 Context 中获取带有本次调用 trace id 的 logger 来打印日志几个进行调用链路跟踪。

**示例 1 普通函数中打印打印带 Trace ID 的日志 [example/context.go](_example/context.go)**

**示例 2 gin 中打印带 Trace ID 的日志 [example/gin.go](_example/gintraceid.go)**:

## 动态修改 logger 日志级别

`logit` 可以在代码中对 AtomicLevel 调用 SetLevel 动态修改日志级别，也可以通过请求 HTTP 接口修改。
创建 logger 时可自定义端口运行 HTTP 服务来接收请求修改日志级别。实际使用中日志级别通常写在配置文件中，
可以通过监听配置文件的修改来动态调用 SetLevel 方法。

**示例 [example/atomiclevel.go](_example/atomiclevel.go)**

## 自定义 logger Encoder 配置

**示例 [example/encoder.go](_example/encoder.go)**

## 日志保存到文件并自动 rotate

使用 lumberjack 将日志保存到文件并 rotate ，采用 zap 的 RegisterSink 方法和 Config.OutputPaths 字段添加自定义的日志输出的方式来使用 lumberjack 。

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:       true,
		AllowGlobalUpdate: false,
		Logger: logit.NewGormLogger(logit.GormLoggerOptions{
			Name:             "gorm",
			LoggerCallerSkip: 3,
			LogLevel:         zap.InfoLevel,
			SlowThreshold:    time.Millisecond * 200,
			TraceWithLevel:   zap.InfoLevel,
		}),
	})
}

```

**示例 [example/gorm.go](_example/gorm.go)**

## 支持 Go-redis 日志打印

使用 go-redis v8 并支持打印 trace id

```go
package main

import (
	"github.com/feymanlee/logit"
	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	// 这里可以添加一写自定义的配置
	logHook := logit.NewRedisLogger(logit.RedisLoggerOptions{})
	client.AddHook(logHook)
}

```

**示例 [example/gorm.go](_example/gorm.go)**

## gin middleware: GinLogger

支持打印 gin 日志

```go
package main

import (
	"fmt"

	"github.com/feymanlee/logit"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	// you can custom the config or use logit.GinLogger() by default config
	conf := logit.GinLoggerConfig{
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
		SkipPaths:     []string{},
		EnableDetails: false,
		TraceIDFunc:   func(c *gin.Context) string { return "my-trace-id" },
	}
	app.Use(logit.GinLoggerWithConfig(conf))
	app.POST("/ping", func(c *gin.Context) {
		// panic("xx")
		// time.Sleep(300 * time.Millisecond)
		c.JSON(200, string(logit.GetGinRequestBody(c)))
	})
	app.Run(":8888")
}

```

示例： [example/ginlogger.go](_example/ginlogger.go)

## 感谢

* 从 [axiaoxin-com/logging](https://github.com/axiaoxin-com/logging) 获得灵感并参考了很多的代码
