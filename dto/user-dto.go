package dto

import "github.com/go-playground/validator/v10"

type UserCreateDTO struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Username  string `json:"username" binding:"required" validate:"email"`
}

type UserUpdateDTO struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password,omitempty"`
}

func (u UserCreateDTO) Validate() error {
	return validator.New().Struct(u)
}
