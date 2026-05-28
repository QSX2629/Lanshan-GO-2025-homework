package handler

import (
	"testing"

	"github.com/aim/aim/internal/pkg/database"
	usermodel "github.com/aim/aim/internal/user/model"
	userrepo "github.com/aim/aim/internal/user/repo"
)

func setupUserHandler(t *testing.T) *UserHandler {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := userrepo.NewUserRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return NewUserHandler(db)
}

func TestUserHandler_RegisterAndLogin(t *testing.T) {
	h := setupUserHandler(t)

	resp, err := h.Register(nil, &RegisterRequest{Username: "alice", Password: "pass", Nickname: "Alice"})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if resp.Username != "alice" {
		t.Errorf("Username = %q, want alice", resp.Username)
	}

	resp2, err := h.Login(nil, &LoginRequest{Username: "alice", Password: "pass"})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if resp2.ID != resp.ID {
		t.Errorf("ID mismatch after login")
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	h := setupUserHandler(t)
	reg, _ := h.Register(nil, &RegisterRequest{Username: "bob", Password: "pw", Nickname: "Bob"})

	profile, err := h.GetProfile(nil, reg.ID)
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if profile.Nickname != "Bob" {
		t.Errorf("Nickname = %q, want Bob", profile.Nickname)
	}
}

func TestUserHandler_SetOnline(t *testing.T) {
	h := setupUserHandler(t)

	if err := h.SetOnline(nil, 1, true); err != nil {
		t.Fatalf("SetOnline() error = %v", err)
	}
}

func TestToUserResponse(t *testing.T) {
	u := setupUserHandler(t)
	reg, _ := u.Register(nil, &RegisterRequest{Username: "test", Password: "pw", Nickname: "Test"})

	r := toUserResponse(&usermodel.User{
		ID: reg.ID, Username: reg.Username, Nickname: reg.Nickname,
		Avatar: reg.Avatar, Email: reg.Email, Phone: reg.Phone,
	})
	if r.ID != reg.ID {
		t.Errorf("ID = %d, want %d", r.ID, reg.ID)
	}
}
