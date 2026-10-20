package utils

func Search[T comparable](slice []T, search T) int {
	for k, v := range slice {
		if search == v {
			return k
		}
	}

	return -1
}
