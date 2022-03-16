package errors

import (
	"errors"
	"fmt"
)

var (
	Config                        = errors.New("configuration")
	ConfigAcceptorNotFound        = fmt.Errorf("%w: acceptor not found", Config)
	ConfigAlreadyExists           = fmt.Errorf("%w: already exists", Config)
	ConfigCanNotBeCreated         = fmt.Errorf("%w: file can not be created", Config)
	ConfigContextMultipleSessions = fmt.Errorf("%w: multiple sessions in initiator context", Config)
	ConfigContextNoSession        = fmt.Errorf("%w: context has no session", Config)
	ConfigContextNotFound         = fmt.Errorf("%w: context not found", Config)
	ConfigDuplicateContextName    = fmt.Errorf("%w: duplicate context name", Config)
	ConfigDuplicateInitiatorName  = fmt.Errorf("%w: duplicate acceptor name", Config)
	ConfigDuplicateSessionName    = fmt.Errorf("%w: duplicate session name", Config)
	ConfigInitiatorNotFound       = fmt.Errorf("%w: initiator not found", Config)
	ConfigSessionNotFound         = fmt.Errorf("%w: session not found", Config)
	ConfigSessionNotInContext     = fmt.Errorf("%w: session name not in context", Config)
	ConnectionTimeout             = errors.New("connection timeout")
	Fix                           = errors.New("FIX")
	FixLogout                     = fmt.Errorf("%w: logout received", Fix)
	FixOrderRejected              = fmt.Errorf("%w: rejected order", Fix)
	FixOrderSideUnknown           = fmt.Errorf("%w: unknown order side", Fix)
	FixOrderStatusUnknown         = fmt.Errorf("%w: unknown order status", Fix)
	FixOrderTypeUnknown           = fmt.Errorf("%w: unknown order type", Fix)
	FixVersionNotImplemented      = fmt.Errorf("%w: version not implemented", Fix)
	NotImplemented                = errors.New("not implemented")
	Options                       = errors.New("options")
	OptionsInvalidMarketPrice     = errors.New("%w: can't give price for market order")
	ResponseTimeout               = errors.New("timeout while waiting for response")
)
