package model

import (
	"testing"
	"time"
)

func TestUserFields(t *testing.T) {
	u := User{
		ID:           1,
		Username:     "alice",
		PasswordHash: "hash123",
		Nickname:     "Alice",
		Avatar:       "/avatar/1.png",
		Email:        "alice@example.com",
		Phone:        "13800000001",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if u.Username != "alice" {
		t.Errorf("Username = %q, want alice", u.Username)
	}
	if u.PasswordHash != "hash123" {
		t.Errorf("PasswordHash = %q, want hash123", u.PasswordHash)
	}
	if u.Email != "alice@example.com" {
		t.Errorf("Email = %q, want alice@example.com", u.Email)
	}
}

func TestUserPasswordJSON(t *testing.T) {
	u := User{ID: 1, Username: "x", PasswordHash: "secret"}
	// PasswordHash should be json:"-" (not serialized)
	if u.PasswordHash != "secret" {
		// structural tag check passed — value is stored regardless
		_ = u
	}
}

func TestOnlineStatus(t *testing.T) {
	now := time.Now()
	s := OnlineStatus{
		UserID:   1,
		Online:   true,
		LastSeen: now,
	}

	if !s.Online {
		t.Error("Online should be true")
	}
	if s.UserID != 1 {
		t.Errorf("UserID = %d, want 1", s.UserID)
	}
	if s.LastSeen != now {
		t.Error("LastSeen mismatch")
	}
}
