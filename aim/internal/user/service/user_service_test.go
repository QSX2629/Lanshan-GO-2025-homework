package service

import (
	"testing"

	"github.com/aim/aim/internal/pkg/database"
	userrepo "github.com/aim/aim/internal/user/repo"
)

func setupUserService(t *testing.T) *UserService {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	// Ensure table exists.
	repo := userrepo.NewUserRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return NewUserService(db)
}

func TestUserService_Register(t *testing.T) {
	svc := setupUserService(t)

	u, err := svc.Register("alice", "pass123", "Alice")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if u.Username != "alice" {
		t.Errorf("Username = %q, want alice", u.Username)
	}
	if u.PasswordHash == "" || u.PasswordHash == "pass123" {
		t.Error("password should be hashed")
	}

	// Duplicate registration.
	_, err = svc.Register("alice", "pass123", "Alice2")
	if err != ErrUserExists {
		t.Errorf("duplicate register error = %v, want ErrUserExists", err)
	}
}

func TestUserService_Login(t *testing.T) {
	svc := setupUserService(t)
	svc.Register("bob", "secret", "Bob")

	u, err := svc.Login("bob", "secret")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if u.Username != "bob" {
		t.Errorf("Username = %q, want bob", u.Username)
	}

	_, err = svc.Login("bob", "wrong")
	if err != ErrWrongPassword {
		t.Errorf("wrong password error = %v, want ErrWrongPassword", err)
	}

	_, err = svc.Login("nobody", "x")
	if err != ErrUserNotFound {
		t.Errorf("not found error = %v, want ErrUserNotFound", err)
	}
}

func TestUserService_Profile(t *testing.T) {
	svc := setupUserService(t)
	u, _ := svc.Register("carol", "pw", "Carol")

	p, err := svc.GetProfile(u.ID)
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if p.Nickname != "Carol" {
		t.Errorf("Nickname = %q, want Carol", p.Nickname)
	}

	p.Nickname = "Caroline"
	if err := svc.UpdateProfile(p); err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
}

func TestUserService_Online(t *testing.T) {
	svc := setupUserService(t)

	if err := svc.SetOnline(1, true); err != nil {
		t.Fatalf("SetOnline() error = %v", err)
	}

	online, err := svc.IsOnline(1)
	if err != nil {
		t.Fatalf("IsOnline() error = %v", err)
	}
	if !online {
		t.Error("user should be online")
	}
}
