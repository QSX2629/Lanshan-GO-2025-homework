package group

import (
	"AIM/internal/storage/mysql/model"
	"AIM/internal/storage/repo"
	"errors"
)

// CreateGroup 创建群
func CreateGroup(name, ownerID, desc string) (*model.Group, error) {
	return repo.CreateGroup(name, ownerID, desc)
}

// JoinGroup 加入群
func JoinGroup(groupID int64, userID string) error {
	return repo.JoinGroup(groupID, userID)
}

// ListMyGroup 获取我的群列表
func ListMyGroup(userID string) ([]model.Group, error) {
	return repo.ListGroup(userID)
}

// SendGroupMsg 发送群消息
func SendGroupMsg(fromUID string, groupID int64, content string) error {
	// 发送前自动校验：禁言 + 是否在群内
	if !repo.IsInGroup(groupID, fromUID) {
		return errors.New("你不在该群内")
	}
	if repo.IsMuted(groupID, fromUID) {
		return errors.New("你已被禁言，无法发言")
	}
	return repo.SendGroupMsg(fromUID, groupID, content)
}

// ======================
// 群高级管理（自动集成）
// ======================

// KickMember 踢人
func KickMember(operatorUID string, groupID int64, targetUID string) error {
	if !repo.IsGroupAdmin(groupID, operatorUID) {
		return errors.New("无权限，仅群主/管理员可操作")
	}
	return repo.KickMember(groupID, targetUID)
}

// MuteMember 禁言
func MuteMember(operatorUID string, groupID int64, targetUID string, minutes int) error {
	if !repo.IsGroupAdmin(groupID, operatorUID) {
		return errors.New("无权限，仅群主/管理员可操作")
	}
	return repo.MuteMember(groupID, targetUID, minutes)
}

// UnMuteMember 解除禁言
func UnMuteMember(operatorUID string, groupID int64, targetUID string) error {
	if !repo.IsGroupAdmin(groupID, operatorUID) {
		return errors.New("无权限，仅群主/管理员可操作")
	}
	return repo.UnMuteMember(groupID, targetUID)
}

// TransferOwner 转让群主
func TransferOwner(operatorUID string, groupID int64, newOwner string) error {
	if !repo.IsGroupOwner(groupID, operatorUID) {
		return errors.New("无权限，仅群主可转让")
	}
	return repo.TransferOwner(groupID, operatorUID, newOwner)
}

// UpdateNotice 更新群公告
func UpdateNotice(operatorUID string, groupID int64, notice string) error {
	if !repo.IsGroupAdmin(groupID, operatorUID) {
		return errors.New("无权限，仅群主/管理员可操作")
	}
	return repo.UpdateGroupNotice(groupID, notice)
}

// GetNotice 获取群公告
func GetNotice(groupID int64) (string, error) {
	return repo.GetGroupNotice(groupID)
}

// SetAdmin 设置管理员
func SetAdmin(operatorUID string, groupID int64, targetUID string) error {
	if !repo.IsGroupOwner(groupID, operatorUID) {
		return errors.New("无权限，仅群主可设置管理员")
	}
	return repo.SetAdmin(groupID, targetUID)
}

// RemoveAdmin 撤销管理员
func RemoveAdmin(operatorUID string, groupID int64, targetUID string) error {
	if !repo.IsGroupOwner(groupID, operatorUID) {
		return errors.New("无权限，仅群主可撤销管理员")
	}
	return repo.RemoveAdmin(groupID, targetUID)
}
