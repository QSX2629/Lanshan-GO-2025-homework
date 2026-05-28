package service

import (
	"testing"
	"time"

	"github.com/aim/aim/internal/pkg/database"
	relrepo "github.com/aim/aim/internal/relation/repo"
)

func setupRelationService(t *testing.T) *RelationService {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := relrepo.NewRelationRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return NewRelationService(db)
}

func TestRelationService_Friend(t *testing.T) {
	svc := setupRelationService(t)

	if err := svc.AddFriend(1, 2, "best"); err != nil {
		t.Fatalf("AddFriend() error = %v", err)
	}

	// Duplicate.
	if err := svc.AddFriend(1, 2, "dup"); err != ErrAlreadyFriend {
		t.Errorf("duplicate error = %v, want ErrAlreadyFriend", err)
	}

	friends, err := svc.ListFriends(1)
	if err != nil {
		t.Fatalf("ListFriends() error = %v", err)
	}
	if len(friends) != 1 {
		t.Errorf("len(friends) = %d, want 1", len(friends))
	}

	if err := svc.UpdateRemark(1, 2, "new remark"); err != nil {
		t.Fatalf("UpdateRemark() error = %v", err)
	}

	if err := svc.DeleteFriend(1, 2); err != nil {
		t.Fatalf("DeleteFriend() error = %v", err)
	}

	if err := svc.DeleteFriend(1, 2); err != ErrNotFriend {
		t.Errorf("delete non-friend error = %v, want ErrNotFriend", err)
	}
}

func TestRelationService_Group(t *testing.T) {
	svc := setupRelationService(t)

	g, err := svc.CreateGroup("dev", 1)
	if err != nil {
		t.Fatalf("CreateGroup() error = %v", err)
	}
	if g.Name != "dev" {
		t.Errorf("Name = %q, want dev", g.Name)
	}

	found, err := svc.GetGroup(g.ID)
	if err != nil {
		t.Fatalf("GetGroup() error = %v", err)
	}
	if found.OwnerID != 1 {
		t.Errorf("OwnerID = %d, want 1", found.OwnerID)
	}

	members, err := svc.GetGroupMembers(g.ID)
	if err != nil {
		t.Fatalf("GetGroupMembers() error = %v", err)
	}
	if len(members) != 1 {
		t.Errorf("len(members) = %d, want 1", len(members))
	}
}

func TestRelationService_InviteAndKick(t *testing.T) {
	svc := setupRelationService(t)
	g, _ := svc.CreateGroup("team", 1)

	// Owner invites.
	if err := svc.InviteMember(g.ID, 2, 1); err != nil {
		t.Fatalf("InviteMember() error = %v", err)
	}

	members, _ := svc.GetGroupMembers(g.ID)
	if len(members) != 2 {
		t.Errorf("len(members) = %d, want 2", len(members))
	}

	// Non-admin cannot kick.
	if err := svc.KickMember(g.ID, 3, 2); err != ErrNotGroupAdmin {
		t.Errorf("non-admin kick error = %v, want ErrNotGroupAdmin", err)
	}

	// Owner kicks member.
	if err := svc.KickMember(g.ID, 2, 1); err != nil {
		t.Fatalf("KickMember() error = %v", err)
	}

	members, _ = svc.GetGroupMembers(g.ID)
	if len(members) != 1 {
		t.Errorf("len(members) after kick = %d, want 1", len(members))
	}
}

func TestRelationService_Mute(t *testing.T) {
	svc := setupRelationService(t)
	g, _ := svc.CreateGroup("team", 1)
	svc.InviteMember(g.ID, 2, 1)

	if err := svc.MuteMember(g.ID, 2, 1, time.Hour); err != nil {
		t.Fatalf("MuteMember() error = %v", err)
	}
}

func TestRelationService_TransferOwner(t *testing.T) {
	svc := setupRelationService(t)
	g, _ := svc.CreateGroup("team", 1)
	svc.InviteMember(g.ID, 2, 1)

	// Non-owner cannot transfer.
	if err := svc.TransferOwner(g.ID, 2, 2); err != ErrNotGroupOwner {
		t.Errorf("non-owner transfer error = %v, want ErrNotGroupOwner", err)
	}

	if err := svc.TransferOwner(g.ID, 2, 1); err != nil {
		t.Fatalf("TransferOwner() error = %v", err)
	}

	g2, _ := svc.GetGroup(g.ID)
	if g2.OwnerID != 2 {
		t.Errorf("OwnerID after transfer = %d, want 2", g2.OwnerID)
	}
}

func TestRelationService_Announcement(t *testing.T) {
	svc := setupRelationService(t)
	g, _ := svc.CreateGroup("team", 1)

	if err := svc.UpdateAnnouncement(g.ID, "hello world", 1); err != nil {
		t.Fatalf("UpdateAnnouncement() error = %v", err)
	}

	g2, _ := svc.GetGroup(g.ID)
	if g2.Announcement != "hello world" {
		t.Errorf("Announcement = %q, want hello world", g2.Announcement)
	}
}
