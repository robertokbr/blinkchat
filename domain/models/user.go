package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/robertokbr/blinkchat/domain/dtos"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TODO: Add validation
func NewUser(data dtos.CreateConnection) *User {
	id := uuid.NewString()

	return &User{
		ID:        id,
		Name:      data.Name,
		ImageURL:  data.ImageURL,
		Email:     data.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
