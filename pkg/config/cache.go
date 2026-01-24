package config

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
)

func GetRedisFromEnvironment() (*redis.Pool, error) {
	redisURL := RedisURL()

	redisPool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialURLContext(ctx, redisURL)
		},
		MaxIdle:         5,
		MaxActive:       500,
		IdleTimeout:     30 * time.Second,
		Wait:            false,
		MaxConnLifetime: 600 * time.Second,
	}

	return redisPool, nil
}
