package dtos

import (
	"encoding/json"

	"github.com/robertokbr/blinkchat/src/domain/logger"
)

type rawMessageData struct {
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
}

type rawMessage struct {
	Data   rawMessageData `json:"data"`
	Action string         `json:"action"`
}

type CreateMessage struct {
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
	Event       string `json:"event"`
}

func NewCreateMessage(message []byte) (*CreateMessage, error) {
	var rawMessage rawMessage

	err := json.Unmarshal(message, &rawMessage)

	if err != nil {
		logger.Errorf("failing to create rawMessage with data %v", message)
		return nil, err
	}

	createMessage := &CreateMessage{
		Content:     rawMessage.Data.Content,
		MessageType: rawMessage.Data.MessageType,
		Event:       rawMessage.Action,
	}

	return createMessage, nil
}
