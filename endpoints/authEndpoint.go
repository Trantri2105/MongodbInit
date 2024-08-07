package endpoints

import (
	requestDto "backend/dto/request"
	responseDto "backend/dto/response"
	"backend/mapper"
	"backend/service"
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator"
)

type AuthEndpoint interface {
	Signup() endpoint.Endpoint
	Login() endpoint.Endpoint
}

type authEndpointImpl struct {
	authService service.AuthService
}

var _ AuthEndpoint = authEndpointImpl{}

func (a authEndpointImpl) Signup() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		signupDto := request.(requestDto.SignupDto)
		validate := validator.New()
		err := validate.Struct(signupDto)
		if err != nil {
			err = err.(validator.ValidationErrors)
			return nil, err
		}

		user := mapper.SignupDtoToUser(signupDto)
		err = a.authService.RegisterUser(ctx, user)
		if err != nil {
			return nil, err
		}

		return responseDto.Response{Message: "User register successfully"}, nil
	}
}

func (a authEndpointImpl) Login() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		loginDto := request.(requestDto.LoginDto)
		validate := validator.New()
		err := validate.Struct(loginDto)
		if err != nil {
			err = err.(validator.ValidationErrors)
			return nil, err
		}

		var token string
		token, err = a.authService.LoginUser(ctx, loginDto)
				if err != nil {
			err = err.(validator.ValidationErrors)
			return nil, err
		}
		
		return responseDto.LoginResponse{Token: token}, nil
	}
}

func NewAuthEndpoint(authService service.AuthService) AuthEndpoint {
	return authEndpointImpl{
		authService: authService,
	}
}
