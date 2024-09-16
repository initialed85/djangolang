package config

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

func GetRedisFromEnvironment() (*redis.Pool, error) {
	redisURL := RedisURL()

	redisPool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialURLContext(ctx, redisURL)
		},
		MaxIdle:         2,
		MaxActive:       100,
		IdleTimeout:     300,
		Wait:            false,
		MaxConnLifetime: 86400,
	}

	return redisPool, nil
}
