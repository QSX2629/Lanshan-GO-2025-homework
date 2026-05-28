package repo

import (
	"testing"

	"github.com/aim/aim/internal/pkg/database"
	usermodel "github.com/aim/aim/internal/user/model"
)

func setupUserRepo(t *testing.T) *UserRepo {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := NewUserRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return repo
}

func TestUserRepo_CreateAndFind(t *testing.T) {
	repo := setupUserRepo(t)

	u := &usermodel.User{Username: "alice", Nickname: "Alice"}
	if err := repo.Create(u); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByUsername("alice")
	if err != nil {
		t.Fatalf("FindByUsername() error = %v", err)
	}
	if found.Nickname != "Alice" {
		t.Errorf("Nickname = %q, want Alice", found.Nickname)
	}

	found2, err := repo.FindByID(found.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found2.Username != "alice" {
		t.Errorf("Username = %q, want alice", found2.Username)
	}
}

func TestUserRepo_NotFound(t *testing.T) {
	repo := setupUserRepo(t)

	_, err := repo.FindByUsername("nobody")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestUserRepo_Update(t *testing.T) {
	repo := setupUserRepo(t)

	u := &usermodel.User{Username: "bob", Nickname: "Bob"}
	repo.Create(u)

	u.Nickname = "Bobby"
	if err := repo.Update(u); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	found, _ := repo.FindByID(u.ID)
	if found.Nickname != "Bobby" {
		t.Errorf("Nickname = %q, want Bobby", found.Nickname)
	}
}

func TestUserRepo_OnlineStatus(t *testing.T) {
	repo := setupUserRepo(t)

	online, err := repo.IsOnline(1)
	if err != nil {
		t.Fatalf("IsOnline() error = %v", err)
	}
	if online {
		t.Error("IsOnline should be false for unknown user")
	}

	if err := repo.UpsertOnlineStatus(1, true); err != nil {
		t.Fatalf("UpsertOnlineStatus() error = %v", err)
	}

	online, err = repo.IsOnline(1)
	if err != nil {
		t.Fatalf("IsOnline() error = %v", err)
	}
	if !online {
		t.Error("IsOnline should be true after upsert")
	}

	repo.UpsertOnlineStatus(1, false)
	online, _ = repo.IsOnline(1)
	if online {
		t.Error("IsOnline should be false after turning off")
	}
}
