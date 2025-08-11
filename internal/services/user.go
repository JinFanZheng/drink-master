package services

import "github.com/ddteam/drink-master/internal/repositories"

// UserService 用户服务
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// TODO: 实现用户服务方法