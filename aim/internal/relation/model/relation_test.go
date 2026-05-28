package model

import (
	"testing"
	"time"
)

func TestFriendFields(t *testing.T) {
	f := Friend{
		UserID:    1,
		FriendID:  2,
		Remark:    "best friend",
		GroupName: "close",
	}

	if f.Remark != "best friend" {
		t.Errorf("Remark = %q, want best friend", f.Remark)
	}
	if f.GroupName != "close" {
		t.Errorf("GroupName = %q, want close", f.GroupName)
	}
}

func TestGroupRoles(t *testing.T) {
	if RoleOwner != "owner" || RoleAdmin != "admin" || RoleMember != "member" {
		t.Error("group role constants mismatch")
	}
}

func TestGroupFields(t *testing.T) {
	g := Group{
		ID:           1,
		Name:         "dev team",
		OwnerID:      100,
		Announcement: "welcome!",
	}

	if g.Name != "dev team" {
		t.Errorf("Name = %q, want dev team", g.Name)
	}
	if g.OwnerID != 100 {
		t.Errorf("OwnerID = %d, want 100", g.OwnerID)
	}
}

func TestGroupMemberFields(t *testing.T) {
	muteUntil := time.Now().Add(time.Hour)
	m := GroupMember{
		GroupID:    1,
		UserID:     2,
		Role:       RoleAdmin,
		Muted:      true,
		MutedUntil: &muteUntil,
	}

	if !m.Muted {
		t.Error("Muted should be true")
	}
	if m.Role != RoleAdmin {
		t.Errorf("Role = %q, want %q", m.Role, RoleAdmin)
	}
	if m.MutedUntil == nil || !m.MutedUntil.Equal(muteUntil) {
		t.Error("MutedUntil mismatch")
	}
}

func TestFriendGroup(t *testing.T) {
	fg := FriendGroup{ID: 1, UserID: 1, Name: "work"}
	if fg.Name != "work" {
		t.Errorf("Name = %q, want work", fg.Name)
	}
}
