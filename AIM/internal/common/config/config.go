package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var Config *AppConfig

type AppConfig struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	AI     AIConfig     `mapstructure:"ai"`
	Log    LogConfig    `mapstructure:"log"`
	App    AppInfo      `mapstructure:"app"`
}
type ServerConfig struct {
	CometPort int `mapstructure:"comet_port"`
	LogicPort int `mapstructure:"logic_port"`
	AiPort    int `mapstructure:"ai_port"`
	ApiPort   int `mapstructure:"api_port"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
type AIConfig struct {
	Enable  bool   `mapstructure:"enable"`
	Model   string `mapstructure:"model"`
	ApiKey  string `mapstructure:"api_key"`
	Timeout int    `mapstructure:"timeout"`
}
type LogConfig struct {
	Level     string `mapstructure:"level"`
	Filename  string `mapstructure:"filename"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxBackup int    `mapstructure:"max_backup"`
}
type AppInfo struct {
	Env     string `mapstructure:"env"`
	Version string `mapstructure:"version"`
}

// Load 加载配置，自动适配当前文件位置
func Load() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// 关键：直接使用当前文件的路径，找到 config.yaml
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("无法获取当前文件路径")
	}
	// 当前文件是 config.go，它的同级目录就是 config.yaml 所在目录
	configDir := filepath.Dir(filename)
	v.AddConfigPath(configDir)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("配置文件加载失败: %v", err)
	}

	Config = &AppConfig{}
	if err := v.Unmarshal(Config); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}

	log.Println("✅ 配置加载成功")
}
