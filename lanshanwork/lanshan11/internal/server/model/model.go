package model

import "gorm.io/gorm"

// User 用户表结构体（与RPC服务端模型一致）
type User struct {
	gorm.Model        // 内置字段：ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string `gorm:"type:varchar(55);uniqueIndex;not null" json:"username"` // 唯一索引，非空
	Password   string `gorm:"type:varchar(255);not null" json:"password"`            // 密码（建议加密存储）
}
