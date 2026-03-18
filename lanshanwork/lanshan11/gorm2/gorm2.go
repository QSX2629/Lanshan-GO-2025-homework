package gorm2

import (
	"context"
	"lanshan11/user/model"
	"log"
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const UserDSN = "root:123456@tcp(127.0.0.1:3306)/原神启动?charset=utf8mb4&parseTime=True&loc=Local"

var UserDB *gorm.DB

func init() {
	newLogger := logger.New(log.New(os.Stdout, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Error,
		})
	db, err := gorm.Open(mysql.Open(UserDSN),
		&gorm.Config{
			Logger: newLogger,
		})
	if err != nil {
		logx.WithContext(context.Background()).Errorf("GORM connect UserDB Error: %+v", err)
	}

	err = db.AutoMigrate(&model.User{}) //这里的model需要你提前声明（例如你可以在api和rpc目录下新建model目录，然后在model目录中去定义）
	if err != nil {
		logx.WithContext(context.Background()).Errorf("GORM AutoMigrate user ERROR:%+v", err)
	}
	UserDB = db
}
