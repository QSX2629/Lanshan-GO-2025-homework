package repo

import (
	"AIM/internal/comet/protocol"
	"AIM/internal/comet/push"
	"AIM/internal/storage/mysql"
	"AIM/internal/storage/mysql/model"
	"time"
)

// ======================
// 基础群操作
// ======================

// CreateGroup 创建群
func CreateGroup(name, ownerID, desc string) (*model.Group, error) {
	group := &model.Group{
		GroupName:  name,
		OwnerUID:   ownerID,
		Intro:      desc,
		CreateTime: time.Now(),
	}

	err := mysql.DB.Create(group).Error
	if err != nil {
		return nil, err
	}

	// 群主自动入群，并且角色是群主
	AddGroupMember(group.ID, ownerID)
	SetGroupOwnerRole(group.ID, ownerID)
	return group, nil
}

// AddGroupMember 添加群成员
func AddGroupMember(groupID int64, userID string) error {
	m := &model.GroupMember{
		GroupID:  groupID,
		UID:      userID,
		JoinTime: time.Now(),
		Role:     model.GroupRoleMember,
	}
	return mysql.DB.Create(m).Error
}

// SetGroupOwnerRole 设置群主角色
func SetGroupOwnerRole(groupID int64, uid string) error {
	return mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND uid = ?", groupID, uid).
		Update("role", model.GroupRoleOwner).Error
}

// JoinGroup 用户加入群
func JoinGroup(groupID int64, userID string) error {
	return mysql.DB.Create(&model.GroupMember{
		GroupID:  groupID,
		UID:      userID,
		JoinTime: time.Now(),
		Role:     model.GroupRoleMember,
	}).Error
}

// ListGroup 查询用户的群列表
func ListGroup(userID string) ([]model.Group, error) {
	var groups []model.Group
	err := mysql.DB.Raw(`
	SELECT g.* FROM chat_group g
	JOIN chat_group_member gm ON g.id = gm.group_id
	WHERE gm.uid = ?
	`, userID).Scan(&groups).Error
	return groups, err
}

// SendGroupMsg 发送群消息 + 推送
func SendGroupMsg(fromUID string, groupID int64, content string) error {
	// 1. 存消息
	msg := &model.Message{
		SendID:     fromUID,
		GroupID:    groupID,
		Content:    content,
		CreateTime: time.Now(),
	}
	if err := mysql.DB.Create(msg).Error; err != nil {
		return err
	}

	// 2. 获取群所有成员
	var uids []string
	mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("uid", &uids)

	// 3. 构建推送消息
	pushMsg := &protocol.Message{
		Op:      protocol.OpGroupChat,
		FromUID: fromUID,
		GroupID: groupID,
		Content: content,
	}

	// 4. 批量推送（排除自己）
	var pushUids []string
	for _, uid := range uids {
		if uid != fromUID {
			pushUids = append(pushUids, uid)
		}
	}

	push.PushToGroup(pushUids, pushMsg)
	return nil
}

// ======================
// 群高级管理
// ======================

// IsGroupOwner 是否群主
func IsGroupOwner(groupID int64, uid string) bool {
	var g model.Group
	err := mysql.DB.Where("id = ? AND owner_uid = ?", groupID, uid).First(&g).Error
	return err == nil
}

// IsGroupAdmin 是否管理员/群主
func IsGroupAdmin(groupID int64, uid string) bool {
	if IsGroupOwner(groupID, uid) {
		return true
	}
	var m model.GroupMember
	err := mysql.DB.Where("group_id = ? AND uid = ? AND role = ?", groupID, uid, model.GroupRoleAdmin).First(&m).Error
	return err == nil
}

// IsInGroup 是否在群里
func IsInGroup(groupID int64, uid string) bool {
	var m model.GroupMember
	err := mysql.DB.Where("group_id = ? AND uid = ?", groupID, uid).First(&m).Error
	return err == nil
}

// IsMuted 是否被禁言
func IsMuted(groupID int64, uid string) bool {
	var m model.GroupMember
	err := mysql.DB.Where("group_id = ? AND uid = ? AND is_muted = ?", groupID, uid, true).First(&m).Error
	return err == nil
}

// KickMember 踢人
func KickMember(groupID int64, targetUID string) error {
	return mysql.DB.Where("group_id = ? AND uid = ?", groupID, targetUID).Delete(&model.GroupMember{}).Error
}

// MuteMember 禁言 minutes 分钟
func MuteMember(groupID int64, targetUID string, minutes int) error {
	until := time.Now().Add(time.Minute * time.Duration(minutes)).Unix()
	return mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND uid = ?", groupID, targetUID).
		Updates(map[string]any{
			"is_muted":    true,
			"muted_until": until,
		}).Error
}

// UnMuteMember 取消禁言
func UnMuteMember(groupID int64, targetUID string) error {
	return mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND uid = ?", groupID, targetUID).
		Updates(map[string]any{
			"is_muted":    false,
			"muted_until": 0,
		}).Error
}

// TransferOwner 转让群主
func TransferOwner(groupID int64, oldOwner, newOwner string) error {
	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tx.Model(&model.GroupMember{}).Where("group_id = ? AND uid = ?", groupID, oldOwner).Update("role", model.GroupRoleMember)
	tx.Model(&model.GroupMember{}).Where("group_id = ? AND uid = ?", groupID, newOwner).Update("role", model.GroupRoleOwner)
	tx.Model(&model.Group{}).Where("id = ?", groupID).Update("owner_uid", newOwner)

	return tx.Commit().Error
}

// UpdateGroupNotice 更新群公告
func UpdateGroupNotice(groupID int64, notice string) error {
	return mysql.DB.Model(&model.Group{}).Where("id = ?", groupID).Update("notice", notice).Error
}

// GetGroupNotice 获取群公告
func GetGroupNotice(groupID int64) (string, error) {
	var g model.Group
	if err := mysql.DB.Where("id = ?", groupID).First(&g).Error; err != nil {
		return "", err
	}
	return g.Notice, nil
}

// SetAdmin 设置管理员
func SetAdmin(groupID int64, targetUID string) error {
	return mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND uid = ?", groupID, targetUID).
		Update("role", model.GroupRoleAdmin).Error
}

// RemoveAdmin 撤销管理员
func RemoveAdmin(groupID int64, targetUID string) error {
	return mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ? AND uid = ?", groupID, targetUID).
		Update("role", model.GroupRoleMember).Error
}

// GetGroupMemberUids 获取群所有成员ID列表
func GetGroupMemberUids(groupID int64, uids *[]string) error {
	return mysql.DB.Model(&model.GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("uid", uids).Error
}
