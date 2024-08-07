package service

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtService interface {
	CreateToken(userId string, role string) (string, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

type jwtServiceImpl struct{}

func (jwtServiceImpl) CreateToken(userId string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"role":   role,
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("failed to generate access token")
	}
	return tokenString, nil
}

func (jwtServiceImpl) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}

func NewJwtService() JwtService {
	return jwtServiceImpl{}
}
