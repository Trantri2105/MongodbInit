package mapper

import (
	dto "backend/dto/request"
	"backend/model"
)

func SignupDtoToUser(signupDto dto.SignupDto) model.User {
	user := model.User{
		FirstName:   signupDto.FirstName,
		LastName:    signupDto.LastName,
		Email:       signupDto.Email,
		Password:    signupDto.Password,
		PhoneNumber: signupDto.PhoneNumber,
		Age:         signupDto.Age,
		Role:        "user",
	}
	return user
}
