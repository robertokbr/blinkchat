package usecases_tests_spies

import "github.com/robertokbr/blinkchat/src/domain/models"

type WebsocketConnection struct {
	Messages             chan string
	MessagesSent         []models.Message
	WriteJSONTimesCalled int
}

func NewWebsocketConnection() *WebsocketConnection {
	return &WebsocketConnection{
		Messages:             make(chan string),
		MessagesSent:         make([]models.Message, 0),
		WriteJSONTimesCalled: 0,
	}
}

func (wsc *WebsocketConnection) Close() error {
	return nil
}

func (wsc *WebsocketConnection) ReadMessage() (int, []byte, error) {
	for {
		select {
		case message := <-wsc.Messages:
			return 0, []byte(message), nil
		}
	}
}

func (wsc *WebsocketConnection) WriteJSON(data interface{}) error {
	wsc.WriteJSONTimesCalled++
	wsc.MessagesSent = append(wsc.MessagesSent, data.(models.Message))
	return nil
}
