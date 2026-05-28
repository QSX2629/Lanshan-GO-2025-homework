package model

import "time"

type User struct {
	ID        int64     `gorm2:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm2:"uniqueIndex;not null" json:"username"` // 用户名唯一且非空
	Password  string    `gorm2:"not null" json:"password"`             // 密码非空
	CreatedAt time.Time `gorm2:"autoCreateTime" json:"created_at"`     // 自动创建时间
	UpdatedAt time.Time `gorm2:"autoUpdateTime" json:"updated_at"`     // 自动更新时间
}
