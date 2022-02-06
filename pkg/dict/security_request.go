package dict

import (
	"fmt"
	"strings"

	"github.com/quickfixgo/enum"
)

var SecurityRequestTypes = map[enum.SecurityRequestType]string{
	enum.SecurityRequestType_REQUEST_SECURITY_IDENTITY_AND_SPECIFICATIONS:              "REQUEST_SECURITY_IDENTITY_AND_SPECIFICATIONS",
	enum.SecurityRequestType_REQUEST_SECURITY_IDENTITY_FOR_THE_SPECIFICATIONS_PROVIDED: "REQUEST_SECURITY_IDENTITY_FOR_THE_SPECIFICATIONS_PROVIDED",
	enum.SecurityRequestType_REQUEST_LIST_SECURITY_TYPES:                               "REQUEST_LIST_SECURITY_TYPES",
	enum.SecurityRequestType_REQUEST_LIST_SECURITIES:                                   "REQUEST_LIST_SECURITIES",
	enum.SecurityRequestType_SYMBOL:                                                    "SYMBOL",
	enum.SecurityRequestType_SECURITYTYPE_AND_OR_CFICODE:                               "SECURITYTYPE_AND_OR_CFICODE",
	enum.SecurityRequestType_PRODUCT:                                                   "PRODUCT",
	enum.SecurityRequestType_TRADINGSESSIONID:                                          "TRADINGSESSIONID",
	enum.SecurityRequestType_ALL_SECURITIES:                                            "ALL_SECURITIES",
	enum.SecurityRequestType_MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID:                 "MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID",
}

var SecurityRequestTypesReversed = map[string]enum.SecurityRequestType{
	"REQUEST_SECURITY_IDENTITY_AND_SPECIFICATIONS":              enum.SecurityRequestType_REQUEST_SECURITY_IDENTITY_AND_SPECIFICATIONS,
	"REQUEST_SECURITY_IDENTITY_FOR_THE_SPECIFICATIONS_PROVIDED": enum.SecurityRequestType_REQUEST_SECURITY_IDENTITY_FOR_THE_SPECIFICATIONS_PROVIDED,
	"REQUEST_LIST_SECURITY_TYPES":                               enum.SecurityRequestType_REQUEST_LIST_SECURITY_TYPES,
	"REQUEST_LIST_SECURITIES":                                   enum.SecurityRequestType_REQUEST_LIST_SECURITIES,
	"SYMBOL":                                                    enum.SecurityRequestType_SYMBOL,
	"SECURITYTYPE_AND_OR_CFICODE":                               enum.SecurityRequestType_SECURITYTYPE_AND_OR_CFICODE,
	"PRODUCT":                                                   enum.SecurityRequestType_PRODUCT,
	"TRADINGSESSIONID":                                          enum.SecurityRequestType_TRADINGSESSIONID,
	"ALL_SECURITIES":                                            enum.SecurityRequestType_ALL_SECURITIES,
	"MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID":                 enum.SecurityRequestType_MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID,
}

func SecurityRequestTypeStringToEnum(side string) (enum.SecurityRequestType, error) {
	side = strings.ToUpper(side)
	if e, ok := SecurityRequestTypesReversed[side]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unkown security request type")
}

var SecurityListRequestTypes = map[enum.SecurityListRequestType]string{
	enum.SecurityListRequestType_SYMBOL:                                    "SYMBOL",
	enum.SecurityListRequestType_SECURITYTYPE_AND_OR_CFICODE:               "SECURITYTYPE_AND_OR_CFICODE",
	enum.SecurityListRequestType_PRODUCT:                                   "PRODUCT",
	enum.SecurityListRequestType_TRADINGSESSIONID:                          "TRADINGSESSIONID",
	enum.SecurityListRequestType_ALL_SECURITIES:                            "ALL_SECURITIES",
	enum.SecurityListRequestType_MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID: "MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID",
}

var SecurityListRequestTypesReversed = map[string]enum.SecurityListRequestType{
	"SYMBOL":                      enum.SecurityListRequestType_SYMBOL,
	"SECURITYTYPE_AND_OR_CFICODE": enum.SecurityListRequestType_SECURITYTYPE_AND_OR_CFICODE,
	"PRODUCT":                     enum.SecurityListRequestType_PRODUCT,
	"TRADINGSESSIONID":            enum.SecurityListRequestType_TRADINGSESSIONID,
	"ALL_SECURITIES":              enum.SecurityListRequestType_ALL_SECURITIES,
	"MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID": enum.SecurityListRequestType_MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID,
}

func SecurityListRequestTypeStringToEnum(side string) (enum.SecurityListRequestType, error) {
	side = strings.ToUpper(side)
	if e, ok := SecurityListRequestTypesReversed[side]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unkown security list request type")
}
