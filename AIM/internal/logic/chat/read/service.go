package read

import (
	"AIM/internal/storage/repo"
)

// MarkMsgRead 标记消息已读
func MarkMsgRead(msgID, userID string) error {
	return repo.MarkRead(msgID, userID)
}

// GetUserUnread 获取用户未读消息数
func GetUserUnread(userID string) (int64, error) {
	return repo.GetUnreadCount(userID)
}
