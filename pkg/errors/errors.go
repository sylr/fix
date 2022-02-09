package errors

import (
	"errors"
	"fmt"
)

var (
	ConnectionTimeout           = errors.New("connection timeout")
	NotImplemented              = errors.New("not implemented")
	Options                     = errors.New("options")
	Fix                         = errors.New("FIX")
	FixLogout                   = fmt.Errorf("%w: logout received", Fix)
	FixOrderStatusUnknown       = fmt.Errorf("%w: unknown order status", Fix)
	FixOrderTypeUnknown         = fmt.Errorf("%w: unknown order type", Fix)
	FixOrderSideUnknown         = fmt.Errorf("%w: unknown order side", Fix)
	FixOrderRejected            = fmt.Errorf("%w: rejected order", Fix)
	FixVersionNotImplemented    = fmt.Errorf("%w: version not implemented", Fix)
	Config                      = errors.New("configuration")
	ConfigContextNotFound       = fmt.Errorf("%w: context not found", Config)
	ConfigAcceptorNotFound      = fmt.Errorf("%w: acceptor not found", Config)
	ConfigSessionNotFound       = fmt.Errorf("%w: session not found", Config)
	ConfigDuplicateContextName  = fmt.Errorf("%w: duplicate context name", Config)
	ConfigDuplicateSessionName  = fmt.Errorf("%w: duplicate session name", Config)
	ConfigDuplicateAcceptorName = fmt.Errorf("%w: duplicate acceptor name", Config)
)
