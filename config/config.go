package config

import (
	"fmt"
	"time"

	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/utils"
)

var (
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
	for k, session := range config.Sessions {
		if session.Name == name {
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
	Initiator   string
	Timeout     time.Duration
	Verbose     int
	Interactive bool
	LogCaller   bool
	FixLogging  bool
	Metrics     bool
	PProf       bool
	HTTPPort    int
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

	err = validateNames(f.Acceptors, errors.ConfigDuplicateAcceptorName)
	if err != nil {
		return err
	}

	for _, session := range f.Sessions {
		if err := session.validate(); err != nil {
			return fmt.Errorf("%w: session %s: %w", errors.Config, session.Name, err)
		}
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

// TLS configures the TLS of an initiator's or an acceptor's connections.
type TLS struct {
	// ServerName is the name an initiator checks the acceptor's certificate
	// against, the host of the address dialed if empty.
	ServerName         string `yaml:"serverName"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
	// CertFile and KeyFile are the certificate presented to the
	// counterparty, required for an acceptor.
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
	// CAFile holds the certificates of the authorities trusted, the
	// system's if empty.
	CAFile string `yaml:"caFile"`
	// ClientAuth makes an acceptor require and verify client certificates.
	ClientAuth bool `yaml:"clientAuth"`
}

// Store configures where sessions keep their sequence numbers and messages.
type Store struct {
	// Type is "memory", the default, or "sqlite3".
	Type string `yaml:"type"`
	// Path is the SQLite database, $HOME/.fix/initiator.db or
	// $HOME/.fix/acceptor.db if empty.
	Path string `yaml:"path"`
}

type Initiator struct {
	Name string `yaml:"name"`
	// Addresses are the acceptor's "host:port" addresses, tried in turn.
	Addresses []string `yaml:"addresses"`
	// TLS, if set, makes connections TLS.
	TLS *TLS `yaml:"tls"`
	// DialTimeout bounds each connection attempt, 5s if 0, --timeout if set.
	DialTimeout time.Duration `yaml:"dialTimeout"`
	// ReconnectMin and ReconnectMax bound the wait between connection
	// attempts, 1s and 30s if 0.
	ReconnectMin time.Duration `yaml:"reconnectMin"`
	ReconnectMax time.Duration `yaml:"reconnectMax"`
	Store        Store         `yaml:"store"`
}

func (i *Initiator) GetName() string {
	return i.Name
}

type Acceptor struct {
	Name string `yaml:"name"`
	// Address is the "host:port" address to listen on.
	Address string `yaml:"address"`
	TLS     *TLS   `yaml:"tls"`
	Store   Store  `yaml:"store"`
}

func (a *Acceptor) GetName() string {
	return a.Name
}

// Dictionaries are QuickFIX XML dictionaries replacing the ones fix is built
// with.
type Dictionaries struct {
	// Session is FIXT11.xml for FIXT.1.1 sessions.
	Session string `yaml:"session"`
	// App is the application dictionary, FIX50SP2.xml for FIX.5.0SP2.
	App string `yaml:"app"`
}

// Schedule is when a session runs: every day in Days, from Start to End,
// both "HH:MM:SS"; End before Start runs across midnight. With StartDay and
// EndDay, it runs once a week instead, from StartDay at Start to EndDay at
// End.
type Schedule struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
	// Days are the week days the session starts on, every day if empty.
	Days []string `yaml:"days"`
	// StartDay and EndDay make the schedule weekly.
	StartDay string `yaml:"startDay"`
	EndDay   string `yaml:"endDay"`
	// Location is the IANA time zone Start and End are in, UTC if empty.
	Location string `yaml:"location"`
}

type Session struct {
	Name string `yaml:"name"`

	BeginString      string `yaml:"beginString"`
	SenderCompID     string `yaml:"senderCompID"`
	SenderSubID      string `yaml:"senderSubID"`
	SenderLocationID string `yaml:"senderLocationID"`
	TargetCompID     string `yaml:"targetCompID"`
	TargetSubID      string `yaml:"targetSubID"`
	TargetLocationID string `yaml:"targetLocationID"`
	// Qualifier tells apart sessions with the same identities on the wire.
	Qualifier string `yaml:"qualifier"`
	// DefaultApplVerID is the application version of FIXT.1.1 sessions,
	// "FIX.5.0SP2" or its ApplVerID "9".
	DefaultApplVerID string `yaml:"defaultApplVerID"`

	Dictionaries Dictionaries `yaml:"dictionaries"`

	// HeartBtInt is the heartbeat interval, whole seconds, 30s if 0.
	HeartBtInt time.Duration `yaml:"heartBtInt"`
	// Username (553) and Password (554) are set on the Logons sent.
	Username string    `yaml:"username"`
	Password string    `yaml:"password"`
	Schedule *Schedule `yaml:"schedule"`

	ResetOnLogon      bool `yaml:"resetOnLogon"`
	ResetOnLogout     bool `yaml:"resetOnLogout"`
	ResetOnDisconnect bool `yaml:"resetOnDisconnect"`

	// LogonTimeout and LogoutTimeout default to 5s, --timeout if set.
	LogonTimeout  time.Duration `yaml:"logonTimeout"`
	LogoutTimeout time.Duration `yaml:"logoutTimeout"`
	// Timestamps is the precision of the timestamps sent: seconds,
	// milliseconds (the default), microseconds or nanoseconds.
	Timestamps string `yaml:"timestamps"`
	// MaxLatency is how far a received SendingTime (52) may be from our
	// clock, 120s if 0; negative disables the check.
	MaxLatency time.Duration `yaml:"maxLatency"`

	NextExpectedMsgSeqNum  bool `yaml:"nextExpectedMsgSeqNum"`
	LastMsgSeqNumProcessed bool `yaml:"lastMsgSeqNumProcessed"`

	// Tolerations name the departures from the FIX specifications the
	// session allows (see Tolerations).
	Tolerations []string `yaml:"tolerations"`
	// SkipChecks name the dictionary validations not run (see Checks).
	// "unknownTags" is skipped with the built-in FIX.5.0SP2 dictionary,
	// which leaves out the components fix does not use.
	SkipChecks []string `yaml:"skipChecks"`
}

func (s *Session) GetName() string {
	return s.Name
}
