package redis

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"

	"AIM/internal/common/config"
	"AIM/internal/common/logger"
)

var RDB *redis.Client
var Ctx = context.Background()

func Init() {
	cfg := config.Config.Redis
	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic("Redis connect failed: " + err.Error())
	}

	RDB = rdb
	logger.Info("Redis connected ✅")
}
