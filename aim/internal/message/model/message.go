package model

import "time"

// Target type constants.
const (
	TargetUser  = "user"
	TargetGroup = "group"
)

// Message status constants.
const (
	MsgSent      = "sent"
	MsgDelivered = "delivered"
	MsgRead      = "read"
)

// Message represents a single chat message.
type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Seq        int64     `gorm:"index;not null" json:"seq"`
	FromID     uint      `gorm:"index;not null" json:"from_id"`
	ToID       uint      `gorm:"index;not null" json:"to_id"`
	TargetType string    `gorm:"size:16;not null" json:"target_type"` // "user" or "group"
	MsgType    string    `gorm:"size:16;not null" json:"msg_type"`    // text/image/file/voice
	Content    string    `gorm:"type:text" json:"content"`
	Status     string    `gorm:"size:16;not null;default:sent" json:"status"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

// Session holds the latest conversation state for a user's chat.
type Session struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex:idx_user_target;not null" json:"user_id"`
	TargetID    uint      `gorm:"uniqueIndex:idx_user_target;not null" json:"target_id"`
	TargetType  string    `gorm:"size:16;not null" json:"target_type"`
	LastMsg     string    `gorm:"type:text" json:"last_msg"`
	LastMsgTime time.Time `json:"last_msg_time"`
	UnreadCount int       `gorm:"not null;default:0" json:"unread_count"`
}

// ReadReceipt tracks which users have read a message.
type ReadReceipt struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	MsgID  uint      `gorm:"index;not null" json:"msg_id"`
	UserID uint      `gorm:"index;not null" json:"user_id"`
	ReadAt time.Time `json:"read_at"`
}
