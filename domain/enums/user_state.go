package enums

type UserState string

const (
	DISCONNECTED UserState = "DISCONNECTED"
	CONNECTED    UserState = "CONNECTED"
	SEARCHING    UserState = "SEARCHING"
	IN_ROOM      UserState = "IN_ROOM"
)
