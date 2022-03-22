package dict

import (
	"github.com/quickfixgo/enum"
)

var SubscriptionRequestTypes = map[enum.SubscriptionRequestType]string{
	enum.SubscriptionRequestType_SNAPSHOT:                                      "SNAPSHOT",
	enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES:                         "SNAPSHOT_PLUS_UPDATES",
	enum.SubscriptionRequestType_DISABLE_PREVIOUS_SNAPSHOT_PLUS_UPDATE_REQUEST: "DISABLE_PREVIOUS_SNAPSHOT_PLUS_UPDATE_REQUEST",
}

var SubscriptionRequestTypesReversed = map[string]enum.SubscriptionRequestType{
	"SNAPSHOT":              enum.SubscriptionRequestType_SNAPSHOT,
	"SNAPSHOT_PLUS_UPDATES": enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES,
	"DISABLE_PREVIOUS_SNAPSHOT_PLUS_UPDATE_REQUEST": enum.SubscriptionRequestType_DISABLE_PREVIOUS_SNAPSHOT_PLUS_UPDATE_REQUEST,
}
