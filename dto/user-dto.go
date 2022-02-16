package dto

import "time"

type UserCreateDTO struct {
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	Username  string    `json:"username" binding:"required" validate:"required,email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type UserUpdateDTO struct {
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
