package repo

import (
	"AIM/internal/storage/mysql"
	"AIM/internal/storage/mysql/model"
)

func CreateUser(user *model.User) error {
	return mysql.DB.Create(user).Error
}
func GetUserByID(id string) (*model.User, error) {
	var user model.User
	err := mysql.DB.Where("id=?", id).First(&user).Error
	return &user, err
}
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := mysql.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
