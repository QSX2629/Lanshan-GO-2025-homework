package repo

import (
	"AIM/internal/storage/mysql"
	"AIM/internal/storage/mysql/model"
)

func IsFriend(user_id, friend_id string) bool {
	var cnt int64
	mysql.DB.Model(&model.Friend{}).
		Where("user_id=?AND friend_id=?", user_id, friend_id).
		Count(&cnt)
	return cnt > 0
}
func AddFriend(userID, friendID, remark string) error {
	item := model.Friend{
		UserID:   userID,
		FriendID: friendID,
		Remark:   remark,
	}
	return mysql.DB.Create(&item).Error
}

// ListFriend 好友列表
func ListFriend(userID string) ([]model.Friend, error) {
	var list []model.Friend
	err := mysql.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

// DelFriend 删除好友
func DelFriend(userID, friendID string) error {
	return mysql.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).
		Delete(&model.Friend{}).Error
}
