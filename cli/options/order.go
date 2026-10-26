package options

import (
	"fmt"

	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
)

// Side parses v, the value of --side.
func Side(v string) (fix50sp2.Side, error) {
	side, ok := Lookup(dict.OrderSides, v)
	if !ok {
		return "", fmt.Errorf("%w: `%s`", errors.ErrOptionOrderSideUnknown, v)
	}

	return side, nil
}

// OrdType parses v, the value of --type.
func OrdType(v string) (fix50sp2.OrdType, error) {
	ordType, ok := Lookup(dict.OrderTypes, v)
	if !ok {
		return "", fmt.Errorf("%w: `%s`", errors.ErrOptionOrderTypeUnknown, v)
	}

	return ordType, nil
}

// TimeInForce parses v, the value of --expiry.
func TimeInForce(v string) (fix50sp2.TimeInForce, error) {
	tif, ok := Lookup(dict.OrderTimeInForces, v)
	if !ok {
		return "", fmt.Errorf("%w: `%s`", errors.ErrOptionOrderExpiryUnknown, v)
	}

	return tif, nil
}

// CheckPrice checks price, the value of --price, fits an order of type
// ordType: none for a market order, one otherwise. Prices may be zero or
// negative, as FIX prices of spreads and some futures are.
func CheckPrice(ordType fix50sp2.OrdType, price Decimal) error {
	switch {
	case ordType == fix50sp2.OrdTypeMarket && price.IsSet():
		return errors.ErrOptionsInvalidMarketPrice
	case ordType == fix50sp2.OrdTypeMarket:
		return nil
	case !price.IsSet():
		return errors.ErrOptionsNoPriceGiven
	}

	return nil
}
