/**
 * @Author: feymanlee@gmail.com
 * @Description:
 * @File:  redis
 * @Date: 2023/4/7 15:17
 */

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
