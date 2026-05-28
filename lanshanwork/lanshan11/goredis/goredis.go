package goredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

var Rdb *redis.Client

func init() {
	option := redis.Options{
		Addr:     "localhost:6379",
		DB:       1,
		PoolSize: 100,
	}
	rdb := redis.NewClient(&option)
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logx.WithContext(context.Background()).Error("Redis connect ERROR: %+v", err)
	}
	Rdb = rdb
}
