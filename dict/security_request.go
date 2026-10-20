package dict

import (
	"fmt"
	"strings"

	"sylr.dev/fix/v2/fix50sp2"
)

// SecurityRequestType is the value of tag 321, which is not part of the
// trimmed FIX 5.0 SP2 dictionary used by fix50sp2.
type SecurityRequestType string

const (
	SecurityRequestTypeRequestSecurityIdentityAndSpecifications            SecurityRequestType = "0"
	SecurityRequestTypeRequestSecurityIdentityForTheSpecificationsProvided SecurityRequestType = "1"
	SecurityRequestTypeRequestListSecurityTypes                            SecurityRequestType = "2"
	SecurityRequestTypeRequestListSecurities                               SecurityRequestType = "3"
	SecurityRequestTypeSymbol                                              SecurityRequestType = "4"
	SecurityRequestTypeSecuritytypeAndOrCficode                            SecurityRequestType = "5"
	SecurityRequestTypeProduct                                             SecurityRequestType = "6"
	SecurityRequestTypeTradingsessionid                                    SecurityRequestType = "7"
	SecurityRequestTypeAllSecurities                                       SecurityRequestType = "8"
	SecurityRequestTypeMarketidOrMarketidPlusMarketsegmentid               SecurityRequestType = "9"
)

var SecurityRequestTypes = map[string]SecurityRequestType{
	"REQUEST_SECURITY_IDENTITY_AND_SPECIFICATIONS":              SecurityRequestTypeRequestSecurityIdentityAndSpecifications,
	"REQUEST_SECURITY_IDENTITY_FOR_THE_SPECIFICATIONS_PROVIDED": SecurityRequestTypeRequestSecurityIdentityForTheSpecificationsProvided,
	"REQUEST_LIST_SECURITY_TYPES":                               SecurityRequestTypeRequestListSecurityTypes,
	"REQUEST_LIST_SECURITIES":                                   SecurityRequestTypeRequestListSecurities,
	"SYMBOL":                                                    SecurityRequestTypeSymbol,
	"SECURITYTYPE_AND_OR_CFICODE":                               SecurityRequestTypeSecuritytypeAndOrCficode,
	"PRODUCT":                                                   SecurityRequestTypeProduct,
	"TRADINGSESSIONID":                                          SecurityRequestTypeTradingsessionid,
	"ALL_SECURITIES":                                            SecurityRequestTypeAllSecurities,
	"MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID":                 SecurityRequestTypeMarketidOrMarketidPlusMarketsegmentid,
}

func SecurityRequestTypeStringToEnum(side string) (SecurityRequestType, error) {
	side = strings.ToUpper(side)
	if e, ok := SecurityRequestTypes[side]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unkown security request type")
}

var SecurityListRequestTypes = map[string]fix50sp2.SecurityListRequestType{
	"SYMBOL":                      fix50sp2.SecurityListRequestTypeSymbol,
	"SECURITYTYPE_AND_OR_CFICODE": fix50sp2.SecurityListRequestTypeSecuritytypeAndOrCficode,
	"PRODUCT":                     fix50sp2.SecurityListRequestTypeProduct,
	"TRADINGSESSIONID":            fix50sp2.SecurityListRequestTypeTradingsessionid,
	"ALL_SECURITIES":              fix50sp2.SecurityListRequestTypeAllSecurities,
	"MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID": fix50sp2.SecurityListRequestTypeMarketidOrMarketidPlusMarketsegmentid,
}

func SecurityListRequestTypeStringToEnum(side string) (fix50sp2.SecurityListRequestType, error) {
	side = strings.ToUpper(side)
	if e, ok := SecurityListRequestTypes[side]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unkown security list request type")
}
