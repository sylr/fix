package dict

import (
	"fmt"
)

func Search[T ~string](dict map[string]T, raw string) (T, error) {
	for k, v := range dict {
		if k == raw {
			return v, nil
		}
	}

	return T(raw), fmt.Errorf("unknown value")
}

func SearchValue[T ~string](dict map[string]T, raw T) (string, error) {
	for k, v := range dict {
		if v == raw {
			return k, nil
		}
	}

	return "", fmt.Errorf("unknown value")
}
