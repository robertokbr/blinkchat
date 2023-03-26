package enums

type UserState string

const (
	IN_A_MATCH        UserState = "IN_A_MATCH"
	NOT_IN_A_MATCH    UserState = "NOT_IN_A_MATCH"
	LOOKING_FOR_MATCH UserState = "LOOKING_FOR_MATCH"
)
