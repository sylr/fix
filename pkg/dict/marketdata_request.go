package dict

import (
	"github.com/quickfixgo/enum"
)

var MDEntryTypes = map[string]enum.MDEntryType{
	"BID":                           enum.MDEntryType_BID,
	"OFFER":                         enum.MDEntryType_OFFER,
	"TRADE":                         enum.MDEntryType_TRADE,
	"INDEX_VALUE":                   enum.MDEntryType_INDEX_VALUE,
	"OPENING_PRICE":                 enum.MDEntryType_OPENING_PRICE,
	"CLOSING_PRICE":                 enum.MDEntryType_CLOSING_PRICE,
	"SETTLEMENT_PRICE":              enum.MDEntryType_SETTLEMENT_PRICE,
	"TRADING_SESSION_HIGH_PRICE":    enum.MDEntryType_TRADING_SESSION_HIGH_PRICE,
	"TRADING_SESSION_LOW_PRICE":     enum.MDEntryType_TRADING_SESSION_LOW_PRICE,
	"VOLUME_WEIGHTED_AVERAGE_PRICE": enum.MDEntryType_VOLUME_WEIGHTED_AVERAGE_PRICE,
	"IMBALANCE":                     enum.MDEntryType_IMBALANCE,
	"TRADE_VOLUME":                  enum.MDEntryType_TRADE_VOLUME,
	"OPEN_INTEREST":                 enum.MDEntryType_OPEN_INTEREST,
	"COMPOSITE_UNDERLYING_PRICE":    enum.MDEntryType_COMPOSITE_UNDERLYING_PRICE,
	"SIMULATED_SELL_PRICE":          enum.MDEntryType_SIMULATED_SELL_PRICE,
	"SIMULATED_BUY_PRICE":           enum.MDEntryType_SIMULATED_BUY_PRICE,
	"MARGIN_RATE":                   enum.MDEntryType_MARGIN_RATE,
	"MID_PRICE":                     enum.MDEntryType_MID_PRICE,
	"EMPTY_BOOK":                    enum.MDEntryType_EMPTY_BOOK,
	"SETTLE_HIGH_PRICE":             enum.MDEntryType_SETTLE_HIGH_PRICE,
	"SETTLE_LOW_PRICE":              enum.MDEntryType_SETTLE_LOW_PRICE,
	"PRIOR_SETTLE_PRICE":            enum.MDEntryType_PRIOR_SETTLE_PRICE,
	"SESSION_HIGH_BID":              enum.MDEntryType_SESSION_HIGH_BID,
	"SESSION_LOW_OFFER":             enum.MDEntryType_SESSION_LOW_OFFER,
	"EARLY_PRICES":                  enum.MDEntryType_EARLY_PRICES,
	"AUCTION_CLEARING_PRICE":        enum.MDEntryType_AUCTION_CLEARING_PRICE,
	"DAILY_VALUE_ADJUSTMENT_FOR_LONG_POSITIONS":       enum.MDEntryType_DAILY_VALUE_ADJUSTMENT_FOR_LONG_POSITIONS,
	"SWAP_VALUE_FACTOR":                               enum.MDEntryType_SWAP_VALUE_FACTOR,
	"CUMULATIVE_VALUE_ADJUSTMENT_FOR_LONG_POSITIONS":  enum.MDEntryType_CUMULATIVE_VALUE_ADJUSTMENT_FOR_LONG_POSITIONS,
	"DAILY_VALUE_ADJUSTMENT_FOR_SHORT_POSITIONS":      enum.MDEntryType_DAILY_VALUE_ADJUSTMENT_FOR_SHORT_POSITIONS,
	"CUMULATIVE_VALUE_ADJUSTMENT_FOR_SHORT_POSITIONS": enum.MDEntryType_CUMULATIVE_VALUE_ADJUSTMENT_FOR_SHORT_POSITIONS,
	"FIXING_PRICE":                     enum.MDEntryType_FIXING_PRICE,
	"CASH_RATE":                        enum.MDEntryType_CASH_RATE,
	"RECOVERY_RATE":                    enum.MDEntryType_RECOVERY_RATE,
	"RECOVERY_RATE_FOR_LONG_POSITIONS": enum.MDEntryType_RECOVERY_RATE_FOR_LONG_POSITIONS,
	"RECOVERY_RATE_FOR_SHORT":          enum.MDEntryType_RECOVERY_RATE_FOR_SHORT_POSITIONS,
	"TRADE_HISTORY":                    "101",
}

var MDEntryTypesReversed = map[enum.MDEntryType]string{
	enum.MDEntryType_BID:                                             "BID",
	enum.MDEntryType_OFFER:                                           "OFFER",
	enum.MDEntryType_TRADE:                                           "TRADE",
	enum.MDEntryType_INDEX_VALUE:                                     "INDEX_VALUE",
	enum.MDEntryType_OPENING_PRICE:                                   "OPENING_PRICE",
	enum.MDEntryType_CLOSING_PRICE:                                   "CLOSING_PRICE",
	enum.MDEntryType_SETTLEMENT_PRICE:                                "SETTLEMENT_PRICE",
	enum.MDEntryType_TRADING_SESSION_HIGH_PRICE:                      "TRADING_SESSION_HIGH_PRICE",
	enum.MDEntryType_TRADING_SESSION_LOW_PRICE:                       "TRADING_SESSION_LOW_PRICE",
	enum.MDEntryType_VOLUME_WEIGHTED_AVERAGE_PRICE:                   "VOLUME_WEIGHTED_AVERAGE_PRICE",
	enum.MDEntryType_IMBALANCE:                                       "IMBALANCE",
	enum.MDEntryType_TRADE_VOLUME:                                    "TRADE_VOLUME",
	enum.MDEntryType_OPEN_INTEREST:                                   "OPEN_INTEREST",
	enum.MDEntryType_COMPOSITE_UNDERLYING_PRICE:                      "COMPOSITE_UNDERLYING_PRICE",
	enum.MDEntryType_SIMULATED_SELL_PRICE:                            "SIMULATED_SELL_PRICE",
	enum.MDEntryType_SIMULATED_BUY_PRICE:                             "SIMULATED_BUY_PRICE",
	enum.MDEntryType_MARGIN_RATE:                                     "MARGIN_RATE",
	enum.MDEntryType_MID_PRICE:                                       "MID_PRICE",
	enum.MDEntryType_EMPTY_BOOK:                                      "EMPTY_BOOK",
	enum.MDEntryType_SETTLE_HIGH_PRICE:                               "SETTLE_HIGH_PRICE",
	enum.MDEntryType_SETTLE_LOW_PRICE:                                "SETTLE_LOW_PRICE",
	enum.MDEntryType_PRIOR_SETTLE_PRICE:                              "PRIOR_SETTLE_PRICE",
	enum.MDEntryType_SESSION_HIGH_BID:                                "SESSION_HIGH_BID",
	enum.MDEntryType_SESSION_LOW_OFFER:                               "SESSION_LOW_OFFER",
	enum.MDEntryType_EARLY_PRICES:                                    "EARLY_PRICES",
	enum.MDEntryType_AUCTION_CLEARING_PRICE:                          "AUCTION_CLEARING_PRICE",
	enum.MDEntryType_DAILY_VALUE_ADJUSTMENT_FOR_LONG_POSITIONS:       "DAILY_VALUE_ADJUSTMENT_FOR_LONG_POSITIONS",
	enum.MDEntryType_SWAP_VALUE_FACTOR:                               "SWAP_VALUE_FACTOR",
	enum.MDEntryType_CUMULATIVE_VALUE_ADJUSTMENT_FOR_LONG_POSITIONS:  "CUMULATIVE_VALUE_ADJUSTMENT_FOR_LONG_POSITIONS",
	enum.MDEntryType_DAILY_VALUE_ADJUSTMENT_FOR_SHORT_POSITIONS:      "DAILY_VALUE_ADJUSTMENT_FOR_SHORT_POSITIONS",
	enum.MDEntryType_CUMULATIVE_VALUE_ADJUSTMENT_FOR_SHORT_POSITIONS: "CUMULATIVE_VALUE_ADJUSTMENT_FOR_SHORT_POSITIONS",
	enum.MDEntryType_FIXING_PRICE:                                    "FIXING_PRICE",
	enum.MDEntryType_CASH_RATE:                                       "CASH_RATE",
	enum.MDEntryType_RECOVERY_RATE:                                   "RECOVERY_RATE",
	enum.MDEntryType_RECOVERY_RATE_FOR_LONG_POSITIONS:                "RECOVERY_RATE_FOR_LONG_POSITIONS",
	enum.MDEntryType_RECOVERY_RATE_FOR_SHORT_POSITIONS:               "RECOVERY_RATE_FOR_SHORT",
	"101": "TRADE_HISTORY",
}

var SubscriptionRequestTypes = map[string]enum.SubscriptionRequestType{
	"SNAPSHOT":              enum.SubscriptionRequestType_SNAPSHOT,
	"SNAPSHOT_PLUS_UPDATES": enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES,
	"DISABLE_PREVIOUS_SNAPSHOT_PLUS_UPDATE_REQUEST": enum.SubscriptionRequestType_DISABLE_PREVIOUS_SNAPSHOT_PLUS_UPDATE_REQUEST,
}

var MDUpdateTypes = map[string]enum.MDUpdateType{
	"FULL_REFRESH":        enum.MDUpdateType_FULL_REFRESH,
	"INCREMENTAL_REFRESH": enum.MDUpdateType_INCREMENTAL_REFRESH,
}
