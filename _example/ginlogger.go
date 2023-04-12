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
		OutputPaths:         []string{"stdout", "lumberjack:"},
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
