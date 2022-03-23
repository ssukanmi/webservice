package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserImage struct {
	ID        uuid.UUID `gorm:"primary_key" json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FileName  string    `gorm:"type:varchar(255)" json:"file_name"`
	URL       string    `gorm:"type:varchar(255)" json:"url"`
	UpdatedAt time.Time `gorm:"column:upload_date" json:"upload_date"`
}

func (image *UserImage) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("id", uuid.New())
	return nil
}
