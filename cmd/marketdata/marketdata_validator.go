//go:build validator

package marketdata

import (
	marketdatavalidator "sylr.dev/fix/v2/cmd/marketdata/validator"
)

func init() {
	MarketDataCmd.AddCommand(marketdatavalidator.MarketDataValidatorCmd)
}
