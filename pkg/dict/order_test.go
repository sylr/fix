package dict

import (
	"testing"
)

func TestOrderSidesMapsSync(t *testing.T) {
	for k, v := range OrderSides {
		if rv, ok := OrderSidesReversed[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.OrderSidesReversed", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.OrderSides[%s] = `%s` AND dict.OrderSidesReversed[`%s`] = `%s`", k, v, v, rv)
		}
	}

	for k, v := range OrderSidesReversed {
		if rv, ok := OrderSides[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.OrderSides", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.OrderSidesReversed[`%s`] = `%s` AND dict.OrderSides[`%s`] = `%s`", k, v, v, rv)
		}
	}
}

func TestOrderTypesMapsSync(t *testing.T) {
	for k, v := range OrderTypes {
		if rv, ok := OrderTypesReversed[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.OrderTypesReversed", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.OrderTypes[%s] = `%s` AND dict.OrderTypesReversed[`%s`] = `%s`", k, v, v, rv)
		}
	}

	for k, v := range OrderTypesReversed {
		if rv, ok := OrderTypes[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.OrderTypes", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.OrderTypesReversed[`%s`] = `%s` AND dict.OrderTypes[`%s`] = `%s`", k, v, v, rv)
		}
	}
}
