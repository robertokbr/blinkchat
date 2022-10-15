package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/robertokbr/blinkchat/domain/enums"
)

type User struct {
	ID           string          `json:"id"`
	ConnectionID string          `json:"connection_id"`
	Name         string          `json:"name"`
	ImageURL     string          `json:"image_url"`
	Email        string          `json:"email"`
	State        enums.UserState `json:"state"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// TODO: Add validation
func NewUser(connectionID, name, imageURL, email string, state enums.UserState, createdAt, updatedAt time.Time) *User {
	id := uuid.NewString()

	return &User{
		ID:           id,
		ConnectionID: connectionID,
		Name:         name,
		ImageURL:     imageURL,
		Email:        email,
		State:        state,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}
