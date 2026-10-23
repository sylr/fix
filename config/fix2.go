package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"maps"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/dictionary/xmldict"
	"github.com/sylr/fix2/fix/session"
	"github.com/sylr/fix2/fix/session/schedule"
	"github.com/sylr/fix2/fix/session/store"
	"github.com/sylr/fix2/fix/session/store/memstore"
	"github.com/sylr/fix2/transport/tcp"

	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/fixt11"
)

const defaultTimeout = 5 * time.Second

// ApplVerIDs maps the names of the FIX application versions to their
// ApplVerID (1128) values.
var ApplVerIDs = map[string]string{
	"FIX.2.7":    "0",
	"FIX.3.0":    "1",
	"FIX.4.0":    "2",
	"FIX.4.1":    "3",
	"FIX.4.2":    "4",
	"FIX.4.3":    "5",
	"FIX.4.4":    "6",
	"FIX.5.0":    "7",
	"FIX.5.0SP1": "8",
	"FIX.5.0SP2": "9",
}

// Tolerations maps the names a session's tolerations are configured with to
// what they set. "quickfix" sets those QuickFIX's acceptance tests need.
var Tolerations = map[string]func(*session.Tolerations){
	"quickfix": func(t *session.Tolerations) {
		q := session.QuickFIXTolerations()
		t.TolerateUnsupportedDefaultApplVerID = t.TolerateUnsupportedDefaultApplVerID || q.TolerateUnsupportedDefaultApplVerID
		t.TolerateLogoutBeforeGapRecovery = t.TolerateLogoutBeforeGapRecovery || q.TolerateLogoutBeforeGapRecovery
		t.TolerateNoHeartbeatAwaitingTestRequest = t.TolerateNoHeartbeatAwaitingTestRequest || q.TolerateNoHeartbeatAwaitingTestRequest
	},
	"logonWithoutEncryptMethod":           func(t *session.Tolerations) { t.TolerateLogonWithoutEncryptMethod = true },
	"unsupportedDefaultApplVerID":         func(t *session.Tolerations) { t.TolerateUnsupportedDefaultApplVerID = true },
	"logoutBeforeGapRecovery":             func(t *session.Tolerations) { t.TolerateLogoutBeforeGapRecovery = true },
	"sequenceResetWithoutOrigSendingTime": func(t *session.Tolerations) { t.TolerateSequenceResetWithoutOrigSendingTime = true },
	"resetLogonSeqNumNotOne":              func(t *session.Tolerations) { t.TolerateResetLogonSeqNumNotOne = true },
	"malformedNextExpectedMsgSeqNum":      func(t *session.Tolerations) { t.TolerateMalformedNextExpectedMsgSeqNum = true },
	"noHeartbeatAwaitingTestRequest":      func(t *session.Tolerations) { t.TolerateNoHeartbeatAwaitingTestRequest = true },
	"tagLeadingZeros":                     func(t *session.Tolerations) { t.TolerateTagLeadingZeros = true },
	"finerTimestamps":                     func(t *session.Tolerations) { t.TolerateFinerTimestamps = true },
	"emptyRequiredGroups":                 func(t *session.Tolerations) { t.TolerateEmptyRequiredGroups = true },
}

// Checks maps the names of the dictionary validations a session may skip to
// their flags.
var Checks = map[string]dictionary.Checks{
	"emptyValues":         dictionary.CheckEmptyValues,
	"formats":             dictionary.CheckFormats,
	"repeatedTags":        dictionary.CheckRepeatedTags,
	"order":               dictionary.CheckOrder,
	"groupCounts":         dictionary.CheckGroupCounts,
	"unknownTags":         dictionary.CheckUnknownTags,
	"userDefinedTags":     dictionary.CheckUserDefinedTags,
	"groupRequired":       dictionary.CheckGroupRequired,
	"groupOrder":          dictionary.CheckGroupOrder,
	"timestampPrecision":  dictionary.CheckTimestampPrecision,
	"emptyRequiredGroups": dictionary.CheckEmptyRequiredGroups,
}

var timestamps = map[string]session.TimestampPrecision{
	"":             session.DefaultTimestamps,
	"seconds":      session.SecondTimestamps,
	"milliseconds": session.MillisecondTimestamps,
	"microseconds": session.MicrosecondTimestamps,
	"nanoseconds":  session.NanosecondTimestamps,
}

func names[V any](m map[string]V) string {
	return strings.Join(slices.Sorted(maps.Keys(m)), ", ")
}

// validate checks what Config.Validate in fix2 cannot: the names and
// formats fix2 has no notion of.
func (s *Session) validate() error {
	for _, t := range s.Tolerations {
		if _, ok := Tolerations[t]; !ok {
			return fmt.Errorf("unknown toleration %q (known: %s)", t, names(Tolerations))
		}
	}

	for _, c := range s.SkipChecks {
		if _, ok := Checks[c]; !ok {
			return fmt.Errorf("unknown check %q (known: %s)", c, names(Checks))
		}
	}

	if _, ok := timestamps[s.Timestamps]; !ok {
		return fmt.Errorf("unknown timestamps %q (known: seconds, milliseconds, microseconds, nanoseconds)", s.Timestamps)
	}

	if _, err := s.schedule(); err != nil {
		return err
	}

	return nil
}

// ID returns the identity of the session.
func (s *Session) ID() session.SessionID {
	return session.SessionID{
		BeginString:      s.BeginString,
		SenderCompID:     s.SenderCompID,
		SenderSubID:      s.SenderSubID,
		SenderLocationID: s.SenderLocationID,
		TargetCompID:     s.TargetCompID,
		TargetSubID:      s.TargetSubID,
		TargetLocationID: s.TargetLocationID,
		Qualifier:        s.Qualifier,
	}
}

// ApplVerID returns the ApplVerID (1128) value of DefaultApplVerID.
func (s *Session) ApplVerID() string {
	if v, ok := ApplVerIDs[s.DefaultApplVerID]; ok {
		return v
	}

	return s.DefaultApplVerID
}

// IsFIX50SP2 reports whether the session carries FIX 5.0 SP2 application
// messages, those of package fix50sp2.
func (s *Session) IsFIX50SP2() bool {
	return s.BeginString == session.FIXT11 && s.ApplVerID() == fix50sp2.Version.ApplVerID
}

func loadDictionary(path string) (*dictionary.Dict, error) {
	d, err := xmldict.LoadFile(os.ExpandEnv(path))
	if err != nil {
		return nil, fmt.Errorf("loading dictionary %s: %w", path, err)
	}

	return d, nil
}

// FIXDictionaries returns the dictionaries the session validates messages
// with, and the one describing every field of its messages, header and
// trailer included, for printing them. fix's own dictionaries are used
// unless Dictionaries name others.
func (s *Session) FIXDictionaries() (session.Dictionaries, *dictionary.Dict, error) {
	var dicts session.Dictionaries
	var err error

	if s.BeginString != session.FIXT11 {
		path := s.Dictionaries.App
		if path == "" {
			path = s.Dictionaries.Session
		}
		if path == "" {
			return dicts, nil, nil
		}

		dicts.Session, err = loadDictionary(path)
		return dicts, dicts.Session, err
	}

	dicts.Session = fixt11.Dictionary()
	if s.Dictionaries.Session != "" {
		if dicts.Session, err = loadDictionary(s.Dictionaries.Session); err != nil {
			return dicts, nil, err
		}
	}

	dicts.DefaultApplVerID = s.ApplVerID()

	var app *dictionary.Dict
	switch {
	case s.Dictionaries.App != "":
		if app, err = loadDictionary(s.Dictionaries.App); err != nil {
			return dicts, nil, err
		}
	case dicts.DefaultApplVerID == fix50sp2.Version.ApplVerID:
		app = fix50sp2.Dictionary()
	default:
		return dicts, dicts.Session, nil
	}

	dicts.Apps = map[string]*dictionary.Dict{dicts.DefaultApplVerID: app}

	return dicts, dictionary.Compose(dicts.Session, app), nil
}

// skipChecks returns the validations the session skips.
func (s *Session) skipChecks() dictionary.Checks {
	var checks dictionary.Checks

	// The built-in FIX 5.0 SP2 dictionary lacks the components fix does not
	// use, whose tags counterparties may still send.
	if s.BeginString == session.FIXT11 && s.Dictionaries.App == "" {
		checks |= dictionary.CheckUnknownTags
	}

	for _, c := range s.SkipChecks {
		checks |= Checks[c]
	}

	return checks
}

func (s *Session) schedule() (schedule.Schedule, error) {
	if s.Schedule == nil {
		return nil, nil
	}

	clock := func(v string) (schedule.ClockTime, error) {
		t, err := time.Parse(time.TimeOnly, v)
		if err != nil {
			return schedule.ClockTime{}, fmt.Errorf("schedule: %q is not HH:MM:SS", v)
		}
		return schedule.Clock(t.Hour(), t.Minute(), t.Second()), nil
	}

	weekday := func(v string) (time.Weekday, error) {
		day, ok := weekdays[strings.ToLower(v)]
		if !ok {
			return 0, fmt.Errorf("schedule: unknown day %q", v)
		}
		return day, nil
	}

	start, err := clock(s.Schedule.Start)
	if err != nil {
		return nil, err
	}
	end, err := clock(s.Schedule.End)
	if err != nil {
		return nil, err
	}

	var location *time.Location
	if s.Schedule.Location != "" {
		if location, err = time.LoadLocation(s.Schedule.Location); err != nil {
			return nil, fmt.Errorf("schedule: %w", err)
		}
	}

	if s.Schedule.StartDay != "" || s.Schedule.EndDay != "" {
		if s.Schedule.StartDay == "" || s.Schedule.EndDay == "" || len(s.Schedule.Days) > 0 {
			return nil, fmt.Errorf("schedule: a weekly schedule takes startDay and endDay, without days")
		}

		weekly := schedule.Weekly{Start: start, End: end, Location: location}
		if weekly.StartDay, err = weekday(s.Schedule.StartDay); err != nil {
			return nil, err
		}
		if weekly.EndDay, err = weekday(s.Schedule.EndDay); err != nil {
			return nil, err
		}

		return weekly, nil
	}

	daily := schedule.Daily{Start: start, End: end, Location: location}
	for _, d := range s.Schedule.Days {
		day, err := weekday(d)
		if err != nil {
			return nil, err
		}
		daily.Days = append(daily.Days, day)
	}

	return daily, nil
}

var weekdays = map[string]time.Weekday{
	"sunday": time.Sunday, "sun": time.Sunday,
	"monday": time.Monday, "mon": time.Monday,
	"tuesday": time.Tuesday, "tue": time.Tuesday,
	"wednesday": time.Wednesday, "wed": time.Wednesday,
	"thursday": time.Thursday, "thu": time.Thursday,
	"friday": time.Friday, "fri": time.Friday,
	"saturday": time.Saturday, "sat": time.Saturday,
}

// Timeout returns --timeout if given, else v if set, else 5s.
func Timeout(v time.Duration) time.Duration {
	switch {
	case options.Timeout > 0:
		return options.Timeout
	case v > 0:
		return v
	}

	return defaultTimeout
}

// FIXConfig returns the fix2 configuration of the session, without its
// Store, which the initiator or acceptor sets.
func (s *Session) FIXConfig(handler session.Handler, logger *slog.Logger, observer session.MessageObserver) (session.Config, error) {
	dicts, _, err := s.FIXDictionaries()
	if err != nil {
		return session.Config{}, err
	}

	sched, err := s.schedule()
	if err != nil {
		return session.Config{}, err
	}

	cfg := session.Config{
		ID:                     s.ID(),
		Dictionaries:           dicts,
		HeartBtInt:             s.HeartBtInt,
		Schedule:               sched,
		Logger:                 logger,
		Observer:               observer,
		Handler:                handler,
		Timestamps:             timestamps[s.Timestamps],
		MaxLatency:             s.MaxLatency,
		ResetOnLogon:           s.ResetOnLogon,
		ResetOnLogout:          s.ResetOnLogout,
		ResetOnDisconnect:      s.ResetOnDisconnect,
		NextExpectedMsgSeqNum:  s.NextExpectedMsgSeqNum,
		LastMsgSeqNumProcessed: s.LastMsgSeqNumProcessed,
		LogonTimeout:           Timeout(s.LogonTimeout),
		LogoutTimeout:          Timeout(s.LogoutTimeout),
	}

	for _, t := range s.Tolerations {
		Tolerations[t](&cfg.Tolerations)
	}
	cfg.TolerateValidationFailures = s.skipChecks()

	return cfg, cfg.Validate()
}

// IsMemory reports whether the store keeps sessions in memory, the default.
func (st Store) IsMemory() bool {
	return st.Type == "" || st.Type == "memory"
}

func (st Store) validate() error {
	if !st.IsMemory() && st.Type != "sqlite3" {
		return fmt.Errorf("unknown store type %q (known: memory, sqlite3)", st.Type)
	}

	return nil
}

// Open opens the store of the session id, in the SQLite database
// $HOME/.fix/<defaultName> unless Path names another.
func (st Store) Open(defaultName string, id session.SessionID) (store.Store, error) {
	if err := st.validate(); err != nil {
		return nil, err
	}

	if st.IsMemory() {
		return memstore.New(), nil
	}

	path, err := st.SQLitePath(defaultName)
	if err != nil {
		return nil, err
	}

	return openSQLite(path, id.String())
}

// SQLitePath returns the SQLite database of a sqlite3 store, creating its
// directory.
func (st Store) SQLitePath(defaultName string) (string, error) {
	path := st.Path
	if path == "" {
		path = filepath.Join("$HOME", ".fix", defaultName)
	}
	path = os.ExpandEnv(path)

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}

	return path, nil
}

func (t *TLS) config() (*tls.Config, error) {
	cfg := &tls.Config{
		ServerName:         t.ServerName,
		InsecureSkipVerify: t.InsecureSkipVerify,
	}

	if t.CertFile != "" || t.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(os.ExpandEnv(t.CertFile), os.ExpandEnv(t.KeyFile))
		if err != nil {
			return nil, fmt.Errorf("tls: %w", err)
		}
		cfg.Certificates = []tls.Certificate{cert}
	}

	if t.CAFile != "" {
		pem, err := os.ReadFile(os.ExpandEnv(t.CAFile))
		if err != nil {
			return nil, fmt.Errorf("tls: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("tls: no certificate in %s", t.CAFile)
		}
		cfg.RootCAs = pool
		cfg.ClientCAs = pool
	}

	if t.ClientAuth {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return cfg, nil
}

// Dialer returns the dialer connecting to the initiator's acceptor.
func (i *Initiator) Dialer() (*tcp.Dialer, error) {
	if len(i.Addresses) == 0 {
		return nil, fmt.Errorf("%w: initiator %s has no addresses", errors.Config, i.Name)
	}

	d := &tcp.Dialer{
		Addresses: i.Addresses,
		Net:       &net.Dialer{Timeout: Timeout(i.DialTimeout)},
	}

	if i.TLS != nil {
		cfg, err := i.TLS.config()
		if err != nil {
			return nil, err
		}
		d.TLS = cfg
	}

	return d, nil
}

// DialTimeoutOrDefault returns how long a connection attempt may take.
func (i *Initiator) DialTimeoutOrDefault() time.Duration {
	return Timeout(i.DialTimeout)
}

// Listen listens on the acceptor's address.
func (a *Acceptor) Listen() (*tcp.Listener, error) {
	if a.Address == "" {
		return nil, fmt.Errorf("%w: acceptor %s has no address", errors.Config, a.Name)
	}

	var cfg *tls.Config
	if a.TLS != nil {
		var err error
		if cfg, err = a.TLS.config(); err != nil {
			return nil, err
		}
		if len(cfg.Certificates) == 0 {
			return nil, fmt.Errorf("%w: acceptor %s: tls requires certFile and keyFile", errors.Config, a.Name)
		}
	}

	return tcp.Listen(a.Address, cfg)
}

// Port returns the port of a "host:port" address, for display.
func Port(addr string) int {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return 0
	}
	p, _ := strconv.Atoi(port)
	return p
}
