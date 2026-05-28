package handler

import (
	"testing"

	"github.com/aim/aim/internal/pkg/database"
	relrepo "github.com/aim/aim/internal/relation/repo"
)

func setupRelationHandler(t *testing.T) *RelationHandler {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := relrepo.NewRelationRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return NewRelationHandler(db)
}

func TestRelationHandler_Friend(t *testing.T) {
	h := setupRelationHandler(t)

	if err := h.AddFriend(nil, 1, 2, "best"); err != nil {
		t.Fatalf("AddFriend() error = %v", err)
	}

	friends, err := h.ListFriends(nil, 1)
	if err != nil {
		t.Fatalf("ListFriends() error = %v", err)
	}
	if len(friends) != 1 {
		t.Errorf("len(friends) = %d, want 1", len(friends))
	}

	if err := h.UpdateFriendRemark(nil, 1, 2, "updated"); err != nil {
		t.Fatalf("UpdateFriendRemark() error = %v", err)
	}

	if err := h.DeleteFriend(nil, 1, 2); err != nil {
		t.Fatalf("DeleteFriend() error = %v", err)
	}
}

func TestRelationHandler_Group(t *testing.T) {
	h := setupRelationHandler(t)

	g, err := h.CreateGroup(nil, "dev", 1)
	if err != nil {
		t.Fatalf("CreateGroup() error = %v", err)
	}

	got, err := h.GetGroup(nil, g.ID)
	if err != nil {
		t.Fatalf("GetGroup() error = %v", err)
	}
	if got.Name != "dev" {
		t.Errorf("Name = %q, want dev", got.Name)
	}
}

func TestRelationHandler_Members(t *testing.T) {
	h := setupRelationHandler(t)
	g, _ := h.CreateGroup(nil, "team", 1)

	if err := h.InviteMember(nil, g.ID, 2, 1); err != nil {
		t.Fatalf("InviteMember() error = %v", err)
	}

	members, err := h.GetGroupMembers(nil, g.ID)
	if err != nil {
		t.Fatalf("GetGroupMembers() error = %v", err)
	}
	if len(members) != 2 {
		t.Errorf("len(members) = %d, want 2", len(members))
	}

	if err := h.KickMember(nil, g.ID, 2, 1); err != nil {
		t.Fatalf("KickMember() error = %v", err)
	}
}

func TestRelationHandler_TransferOwner(t *testing.T) {
	h := setupRelationHandler(t)
	g, _ := h.CreateGroup(nil, "team", 1)
	h.InviteMember(nil, g.ID, 2, 1)

	if err := h.TransferOwner(nil, g.ID, 2, 1); err != nil {
		t.Fatalf("TransferOwner() error = %v", err)
	}
}

func TestRelationHandler_Announcement(t *testing.T) {
	h := setupRelationHandler(t)
	g, _ := h.CreateGroup(nil, "team", 1)

	if err := h.UpdateAnnouncement(nil, g.ID, "hello", 1); err != nil {
		t.Fatalf("UpdateAnnouncement() error = %v", err)
	}
}
