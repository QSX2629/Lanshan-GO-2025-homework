package handler

import (
	"context"
	"time"

	"github.com/aim/aim/internal/pkg/database"
	relmodel "github.com/aim/aim/internal/relation/model"
	"github.com/aim/aim/internal/relation/service"
)

// RelationHandler provides the handler for friend and group operations.
type RelationHandler struct {
	svc *service.RelationService
}

// NewRelationHandler creates a new RelationHandler.
func NewRelationHandler(db *database.DB) *RelationHandler {
	return &RelationHandler{svc: service.NewRelationService(db)}
}

// --- Friends ---

// AddFriend adds a friend.
func (h *RelationHandler) AddFriend(_ context.Context, userID, friendID uint, remark string) error {
	return h.svc.AddFriend(userID, friendID, remark)
}

// DeleteFriend removes a friend.
func (h *RelationHandler) DeleteFriend(_ context.Context, userID, friendID uint) error {
	return h.svc.DeleteFriend(userID, friendID)
}

// ListFriends returns all friends.
func (h *RelationHandler) ListFriends(_ context.Context, userID uint) ([]relmodel.Friend, error) {
	return h.svc.ListFriends(userID)
}

// UpdateFriendRemark updates a friend's remark.
func (h *RelationHandler) UpdateFriendRemark(_ context.Context, userID, friendID uint, remark string) error {
	return h.svc.UpdateRemark(userID, friendID, remark)
}

// --- Groups ---

// CreateGroup creates a new group.
func (h *RelationHandler) CreateGroup(_ context.Context, name string, ownerID uint) (*relmodel.Group, error) {
	return h.svc.CreateGroup(name, ownerID)
}

// InviteMember invites a user to a group.
func (h *RelationHandler) InviteMember(_ context.Context, groupID, userID, operatorID uint) error {
	return h.svc.InviteMember(groupID, userID, operatorID)
}

// KickMember removes a user from a group.
func (h *RelationHandler) KickMember(_ context.Context, groupID, userID, operatorID uint) error {
	return h.svc.KickMember(groupID, userID, operatorID)
}

// MuteMember mutes a group member.
func (h *RelationHandler) MuteMember(_ context.Context, groupID, userID, operatorID uint, duration time.Duration) error {
	return h.svc.MuteMember(groupID, userID, operatorID, duration)
}

// TransferOwner transfers group ownership.
func (h *RelationHandler) TransferOwner(_ context.Context, groupID, newOwnerID, currentOwnerID uint) error {
	return h.svc.TransferOwner(groupID, newOwnerID, currentOwnerID)
}

// UpdateAnnouncement sets the group announcement.
func (h *RelationHandler) UpdateAnnouncement(_ context.Context, groupID uint, announcement string, operatorID uint) error {
	return h.svc.UpdateAnnouncement(groupID, announcement, operatorID)
}

// GetGroupMembers returns group members.
func (h *RelationHandler) GetGroupMembers(_ context.Context, groupID uint) ([]relmodel.GroupMember, error) {
	return h.svc.GetGroupMembers(groupID)
}

// GetGroup returns group info.
func (h *RelationHandler) GetGroup(_ context.Context, groupID uint) (*relmodel.Group, error) {
	return h.svc.GetGroup(groupID)
}
