package service

import (
	"errors"

	"github.com/aim/aim/internal/pkg/database"
	usermodel "github.com/aim/aim/internal/user/model"
	userrepo "github.com/aim/aim/internal/user/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

// UserService handles user business logic.
type UserService struct {
	repo *userrepo.UserRepo
}

// NewUserService creates a new UserService.
func NewUserService(db *database.DB) *UserService {
	return &UserService{repo: userrepo.NewUserRepo(db)}
}

// Register creates a new user account with hashed password.
func (s *UserService) Register(username, password, nickname string) (*usermodel.User, error) {
	existing, err := s.repo.FindByUsername(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &usermodel.User{
		Username:     username,
		PasswordHash: string(hash),
		Nickname:     nickname,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

// Login verifies credentials and returns the user.
func (s *UserService) Login(username, password string) (*usermodel.User, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrWrongPassword
	}
	return u, nil
}

// GetProfile returns a user's profile by ID.
func (s *UserService) GetProfile(userID uint) (*usermodel.User, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

// UpdateProfile updates user fields.
func (s *UserService) UpdateProfile(u *usermodel.User) error {
	return s.repo.Update(u)
}

// SetOnline sets the online status of a user.
func (s *UserService) SetOnline(userID uint, online bool) error {
	return s.repo.UpsertOnlineStatus(userID, online)
}

// IsOnline checks if a user is online.
func (s *UserService) IsOnline(userID uint) (bool, error) {
	return s.repo.IsOnline(userID)
}
