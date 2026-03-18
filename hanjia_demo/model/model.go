package model

import (
	"gorm.io/gorm"
)

const (
	ArticleStatusDraft     = "draft"
	ArticleStatusPublished = "published"
	ArticleStatusDeleted   = "deleted"
)

type User struct {
	gorm.Model
	Email       string `gorm2:"size:20;DEFAULT:'0'"json:"email"`
	ID          uint   `gorm2:"primary_key" json:"id"`
	Username    string `gorm2:"size:255;uniqueIndex;fulltext" json:"username"`
	Password    string `gorm2:"size:255" json:"password"`
	Common_user bool   `gorm2:"not null;default:false" json:"common_user"`
	Admin_user  bool   `gorm2:"not null;default:false" json:"admin_user"`
}
type Article struct {
	gorm.Model
	Title     string `gorm2:"size:255;uniqueIndex;fulltext" json:"title"`
	Content   string `gorm2:"size:255" json:"content"`
	ArticleID uint   `gorm2: json:"article_id"`
	UserId    uint   `gorm2:"index" json:"user_id"`
	Status    string `gorm2:"size:255" json:"status;default:'draft'"`
}
type Comment struct {
	gorm.Model
	CommentID uint   `gorm2: json:"id"`
	Content   string `gorm2:"size:255" json:"content"`
	ArticleId uint   `gorm2:"index" json:"article_id"`
}
type Follow struct {
	gorm.Model
	ID         uint `gorm2:"primary" json:"id"`
	UserID     uint `gorm2:"not null;index:idx_user_followed" json:"user_id"` // 关注者ID
	FollowedID uint `gorm2:"not null;index:idx_user_followed" json:"followed_id"`
} // 被关注者ID
