package endpoints

import (
	requestDto "backend/dto/request"
	responseDto "backend/dto/response"
	"backend/mapper"
	"backend/middleware"
	"backend/service"
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type UserEndpoint interface {
	GetUserInfo() endpoint.Endpoint
	UpdateUserInfo() endpoint.Endpoint
	DeleteUser() endpoint.Endpoint
}

type userEndpointImpl struct {
	userService service.UserService
}

var _ UserEndpoint = userEndpointImpl{}

func (u userEndpointImpl) GetUserInfo() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userId := ctx.Value(middleware.UserIdContextKey{}).(string)
		user, err := u.userService.GetUserInfoById(ctx, userId)
		if err != nil {
			return nil, err
		}
		userDto := mapper.UserToUserDto(user)
		return userDto, nil
	}
}

func (u userEndpointImpl) UpdateUserInfo() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userId := ctx.Value(middleware.UserIdContextKey{}).(string)
		updateUserDto := request.(requestDto.UpdateUserDto)
		err := u.userService.UpdateUserInfo(ctx, userId, mapper.UpdateUserDtoToUser(updateUserDto))
		if err != nil {
			return nil, errors.New("failed to update user")
		}
		return responseDto.Response{Message: "user updated"}, nil
	}
}

func (u userEndpointImpl) DeleteUser() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userId := request.(string)
		err := u.userService.DeleteUserById(ctx, userId)
		if err != nil {
			return nil, errors.New("failed to delete user")
		}
		return responseDto.Response{Message: "user deleted"}, nil
	}
}

func NewUserEndpoint(userService service.UserService) UserEndpoint {
	return userEndpointImpl{
		userService: userService,
	}
}
