package repo

import (
	"AIM/internal/comet/connection"
	"AIM/internal/comet/protocol"
	"AIM/internal/storage/mysql"
	"AIM/internal/storage/mysql/model"
	"encoding/json"
	"time"
)

// SaveMessage 保存消息
func SaveMessage(message *model.Message) error {
	return mysql.DB.Create(message).Error
}

// GetChatHistory 获取两人聊天记录（修复了你原有字段错误）
func GetChatHistory(uid1, uid2 string, limit int) ([]model.Message, error) {
	var msgs []model.Message
	// 修复：你的字段是 from_uid / to_uid，不是 send_id / recv_id
	err := mysql.DB.
		Where("(from_uid = ? AND to_uid = ?) OR (from_uid = ? AND to_uid = ?)",
			uid1, uid2, uid2, uid1).
		Order("create_time asc").
		Limit(limit).
		Find(&msgs).Error
	return msgs, err
}

// MarkRead 标记单条消息已读
func MarkRead(msgID, userID string) error {
	return mysql.DB.Model(&model.MessageRead{}).
		Where("message_id = ? AND user_id = ?", msgID, userID).
		Update("is_read", true).Error
}

// GetUnreadCount 获取未读消息总数
func GetUnreadCount(userID string) (int64, error) {
	var cnt int64
	err := mysql.DB.Model(&model.MessageRead{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&cnt).Error
	return cnt, err
}

// SetMessageRead 把两个人之间的消息设为已读
func SetMessageRead(fromUID, toUID string) error {
	return mysql.DB.Model(&model.Message{}).
		Where("from_uid = ? AND to_uid = ? AND is_read = 0", fromUID, toUID).
		Updates(map[string]any{
			"is_read":   1,
			"read_time": time.Now().Unix(),
		}).Error
}

// PushTypingNotify 推送正在输入状态
func PushTypingNotify(fromUID, toUID string) error {
	conn := connection.GetConn(toUID)
	if conn == nil {
		return nil
	}
	msg := protocol.Message{
		Op:      protocol.OpTypingStatus,
		FromUID: fromUID,
		ToUID:   toUID,
		Content: "typing",
	}
	data, _ := json.Marshal(msg)
	conn.SendMessage(data)
	return nil
}

// CreateMessage 创建消息（支持多类型）
func CreateMessage(fromUID, toUID string, groupID int64, content string, msgType int) error {
	msg := model.Message{
		FromUID: fromUID,
		ToUID:   toUID,
		GroupID: groupID,
		Content: content,
		MsgType: msgType,
		IsRead:  0,
	}
	return mysql.DB.Create(&msg).Error
}

// ======================
// 以下是【消息漫游 + 搜索】（我给你加上的）
// ======================

// List 云端漫游消息（和你原有逻辑兼容）
func List(uid1, uid2 string, limit int) ([]model.Message, error) {
	var list []model.Message
	err := mysql.DB.
		Where("(from_uid = ? AND to_uid = ?) OR (from_uid = ? AND to_uid = ?)",
			uid1, uid2, uid2, uid1).
		Order("create_time desc").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// SearchByKeyword 关键词搜索消息
func SearchByKeyword(uid1, uid2, keyword string, limit int) ([]model.Message, error) {
	var list []model.Message
	err := mysql.DB.
		Where("(from_uid = ? AND to_uid = ?) OR (from_uid = ? AND to_uid = ?)",
			uid1, uid2, uid2, uid1).
		Where("content LIKE ?", "%"+keyword+"%").
		Order("create_time desc").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// SearchByTimeRange 按时间范围搜索消息
func SearchByTimeRange(uid1, uid2 string, startTime, endTime int64, limit int) ([]model.Message, error) {
	var list []model.Message
	err := mysql.DB.
		Where("(from_uid = ? AND to_uid = ?) OR (from_uid = ? AND to_uid = ?)",
			uid1, uid2, uid2, uid1).
		Where("create_time BETWEEN ? AND ?", startTime, endTime).
		Order("create_time desc").
		Limit(limit).
		Find(&list).Error
	return list, err
}
