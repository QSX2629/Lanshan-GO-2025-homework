package redis

import (
	"context"

	"practice/log"
	"practice/viper"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 全局Redis客户端和上下文
var (
	RedisClient *redis.Client
	Ctx         = context.Background() // 定义全局上下文
)

// InitRedis 初始化Redis
func InitRedis() {

	cfg := viper.GetConfig()
	if cfg == nil {
		log.Panic("初始化Redis失败：配置未加载")
	}
	redisCfg := cfg.Redis

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,     // 从配置读取地址
		Password: redisCfg.Password, // 从配置读取密码
		DB:       redisCfg.DB,       // 从配置读取数据库编号
		PoolSize: redisCfg.PoolSize, // 新增连接池大小（可选，提升性能）
	})

	// 3. 测试Redis连接
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Panic("Redis初始化失败",
			zap.String("addr", redisCfg.Addr),
			zap.Error(err),
		)
	}

	// 4. 打印成功日志
	log.Info("Redis初始化成功",
		zap.String("addr", redisCfg.Addr),
		zap.Int("db", redisCfg.DB),
		zap.Int("pool_size", redisCfg.PoolSize),
	)
}
