package repo

import (
	"github.com/aim/aim/internal/pkg/database"
	usermodel "github.com/aim/aim/internal/user/model"
	"gorm.io/gorm"
)

// UserRepo handles user data access.
type UserRepo struct {
	db *database.DB
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(db *database.DB) *UserRepo {
	return &UserRepo{db: db}
}

// AutoMigrate creates the user tables.
func (r *UserRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&usermodel.User{}, &usermodel.OnlineStatus{})
}

// Create inserts a new user.
func (r *UserRepo) Create(u *usermodel.User) error {
	return r.db.Create(u).Error
}

// FindByUsername looks up a user by username.
func (r *UserRepo) FindByUsername(username string) (*usermodel.User, error) {
	var u usermodel.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByID looks up a user by ID.
func (r *UserRepo) FindByID(id uint) (*usermodel.User, error) {
	var u usermodel.User
	err := r.db.First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Update updates user fields.
func (r *UserRepo) Update(u *usermodel.User) error {
	return r.db.Save(u).Error
}

// UpsertOnlineStatus sets or updates the online status for a user.
func (r *UserRepo) UpsertOnlineStatus(userID uint, online bool) error {
	status := usermodel.OnlineStatus{UserID: userID}
	r.db.Where("user_id = ?", userID).FirstOrCreate(&status)
	status.Online = online
	return r.db.Save(&status).Error
}

// IsOnline checks if a user is online.
func (r *UserRepo) IsOnline(userID uint) (bool, error) {
	var status usermodel.OnlineStatus
	err := r.db.Where("user_id = ?", userID).First(&status).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return status.Online, nil
}
