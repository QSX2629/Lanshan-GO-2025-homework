package test

import (
	"AIM/internal/common/config"
	"AIM/internal/common/logger"
	"AIM/internal/storage/mysql"
	"testing"
)

// TestMain 全局只初始化一次
func TestMain(m *testing.M) {
	// 加载配置
	config.Load()
	// 初始化日志
	logger.Init()
	// 初始化数据库
	mysql.Init()

	// 运行所有测试
	m.Run()
}
