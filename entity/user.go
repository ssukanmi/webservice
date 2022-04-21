package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key" json:"id"`
	FirstName string    `gorm:"type:varchar(255)" json:"first_name"`
	LastName  string    `gorm:"type:varchar(255)" json:"last_name"`
	Password  string    `gorm:"->;<-;not null" json:"-"`
	Username  string    `gorm:"unique;type:varchar(255)" json:"username"`
	CreatedAt time.Time `gorm:"column:account_created" json:"account_created"`
	UpdatedAt time.Time `gorm:"column:account_updated" json:"account_updated"`
	Verified  bool      `gorm:"type:bool" json:"verified"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("id", uuid.New())
	return nil
}
