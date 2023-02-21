package enums

type ActionType string

const (
	DISCONNECTED ActionType = "DISCONNECTED"
	CONNECTED    ActionType = "CONNECTED"
	SEARCHING    ActionType = "SEARCHING"
	MATCHED      ActionType = "MATCHED"
	MESSAGE      ActionType = "MESSAGE"
)
