package repo

import (
	"testing"

	"github.com/aim/aim/internal/pkg/database"
	relmodel "github.com/aim/aim/internal/relation/model"
)

func setupRelationRepo(t *testing.T) *RelationRepo {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := NewRelationRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return repo
}

func TestRelationRepo_Friend(t *testing.T) {
	repo := setupRelationRepo(t)

	f := &relmodel.Friend{UserID: 1, FriendID: 2, Remark: "best"}
	if err := repo.AddFriend(f); err != nil {
		t.Fatalf("AddFriend() error = %v", err)
	}

	isFriend, err := repo.IsFriend(1, 2)
	if err != nil {
		t.Fatalf("IsFriend() error = %v", err)
	}
	if !isFriend {
		t.Error("IsFriend should be true")
	}

	friends, err := repo.ListFriends(1)
	if err != nil {
		t.Fatalf("ListFriends() error = %v", err)
	}
	if len(friends) != 1 {
		t.Errorf("len(friends) = %d, want 1", len(friends))
	}

	if err := repo.UpdateFriendRemark(1, 2, "updated"); err != nil {
		t.Fatalf("UpdateFriendRemark() error = %v", err)
	}

	if err := repo.DeleteFriend(1, 2); err != nil {
		t.Fatalf("DeleteFriend() error = %v", err)
	}

	friends, _ = repo.ListFriends(1)
	if len(friends) != 0 {
		t.Errorf("len(friends) after delete = %d, want 0", len(friends))
	}
}

func TestRelationRepo_Group(t *testing.T) {
	repo := setupRelationRepo(t)

	g := &relmodel.Group{Name: "dev", OwnerID: 1}
	if err := repo.CreateGroup(g); err != nil {
		t.Fatalf("CreateGroup() error = %v", err)
	}

	found, err := repo.FindGroupByID(g.ID)
	if err != nil {
		t.Fatalf("FindGroupByID() error = %v", err)
	}
	if found.Name != "dev" {
		t.Errorf("Name = %q, want dev", found.Name)
	}

	g.Announcement = "hello team"
	if err := repo.UpdateGroup(g); err != nil {
		t.Fatalf("UpdateGroup() error = %v", err)
	}
}

func TestRelationRepo_GroupMember(t *testing.T) {
	repo := setupRelationRepo(t)

	g := &relmodel.Group{Name: "dev", OwnerID: 1}
	repo.CreateGroup(g)

	m := &relmodel.GroupMember{GroupID: g.ID, UserID: 2, Role: relmodel.RoleMember}
	if err := repo.AddGroupMember(m); err != nil {
		t.Fatalf("AddGroupMember() error = %v", err)
	}

	isMember, err := repo.IsGroupMember(g.ID, 2)
	if err != nil {
		t.Fatalf("IsGroupMember() error = %v", err)
	}
	if !isMember {
		t.Error("IsGroupMember should be true")
	}

	members, err := repo.ListGroupMembers(g.ID)
	if err != nil {
		t.Fatalf("ListGroupMembers() error = %v", err)
	}
	if len(members) != 1 {
		t.Errorf("len(members) = %d, want 1", len(members))
	}

	member, err := repo.GetGroupMember(g.ID, 2)
	if err != nil {
		t.Fatalf("GetGroupMember() error = %v", err)
	}
	member.Role = relmodel.RoleAdmin
	if err := repo.UpdateGroupMember(member); err != nil {
		t.Fatalf("UpdateGroupMember() error = %v", err)
	}

	groups, err := repo.ListUserGroups(2)
	if err != nil {
		t.Fatalf("ListUserGroups() error = %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("len(groups) = %d, want 1", len(groups))
	}

	if err := repo.RemoveGroupMember(g.ID, 2); err != nil {
		t.Fatalf("RemoveGroupMember() error = %v", err)
	}

	isMember, _ = repo.IsGroupMember(g.ID, 2)
	if isMember {
		t.Error("IsGroupMember should be false after removal")
	}
}
