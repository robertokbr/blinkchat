package interfaces

type WebsocketConnection interface {
	Close() error
	ReadMessage() (messageType int, p []byte, err error)
	WriteJSON(interface{}) error
}
