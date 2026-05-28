package model

import (
	"time"
)

type Group struct {
	ID        int64  `gorm:"primarykey"`
	GroupName string `json:"group_name"`
	OwnerUID  string `gorm:"index"` // 群主UID
	Avatar    string `json:"avatar"`

	Intro      string `json:"intro"`
	Notice     string `gorm:"size:512"` // 群公告
	CreateTime time.Time
}

func (Group) TableName() string {
	return "group"
}

// GroupMember 群成员
type GroupMember struct {
	ID         int64 `gorm:"primarykey"`
	GroupID    int64
	UID        string
	JoinTime   time.Time
	Role       int   `gorm:"default:1"`     // 1=普通成员 2=管理员 3=群主
	IsMuted    bool  `gorm:"default:false"` // 是否禁言
	MutedUntil int64 // 禁言到期时间
}

func (GroupMember) TableName() string {
	return "group_member"
}

const (
	GroupRoleMember = 1 // 普通成员
	GroupRoleAdmin  = 2 // 管理员
	GroupRoleOwner  = 3 // 群主
)
