package dict

import (
	"fmt"
	"strings"

	"github.com/quickfixgo/enum"
)

var SecurityRequestTypes = map[string]enum.SecurityRequestType{
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
	if e, ok := SecurityRequestTypes[side]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unkown security request type")
}

var SecurityListRequestTypes = map[string]enum.SecurityListRequestType{
	"SYMBOL":                      enum.SecurityListRequestType_SYMBOL,
	"SECURITYTYPE_AND_OR_CFICODE": enum.SecurityListRequestType_SECURITYTYPE_AND_OR_CFICODE,
	"PRODUCT":                     enum.SecurityListRequestType_PRODUCT,
	"TRADINGSESSIONID":            enum.SecurityListRequestType_TRADINGSESSIONID,
	"ALL_SECURITIES":              enum.SecurityListRequestType_ALL_SECURITIES,
	"MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID": enum.SecurityListRequestType_MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID,
}

func SecurityListRequestTypeStringToEnum(side string) (enum.SecurityListRequestType, error) {
	side = strings.ToUpper(side)
	if e, ok := SecurityListRequestTypes[side]; ok {
		return e, nil
	}

	return "", fmt.Errorf("unkown security list request type")
}
