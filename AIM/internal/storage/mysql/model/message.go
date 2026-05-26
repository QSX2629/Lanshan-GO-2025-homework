package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         string    `gorm:"column:id;primaryKey"`
	SendID     string    `gorm:"column:send_id"`
	RecvID     string    `gorm:"column:recv_id"`
	FromUID    string    `gorm:"index"`
	ToUID      string    `gorm:"index"`
	Content    string    `gorm:"column:content"`
	CreateTime time.Time `gorm:"column:create_time"`
	GroupID    int64     `gorm:"column:group_id"`
	IsRead     int       `gorm:"default:0"` // 0=未读 1=已读
	ReadTime   int64     // 已读时间
	MsgType    int       `gorm:"default:1"` // 1=文本 2=图片 3=文件 4=语音
}

func (m *Message) TableName() string {
	return "message"
}

type MessageRead struct {
	gorm.Model
	MessageID string `gorm:"size:64;index"`
	UserID    string `gorm:"size:64;index"`
	IsRead    bool   `gorm:"default:false"`
}
