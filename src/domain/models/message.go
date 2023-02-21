package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/robertokbr/blinkchat/src/domain/enums"
)

type Message struct {
	ID     string           `json:"id"`
	Action enums.ActionType `json:"action"`
	Data   MessageData      `json:"data"`
}

type MessageData struct {
	Content   string            `json:"content"`
	Type      enums.MessageType `json:"type"`
	From      *User             `json:"from"`
	CreatedAt time.Time         `json:"created_at"`
}

func NewMessage(
	content string,
	from *User,
	messageType enums.MessageType,
	action enums.ActionType,
) *Message {
	id := uuid.NewString()

	message := MessageData{
		Content:   content,
		Type:      messageType,
		From:      from,
		CreatedAt: time.Now(),
	}

	return &Message{
		ID:     id,
		Action: action,
		Data:   message,
	}
}
