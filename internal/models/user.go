package models

const UsersTable = "users"

type User struct {
	Base
	People
	PassportNumber string `gorm:"unique;not null" json:"passport_number" validate:"len=11"`
	Task           []Task `gorm:"foreignKey:UserId" json:"task"`
}

func (User) TableName() string {
	return UsersTable
}

type UserUpdateRequest struct {
	Name           string `json:"name" validate:"str_min_max=2 255"`
	Surname        string `json:"surname" validate:"str_min_max=2 255"`
	Patronymic     string `json:"patronymic" validate:"str_min_max=2 255"`
	Address        string `json:"address" validate:"str_min_max=3 255"`
	PassportNumber string `json:"passport_number" validate:"str_min_max=11"`
}

type UserCreateRequest struct {
	PassportNumber string `json:"passport_number" validate:"len=11"`
}

type UsersGetFilter struct {
	Name           string `json:"name" validate:"max=255"`
	Surname        string `json:"surname" validate:"max=255"`
	Patronymic     string `json:"patronymic" validate:"max=255"`
	PassportNumber string `json:"passport_number" validate:"max=11"`
	Address        string `json:"address" validate:"max=255"`
}

type UsersGetRequest struct {
	Filter UsersGetFilter `json:"filter"`
	Page   Page           `json:"page"`
}

type UsersResponse struct {
	Response[[]User]
	Page Page `json:"page"`
}
