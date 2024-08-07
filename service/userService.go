package service

import (
	"backend/model"
	"backend/repository"
	"context"
)

type UserService interface {
	GetUserInfoById(ctx context.Context, userId string) (model.User, error)
	UpdateUserInfo(ctx context.Context, userId string, user model.User) error
	DeleteUserById(ctx context.Context, userId string) error
}

type userServiceImpl struct {
	userRepository repository.UserRepository
}

func (u userServiceImpl) GetUserInfoById(ctx context.Context, userId string) (model.User, error) {
	return u.userRepository.FindUserById(ctx, userId)
}

func (u userServiceImpl) UpdateUserInfo(ctx context.Context, userId string, user model.User) error {
	return u.userRepository.UpdateUserById(ctx, userId, user)
}

func (u userServiceImpl) DeleteUserById(ctx context.Context, userId string) error {
	return u.userRepository.DeleteUserById(ctx, userId)
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return userServiceImpl{
		userRepository: userRepository,
	}
}
