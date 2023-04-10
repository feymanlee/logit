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
