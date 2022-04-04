package utils

func MustNot[T any](v T, err error) T {
	return v
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
