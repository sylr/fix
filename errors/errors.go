package errors

import (
	"errors"
	"fmt"
)

var (
	Is  = errors.Is
	As  = errors.As
	New = errors.New
)

var (
	ErrConfig                          = errors.New("configuration")
	ErrConfigAcceptorNotFound          = fmt.Errorf("%w: acceptor not found", ErrConfig)
	ErrConfigAlreadyExists             = fmt.Errorf("%w: already exists", ErrConfig)
	ErrConfigCanNotBeCreated           = fmt.Errorf("%w: file can not be created", ErrConfig)
	ErrConfigContextMultipleSessions   = fmt.Errorf("%w: multiple sessions in initiator context", ErrConfig)
	ErrConfigContextNoSession          = fmt.Errorf("%w: context has no session", ErrConfig)
	ErrConfigContextNotFound           = fmt.Errorf("%w: context not found", ErrConfig)
	ErrConfigDuplicateContextName      = fmt.Errorf("%w: duplicate context name", ErrConfig)
	ErrConfigDuplicateAcceptorName     = fmt.Errorf("%w: duplicate acceptor name", ErrConfig)
	ErrConfigDuplicateInitiatorName    = fmt.Errorf("%w: duplicate initiator name", ErrConfig)
	ErrConfigDuplicateSessionName      = fmt.Errorf("%w: duplicate session name", ErrConfig)
	ErrConfigInitiatorNotFound         = fmt.Errorf("%w: initiator not found", ErrConfig)
	ErrConfigSessionNotFound           = fmt.Errorf("%w: session not found", ErrConfig)
	ErrConfigSessionNotInContext       = fmt.Errorf("%w: session name not in context", ErrConfig)
	ErrConnectionTimeout               = errors.New("connection timeout")
	ErrFix                             = errors.New("FIX")
	ErrFixLogout                       = fmt.Errorf("%w: logout received", ErrFix)
	ErrFixLogonRefused                 = fmt.Errorf("%w: logon refused", ErrFix)
	ErrFixRejected                     = fmt.Errorf("%w: message rejected", ErrFix)
	ErrFixOrderCanceled                = fmt.Errorf("%w: canceled order", ErrFix)
	ErrFixOrderRejected                = fmt.Errorf("%w: rejected order", ErrFix)
	ErrFixVersionNotImplemented        = fmt.Errorf("%w: version not implemented", ErrFix)
	ErrFixOrderStatusUnknown           = fmt.Errorf("%w: unknown order status", ErrFix)
	ErrFixQuoteRejected                = fmt.Errorf("%w: rejected quote", ErrFix)
	ErrOptions                         = errors.New("options")
	ErrOptionsInvalidMarketPrice       = fmt.Errorf("%w: can't give price for market order", ErrOptions)
	ErrOptionsNoSymbolGiven            = fmt.Errorf("%w: no symbol given", ErrOptions)
	ErrOptionsNoTypeGiven              = fmt.Errorf("%w: no type given", ErrOptions)
	ErrOptionsNoPriceGiven             = fmt.Errorf("%w: no price given", ErrOptions)
	ErrOptionsInconsistentValues       = fmt.Errorf("%w: inconsistent values", ErrOptions)
	ErrOptionOrderSideUnknown          = fmt.Errorf("%w: unknown order side", ErrOptions)
	ErrOptionOrderTypeUnknown          = fmt.Errorf("%w: unknown order type", ErrOptions)
	ErrOptionOrderOriginationUnknown   = fmt.Errorf("%w: unknown order origination", ErrOptions)
	ErrOptionOrderExpiryUnknown        = fmt.Errorf("%w: unknown order expiry", ErrOptions)
	ErrOptionOrderAttributeTypeUnknown = fmt.Errorf("%w: unknown order attribute type", ErrOptions)
	ErrOptionOrderRoleUnknown          = fmt.Errorf("%w: unknown order role", ErrOptions)
	ErrOptionOrderRoleQualifierUnknown = fmt.Errorf("%w: unknown order role qualifier", ErrOptions)
	ErrOptionOrderIDSourceUnknown      = fmt.Errorf("%w: unknown order id source", ErrOptions)
	ErrOptionPartySubIDTypeUnknown     = fmt.Errorf("%w: unknown party sub id type", ErrOptions)
	ErrResponseTimeout                 = errors.New("timeout while waiting for response")
)
