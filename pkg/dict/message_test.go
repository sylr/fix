package dict

import (
	"testing"
)

func TestMessageTypesMapsSync(t *testing.T) {
	for k, v := range MessageTypes {
		if rv, ok := MessageTypesReversed[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.MessageTypesReversed", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.MessageTypes[%s] = `%s` AND dict.MessageTypesReversed[`%s`] = `%s`", k, v, v, rv)
		}
	}

	for k, v := range MessageTypesReversed {
		if rv, ok := MessageTypes[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.MessageTypes", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.MessageTypesReversed[`%s`] = `%s` AND dict.MessageTypes[`%s`] = `%s`", k, v, v, rv)
		}
	}
}
