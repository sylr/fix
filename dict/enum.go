package dict

import (
	"fmt"

	"github.com/sylr/fix2/fix/dictionary"
)

// EnumName names value of tag after d, as NAME(value), or returns value if
// d does not name it.
func EnumName(d *dictionary.Dict, tag int, value string) string {
	if d != nil {
		if def := d.Fields[tag]; def != nil {
			if name, ok := def.Enums[value]; ok {
				return fmt.Sprintf("%s(%s)", name, value)
			}
		}
	}

	return value
}
