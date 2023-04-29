package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID    string `json:"id" gorm:"primaryKey ;type:uuid;"`
	Token string `json:"token" gorm:"type:uuid;"`
	// references the user table
	UserID    string `json:"user_id" gorm:"type:uuid; references:users(id); on delete cascade;"`
	CreatedAt int64  `json:"expires_at"`
}

func NewSession(userID string) *Session {
	return &Session{
		ID:        uuid.NewString(),
		Token:     uuid.NewString(),
		UserID:    userID,
		CreatedAt: time.Now().UnixMilli(),
	}
}

func (s *Session) IsExpired() bool {
	return s.CreatedAt < time.Now().UnixMilli()-1000*60*60*3
}
