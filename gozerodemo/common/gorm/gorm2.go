package gorm2

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const UserDSN = "root:Qsx20061212@tcp(127.0.0.1:3306)/gozerodemo?charset=utf8mb4&parseTime=True&loc=Local"

var UserDB *gorm.DB

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Error,
		},
	)
	db, err := gorm.Open(mysql.Open(UserDSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logx.WithContext(context.Background()).Errorf("GORM connect UserDB Error: %+v", err)
		return
	}
	UserDB = db
}
