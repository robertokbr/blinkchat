package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/robertokbr/blinkchat/src/domain/enums"
)

type Message struct {
	ID     string               `json:"id"`
	Action enums.WebsocketEvent `json:"action"`
	Data   MessageData          `json:"data"`
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
	action enums.WebsocketEvent,
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

func NewPingMessage(user *User) *Message {
	return NewMessage(
		"Ping",
		user,
		enums.TEXT,
		enums.BROADCASTING,
	)
}

func NewUserDisconnectedMessage(user *User) *Message {
	return NewMessage(
		"User has disconnected",
		user,
		enums.TEXT,
		enums.DISCONNECTED,
	)
}

func NewUserConnectedMessage(user *User) *Message {
	return NewMessage(
		"New user joined",
		user,
		enums.TEXT,
		enums.CONNECTED,
	)
}

func NewUserUnmatchedMessage(user *User) *Message {
	return NewMessage(
		"User has been unmatched",
		user,
		enums.TEXT,
		enums.UNMATCHED,
	)
}

func NewUserMatchedMessage(user *User) *Message {
	return NewMessage(
		"You have been matched with a new user",
		user,
		enums.TEXT,
		enums.MATCHED,
	)
}
