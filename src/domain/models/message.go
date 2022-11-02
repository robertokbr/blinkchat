package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/robertokbr/blinkchat/src/domain/enums"
)

type Message struct {
	ID        string            `json:"id"`
	Content   string            `json:"content"`
	Type      enums.MessageType `json:"type"`
	From      User              `json:"from"`
	CreatedAt time.Time         `json:"created_at"`
}

func NewMessage(content string, from User) *Message {
	id := uuid.NewString()

	return &Message{
		ID:        id,
		Content:   content,
		From:      from,
		CreatedAt: time.Now(),
	}
}
