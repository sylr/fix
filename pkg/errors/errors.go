package errors

import (
	"errors"
	"fmt"
)

var (
	Is = errors.Is
	As = errors.As
)

var (
	Config                          = errors.New("configuration")
	ConfigAcceptorNotFound          = fmt.Errorf("%w: acceptor not found", Config)
	ConfigAlreadyExists             = fmt.Errorf("%w: already exists", Config)
	ConfigCanNotBeCreated           = fmt.Errorf("%w: file can not be created", Config)
	ConfigContextMultipleSessions   = fmt.Errorf("%w: multiple sessions in initiator context", Config)
	ConfigContextNoSession          = fmt.Errorf("%w: context has no session", Config)
	ConfigContextNotFound           = fmt.Errorf("%w: context not found", Config)
	ConfigDuplicateContextName      = fmt.Errorf("%w: duplicate context name", Config)
	ConfigDuplicateInitiatorName    = fmt.Errorf("%w: duplicate acceptor name", Config)
	ConfigDuplicateSessionName      = fmt.Errorf("%w: duplicate session name", Config)
	ConfigInitiatorNotFound         = fmt.Errorf("%w: initiator not found", Config)
	ConfigSessionNotFound           = fmt.Errorf("%w: session not found", Config)
	ConfigSessionNotInContext       = fmt.Errorf("%w: session name not in context", Config)
	ConnectionTimeout               = errors.New("connection timeout")
	Fix                             = errors.New("FIX")
	FixLogout                       = fmt.Errorf("%w: logout received", Fix)
	FixOrderCanceled                = fmt.Errorf("%w: canceled order", Fix)
	FixOrderRejected                = fmt.Errorf("%w: rejected order", Fix)
	FixVersionNotImplemented        = fmt.Errorf("%w: version not implemented", Fix)
	FixOrderStatusUnknown           = fmt.Errorf("%w: unknown order status", Fix)
	NotImplemented                  = errors.New("not implemented")
	Options                         = errors.New("options")
	OptionsInvalidMarketPrice       = fmt.Errorf("%w: can't give price for market order", Options)
	OptionsNoSymbolGiven            = fmt.Errorf("%w: no symbol given", Options)
	OptionsNoTypeGiven              = fmt.Errorf("%w: no type given", Options)
	OptionsNoPriceGiven             = fmt.Errorf("%w: no price given", Options)
	OptionsInconsistentValues       = fmt.Errorf("%w: inconsistent values", Options)
	OptionOrderSideUnknown          = fmt.Errorf("%w: unknown order side", Options)
	OptionOrderTypeUnknown          = fmt.Errorf("%w: unknown order type", Options)
	OptionOrderOriginationUnknown   = fmt.Errorf("%w: unknown order origination", Options)
	OptionOrderAttributeTypeUnkonwn = fmt.Errorf("%w: unknown order attribute type", Options)
	OptionOrderRoleUnknown          = fmt.Errorf("%w: unknown order role", Options)
	OptionOrderRoleQualifierUnknown = fmt.Errorf("%w: unknown order role qualifier", Options)
	OptionOrderIDSourceUnknown      = fmt.Errorf("%w: unknown order id source", Options)
	OptionPartySubIDTypeUnknown     = fmt.Errorf("%w: unknown party sub id type", Options)
	ResponseTimeout                 = errors.New("timeout while waiting for response")
)
