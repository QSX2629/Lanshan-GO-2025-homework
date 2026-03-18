package viper

import (
	"fmt"
	"practice/models"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	config     *models.AppConfig
	configLock sync.RWMutex
)

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败")
	}
	configLock.Lock()
	defer configLock.Unlock()
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置失败")
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
		configLock.Lock()
		defer configLock.Unlock()
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[配置热更新] 生效 - 日志级别: %s | JWT过期时间: %d小时 | Redis地址: %s\n",
			config.Log.Level, config.JWT.ExpireTime, config.JWT.Secret, config.Redis.Addr)
	})
	fmt.Printf("配置成功")
	return nil
}
func GetConfig() *models.AppConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}
func GetEnv() string {
	cfg := GetConfig()
	if cfg == nil {
		return "dev"
	}
	return cfg.App.Env
}
func IsDev() bool {
	return GetEnv() == "dev"
}
