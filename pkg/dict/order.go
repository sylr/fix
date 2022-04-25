package dict

import (
	"fmt"
	"strings"

	"github.com/quickfixgo/enum"
)

var OrderSides = map[string]enum.Side{
	"BUY":                enum.Side_BUY,
	"SELL":               enum.Side_SELL,
	"BUY_MINUS":          enum.Side_BUY_MINUS,
	"SELL_PLUS":          enum.Side_SELL_PLUS,
	"SELL_SHORT":         enum.Side_SELL_SHORT,
	"SELL_SHORT_EXEMPT":  enum.Side_SELL_SHORT_EXEMPT,
	"UNDISCLOSED":        enum.Side_UNDISCLOSED,
	"CROSS":              enum.Side_CROSS,
	"CROSS_SHORT":        enum.Side_CROSS_SHORT,
	"CROSS_SHORT_EXEMPT": enum.Side_CROSS_SHORT_EXEMPT,
	"AS_DEFINED":         enum.Side_AS_DEFINED,
	"OPPOSITE":           enum.Side_OPPOSITE,
	"SUBSCRIBE":          enum.Side_SUBSCRIBE,
	"REDEEM":             enum.Side_REDEEM,
	"LEND":               enum.Side_LEND,
	"BORROW":             enum.Side_BORROW,
}

var OrderSidesReversed = map[enum.Side]string{
	enum.Side_BUY:                "BUY",
	enum.Side_SELL:               "SELL",
	enum.Side_BUY_MINUS:          "BUY_MINUS",
	enum.Side_SELL_PLUS:          "SELL_PLUS",
	enum.Side_SELL_SHORT:         "SELL_SHORT",
	enum.Side_SELL_SHORT_EXEMPT:  "SELL_SHORT_EXEMPT",
	enum.Side_UNDISCLOSED:        "UNDISCLOSED",
	enum.Side_CROSS:              "CROSS",
	enum.Side_CROSS_SHORT:        "CROSS_SHORT",
	enum.Side_CROSS_SHORT_EXEMPT: "CROSS_SHORT_EXEMPT",
	enum.Side_AS_DEFINED:         "AS_DEFINED",
	enum.Side_OPPOSITE:           "OPPOSITE",
	enum.Side_SUBSCRIBE:          "SUBSCRIBE",
	enum.Side_REDEEM:             "REDEEM",
	enum.Side_LEND:               "LEND",
	enum.Side_BORROW:             "BORROW",
}

func OrderSideStringToEnum(side string) (enum.Side, error) {
	side = strings.ToUpper(side)
	if e, ok := OrderSides[side]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unkown order side")
}

var OrderTypes = map[string]enum.OrdType{
	"MARKET":                         enum.OrdType_MARKET,
	"LIMIT":                          enum.OrdType_LIMIT,
	"STOP":                           enum.OrdType_STOP_STOP_LOSS,
	"STOP_LIMIT":                     enum.OrdType_STOP_LIMIT,
	"MARKET_ON_CLOSE":                enum.OrdType_MARKET_ON_CLOSE,
	"WITH_OR_WITHOUT":                enum.OrdType_WITH_OR_WITHOUT,
	"LIMIT_OR_BETTER":                enum.OrdType_LIMIT_OR_BETTER,
	"LIMIT_WITH_OR_WITHOUT":          enum.OrdType_LIMIT_WITH_OR_WITHOUT,
	"ON_BASIS":                       enum.OrdType_ON_BASIS,
	"ON_CLOSE":                       enum.OrdType_ON_CLOSE,
	"LIMIT_ON_CLOSE":                 enum.OrdType_LIMIT_ON_CLOSE,
	"FOREX_MARKET":                   enum.OrdType_FOREX_MARKET,
	"PREVIOUSLY_QUOTED":              enum.OrdType_PREVIOUSLY_QUOTED,
	"PREVIOUSLY_INDICATED":           enum.OrdType_PREVIOUSLY_INDICATED,
	"FOREX_LIMIT":                    enum.OrdType_FOREX_LIMIT,
	"FOREX_SWAP":                     enum.OrdType_FOREX_SWAP,
	"FOREX_PREVIOUSLY_QUOTED":        enum.OrdType_FOREX_PREVIOUSLY_QUOTED,
	"FUNARI":                         enum.OrdType_FUNARI,
	"MARKET_IF_TOUCHED":              enum.OrdType_MARKET_IF_TOUCHED,
	"MARKET_WITH_LEFT_OVER_AS_LIMIT": enum.OrdType_MARKET_WITH_LEFT_OVER_AS_LIMIT,
	"PREVIOUS_FUND_VALUATION_POINT":  enum.OrdType_PREVIOUS_FUND_VALUATION_POINT,
	"NEXT_FUND_VALUATION_POINT":      enum.OrdType_NEXT_FUND_VALUATION_POINT,
	"PEGGED":                         enum.OrdType_PEGGED,
	"COUNTER_ORDER_SELECTION":        enum.OrdType_COUNTER_ORDER_SELECTION,
}

func OrderTypeStringToEnum(order string) (enum.OrdType, error) {
	order = strings.ToUpper(order)
	if e, ok := OrderTypes[order]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unkown order type")
}

var OrderTimeInForces = map[string]enum.TimeInForce{
	"DAY":                   enum.TimeInForce_DAY,
	"GOOD_TILL_CANCEL":      enum.TimeInForce_GOOD_TILL_CANCEL,
	"AT_THE_OPENING":        enum.TimeInForce_AT_THE_OPENING,
	"IMMEDIATE_OR_CANCEL":   enum.TimeInForce_IMMEDIATE_OR_CANCEL,
	"FILL_OR_KILL":          enum.TimeInForce_FILL_OR_KILL,
	"GOOD_TILL_CROSSING":    enum.TimeInForce_GOOD_TILL_CROSSING,
	"GOOD_TILL_DATE":        enum.TimeInForce_GOOD_TILL_DATE,
	"AT_THE_CLOSE":          enum.TimeInForce_AT_THE_CLOSE,
	"GOOD_THROUGH_CROSSING": enum.TimeInForce_GOOD_THROUGH_CROSSING,
	"AT_CROSSING":           enum.TimeInForce_AT_CROSSING,
}

func OrderTimeInForceStringToEnum(t string) (enum.TimeInForce, error) {
	t = strings.ToUpper(t)
	if e, ok := OrderTimeInForces[t]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unkown time in force")
}

var OrderOriginations = map[string]enum.OrderOrigination{
	"CUSTOMER":              enum.OrderOrigination_ORDER_RECEIVED_FROM_A_CUSTOMER,
	"WITHIN_THE_FIRM":       enum.OrderOrigination_ORDER_RECEIVED_FROM_WITHIN_THE_FIRM,
	"ANOTHER_BROKER_DEALER": enum.OrderOrigination_ORDER_RECEIVED_FROM_ANOTHER_BROKER_DEALER,
	"A_CUSTOMER_OR_ORIGINATED_FROM_WITHIN_THE_FIRM": enum.OrderOrigination_ORDER_RECEIVED_FROM_A_CUSTOMER_OR_ORIGINATED_FROM_WITHIN_THE_FIRM,
	"A_DIRECT_ACCESS_OR_SPONSORED_ACCESS_CUSTOMER":  enum.OrderOrigination_ORDER_RECEIVED_FROM_A_DIRECT_ACCESS_OR_SPONSORED_ACCESS_CUSTOMER,
	"A_FOREIGN_DEALER_EQUIVALENT":                   enum.OrderOrigination_ORDER_RECEIVED_FROM_A_FOREIGN_DEALER_EQUIVALENT,
	"AN_EXECUTION_ONLY_SERVICE":                     enum.OrderOrigination_ORDER_RECEIVED_FROM_AN_EXECUTION_ONLY_SERVICE,
}

var OrderAttributeTypes = map[string]enum.OrderAttributeType{
	"AGGREGATED_ORDER":                                           enum.OrderAttributeType_AGGREGATED_ORDER,
	"ORDER_PENDING_ALLOCATION":                                   enum.OrderAttributeType_ORDER_PENDING_ALLOCATION,
	"SUBJECT_TO_EU_SHARE_TRADING_OBLIGATION":                     enum.OrderAttributeType_SUBJECT_TO_EU_SHARE_TRADING_OBLIGATION,
	"SUBJECT_TO_UK_SHARE_TRADING_OBLIGATION":                     enum.OrderAttributeType_SUBJECT_TO_UK_SHARE_TRADING_OBLIGATION,
	"REPRESENTATIVE_ORDER":                                       enum.OrderAttributeType_REPRESENTATIVE_ORDER,
	"LINKAGE_TYPE":                                               enum.OrderAttributeType_LINKAGE_TYPE,
	"EXEMPT_FROM_SHARE_TRADING_OBLIGATION":                       enum.OrderAttributeType_EXEMPT_FROM_SHARE_TRADING_OBLIGATION,
	"LIQUIDITY_PROVISION_ACTIVITY_ORDER":                         enum.OrderAttributeType_LIQUIDITY_PROVISION_ACTIVITY_ORDER,
	"RISK_REDUCTION_ORDER":                                       enum.OrderAttributeType_RISK_REDUCTION_ORDER,
	"ALGORITHMIC_ORDER":                                          enum.OrderAttributeType_ALGORITHMIC_ORDER,
	"SYSTEMIC_INTERNALISER_ORDER":                                enum.OrderAttributeType_SYSTEMIC_INTERNALISER_ORDER,
	"ALL_EXECUTIONS_FOR_THE_ORDER_ARE_TO_BE_SUBMITTED_TO_AN_APA": enum.OrderAttributeType_ALL_EXECUTIONS_FOR_THE_ORDER_ARE_TO_BE_SUBMITTED_TO_AN_APA,
	"ORDER_EXECUTION_INSTRUCTED_BY_CLIENT":                       enum.OrderAttributeType_ORDER_EXECUTION_INSTRUCTED_BY_CLIENT,
	"LARGE_IN_SCALE_ORDER":                                       enum.OrderAttributeType_LARGE_IN_SCALE_ORDER,
	"HIDDEN_ORDER":                                               enum.OrderAttributeType_HIDDEN_ORDER,
}
