/**
 * @Author: feymanlee@gmail.com
 * @Description:
 * @File:  redis_test
 * @Date: 2023/4/10 13:42
 */

package logit

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestRedisLogger(t *testing.T) {
	logger, err := NewRedisLogger(RedisLoggerOptions{
		Name:          "redis",
		CallerSkip:    4,
		SlowThreshold: 10 * time.Millisecond,
		OutputPaths:   []string{"stdout"},
		InitialFields: map[string]interface{}{
			"key1": "value1",
		},
		DisableCaller:     false,
		DisableStacktrace: false,
		EncoderConfig:     &defaultEncoderConfig,
		LumberjackSink:    nil,
	})
	if err != nil {
		t.Errorf("new gorm logger failed: %v", err)
	}
	if logger == (RedisLogger{}) {
		t.Error("CtxGormLogger return empty GormLogger")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	redisClient.AddHook(logger)
	ctx := context.Background()
	// redisClient.Get(ctx, "a")
	// redisClient.Set(ctx, "b", 1, -1)
	// redisClient.Get(ctx, "b")
	pipeline := redisClient.Pipeline()
	pipeline.Set(ctx, "c", 1, -1)
	pipeline.Get(ctx, "a")
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return
	}
	redisClient.Del(ctx, "b", "c")
}
