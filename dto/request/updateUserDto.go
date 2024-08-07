package dto

type UpdateUserDto struct {
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	Age 		int `json:"age"`
}
