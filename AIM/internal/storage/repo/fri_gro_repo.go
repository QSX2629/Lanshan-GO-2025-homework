package repo

import (
	"AIM/internal/storage/mysql"
	"AIM/internal/storage/mysql/model"
)

// ====================== 分组管理 ======================

// CreateGroup 创建分组
func CreateFriendGroup(uid string, groupName string) error {
	return mysql.DB.Create(&model.FriendGroup{
		UserID:    uid,
		GroupName: groupName,
	}).Error
}

// ListGroup 获取用户所有分组
func ListFriendGroup(uid string) ([]model.FriendGroup, error) {
	var list []model.FriendGroup
	err := mysql.DB.Where("uid = ?", uid).Find(&list).Error
	return list, err
}

// DeleteGroup 删除分组
func DeleteFriendGroup(uid string, groupID int64) error {
	return mysql.DB.Where("id = ? AND uid = ?", groupID, uid).Delete(&model.FriendGroup{}).Error
}

// ====================== 好友移动分组 ======================

// MoveFriendToGroup 移动好友到分组
func MoveFriendToGroup(uid, friendID string, groupID int64) error {
	return mysql.DB.Model(&model.FriendGroup{}).
		Where("uid = ? AND friend_id = ?", uid, friendID).
		Update("group_id", groupID).Error
}

// ListFriendByGroup 按分组获取好友列表
func ListFriendByGroup(uid string, groupID int64) ([]model.FriendGroup, error) {
	var list []model.FriendGroup
	err := mysql.DB.Where("uid = ? AND group_id = ?", uid, groupID).Find(&list).Error
	return list, err
}
