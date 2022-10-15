package enums

type UserState string

const (
	DISCONNECTED UserState = "DISCONNECTED"
	SEARCHING    UserState = "searching"
	CONNECTED    UserState = "CONNECTED"
)
