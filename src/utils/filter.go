package utils

func Filter[A []V, V any](a *A, f func(v V) bool) {
	filtered := make([]V, 0)

	for _, data := range *a {
		if f(data) {
			filtered = append(filtered, data)
		}
	}

	*a = filtered
}
