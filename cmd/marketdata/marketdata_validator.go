//go:build validator
// +build validator

package marketdata

import (
	markedatavalidator "sylr.dev/fix/v2/cmd/marketdata/validator"
	"sylr.dev/fix/v2/initiator"
)

func init() {
	initiator.AddPersistentFlags(markedatavalidator.MarketDataValidatorCmd)

	if err := initiator.AddPersistentFlagCompletions(markedatavalidator.MarketDataValidatorCmd); err != nil {
		panic(err)
	}

	MarketDataCmd.AddCommand(markedatavalidator.MarketDataValidatorCmd)
}
