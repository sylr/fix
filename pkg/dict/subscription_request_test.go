package dict

import (
	"testing"
)

func TestSubscriptionRequestTypesMapsSync(t *testing.T) {
	for k, v := range SubscriptionRequestTypes {
		if rv, ok := SubscriptionRequestTypesReversed[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.SubscriptionRequestTypesReversed", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.SubscriptionRequestTypes[%s] = `%s` AND dict.SubscriptionRequestTypesReversed[`%s`] = `%s`", k, v, v, rv)
		}
	}

	for k, v := range SubscriptionRequestTypesReversed {
		if rv, ok := SubscriptionRequestTypes[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.SubscriptionRequestTypes", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.SubscriptionRequestTypesReversed[`%s`] = `%s` AND dict.SubscriptionRequestTypes[`%s`] = `%s`", k, v, v, rv)
		}
	}
}
