package models

import (
	"github.com/google/uuid"
	"github.com/robertokbr/blinkchat/src/domain/dtos"
)

type User struct {
	ID        string `json:"id" gorm:"primaryKey type:uuid"`
	Name      string `json:"name" gorm:"type:varchar(255)"`
	ImageURL  string `json:"image_url" gorm:"type:varchar(255)"`
	Email     string `json:"email" gorm:"type:varchar(255) unique"`
	CreatedAt string `json:"created_at" gorm:"type:timestampz"`
	UpdatedAt string `json:"updated_at" gorm:"type:timestampz"`
}

// TODO: Add validation
func NewUser(data dtos.CreateConnection) *User {
	id := uuid.NewString()

	return &User{
		ID:       id,
		Name:     data.Name,
		ImageURL: data.ImageURL,
		Email:    data.Email,
	}
}
