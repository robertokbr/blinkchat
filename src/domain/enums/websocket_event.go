package enums

type WebsocketEvent string

const (
	DISCONNECTED WebsocketEvent = "disconnected"
	CONNECTED    WebsocketEvent = "connected"
	SEARCHING    WebsocketEvent = "searching"
	MATCHING     WebsocketEvent = "matching"
	UNMATCHING   WebsocketEvent = "unmatching"
	BROADCASTING WebsocketEvent = "broadcasting"
)
