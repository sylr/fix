package enum

// OrderRestrictions field enumeration values.
type OrderOrigination string

const (
	OrderOriginations_CUSTOMER                                    OrderOrigination = "1"
	OrderOriginations_WITHIN_THE_FIRM                             OrderOrigination = "2"
	OrderOriginations_ANOTHER_BROKER_DEALER                       OrderOrigination = "3"
	OrderOriginations_CUSTOMER_OR_ORIGINATED_FROM_WITHIN_THE_FIRM OrderOrigination = "4"
	OrderOriginations_DIRECT_ACCESS_OR_SPONSORED_ACCESS_CUSTOMER  OrderOrigination = "5"
	OrderOriginations_FOREIGN_DEALER_EQUIVALENT                   OrderOrigination = "6"
	OrderOriginations_EXECUTION_ONLY_SERVICE                      OrderOrigination = "7"
)
