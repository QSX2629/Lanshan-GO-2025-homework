package mysql

import (
	"AIM/internal/storage/mysql/model"
	"fmt"
	"log"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"AIM/internal/common/config"
	"AIM/internal/common/logger"
)

var DB *gorm.DB

func Init() {
	conf := config.Config.MySQL

	// 1. 先连接 MySQL 服务本身（不指定数据库），用来创建库
	rootDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
	)

	rootDB, err := gorm.Open(mysql.Open(rootDsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		log.Fatal("连接 MySQL 服务失败：", err)
	}

	// 2. 检查并创建数据库
	createDBsql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", conf.Database)
	if err := rootDB.Exec(createDBsql).Error; err != nil {
		log.Fatal("创建数据库失败：", err)
	}
	logger.Info("数据库已创建或已存在", zap.String("db", conf.Database))

	// 3. 连接到目标数据库
	targetDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	)

	db, err := gorm.Open(mysql.Open(targetDsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		log.Fatal("连接目标数据库失败：", err)
	}

	DB = db
	logger.Info("MySQL 连接成功 ✅")
	err = DB.AutoMigrate(
		&model.User{},
		&model.Message{},
		&model.Friend{},
		&model.Group{},
		&model.GroupMember{},
		&model.UserAPIConfig{}, // 从 model 包导入，无循环依赖
	)
	if err != nil {
		log.Fatal("建表失败：", err)
	}

	logger.Info("✅ MySQL 表自动创建完成", zap.String("db", conf.Database))
}
