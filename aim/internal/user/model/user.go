package model

import "time"

// User represents a registered account.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash string    `gorm:"size:256;not null" json:"-"`
	Nickname     string    `gorm:"size:128" json:"nickname"`
	Avatar       string    `gorm:"size:512" json:"avatar"`
	Email        string    `gorm:"size:128" json:"email"`
	Phone        string    `gorm:"size:32" json:"phone"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// OnlineStatus tracks the current online state of a user.
type OnlineStatus struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Online   bool      `gorm:"not null;default:false" json:"online"`
	LastSeen time.Time `json:"last_seen"`
}
