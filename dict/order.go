package dict

import (
	"fmt"
	"strings"

	"sylr.dev/fix/v2/fix50sp2"
)

var OrderSides = map[string]fix50sp2.Side{
	"BUY":                fix50sp2.SideBuy,
	"SELL":               fix50sp2.SideSell,
	"BUY_MINUS":          fix50sp2.SideBuyMinus,
	"SELL_PLUS":          fix50sp2.SideSellPlus,
	"SELL_SHORT":         fix50sp2.SideSellShort,
	"SELL_SHORT_EXEMPT":  fix50sp2.SideSellShortExempt,
	"UNDISCLOSED":        fix50sp2.SideUndisclosed,
	"CROSS":              fix50sp2.SideCross,
	"CROSS_SHORT":        fix50sp2.SideCrossShort,
	"CROSS_SHORT_EXEMPT": fix50sp2.SideCrossShortExempt,
	"AS_DEFINED":         fix50sp2.SideAsDefined,
	"OPPOSITE":           fix50sp2.SideOpposite,
	"SUBSCRIBE":          fix50sp2.SideSubscribe,
	"REDEEM":             fix50sp2.SideRedeem,
	"LEND":               fix50sp2.SideLend,
	"BORROW":             fix50sp2.SideBorrow,
}

var OrderSidesReversed = map[fix50sp2.Side]string{
	fix50sp2.SideBuy:              "BUY",
	fix50sp2.SideSell:             "SELL",
	fix50sp2.SideBuyMinus:         "BUY_MINUS",
	fix50sp2.SideSellPlus:         "SELL_PLUS",
	fix50sp2.SideSellShort:        "SELL_SHORT",
	fix50sp2.SideSellShortExempt:  "SELL_SHORT_EXEMPT",
	fix50sp2.SideUndisclosed:      "UNDISCLOSED",
	fix50sp2.SideCross:            "CROSS",
	fix50sp2.SideCrossShort:       "CROSS_SHORT",
	fix50sp2.SideCrossShortExempt: "CROSS_SHORT_EXEMPT",
	fix50sp2.SideAsDefined:        "AS_DEFINED",
	fix50sp2.SideOpposite:         "OPPOSITE",
	fix50sp2.SideSubscribe:        "SUBSCRIBE",
	fix50sp2.SideRedeem:           "REDEEM",
	fix50sp2.SideLend:             "LEND",
	fix50sp2.SideBorrow:           "BORROW",
}

func OrderSideStringToEnum(side string) (fix50sp2.Side, error) {
	side = strings.ToUpper(side)
	if e, ok := OrderSides[side]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unknown order side")
}

var OrderTypes = map[string]fix50sp2.OrdType{
	"MARKET":                         fix50sp2.OrdTypeMarket,
	"LIMIT":                          fix50sp2.OrdTypeLimit,
	"STOP":                           fix50sp2.OrdTypeStopStopLoss,
	"STOP_LIMIT":                     fix50sp2.OrdTypeStopLimit,
	"MARKET_ON_CLOSE":                fix50sp2.OrdTypeMarketOnClose,
	"WITH_OR_WITHOUT":                fix50sp2.OrdTypeWithOrWithout,
	"LIMIT_OR_BETTER":                fix50sp2.OrdTypeLimitOrBetter,
	"LIMIT_WITH_OR_WITHOUT":          fix50sp2.OrdTypeLimitWithOrWithout,
	"ON_BASIS":                       fix50sp2.OrdTypeOnBasis,
	"ON_CLOSE":                       fix50sp2.OrdTypeOnClose,
	"LIMIT_ON_CLOSE":                 fix50sp2.OrdTypeLimitOnClose,
	"FOREX_MARKET":                   fix50sp2.OrdTypeForexMarket,
	"PREVIOUSLY_QUOTED":              fix50sp2.OrdTypePreviouslyQuoted,
	"PREVIOUSLY_INDICATED":           fix50sp2.OrdTypePreviouslyIndicated,
	"FOREX_LIMIT":                    fix50sp2.OrdTypeForexLimit,
	"FOREX_SWAP":                     fix50sp2.OrdTypeForexSwap,
	"FOREX_PREVIOUSLY_QUOTED":        fix50sp2.OrdTypeForexPreviouslyQuoted,
	"FUNARI":                         fix50sp2.OrdTypeFunari,
	"MARKET_IF_TOUCHED":              fix50sp2.OrdTypeMarketIfTouched,
	"MARKET_WITH_LEFT_OVER_AS_LIMIT": fix50sp2.OrdTypeMarketWithLeftOverAsLimit,
	"PREVIOUS_FUND_VALUATION_POINT":  fix50sp2.OrdTypePreviousFundValuationPoint,
	"NEXT_FUND_VALUATION_POINT":      fix50sp2.OrdTypeNextFundValuationPoint,
	"PEGGED":                         fix50sp2.OrdTypePegged,
	"COUNTER_ORDER_SELECTION":        fix50sp2.OrdTypeCounterOrderSelection,
}

var OrderTypesReversed = map[fix50sp2.OrdType]string{
	fix50sp2.OrdTypeMarket:                     "MARKET",
	fix50sp2.OrdTypeLimit:                      "LIMIT",
	fix50sp2.OrdTypeStopStopLoss:               "STOP",
	fix50sp2.OrdTypeStopLimit:                  "STOP_LIMIT",
	fix50sp2.OrdTypeMarketOnClose:              "MARKET_ON_CLOSE",
	fix50sp2.OrdTypeWithOrWithout:              "WITH_OR_WITHOUT",
	fix50sp2.OrdTypeLimitOrBetter:              "LIMIT_OR_BETTER",
	fix50sp2.OrdTypeLimitWithOrWithout:         "LIMIT_WITH_OR_WITHOUT",
	fix50sp2.OrdTypeOnBasis:                    "ON_BASIS",
	fix50sp2.OrdTypeOnClose:                    "ON_CLOSE",
	fix50sp2.OrdTypeLimitOnClose:               "LIMIT_ON_CLOSE",
	fix50sp2.OrdTypeForexMarket:                "FOREX_MARKET",
	fix50sp2.OrdTypePreviouslyQuoted:           "PREVIOUSLY_QUOTED",
	fix50sp2.OrdTypePreviouslyIndicated:        "PREVIOUSLY_INDICATED",
	fix50sp2.OrdTypeForexLimit:                 "FOREX_LIMIT",
	fix50sp2.OrdTypeForexSwap:                  "FOREX_SWAP",
	fix50sp2.OrdTypeForexPreviouslyQuoted:      "FOREX_PREVIOUSLY_QUOTED",
	fix50sp2.OrdTypeFunari:                     "FUNARI",
	fix50sp2.OrdTypeMarketIfTouched:            "MARKET_IF_TOUCHED",
	fix50sp2.OrdTypeMarketWithLeftOverAsLimit:  "MARKET_WITH_LEFT_OVER_AS_LIMIT",
	fix50sp2.OrdTypePreviousFundValuationPoint: "PREVIOUS_FUND_VALUATION_POINT",
	fix50sp2.OrdTypeNextFundValuationPoint:     "NEXT_FUND_VALUATION_POINT",
	fix50sp2.OrdTypePegged:                     "PEGGED",
	fix50sp2.OrdTypeCounterOrderSelection:      "COUNTER_ORDER_SELECTION",
}

func OrderTypeStringToEnum(order string) (fix50sp2.OrdType, error) {
	order = strings.ToUpper(order)
	if e, ok := OrderTypes[order]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unknown order type")
}

var OrderTimeInForces = map[string]fix50sp2.TimeInForce{
	"DAY":                   fix50sp2.TimeInForceDay,
	"GOOD_TILL_CANCEL":      fix50sp2.TimeInForceGoodTillCancel,
	"AT_THE_OPENING":        fix50sp2.TimeInForceAtTheOpening,
	"IMMEDIATE_OR_CANCEL":   fix50sp2.TimeInForceImmediateOrCancel,
	"FILL_OR_KILL":          fix50sp2.TimeInForceFillOrKill,
	"GOOD_TILL_CROSSING":    fix50sp2.TimeInForceGoodTillCrossing,
	"GOOD_TILL_DATE":        fix50sp2.TimeInForceGoodTillDate,
	"AT_THE_CLOSE":          fix50sp2.TimeInForceAtTheClose,
	"GOOD_THROUGH_CROSSING": fix50sp2.TimeInForceGoodThroughCrossing,
	"AT_CROSSING":           fix50sp2.TimeInForceAtCrossing,
}

func OrderTimeInForceStringToEnum(t string) (fix50sp2.TimeInForce, error) {
	t = strings.ToUpper(t)
	if e, ok := OrderTimeInForces[t]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unknown time in force")
}

var OrderOriginations = map[string]fix50sp2.OrderOrigination{
	"CUSTOMER":              fix50sp2.OrderOriginationOrderReceivedFromACustomer,
	"WITHIN_THE_FIRM":       fix50sp2.OrderOriginationOrderReceivedFromWithinTheFirm,
	"ANOTHER_BROKER_DEALER": fix50sp2.OrderOriginationOrderReceivedFromAnotherBrokerDealer,
	"A_CUSTOMER_OR_ORIGINATED_FROM_WITHIN_THE_FIRM": fix50sp2.OrderOriginationOrderReceivedFromACustomerOrOriginatedFromWithinTheFirm,
	"A_DIRECT_ACCESS_OR_SPONSORED_ACCESS_CUSTOMER":  fix50sp2.OrderOriginationOrderReceivedFromADirectAccessOrSponsoredAccessCustomer,
	"A_FOREIGN_DEALER_EQUIVALENT":                   fix50sp2.OrderOriginationOrderReceivedFromAForeignDealerEquivalent,
	"AN_EXECUTION_ONLY_SERVICE":                     fix50sp2.OrderOriginationOrderReceivedFromAnExecutionOnlyService,
}

var OrderAttributeTypes = map[string]fix50sp2.OrderAttributeType{
	"AGGREGATED_ORDER":                                           fix50sp2.OrderAttributeTypeAggregatedOrder,
	"ORDER_PENDING_ALLOCATION":                                   fix50sp2.OrderAttributeTypeOrderPendingAllocation,
	"SUBJECT_TO_EU_SHARE_TRADING_OBLIGATION":                     fix50sp2.OrderAttributeTypeSubjectToEuShareTradingObligation,
	"SUBJECT_TO_UK_SHARE_TRADING_OBLIGATION":                     fix50sp2.OrderAttributeTypeSubjectToUkShareTradingObligation,
	"REPRESENTATIVE_ORDER":                                       fix50sp2.OrderAttributeTypeRepresentativeOrder,
	"LINKAGE_TYPE":                                               fix50sp2.OrderAttributeTypeLinkageType,
	"EXEMPT_FROM_SHARE_TRADING_OBLIGATION":                       fix50sp2.OrderAttributeTypeExemptFromShareTradingObligation,
	"LIQUIDITY_PROVISION_ACTIVITY_ORDER":                         fix50sp2.OrderAttributeTypeLiquidityProvisionActivityOrder,
	"RISK_REDUCTION_ORDER":                                       fix50sp2.OrderAttributeTypeRiskReductionOrder,
	"ALGORITHMIC_ORDER":                                          fix50sp2.OrderAttributeTypeAlgorithmicOrder,
	"SYSTEMIC_INTERNALISER_ORDER":                                fix50sp2.OrderAttributeTypeSystemicInternaliserOrder,
	"ALL_EXECUTIONS_FOR_THE_ORDER_ARE_TO_BE_SUBMITTED_TO_AN_APA": fix50sp2.OrderAttributeTypeAllExecutionsForTheOrderAreToBeSubmittedToAnApa,
	"ORDER_EXECUTION_INSTRUCTED_BY_CLIENT":                       fix50sp2.OrderAttributeTypeOrderExecutionInstructedByClient,
	"LARGE_IN_SCALE_ORDER":                                       fix50sp2.OrderAttributeTypeLargeInScaleOrder,
	"HIDDEN_ORDER":                                               fix50sp2.OrderAttributeTypeHiddenOrder,
}
