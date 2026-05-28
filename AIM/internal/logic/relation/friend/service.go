package friend

import (
	"AIM/internal/storage/mysql/model"
	"AIM/internal/storage/repo"
)

func AddFriend(userID, friendID, remark string) error {
	return repo.AddFriend(userID, friendID, remark)
}

func ListFriend(userID string) ([]model.Friend, error) {
	return repo.ListFriend(userID)
}

func DelFriend(userID, friendID string) error {
	return repo.DelFriend(userID, friendID)
}
func IsFriend(userID, friendID string) bool {
	return repo.IsFriend(userID, friendID)
}

// ====================== 好友分组 ======================

// CreateGroup 创建好友分组
func CreateGroup(uid, groupName string) error {
	return repo.CreateFriendGroup(uid, groupName)
}

// ListGroup 获取我的所有分组
func ListGroup(uid string) ([]model.FriendGroup, error) {
	return repo.ListFriendGroup(uid)
}

// DeleteGroup 删除分组
func DeleteGroup(uid string, groupID int64) error {
	return repo.DeleteFriendGroup(uid, groupID)
}

// MoveFriend 移动好友到分组
func MoveFriend(uid, friendID string, groupID int64) error {
	return repo.MoveFriendToGroup(uid, friendID, groupID)
}

// ListByGroup 按分组获取好友
func ListByGroup(uid string, groupID int64) ([]model.FriendGroup, error) {
	return repo.ListFriendByGroup(uid, groupID)
}
