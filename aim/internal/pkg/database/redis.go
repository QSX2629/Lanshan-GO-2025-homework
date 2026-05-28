package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// DefaultRedisConfig returns a localhost Redis config for development.
func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
}

// NewRedis creates a new Redis client and verifies connectivity.
func NewRedis(cfg RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
