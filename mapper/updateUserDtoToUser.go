package mapper

import (
	dto "backend/dto/request"
	"backend/model"
)

func UpdateUserDtoToUser(updateUserDto dto.UpdateUserDto) model.User {
	return model.User{
		FirstName:   updateUserDto.FirstName,
		LastName:    updateUserDto.LastName,
		Password:    updateUserDto.Password,
		PhoneNumber: updateUserDto.PhoneNumber,
		Age:         updateUserDto.Age,
	}
}
