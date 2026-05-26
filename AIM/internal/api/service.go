package api //业务逻辑层
import (
	"AIM/internal/common/utils"
	"AIM/internal/storage/mysql/model"
	"AIM/internal/storage/repo"
	"errors"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
func (s *Service) Register(username, password string) (*model.User, error) {
	existUser, err := repo.GetUserByUsername(username)
	if err == nil && existUser != nil {
		return nil, errors.New("用户名已存在")
	}
	user := &model.User{
		Username: username,
		Password: password,
		ID:       utils.UUID(),
		Avatar:   "",
	}
	err = repo.CreateUser(user)
	if err != nil {
		return nil, errors.New("注册失败")

	}
	return user, nil

}
func (s *Service) Login(username, password string) (*model.User, error) {
	user, err := repo.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.Password != password {
		return nil, errors.New("密码错误")
	}
	return user, nil

}
