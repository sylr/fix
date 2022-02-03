package config

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/datadictionary"
)

var (
	fixDict             = make(map[string]*datadictionary.DataDictionary)
	options             = cliOptions{}
	config              = fixConfig{}
	ErrBadConfig        = errors.New("configuration error")
	ErrContextNotFound  = errors.New("context not found")
	ErrAcceptorNotFound = errors.New("acceptor not found")
	ErrSessionNotFound  = errors.New("session not found")
)

func GetOptions() *cliOptions {
	return &options
}

func GetConfig() *fixConfig {
	return &config
}

func GetCurrentContext() (*Context, error) {
	currentContext := options.Context

	return GetContext(currentContext)
}

func GetContext(name string) (*Context, error) {
	for k, context := range config.Contexts {
		if context.Name == name {
			return &config.Contexts[k], nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrContextNotFound, name)
}

func GetContexts() []Context {
	return config.Contexts
}

func GetAcceptor(name string) (*Acceptor, error) {
	for k, acceptor := range config.Acceptors {
		if acceptor.Name == name {
			return &config.Acceptors[k], nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrAcceptorNotFound, name)
}

func GetSession(name string) (*Session, error) {
	for k, acceptor := range config.Sessions {
		if acceptor.Name == name {
			return &config.Sessions[k], nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrAcceptorNotFound, name)
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
	Contexts  []Context  `yaml:"contexts"`
	Acceptors []Acceptor `yaml:"acceptors"`
	Sessions  []Session  `yaml:"sessions"`
}

type Context struct {
	Name     string `yaml:"name"`
	Acceptor string `yaml:"acceptor"`
	Session  string `yaml:"session"`
}

type Acceptor struct {
	Name                     string        `yaml:"name"`
	SocketConnectHost        string        `yaml:"SocketConnectHost"`
	SocketConnectPort        int           `yaml:"SocketConnectPort"`
	SocketInsecureSkipVerify bool          `yaml:"SocketInsecureSkipVerify"`
	SocketServerName         string        `yaml:"SocketServerName"`
	SocketTimeout            time.Duration `yaml:"SocketTimeout"`
	SocketUseSSL             bool          `yaml:"SocketUseSSL"`
}

type Session struct {
	Name                    string `yaml:"name"`
	BeginString             string `yaml:"BeginString"`
	DefaultApplVerID        string `yaml:"DefaultApplVerID"`
	HeartBtInt              int    `yaml:"HeartBtInt"`
	SenderCompID            string `yaml:"SenderCompID"`
	TargetCompID            string `yaml:"TargetCompID"`
	Username                string `yaml:"Username"`
	TransportDataDictionary string `yaml:"TransportDataDictionary"`
	AppDataDictionary       string `yaml:"AppDataDictionary"`
}

func (c Context) GetAcceptor() (*Acceptor, error) {
	return GetAcceptor(c.Acceptor)
}

func (c Context) GetSession() (*Session, error) {
	return GetSession(c.Session)
}

func (c Context) ToQuickFixSettings() (*quickfix.Settings, error) {
	qfSession := quickfix.NewSessionSettings()

	acceptor, err := GetAcceptor(c.Acceptor)
	if err != nil {
		return nil, err
	}

	session, err := GetSession(c.Session)
	if err != nil {
		return nil, err
	}

	// Acceptor settings
	qfSession.Set("SocketConnectHost", acceptor.SocketConnectHost)
	qfSession.Set("SocketConnectPort", FixIntString(acceptor.SocketConnectPort))
	qfSession.Set("SocketServerName", acceptor.SocketServerName)
	qfSession.Set("SocketUseSSL", FixBoolString(acceptor.SocketUseSSL))
	qfSession.Set("SocketInsecureSkipVerify", FixBoolString(acceptor.SocketInsecureSkipVerify))

	if options.Timeout != time.Duration(0) {
		qfSession.Set("SocketTimeout", options.Timeout.String())
	} else if acceptor.SocketTimeout != time.Duration(0) {
		qfSession.Set("SocketTimeout", acceptor.SocketTimeout.String())
	} else {
		qfSession.Set("SocketTimeout", "5s")
	}

	// Session settings
	qfSession.Set("HeartBtInt", FixIntString(session.HeartBtInt))
	qfSession.Set("BeginString", session.BeginString)
	qfSession.Set("DefaultApplVerID", session.DefaultApplVerID)
	qfSession.Set("Username", session.Username)
	qfSession.Set("TargetCompID", session.TargetCompID)
	qfSession.Set("SenderCompID", session.SenderCompID)
	qfSession.Set("BeginString", session.BeginString)

	if len(session.TransportDataDictionary) > 0 {
		qfSession.Set("TransportDataDictionary", session.TransportDataDictionary)
	}
	if len(session.TransportDataDictionary) > 0 {
		qfSession.Set("AppDataDictionary", session.AppDataDictionary)
	}

	if options.Timeout != time.Duration(0) {
		qfSession.Set("LogonTimeout", FixIntString(int(options.Timeout.Seconds())))
		qfSession.Set("LogonTimeout", FixIntString(int(options.Timeout.Seconds())))
	} else if acceptor.SocketTimeout != time.Duration(0) {
		qfSession.Set("LogonTimeout", FixIntString(int(acceptor.SocketTimeout.Seconds())))
		qfSession.Set("LogoutTimeout", FixIntString(int(acceptor.SocketTimeout.Seconds())))
	} else {
		qfSession.Set("LogonTimeout", "5")
		qfSession.Set("LogoutTimeout", "5")
	}

	settings := quickfix.NewSettings()
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
