package repo

import (
	"github.com/aim/aim/internal/pkg/database"
	relmodel "github.com/aim/aim/internal/relation/model"
	"gorm.io/gorm"
)

// RelationRepo handles friend and group data access.
type RelationRepo struct {
	db *database.DB
}

// NewRelationRepo creates a new RelationRepo.
func NewRelationRepo(db *database.DB) *RelationRepo {
	return &RelationRepo{db: db}
}

// AutoMigrate creates the relation tables.
func (r *RelationRepo) AutoMigrate() error {
	return r.db.AutoMigrate(
		&relmodel.Friend{},
		&relmodel.FriendGroup{},
		&relmodel.Group{},
		&relmodel.GroupMember{},
	)
}

// --- Friends ---

// AddFriend creates a friendship record.
func (r *RelationRepo) AddFriend(f *relmodel.Friend) error {
	return r.db.Create(f).Error
}

// DeleteFriend removes a friendship.
func (r *RelationRepo) DeleteFriend(userID, friendID uint) error {
	return r.db.Where("user_id = ? AND friend_id = ?", userID, friendID).Delete(&relmodel.Friend{}).Error
}

// IsFriend checks if two users are friends.
func (r *RelationRepo) IsFriend(userID, friendID uint) (bool, error) {
	var count int64
	err := r.db.Model(&relmodel.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Count(&count).Error
	return count > 0, err
}

// ListFriends returns all friends for a user.
func (r *RelationRepo) ListFriends(userID uint) ([]relmodel.Friend, error) {
	var friends []relmodel.Friend
	err := r.db.Where("user_id = ?", userID).Find(&friends).Error
	return friends, err
}

// UpdateFriendRemark updates the remark for a friend.
func (r *RelationRepo) UpdateFriendRemark(userID, friendID uint, remark string) error {
	return r.db.Model(&relmodel.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Update("remark", remark).Error
}

// --- Groups ---

// CreateGroup inserts a new group.
func (r *RelationRepo) CreateGroup(g *relmodel.Group) error {
	return r.db.Create(g).Error
}

// FindGroupByID looks up a group.
func (r *RelationRepo) FindGroupByID(id uint) (*relmodel.Group, error) {
	var g relmodel.Group
	err := r.db.First(&g, id).Error
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// UpdateGroup updates group fields.
func (r *RelationRepo) UpdateGroup(g *relmodel.Group) error {
	return r.db.Save(g).Error
}

// ListUserGroups returns all groups a user belongs to.
func (r *RelationRepo) ListUserGroups(userID uint) ([]relmodel.Group, error) {
	var groupIDs []uint
	r.db.Model(&relmodel.GroupMember{}).
		Where("user_id = ?", userID).
		Pluck("group_id", &groupIDs)
	if len(groupIDs) == 0 {
		return nil, nil
	}
	var groups []relmodel.Group
	err := r.db.Where("id IN ?", groupIDs).Find(&groups).Error
	return groups, err
}

// --- Group Members ---

// AddGroupMember adds a user to a group.
func (r *RelationRepo) AddGroupMember(m *relmodel.GroupMember) error {
	return r.db.Create(m).Error
}

// RemoveGroupMember removes a user from a group.
func (r *RelationRepo) RemoveGroupMember(groupID, userID uint) error {
	return r.db.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&relmodel.GroupMember{}).Error
}

// GetGroupMember returns a member record.
func (r *RelationRepo) GetGroupMember(groupID, userID uint) (*relmodel.GroupMember, error) {
	var m relmodel.GroupMember
	err := r.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ListGroupMembers returns all members of a group.
func (r *RelationRepo) ListGroupMembers(groupID uint) ([]relmodel.GroupMember, error) {
	var members []relmodel.GroupMember
	err := r.db.Where("group_id = ?", groupID).Find(&members).Error
	return members, err
}

// UpdateGroupMember updates a member's role or mute status.
func (r *RelationRepo) UpdateGroupMember(m *relmodel.GroupMember) error {
	return r.db.Save(m).Error
}

// IsGroupMember checks if a user is in a group.
func (r *RelationRepo) IsGroupMember(groupID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&relmodel.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	return count > 0, err
}

// Ensure import is used.
var _ = gorm.ErrRecordNotFound
