package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

func GetInitiator(name string) (*Initiator, error) {
	if len(name) == 0 {
		panic("empty initiator name")
	}

	for k, initiator := range config.Initiators {
		if initiator.Name == name {
			return config.Initiators[k], nil
		}
	}

	return nil, fmt.Errorf("%w: %s", errors.ConfigInitiatorNotFound, name)
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
	Config          string
	Context         string
	Session         string
	Acceptor        string
	Initiator       string
	Timeout         time.Duration
	Verbose         int
	Interactive     bool
	LogCaller       bool
	QuickFixLogging bool
}

type fixConfig struct {
	Contexts       []*Context   `yaml:"contexts"`
	Acceptors      []*Acceptor  `yaml:"acceptors"`
	Initiators     []*Initiator `yaml:"initiators"`
	Sessions       []*Session   `yaml:"sessions"`
	CurrentContext string       `yaml:"current-context"`
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

	err = validateNames(f.Initiators, errors.ConfigDuplicateInitiatorName)
	if err != nil {
		return err
	}

	return nil
}

type HasName interface {
	GetName() string
}

func validateNames[T HasName](slice []T, err error) error {
	names := make([]string, 0, len(slice))
	for _, item := range slice {
		name := item.GetName()
		if utils.Search(names, name) != -1 {
			return fmt.Errorf("%w: %s", err, name)
		} else {
			names = append(names, name)
		}
	}

	return nil
}

type Context struct {
	Name      string   `yaml:"name"`
	Initiator string   `yaml:"initiator"`
	Acceptor  string   `yaml:"acceptor"`
	Sessions  []string `yaml:"sessions"`
}

func (c *Context) GetName() string {
	return c.Name
}

type SQLStoreConfig interface {
	GetSQLStoreDriver() string
	GetSQLStoreDataSourceName() string
}

type common struct {
	Name                     string        `yaml:"name"`
	SocketUseSSL             bool          `yaml:"SocketUseSSL"`
	SocketInsecureSkipVerify bool          `yaml:"SocketInsecureSkipVerify"`
	SocketPrivateKeyFile     string        `yaml:"SocketPrivateKeyFile"`
	SocketCertificateFile    string        `yaml:"SocketCertificateFile"`
	SocketCAFile             string        `yaml:"SocketCAFile"`
	SocketTimeout            time.Duration `yaml:"SocketTimeout"`
	ResetOnDisconnect        bool          `yaml:"ResetOnDisconnect"`
	SQLStoreDriver           string        `yaml:"SQLStoreDriver"`
	SQLStoreDataSourceName   string        `yaml:"SQLStoreDataSourceName"`
}

func (c *common) GetName() string {
	return c.Name
}

func (c *common) GetSQLStoreDriver() string {
	return c.SQLStoreDriver
}

func (c *common) GetSQLStoreDataSourceName() string {
	return c.SQLStoreDataSourceName
}

func (c *common) setQuickFixGlobalSettings(settings *quickfix.SessionSettings, session *quickfix.SessionSettings) {
	settings.Set(qconfig.ResetOnDisconnect, FixBoolString(c.ResetOnDisconnect))
	session.Set(qconfig.SocketUseSSL, FixBoolString(c.SocketUseSSL))
	session.Set(qconfig.SocketInsecureSkipVerify, FixBoolString(c.SocketInsecureSkipVerify))

	if len(c.SocketPrivateKeyFile) != 0 {
		session.Set(qconfig.SocketPrivateKeyFile, c.SocketPrivateKeyFile)
	}
	if len(c.SocketCertificateFile) != 0 {
		session.Set(qconfig.SocketCertificateFile, c.SocketCertificateFile)
	}
	if len(c.SocketCAFile) != 0 {
		session.Set(qconfig.SocketCAFile, c.SocketCAFile)
	}

	if options.Timeout != time.Duration(0) {
		session.Set(qconfig.SocketTimeout, options.Timeout.String())
	} else if c.SocketTimeout != time.Duration(0) {
		session.Set(qconfig.SocketTimeout, c.SocketTimeout.String())
	} else {
		session.Set(qconfig.SocketTimeout, "5s")
	}
}

type Acceptor struct {
	common `yaml:",inline"`

	SocketAcceptHost string `yaml:"SocketAcceptHost"`
	SocketAcceptPort int    `yaml:"SocketAcceptPort"`
	UseTCPProxy      string `yaml:"UseTCPProxy"`
}

func (a *Acceptor) GetSQLStoreDriver() string {
	return a.common.GetSQLStoreDriver()
}

func (a *Acceptor) GetSQLStoreDataSourceName() string {
	return a.common.GetSQLStoreDataSourceName()
}

type Initiator struct {
	common `yaml:",inline"`

	SocketConnectHost string `yaml:"SocketConnectHost"`
	SocketConnectPort int    `yaml:"SocketConnectPort"`
	SocketServerName  string `yaml:"SocketServerName"`
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
	Password                string `yaml:"Password"`
	TransportDataDictionary string `yaml:"TransportDataDictionary"`
	AppDataDictionary       string `yaml:"AppDataDictionary"`
}

func (s *Session) GetName() string {
	return s.Name
}

func (c Context) GetInitiator() (*Initiator, error) {
	return GetInitiator(c.Initiator)
}

func (c Context) GetAcceptor() (*Acceptor, error) {
	return GetAcceptor(c.Acceptor)
}

func (c Context) GetSession(name string) (*Session, error) {
	if utils.Search(c.Sessions, name) < 0 {
		return nil, errors.ConfigSessionNotInContext
	}
	return GetSession(name)
}

func (c Context) GetSessions() ([]*Session, error) {
	sessions := make([]*Session, len(c.Sessions))
	for i, name := range c.Sessions {
		if session, err := GetSession(name); err != nil {
			return nil, err
		} else {
			sessions[i] = session
		}
	}
	return sessions, nil
}

func (c Context) quickFixGlobalSettings(*quickfix.Settings) error {
	return nil
}

func (c Context) ToQuickFixInitiatorSettings() (*quickfix.Settings, error) {
	settings := quickfix.NewSettings()
	qfSession := quickfix.NewSessionSettings()

	initiator, err := GetInitiator(c.Initiator)
	if err != nil {
		return nil, err
	}

	global := settings.GlobalSettings()
	initiator.setQuickFixGlobalSettings(global, qfSession)

	if len(initiator.SQLStoreDriver) > 0 {
		if initiator.SQLStoreDriver == "sqlite3" {
			if len(initiator.SQLStoreDataSourceName) == 0 {
				initiator.SQLStoreDataSourceName = os.ExpandEnv(strings.Join([]string{"$HOME", ".fix", "initiator.db"}, string(os.PathSeparator)))
			}
			global.Set(qconfig.SQLStoreDriver, initiator.SQLStoreDriver)
			global.Set(qconfig.SQLStoreDataSourceName, initiator.SQLStoreDataSourceName)
		}
	}

	sessions, err := c.GetSessions()
	if err != nil {
		return nil, err
	} else if len(sessions) != 1 {
		return nil, errors.ConfigContextMultipleSessions
	}

	// Session settings
	session := sessions[0]
	setSessionSetting(qfSession, qconfig.SocketConnectHost, initiator.SocketConnectHost)
	setSessionSetting(qfSession, qconfig.SocketConnectPort, initiator.SocketConnectPort)
	setSessionSetting(qfSession, qconfig.SocketServerName, initiator.SocketServerName)
	setSessionSetting(qfSession, qconfig.HeartBtInt, session.HeartBtInt)
	setSessionSetting(qfSession, qconfig.BeginString, session.BeginString)
	setSessionSetting(qfSession, qconfig.DefaultApplVerID, session.DefaultApplVerID)
	setSessionSetting(qfSession, qconfig.SenderCompID, session.SenderCompID)
	setSessionSetting(qfSession, qconfig.SenderSubID, session.SenderSubID)
	setSessionSetting(qfSession, qconfig.TargetCompID, session.TargetCompID)
	setSessionSetting(qfSession, qconfig.TargetSubID, session.TargetSubID)
	setSessionSetting(qfSession, qconfig.BeginString, session.BeginString)
	setSessionSetting(qfSession, "Username", session.Username)
	setSessionSetting(qfSession, "Password", session.Password)
	setSessionSetting(qfSession, qconfig.TransportDataDictionary, os.ExpandEnv(session.TransportDataDictionary))
	setSessionSetting(qfSession, qconfig.AppDataDictionary, os.ExpandEnv(session.AppDataDictionary))

	if options.Timeout != time.Duration(0) {
		qfSession.Set(qconfig.LogonTimeout, FixIntString(int(options.Timeout.Seconds())))
		qfSession.Set(qconfig.LogonTimeout, FixIntString(int(options.Timeout.Seconds())))
	} else if initiator.SocketTimeout != time.Duration(0) {
		qfSession.Set(qconfig.LogonTimeout, FixIntString(int(initiator.SocketTimeout.Seconds())))
		qfSession.Set(qconfig.LogoutTimeout, FixIntString(int(initiator.SocketTimeout.Seconds())))
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

func setSessionSetting(session *quickfix.SessionSettings, key string, value interface{}) {
	switch v := value.(type) {
	case bool:
		session.Set(key, FixBoolString(v))
	case time.Duration:
		if v.Seconds() > 0 {
			session.Set(key, FixIntString(int(v.Seconds())))
		}
	case int:
		session.Set(key, FixIntString(v))
	case string:
		if len(v) > 0 {
			session.Set(key, v)
		}
	default:
		panic("unknown type")
	}
}

func (c Context) ToQuickFixAcceptorSettings() (*quickfix.Settings, error) {
	settings := quickfix.NewSettings()
	qfSession := quickfix.NewSessionSettings()

	acceptor, err := GetAcceptor(c.Acceptor)
	if err != nil {
		return nil, err
	}

	global := settings.GlobalSettings()
	acceptor.setQuickFixGlobalSettings(global, qfSession)

	if len(acceptor.SQLStoreDriver) > 0 {
		if acceptor.SQLStoreDriver == "sqlite3" {
			if len(acceptor.SQLStoreDataSourceName) == 0 {
				acceptor.SQLStoreDataSourceName = os.ExpandEnv(strings.Join([]string{"$HOME", ".fix", "acceptor.db"}, string(os.PathSeparator)))
			}
			global.Set(qconfig.SQLStoreDriver, acceptor.SQLStoreDriver)
			global.Set(qconfig.SQLStoreDataSourceName, acceptor.SQLStoreDataSourceName)
		}
	}

	global.Set(qconfig.SocketAcceptHost, acceptor.SocketAcceptHost)
	global.Set(qconfig.SocketAcceptPort, strconv.Itoa(acceptor.SocketAcceptPort))

	sessions, err := c.GetSessions()
	if err != nil {
		return nil, err
	}

	for _, session := range sessions {
		setSessionSetting(qfSession, qconfig.HeartBtInt, session.HeartBtInt)
		setSessionSetting(qfSession, qconfig.BeginString, session.BeginString)
		setSessionSetting(qfSession, qconfig.DefaultApplVerID, session.DefaultApplVerID)
		setSessionSetting(qfSession, qconfig.SenderCompID, session.SenderCompID)
		setSessionSetting(qfSession, qconfig.SenderSubID, session.SenderSubID)
		setSessionSetting(qfSession, qconfig.TargetCompID, session.TargetCompID)
		setSessionSetting(qfSession, qconfig.TargetSubID, session.TargetSubID)
		setSessionSetting(qfSession, qconfig.BeginString, session.BeginString)
		setSessionSetting(qfSession, qconfig.TransportDataDictionary, os.ExpandEnv(session.TransportDataDictionary))
		setSessionSetting(qfSession, qconfig.AppDataDictionary, os.ExpandEnv(session.AppDataDictionary))

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
	}

	return settings, nil
}

func (s Session) GetFIXDictionaries() (*datadictionary.DataDictionary, *datadictionary.DataDictionary, error) {
	var err error
	var ok bool

	if len(s.TransportDataDictionary) > 0 {
		if _, ok = fixDict[s.TransportDataDictionary]; !ok {
			path := os.ExpandEnv(s.TransportDataDictionary)
			fixDict[s.TransportDataDictionary], err = datadictionary.Parse(path)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	if len(s.AppDataDictionary) > 0 {
		if _, ok = fixDict[s.AppDataDictionary]; !ok {
			path := os.ExpandEnv(s.AppDataDictionary)
			fixDict[s.AppDataDictionary], err = datadictionary.Parse(path)
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
