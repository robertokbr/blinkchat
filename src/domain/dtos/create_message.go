package dtos

type RawMessage struct {
	Data   string `json:"data"`
	Action string `json:"action"`
}

type CreateMessage struct {
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
	Event       string `json:"event"`
}

func NewCreateMessage(content, messageType, event string) *CreateMessage {
	return &CreateMessage{
		Content:     content,
		MessageType: messageType,
		Event:       event,
	}
}
