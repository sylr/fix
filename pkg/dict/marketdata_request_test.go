package dict

import (
	"testing"
)

func TestMDEntryTypesMapsSync(t *testing.T) {
	for k, v := range MDEntryTypes {
		if rv, ok := MDEntryTypesReversed[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.MDEntryTypesReversed", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.MDEntryTypes[%s] = `%s` AND dict.MDEntryTypesReversed[`%s`] = `%s`", k, v, v, rv)
		}
	}

	for k, v := range MDEntryTypesReversed {
		if rv, ok := MDEntryTypes[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.MDEntryTypes", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.MDEntryTypesReversed[`%s`] = `%s` AND dict.MDEntryTypes[`%s`] = `%s`", k, v, v, rv)
		}
	}
}
