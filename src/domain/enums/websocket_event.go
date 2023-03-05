package enums

type WebsocketEvent string

const (
	DISCONNECTED WebsocketEvent = "DISCONNECTED"
	CONNECTED    WebsocketEvent = "CONNECTED"
	SEARCHING    WebsocketEvent = "SEARCHING"
	MATCHING     WebsocketEvent = "MATCHING"
	BROADCASTING WebsocketEvent = "BROADCASTING"
)
