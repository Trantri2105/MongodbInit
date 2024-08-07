package service

import (
	dto "backend/dto/request"
	"backend/model"
	"backend/repository"
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user model.User) error
	LoginUser(ctx context.Context, loginDto dto.LoginDto) (string, error)
}

type authServiceImpl struct {
	userRepository repository.UserRepository
	jwtService     JwtService
}

func (a authServiceImpl) RegisterUser(ctx context.Context, user model.User) error {
	_, err := a.userRepository.FindUserByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("user already existed")
	}

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hash)

	err = a.userRepository.InsertUser(ctx, user)
	if err != nil {
		log.Println(err.Error())
		return errors.New("failed to created user")
	}

	return nil
}

func (a authServiceImpl) LoginUser(ctx context.Context, loginDto dto.LoginDto) (string, error) {
	user, err := a.userRepository.FindUserByEmail(ctx, loginDto.Email)
	if err != nil {
		return "", errors.New("user not found with email: " + loginDto.Email)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))
	if err != nil {
		return "", errors.New("password incorrect")
	}
	return a.jwtService.CreateToken(user.ID.Hex(), user.Role)
}

func NewAuthService(userRepository repository.UserRepository, jwtService JwtService) AuthService {
	return authServiceImpl{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}
