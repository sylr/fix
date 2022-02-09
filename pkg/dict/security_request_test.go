package dict

import (
	"testing"
)

func TestSecurityRequestTypesMapsSync(t *testing.T) {
	for k, v := range SecurityRequestTypes {
		if rv, ok := SecurityRequestTypesReversed[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.SecurityRequestTypesReversed", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.SecurityRequestTypes[%s] = `%s` AND dict.SecurityRequestTypesReversed[`%s`] = `%s`", k, v, v, rv)
		}
	}

	for k, v := range SecurityRequestTypesReversed {
		if rv, ok := SecurityRequestTypes[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.SecurityRequestTypes", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.SecurityRequestTypesReversed[`%s`] = `%s` AND dict.SecurityRequestTypes[`%s`] = `%s`", k, v, v, rv)
		}
	}
}

func TestSecurityListRequestTypesMapsSync(t *testing.T) {
	for k, v := range SecurityListRequestTypes {
		if rv, ok := SecurityListRequestTypesReversed[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.SecurityListRequestTypesReversed", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.SecurityListRequestTypes[%s] = `%s` AND dict.SecurityListRequestTypesReversed[`%s`] = `%s`", k, v, v, rv)
		}
	}

	for k, v := range SecurityListRequestTypesReversed {
		if rv, ok := SecurityListRequestTypes[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.SecurityListRequestTypes", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.SecurityListRequestTypesReversed[`%s`] = `%s` AND dict.SecurityListRequestTypes[`%s`] = `%s`", k, v, v, rv)
		}
	}
}
