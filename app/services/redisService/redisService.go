package redisService

import (
	"context"
	"time"
	"usercenter/config/redis"
)

var (
	ctx = context.Background()
)

func SetRedis(key string, value string) bool {
	var t int64
	t = 900
	expire := time.Duration(t) * time.Second
	if err := redis.RedisClient.Set(ctx, key, value, expire).Err(); err != nil {
		return false
	}
	return true
}

func GetRedis(key string) string {
	result, err := redis.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}
