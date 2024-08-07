package mapper

import (
	dto "backend/dto/response"
	"backend/model"
)

func UserToUserDto(user model.User) dto.UserDto {
	userDto := dto.UserDto{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Age:         user.Age,
		Role:        user.Role,
	}
	return userDto
}
