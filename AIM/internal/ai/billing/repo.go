package billing

import (
	"AIM/internal/storage/mysql"
	"AIM/internal/storage/mysql/model"
	"errors"

	"gorm.io/gorm"
)

// SaveUserConfig 保存用户AI配置
func SaveUserConfig(uid, platform, apiKey string) error {
	return mysql.DB.Save(&model.UserAPIConfig{
		UID:        uid,
		Platform:   platform,
		APIKey:     apiKey,
		UsedTokens: 0,
		Limit:      100000,
	}).Error
}

// GetUserConfig 获取用户AI配置
func GetUserConfig(uid string) (*model.UserAPIConfig, error) {
	var cfg model.UserAPIConfig
	if err := mysql.DB.Where("uid = ?", uid).First(&cfg).Error; err != nil {
		return nil, errors.New("no config found")
	}
	return &cfg, nil
}

// AddUsage 增加已用Token数
func AddUsage(uid string, tokens int64) error {
	return mysql.DB.Model(&model.UserAPIConfig{}).
		Where("uid = ?", uid).
		Update("used_tokens", gorm.Expr("used_tokens + ?", tokens)).Error
}

// CheckLimit 检查是否超出限额
func CheckLimit(uid string) (bool, error) {
	cfg, err := GetUserConfig(uid)
	if err != nil {
		return false, err
	}
	return cfg.UsedTokens < cfg.Limit, nil
}
