package utils

import (
	"sort"
	"strings"
)

func PrettyOptionValues[T any](input map[string]T) []string {
	output := make([]string, 0, len(input))
	for k := range input {
		output = append(output, strings.ToLower(k))
	}
	sort.Strings(output)
	return output
}
