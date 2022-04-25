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
