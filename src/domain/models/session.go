package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	UserID    string `json:"user_id"`
	CreatedAt int64  `json:"expires_at"`
}

func NewSession(userID string) *Session {
	return &Session{
		ID:        uuid.NewString(),
		Code:      uuid.NewString(),
		UserID:    userID,
		CreatedAt: time.Now().UnixMilli(),
	}
}

func (s *Session) IsExpired() bool {
	return s.CreatedAt < time.Now().UnixMilli()-1000*60*60*3
}
