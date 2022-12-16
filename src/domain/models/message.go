package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/robertokbr/blinkchat/src/domain/enums"
)

type Message struct {
	Action string      `json:"action"`
	Data   MessageData `json:"data"`
}

type MessageData struct {
	ID        string            `json:"id"`
	Content   string            `json:"content"`
	Type      enums.MessageType `json:"type"`
	From      User              `json:"from"`
	CreatedAt time.Time         `json:"created_at"`
}

func NewMessage(content string, from User, messageType enums.MessageType, action string) *Message {
	id := uuid.NewString()

	message := MessageData{
		ID:        id,
		Content:   content,
		Type:      messageType,
		From:      from,
		CreatedAt: time.Now(),
	}

	return &Message{
		Action: action,
		Data:   message,
	}
}
