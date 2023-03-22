package utils

func If[T any](condition bool, truthy T, falsy T) T {
	if condition {
		return truthy
	} else {
		return falsy
	}
}
