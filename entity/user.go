package entity

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primary_key:auto_increment" json:"id"`
	FirstName string    `gorm:"type:varchar(255)" json:"first_name"`
	LastName  string    `gorm:"type:varchar(255)" json:"last_name"`
	Password  string    `gorm:"->;<-;not null" json:"-"`
	Username  string    `gorm:"uniqueIndex;type:varchar(255)" json:"username"`
	CreatedAt time.Time `gorm:"column:account_created" json:"account_created"`
	UpdatedAt time.Time `gorm:"column:account_updated" json:"account_updated"`
}
