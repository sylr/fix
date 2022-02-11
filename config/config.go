package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/quickfixgo/quickfix"
	qconfig "github.com/quickfixgo/quickfix/config"
	"github.com/quickfixgo/quickfix/datadictionary"

	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/utils"
)

var (
	fixDict = make(map[string]*datadictionary.DataDictionary)
	options = cliOptions{}
	config  = fixConfig{}
)

func GetOptions() *cliOptions {
	return &options
}

func GetConfig() *fixConfig {
	return &config
}

func GetCurrentContext() (*Context, error) {
	currentContext := config.CurrentContext

	if len(options.Context) > 0 {
		currentContext = options.Context
	}

	if len(currentContext) == 0 {
		return nil, fmt.Errorf("%w:no current-context set and no --context given", errors.Config)
	}

	return GetContext(currentContext)
}

func GetContext(name string) (*Context, error) {
	for k, context := range config.Contexts {
		if context.Name == name {
			return config.Contexts[k], nil
		}
	}

	return nil, fmt.Errorf("%w: %s", errors.ConfigContextNotFound, name)
}

func GetContexts() []*Context {
	return config.Contexts
}

func GetAcceptor(name string) (*Acceptor, error) {
	for k, acceptor := range config.Acceptors {
		if acceptor.Name == name {
			return config.Acceptors[k], nil
		}
	}

	return nil, fmt.Errorf("%w: %s", errors.ConfigAcceptorNotFound, name)
}

func GetSession(name string) (*Session, error) {
	for k, acceptor := range config.Sessions {
		if acceptor.Name == name {
			return config.Sessions[k], nil
		}
	}

	return nil, fmt.Errorf("%w: %s", errors.ConfigSessionNotFound, name)
}

type cliOptions struct {
	Config      string
	Context     string
	Session     string
	Acceptor    string
	Timeout     time.Duration
	Verbose     int
	Interactive bool
	LogCaller   bool
}

type fixConfig struct {
	Contexts       []*Context  `yaml:"contexts"`
	Acceptors      []*Acceptor `yaml:"acceptors"`
	Sessions       []*Session  `yaml:"sessions"`
	CurrentContext string      `yaml:"current-context"`
}

func (f *fixConfig) Validate() error {
	err := validateNames(f.Contexts, errors.ConfigDuplicateContextName)
	if err != nil {
		return err
	}

	err = validateNames(f.Sessions, errors.ConfigDuplicateSessionName)
	if err != nil {
		return err
	}

	err = validateNames(f.Acceptors, errors.ConfigAcceptorNotFound)
	if err != nil {
		return err
	}

	return nil
}

type HasName interface {
	GetName() string
}

func validateNames[T HasName](slice []T, err error) error {
	contextNames := make([]string, 0, len(slice))

	for _, item := range slice {
		name := item.GetName()
		if utils.Search(contextNames, name) != -1 {
			return fmt.Errorf("%w: %s", err, name)
		} else {
			contextNames = append(contextNames, name)
		}
	}

	return nil
}

type Context struct {
	Name     string `yaml:"name"`
	Acceptor string `yaml:"acceptor"`
	Session  string `yaml:"session"`
}

func (c *Context) GetName() string {
	return c.Name
}

type Acceptor struct {
	Name string `yaml:"name"`

	// Initiator
	SocketConnectHost string        `yaml:"SocketConnectHost"`
	SocketConnectPort int           `yaml:"SocketConnectPort"`
	SocketTimeout     time.Duration `yaml:"SocketTimeout"`

	// Acceptor
	SocketAcceptHost string `yaml:"SocketAcceptHost"`
	SocketAcceptPort int    `yaml:"SocketAcceptPort"`
	UseTCPProxy      string `yaml:"UseTCPProxy"`

	// Common
	SocketUseSSL             bool   `yaml:"SocketUseSSL"`
	SocketServerName         string `yaml:"SocketServerName"`
	SocketInsecureSkipVerify bool   `yaml:"SocketInsecureSkipVerify"`
	SocketPrivateKeyFile     string `yaml:"SocketPrivateKeyFile"`
	SocketCertificateFile    string `yaml:"SocketCertificateFile"`
	SocketCAFile             string `yaml:"SocketCAFile"`
}

func (a *Acceptor) GetName() string {
	return a.Name
}

type Session struct {
	Name                    string `yaml:"name"`
	BeginString             string `yaml:"BeginString"`
	DefaultApplVerID        string `yaml:"DefaultApplVerID"`
	HeartBtInt              int    `yaml:"HeartBtInt"`
	SenderCompID            string `yaml:"SenderCompID"`
	SenderSubID             string `yaml:"SenderSubID"`
	TargetCompID            string `yaml:"TargetCompID"`
	TargetSubID             string `yaml:"TargetSubID"`
	Username                string `yaml:"Username"`
	TransportDataDictionary string `yaml:"TransportDataDictionary"`
	AppDataDictionary       string `yaml:"AppDataDictionary"`
}

func (s *Session) GetName() string {
	return s.Name
}

func (c Context) GetAcceptor() (*Acceptor, error) {
	return GetAcceptor(c.Acceptor)
}

func (c Context) GetSession() (*Session, error) {
	return GetSession(c.Session)
}

func (c Context) ToQuickFixSettings() (*quickfix.Settings, error) {
	settings := quickfix.NewSettings()
	qfSession := quickfix.NewSessionSettings()

	acceptor, err := GetAcceptor(c.Acceptor)
	if err != nil {
		return nil, err
	}

	session, err := GetSession(c.Session)
	if err != nil {
		return nil, err
	}

	global := settings.GlobalSettings()

	// Global settings
	global.Set(qconfig.SocketAcceptHost, acceptor.SocketAcceptHost)
	global.Set(qconfig.SocketAcceptPort, FixIntString(acceptor.SocketAcceptPort))
	global.Set(qconfig.ResetOnDisconnect, "Y")

	// Acceptor settings
	qfSession.Set(qconfig.SocketConnectHost, acceptor.SocketConnectHost)
	qfSession.Set(qconfig.SocketConnectPort, FixIntString(acceptor.SocketConnectPort))
	qfSession.Set(qconfig.SocketServerName, acceptor.SocketServerName)
	qfSession.Set(qconfig.SocketUseSSL, FixBoolString(acceptor.SocketUseSSL))
	qfSession.Set(qconfig.SocketInsecureSkipVerify, FixBoolString(acceptor.SocketInsecureSkipVerify))

	if options.Timeout != time.Duration(0) {
		qfSession.Set(qconfig.SocketTimeout, options.Timeout.String())
	} else if acceptor.SocketTimeout != time.Duration(0) {
		qfSession.Set(qconfig.SocketTimeout, acceptor.SocketTimeout.String())
	} else {
		qfSession.Set(qconfig.SocketTimeout, "5s")
	}

	// Session settings
	qfSession.Set(qconfig.HeartBtInt, FixIntString(session.HeartBtInt))
	qfSession.Set(qconfig.BeginString, session.BeginString)
	qfSession.Set(qconfig.DefaultApplVerID, session.DefaultApplVerID)
	qfSession.Set(qconfig.SenderCompID, session.SenderCompID)
	qfSession.Set(qconfig.SenderSubID, session.SenderSubID)
	qfSession.Set(qconfig.TargetCompID, session.TargetCompID)
	qfSession.Set(qconfig.TargetSubID, session.TargetSubID)
	qfSession.Set(qconfig.BeginString, session.BeginString)
	qfSession.Set("Username", session.Username)

	if len(session.TransportDataDictionary) > 0 {
		qfSession.Set(qconfig.TransportDataDictionary, session.TransportDataDictionary)
	}
	if len(session.TransportDataDictionary) > 0 {
		qfSession.Set(qconfig.AppDataDictionary, session.AppDataDictionary)
	}

	if options.Timeout != time.Duration(0) {
		qfSession.Set(qconfig.LogonTimeout, FixIntString(int(options.Timeout.Seconds())))
		qfSession.Set(qconfig.LogonTimeout, FixIntString(int(options.Timeout.Seconds())))
	} else if acceptor.SocketTimeout != time.Duration(0) {
		qfSession.Set(qconfig.LogonTimeout, FixIntString(int(acceptor.SocketTimeout.Seconds())))
		qfSession.Set(qconfig.LogoutTimeout, FixIntString(int(acceptor.SocketTimeout.Seconds())))
	} else {
		qfSession.Set(qconfig.LogonTimeout, "5")
		qfSession.Set(qconfig.LogoutTimeout, "5")
	}

	_, err = settings.AddSession(qfSession)

	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (s Session) GetFIXDictionaries() (*datadictionary.DataDictionary, *datadictionary.DataDictionary, error) {
	var err error
	var ok bool

	if len(s.TransportDataDictionary) > 0 {
		if _, ok = fixDict[s.TransportDataDictionary]; !ok {
			fixDict[s.TransportDataDictionary], err = datadictionary.Parse(s.TransportDataDictionary)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	if len(s.AppDataDictionary) > 0 {
		if _, ok = fixDict[s.AppDataDictionary]; !ok {
			fixDict[s.AppDataDictionary], err = datadictionary.Parse(s.AppDataDictionary)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return fixDict[s.TransportDataDictionary], fixDict[s.AppDataDictionary], nil
}

func FixBoolString(b bool) string {
	if b {
		return "Y"
	}

	return "N"
}

func FixIntString(i int) string {
	return strconv.Itoa(i)
}
