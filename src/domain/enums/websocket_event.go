package enums

type WebsocketEvent string

const (
	DISCONNECTED WebsocketEvent = "disconnected"
	CONNECTED    WebsocketEvent = "connected"
	MATCHING     WebsocketEvent = "matching"
	UNMATCHING   WebsocketEvent = "unmatching"
	MATCHED      WebsocketEvent = "matched"
	UNMATCHED    WebsocketEvent = "unmatched"
	BROADCASTING WebsocketEvent = "broadcasting"
)
