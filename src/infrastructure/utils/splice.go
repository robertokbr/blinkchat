package utils

func Splice[A []V, V any](s *A, index int, params ...V) {
	arr := *s

	if len(params) > 0 {
		arr = append(arr[:index], append(params, arr[index+1:]...)...)

		*s = arr
		return
	}

	*s = append(arr[:index], arr[index+1:]...)
}
