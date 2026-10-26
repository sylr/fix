package dict

import (
	"sylr.dev/fix/v2/fix50sp2"
)

var SecurityListRequestTypes = map[string]fix50sp2.SecurityListRequestType{
	"SYMBOL":                      fix50sp2.SecurityListRequestTypeSymbol,
	"SECURITYTYPE_AND_OR_CFICODE": fix50sp2.SecurityListRequestTypeSecuritytypeAndOrCficode,
	"PRODUCT":                     fix50sp2.SecurityListRequestTypeProduct,
	"TRADINGSESSIONID":            fix50sp2.SecurityListRequestTypeTradingsessionid,
	"ALL_SECURITIES":              fix50sp2.SecurityListRequestTypeAllSecurities,
	"MARKETID_OR_MARKETID_PLUS_MARKETSEGMENTID": fix50sp2.SecurityListRequestTypeMarketidOrMarketidPlusMarketsegmentid,
}
