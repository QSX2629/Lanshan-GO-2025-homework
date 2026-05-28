package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf                    // go-zero内置服务配置（端口、日志等）
	RestConf            rest.RestConf      `json:"restConf"` // REST API 配置（包含 Host、Port）
	UserRPC             zrpc.RpcClientConf `json:"userRpc"`

	// MySQL配置（可选，若直接用gorm2的硬编码可省略）
	Mysql struct {
		DSN string `json:"dsn"`
	}

	// Redis配置
	Redis redis.RedisConf `json:"redis"`

	// User RPC服务配置（etcd注册发现）

}
