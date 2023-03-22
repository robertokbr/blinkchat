package dtos

import (
	"encoding/json"
	"log"
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

func NewCreateMessage(message string) (*CreateMessage, error) {
	var rawMessage rawMessage

	err := json.Unmarshal([]byte(message), &rawMessage)

	if err != nil {
		log.Printf("[error]: failing to create rawMessage with data %v", message)
		return nil, err
	}

	createMessage := &CreateMessage{
		Content:     rawMessage.Data.Content,
		MessageType: rawMessage.Data.MessageType,
		Event:       rawMessage.Action,
	}

	return createMessage, nil
}
