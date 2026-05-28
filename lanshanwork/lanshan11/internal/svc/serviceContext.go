package svc

import (
	pb "lanshan11/desc"
	"lanshan11/goredis"
	"lanshan11/gorm2"
	"lanshan11/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	MysqlDB *gorm.DB
	RedisDB *redis.Client
	UserRPC pb.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建 zrpc 客户端
	rpcClient, err := zrpc.NewClient(c.UserRPC)
	if err != nil {
		panic("failed to create user rpc client: " + err.Error())
	}

	return &ServiceContext{
		Config:  c,
		MysqlDB: gorm2.UserDB,
		RedisDB: goredis.Rdb,
		// ✅ 传入底层连接，实现 grpc.ClientConnInterface
		UserRPC: pb.NewUserClient(rpcClient.Conn()),
	}
}
