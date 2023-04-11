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
		OutputPaths:       []string{"stdout", "lumberjack:"},
		DisableCaller:     false, // 禁用 caller 打印
		DisableStacktrace: false, // 禁用 Stacktrace
		EncoderConfig:     nil,
	})
	if err != nil {
		panic(err)
	}
	client.AddHook(logHook)
}
