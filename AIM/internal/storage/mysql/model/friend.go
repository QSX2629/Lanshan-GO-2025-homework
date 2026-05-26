package model

import (
	"time"

	"gorm.io/gorm"
)

func (Friend) TableName() string {
	return "friend"
}

type Friend struct {
	gorm.Model
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	UserID   string `gorm:"size:64;index"`
	FriendID string `gorm:"size:64;index"`
	Remark   string `gorm:"size:128"`
	GroupID  int64  `gorm:"index;comment:分组ID"`
	Status   int    `gorm:"comment:状态 0-待确认 1-已添加 2-已拒绝 3-已删除"`
}
type FriendGroup struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	UserID     string    `gorm:"size:32;index;comment:所属用户ID"`
	Name       string    `gorm:"size:32;comment:分组名称"`
	Sort       int       `gorm:"comment:排序"`
	CreateTime time.Time `gorm:"autoCreateTime"`
	GroupName  string    `gorm:"size:32"` // 分组名（如：家人、同事
}
