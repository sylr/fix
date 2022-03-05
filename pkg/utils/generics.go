package utils

func MustNot[T any, T2 any](v T, err T2) T {
	return v
}
