package model

import "time"

// Group member role constants.
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)

// Friend represents a friendship between two users.
type Friend struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	FriendID  uint      `gorm:"index;not null" json:"friend_id"`
	Remark    string    `gorm:"size:128" json:"remark"`
	GroupName string    `gorm:"size:64" json:"group_name"` // friend group
	CreatedAt time.Time `json:"created_at"`
}

// FriendGroup is a user-defined category for organizing friends.
type FriendGroup struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `gorm:"index;not null" json:"user_id"`
	Name   string `gorm:"size:64;not null" json:"name"`
}

// Group represents a chat group.
type Group struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:128;not null" json:"name"`
	Avatar       string    `gorm:"size:512" json:"avatar"`
	OwnerID      uint      `gorm:"not null" json:"owner_id"`
	Announcement string    `gorm:"type:text" json:"announcement"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GroupMember represents a user's membership in a group.
type GroupMember struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	GroupID    uint       `gorm:"index;not null" json:"group_id"`
	UserID     uint       `gorm:"index;not null" json:"user_id"`
	Role       string     `gorm:"size:16;not null;default:member" json:"role"`
	Muted      bool       `gorm:"not null;default:false" json:"muted"`
	MutedUntil *time.Time `json:"muted_until,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}
