package utils

func Find[A []V, V any](a A, f func(v V) bool) (V, int) {
	var v V

	for i, value := range a {
		if f(value) {
			return value, i
		}
	}

	return v, -1
}
