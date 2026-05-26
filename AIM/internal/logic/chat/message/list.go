package message

import (
	"AIM/internal/logic/relation/friend"
	"AIM/internal/storage/mysql/model"
	"AIM/internal/storage/repo"
	"errors"
)

// List 获取聊天历史（严格匹配你的 3 个参数）
func List(uid1, uid2 string, limit int) ([]model.Message, error) {
	// 直接调用你真实的 repo 函数
	return repo.GetChatHistory(uid1, uid2, limit)
}
func SearchKeyword(uid1, uid2, keyword string) ([]model.Message, error) {
	if !friend.IsFriend(uid1, uid2) {
		return nil, errors.New("非好友无法搜索")
	}
	return repo.SearchByKeyword(uid1, uid2, keyword, 100)
}

// ======================
// 【新增】时间范围搜索
// ======================
func SearchTime(uid1, uid2 string, startTime, endTime int64) ([]model.Message, error) {
	if !friend.IsFriend(uid1, uid2) {
		return nil, errors.New("非好友无法搜索")
	}
	return repo.SearchByTimeRange(uid1, uid2, startTime, endTime, 100)
}
