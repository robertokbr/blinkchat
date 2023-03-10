package enums

type WebsocketEvent string

const (
	DISCONNECTED WebsocketEvent = "disconnected"
	CONNECTED    WebsocketEvent = "connected"
	MATCHING     WebsocketEvent = "matching"
	UNMATCHING   WebsocketEvent = "unmatching"
	BROADCASTING WebsocketEvent = "broadcasting"
)
