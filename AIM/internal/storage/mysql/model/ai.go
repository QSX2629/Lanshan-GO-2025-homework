// user_api_config.go
package model

import "time"

// UserAPIConfig 用户AI配置表
type UserAPIConfig struct {
	UID        string `gorm:"primaryKey;size:64"`
	Platform   string `gorm:"size:32"`
	APIKey     string `gorm:"size:255"`
	UsedTokens int64  `gorm:"default:0"`
	Limit      int64  `gorm:"default:100000"`
	UpdatedAt  time.Time
}
