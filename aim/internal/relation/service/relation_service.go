package service

import (
	"errors"
	"time"

	"github.com/aim/aim/internal/pkg/database"
	relmodel "github.com/aim/aim/internal/relation/model"
	relrepo "github.com/aim/aim/internal/relation/repo"
)

var (
	ErrAlreadyFriend  = errors.New("already a friend")
	ErrNotFriend      = errors.New("not a friend")
	ErrGroupNotFound  = errors.New("group not found")
	ErrNotGroupMember = errors.New("not a group member")
	ErrNotGroupOwner  = errors.New("not the group owner")
	ErrNotGroupAdmin  = errors.New("not a group admin")
)

// RelationService handles friend and group business logic.
type RelationService struct {
	repo *relrepo.RelationRepo
}

// NewRelationService creates a new RelationService.
func NewRelationService(db *database.DB) *RelationService {
	return &RelationService{repo: relrepo.NewRelationRepo(db)}
}

// --- Friends ---

// AddFriend adds a new friend.
func (s *RelationService) AddFriend(userID, friendID uint, remark string) error {
	isFriend, err := s.repo.IsFriend(userID, friendID)
	if err != nil {
		return err
	}
	if isFriend {
		return ErrAlreadyFriend
	}

	f := &relmodel.Friend{
		UserID:    userID,
		FriendID:  friendID,
		Remark:    remark,
		CreatedAt: time.Now(),
	}
	return s.repo.AddFriend(f)
}

// DeleteFriend removes a friend.
func (s *RelationService) DeleteFriend(userID, friendID uint) error {
	isFriend, err := s.repo.IsFriend(userID, friendID)
	if err != nil {
		return err
	}
	if !isFriend {
		return ErrNotFriend
	}
	return s.repo.DeleteFriend(userID, friendID)
}

// ListFriends returns all friends for a user.
func (s *RelationService) ListFriends(userID uint) ([]relmodel.Friend, error) {
	return s.repo.ListFriends(userID)
}

// UpdateRemark updates the remark for a friend.
func (s *RelationService) UpdateRemark(userID, friendID uint, remark string) error {
	return s.repo.UpdateFriendRemark(userID, friendID, remark)
}

// --- Groups ---

// CreateGroup creates a new group with the creator as owner.
func (s *RelationService) CreateGroup(name string, ownerID uint) (*relmodel.Group, error) {
	g := &relmodel.Group{
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreateGroup(g); err != nil {
		return nil, err
	}

	// Add creator as owner member.
	m := &relmodel.GroupMember{
		GroupID:   g.ID,
		UserID:    ownerID,
		Role:      relmodel.RoleOwner,
		CreatedAt: time.Now(),
	}
	if err := s.repo.AddGroupMember(m); err != nil {
		return nil, err
	}

	return g, nil
}

// InviteMember adds a user to a group. Caller must be admin/owner.
func (s *RelationService) InviteMember(groupID, userID, operatorID uint) error {
	if err := s.requireAdmin(groupID, operatorID); err != nil {
		return err
	}

	isMember, err := s.repo.IsGroupMember(groupID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return nil // Already a member, no error.
	}

	m := &relmodel.GroupMember{
		GroupID:   groupID,
		UserID:    userID,
		Role:      relmodel.RoleMember,
		CreatedAt: time.Now(),
	}
	return s.repo.AddGroupMember(m)
}

// KickMember removes a user from a group.
func (s *RelationService) KickMember(groupID, userID, operatorID uint) error {
	if err := s.requireAdmin(groupID, operatorID); err != nil {
		return err
	}

	isMember, err := s.repo.IsGroupMember(groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrNotGroupMember
	}

	// Cannot kick the owner.
	member, _ := s.repo.GetGroupMember(groupID, userID)
	if member != nil && member.Role == relmodel.RoleOwner {
		return errors.New("cannot kick the group owner")
	}

	return s.repo.RemoveGroupMember(groupID, userID)
}

// MuteMember mutes a group member.
func (s *RelationService) MuteMember(groupID, userID, operatorID uint, duration time.Duration) error {
	if err := s.requireAdmin(groupID, operatorID); err != nil {
		return err
	}

	member, err := s.repo.GetGroupMember(groupID, userID)
	if err != nil {
		return ErrNotGroupMember
	}

	member.Muted = true
	muteUntil := time.Now().Add(duration)
	member.MutedUntil = &muteUntil
	return s.repo.UpdateGroupMember(member)
}

// TransferOwner transfers group ownership.
func (s *RelationService) TransferOwner(groupID, newOwnerID, currentOwnerID uint) error {
	g, err := s.repo.FindGroupByID(groupID)
	if err != nil {
		return ErrGroupNotFound
	}
	if g.OwnerID != currentOwnerID {
		return ErrNotGroupOwner
	}

	isMember, err := s.repo.IsGroupMember(groupID, newOwnerID)
	if err != nil || !isMember {
		return ErrNotGroupMember
	}

	// Update new owner's role.
	newOwner, _ := s.repo.GetGroupMember(groupID, newOwnerID)
	newOwner.Role = relmodel.RoleOwner
	s.repo.UpdateGroupMember(newOwner)

	// Demote current owner to admin.
	oldOwner, _ := s.repo.GetGroupMember(groupID, currentOwnerID)
	if oldOwner != nil {
		oldOwner.Role = relmodel.RoleAdmin
		s.repo.UpdateGroupMember(oldOwner)
	}

	g.OwnerID = newOwnerID
	return s.repo.UpdateGroup(g)
}

// UpdateAnnouncement sets the group announcement.
func (s *RelationService) UpdateAnnouncement(groupID uint, announcement string, operatorID uint) error {
	if err := s.requireAdmin(groupID, operatorID); err != nil {
		return err
	}

	g, err := s.repo.FindGroupByID(groupID)
	if err != nil {
		return ErrGroupNotFound
	}
	g.Announcement = announcement
	return s.repo.UpdateGroup(g)
}

// GetGroupMembers returns all members of a group.
func (s *RelationService) GetGroupMembers(groupID uint) ([]relmodel.GroupMember, error) {
	return s.repo.ListGroupMembers(groupID)
}

// GetGroup returns group info.
func (s *RelationService) GetGroup(groupID uint) (*relmodel.Group, error) {
	g, err := s.repo.FindGroupByID(groupID)
	if err != nil {
		return nil, ErrGroupNotFound
	}
	return g, nil
}

// requireAdmin checks if a user is admin or owner of the group.
func (s *RelationService) requireAdmin(groupID, userID uint) error {
	member, err := s.repo.GetGroupMember(groupID, userID)
	if err != nil {
		return ErrNotGroupMember
	}
	if member.Role != relmodel.RoleOwner && member.Role != relmodel.RoleAdmin {
		return ErrNotGroupAdmin
	}
	return nil
}
