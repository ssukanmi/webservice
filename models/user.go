package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
	CreatedAt time.Time `gorm:"column:account_created" json:"account_created"`
	UpdatedAt time.Time `gorm:"column:account_updated" json:"account_updated"`
}

// type UserOutput struct {
// 	ID        string    `json:"id"`
// 	FirstName string    `json:"first_name"`
// 	LastName  string    `json:"last_name"`
// 	Username  string    `json:"username"`
// 	CreatedAt time.Time `json:"account_created"`
// 	UpdatedAt time.Time `json:"account_updated"`
// }

// func (u *User) CreateUserOutput() (uo UserOutput) {
// 	uo.ID = u.ID
// 	uo.FirstName = u.FirstName
// 	uo.LastName = u.LastName
// 	uo.Username = u.Username
// 	uo.CreatedAt = u.CreatedAt
// 	uo.UpdatedAt = u.UpdatedAt
// 	return
// }
