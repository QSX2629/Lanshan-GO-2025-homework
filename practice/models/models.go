package models

import "github.com/golang-jwt/jwt/v5"

type AppConfigInformation struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port string `mapstructure:"port"`
}
type MysqlConfig struct {
	Dsn             string `mapstructure:"dsn"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifeTime int    `mapstructure:"conn_max_life_time"`
}
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
}
type LogConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	MaxSize  int    `mapstructure:"max_size"`
	MaxAge   int    `mapstructure:"max_age"`
	Compress bool   `mapstructure:"compress"`
}
type AppConfig struct {
	App   AppConfigInformation `mapstructure:"app"`
	Mysql MysqlConfig          `mapstructure:"mysql"`
	Redis RedisConfig          `mapstructure:"redis2"`
	JWT   JWTConfig            `mapstructure:"jwt"`
	Log   LogConfig            `mapstructure:"log"`
}
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}
