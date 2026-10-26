package config

import (
	"cmp"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"maps"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/sylr/fix2/fix/dictionary"
	"github.com/sylr/fix2/fix/dictionary/xmldict"
	"github.com/sylr/fix2/fix/session"
	"github.com/sylr/fix2/fix/session/schedule"
	"github.com/sylr/fix2/fix/session/store"
	"github.com/sylr/fix2/fix/session/store/memstore"
	"github.com/sylr/fix2/transport/tcp"

	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix43"
	"sylr.dev/fix/v2/fix44"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/fixt11"
	"sylr.dev/fix/v2/internal/fixmsg"
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

// validate checks the session's configuration: what Config.Validate in
// fix2 cannot, the names and formats fix2 has no notion of, and what is
// better reported before connecting.
func (s *Session) validate() error {
	if s.SenderCompID == "" || s.TargetCompID == "" {
		return fmt.Errorf("senderCompID and targetCompID are required")
	}

	switch _, fix4 := fix4Dictionaries[s.BeginString]; {
	case s.BeginString == "":
		return fmt.Errorf("beginString is required")
	case s.BeginString == session.FIXT11:
		if !knownApplVerID(s.DefaultApplVerID) {
			return fmt.Errorf("FIXT.1.1 requires defaultApplVerID, a FIX version or its ApplVerID, not %q (known: %s)", s.DefaultApplVerID, names(ApplVerIDs))
		}
	case fix4:
		if s.DefaultApplVerID != "" {
			return fmt.Errorf("defaultApplVerID is for FIXT.1.1 sessions, not %s", s.BeginString)
		}
		if s.Dictionaries.Session != "" && s.Dictionaries.App != "" {
			return fmt.Errorf("%s has one dictionary: set dictionaries.app or dictionaries.session, not both", s.BeginString)
		}
	}

	if err := s.CheckSupported(); err != nil {
		return err
	}

	if s.HeartBtInt < 0 || s.HeartBtInt%time.Second != 0 {
		return fmt.Errorf("heartBtInt must be whole seconds >= 0, not %v", s.HeartBtInt)
	}

	if s.LogonTimeout < 0 || s.LogoutTimeout < 0 {
		return fmt.Errorf("logonTimeout and logoutTimeout must be >= 0")
	}

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

// knownApplVerID reports whether v is a FIX version ApplVerIDs names, or
// its ApplVerID.
func knownApplVerID(v string) bool {
	for name, id := range ApplVerIDs {
		if v == name || v == id {
			return true
		}
	}

	return false
}

// ApplVerID returns the ApplVerID (1128) value of DefaultApplVerID.
func (s *Session) ApplVerID() string {
	if v, ok := ApplVerIDs[s.DefaultApplVerID]; ok {
		return v
	}

	return s.DefaultApplVerID
}

// fix4Dictionaries are the built-in dictionaries of the FIX 4.x versions fix
// supports, by BeginString.
var fix4Dictionaries = map[string]func() *dictionary.Dict{
	fix43.BeginString: fix43.Dictionary,
	fix44.BeginString: fix44.Dictionary,
}

// Version returns the FIX version of the session's application messages:
// its BeginString, or for FIXT.1.1 its DefaultApplVerID.
func (s *Session) Version() string {
	if s.BeginString != session.FIXT11 {
		return s.BeginString
	}

	for name, id := range ApplVerIDs {
		if id == s.ApplVerID() {
			return name
		}
	}

	return s.DefaultApplVerID
}

// IsSupported reports whether fix sends and reads the session's application
// messages: FIX 5.0 SP2 over FIXT.1.1, FIX 4.3 and FIX 4.4. Commands build
// them with package fix50sp2, whose tags and values FIX 4.3 and 4.4 share.
func (s *Session) IsSupported() bool {
	if s.BeginString == session.FIXT11 {
		return s.ApplVerID() == fix50sp2.Version.ApplVerID
	}

	_, ok := fix4Dictionaries[s.BeginString]
	return ok
}

// CheckSupported returns an error naming the session's FIX version unless
// IsSupported.
func (s *Session) CheckSupported() error {
	if !s.IsSupported() {
		return fmt.Errorf("%w: %s (supported: FIX.5.0SP2 over FIXT.1.1, FIX.4.3, FIX.4.4)", errors.ErrFixVersionNotImplemented, s.Version())
	}

	return nil
}

// Defines reports whether the session's dictionary defines tag in messages
// of type msgType, as FIX 4.x does not QuoteMsgID (1166) in a QuoteCancel;
// true without a dictionary, and with fix's built-in FIX 5.0 SP2 one, which
// lacks what fix does not use (see trimmedDictionary).
func (s *Session) Defines(msgType string, tag int) bool {
	_, d, err := s.FIXDictionaries()
	if err != nil || d == nil || s.trimmedDictionary() {
		return true
	}

	return d.Defines(msgType, tag)
}

// trimmedDictionary reports whether the session uses fix's built-in FIX 5.0
// SP2 dictionary, which leaves out the components fix does not use, whose
// tags counterparties may still send.
func (s *Session) trimmedDictionary() bool {
	return s.BeginString == session.FIXT11 && s.Dictionaries.App == "" && s.ApplVerID() == fix50sp2.Version.ApplVerID
}

// Requires reports whether the session's dictionary requires tag in
// messages of type msgType, as FIX 4.3 does HandlInst (21) in a
// NewOrderSingle; false without a dictionary.
func (s *Session) Requires(msgType string, tag int) bool {
	_, d, err := s.FIXDictionaries()
	if err != nil || d == nil || d.Messages[msgType] == nil {
		return false
	}

	for _, f := range d.Messages[msgType].Fields {
		if f.Tag == tag {
			return f.Required
		}
	}

	return false
}

func loadDictionary(path string) (*dictionary.Dict, error) {
	d, err := xmldict.LoadFile(expandPath(path))
	if err != nil {
		return nil, fmt.Errorf("loading dictionary %s: %w", path, err)
	}

	return d, nil
}

// loadedDictionaries holds what FIXDictionaries loaded for a session.
type loadedDictionaries struct {
	once    sync.Once
	dicts   session.Dictionaries
	display *dictionary.Dict
	err     error
}

// loadedDictionariesMu guards Session.loaded.
var loadedDictionariesMu sync.Mutex

// FIXDictionaries returns the dictionaries the session validates messages
// with, and the one describing every field of its messages, header and
// trailer included, for printing them. fix's own dictionaries are used
// unless Dictionaries name others. They are loaded once, on first use:
// the session's configuration must not change after.
func (s *Session) FIXDictionaries() (session.Dictionaries, *dictionary.Dict, error) {
	loadedDictionariesMu.Lock()
	if s.loaded == nil {
		s.loaded = &loadedDictionaries{}
	}
	l := s.loaded
	loadedDictionariesMu.Unlock()

	l.once.Do(func() { l.dicts, l.display, l.err = s.loadDictionaries() })

	return l.dicts, l.display, l.err
}

// loadDictionaries loads the dictionaries FIXDictionaries returns.
func (s *Session) loadDictionaries() (session.Dictionaries, *dictionary.Dict, error) {
	var dicts session.Dictionaries
	var err error

	if s.BeginString != session.FIXT11 {
		path := cmp.Or(s.Dictionaries.App, s.Dictionaries.Session)

		switch builtin, ok := fix4Dictionaries[s.BeginString]; {
		case path != "":
			dicts.Session, err = loadDictionary(path)
		case ok:
			dicts.Session = builtin()
		}

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

// Checks returns the dictionary validations the session runs on the
// messages it receives.
func (s *Session) Checks() dictionary.Checks {
	return dictionary.CheckAll &^ s.skipChecks()
}

// Check refuses m if the session's dictionary does not define it or what it
// holds, or it lacks a field the dictionary requires, as fixmsg.Check
// tells: rather than have the counterparty reject it.
func (s *Session) Check(m session.Outgoing) error {
	_, d, err := s.FIXDictionaries()
	if err != nil || d == nil {
		return err
	}

	msg, err := fixmsg.Message(m, s.BeginString, d)
	if err != nil {
		return err
	}

	return fixmsg.Check(msg, d, s.Checks())
}

// skipChecks returns the validations the session skips.
func (s *Session) skipChecks() dictionary.Checks {
	var checks dictionary.Checks

	if s.trimmedDictionary() {
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

// Validate checks the store's type, and that fix is built with it.
func (st Store) Validate() error {
	switch {
	case st.IsMemory():
	case st.Type != "sqlite3":
		return fmt.Errorf("unknown store type %q (known: memory, sqlite3)", st.Type)
	case !SQLiteStore:
		return fmt.Errorf("the sqlite3 store requires fix built with cgo")
	}

	return nil
}

// Open opens the store of the session id, in the SQLite database
// ~/.fix/<defaultName> unless Path names another.
func (st Store) Open(defaultName string, id session.SessionID) (store.Store, error) {
	if err := st.Validate(); err != nil {
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
	path := expandPath(st.Path)
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, ".fix", defaultName)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}

	return path, nil
}

// expandPath expands the environment variables in path, then a leading "~/"
// to the user's home directory, if known.
func expandPath(path string) string {
	path = os.ExpandEnv(path)

	rest, ok := strings.CutPrefix(filepath.ToSlash(path), "~/")
	if !ok {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return filepath.Join(home, filepath.FromSlash(rest))
}

// validate checks the TLS configuration of an acceptor, or of an initiator.
func (t *TLS) validate(acceptor bool) error {
	if (t.CertFile == "") != (t.KeyFile == "") {
		return fmt.Errorf("tls: certFile and keyFile go together")
	}

	if !acceptor {
		if t.ClientAuth {
			return fmt.Errorf("tls: clientAuth is for acceptors")
		}
		return nil
	}

	switch {
	case t.CertFile == "":
		return fmt.Errorf("tls: an acceptor requires certFile and keyFile")
	case t.ServerName != "" || t.InsecureSkipVerify:
		return fmt.Errorf("tls: serverName and insecureSkipVerify are for initiators")
	case t.ClientAuth && t.CAFile == "":
		return fmt.Errorf("tls: clientAuth requires caFile, the authorities client certificates are verified against")
	case !t.ClientAuth && t.CAFile != "":
		return fmt.Errorf("tls: an acceptor's caFile verifies client certificates, which requires clientAuth")
	}

	return nil
}

// config returns the TLS configuration of an acceptor, or of an initiator.
func (t *TLS) config(acceptor bool) (*tls.Config, error) {
	if err := t.validate(acceptor); err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		ServerName:         t.ServerName,
		InsecureSkipVerify: t.InsecureSkipVerify,
	}

	if t.CertFile != "" {
		cert, err := tls.LoadX509KeyPair(expandPath(t.CertFile), expandPath(t.KeyFile))
		if err != nil {
			return nil, fmt.Errorf("tls: %w", err)
		}
		cfg.Certificates = []tls.Certificate{cert}
	}

	if t.CAFile != "" {
		pem, err := os.ReadFile(expandPath(t.CAFile))
		if err != nil {
			return nil, fmt.Errorf("tls: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("tls: no certificate in %s", t.CAFile)
		}
		if acceptor {
			cfg.ClientCAs = pool
		} else {
			cfg.RootCAs = pool
		}
	}

	if t.ClientAuth {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return cfg, nil
}

// validateAddress checks that addr is "host:port".
func validateAddress(addr string) error {
	if _, port, err := net.SplitHostPort(addr); err != nil || port == "" {
		return fmt.Errorf("address %q is not host:port", addr)
	}

	return nil
}

// validate checks the initiator's configuration.
func (i *Initiator) validate() error {
	if len(i.Addresses) == 0 {
		return fmt.Errorf("addresses are required")
	}

	for _, addr := range i.Addresses {
		if err := validateAddress(addr); err != nil {
			return err
		}
	}

	if i.DialTimeout < 0 || i.ReconnectMin < 0 || i.ReconnectMax < 0 {
		return fmt.Errorf("dialTimeout, reconnectMin and reconnectMax must be >= 0")
	}

	// With the defaults of fix2's initiator.
	if lo, hi := cmp.Or(i.ReconnectMin, time.Second), cmp.Or(i.ReconnectMax, 30*time.Second); lo > hi {
		return fmt.Errorf("reconnectMin %v is above reconnectMax %v", lo, hi)
	}

	if i.TLS != nil {
		if err := i.TLS.validate(false); err != nil {
			return err
		}
	}

	return i.Store.Validate()
}

// Dialer returns the dialer connecting to the initiator's acceptor.
func (i *Initiator) Dialer() (*tcp.Dialer, error) {
	if len(i.Addresses) == 0 {
		return nil, fmt.Errorf("%w: initiator %s has no addresses", errors.ErrConfig, i.Name)
	}

	d := &tcp.Dialer{
		Addresses: i.Addresses,
		Net:       &net.Dialer{Timeout: Timeout(i.DialTimeout)},
	}

	if i.TLS != nil {
		cfg, err := i.TLS.config(false)
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

// validate checks the acceptor's configuration.
func (a *Acceptor) validate() error {
	if a.Address == "" {
		return fmt.Errorf("address is required")
	}

	if err := validateAddress(a.Address); err != nil {
		return err
	}

	if a.TLS != nil {
		if err := a.TLS.validate(true); err != nil {
			return err
		}
	}

	return a.Store.Validate()
}

// Listen listens on the acceptor's address.
func (a *Acceptor) Listen() (*tcp.Listener, error) {
	if a.Address == "" {
		return nil, fmt.Errorf("%w: acceptor %s has no address", errors.ErrConfig, a.Name)
	}

	var cfg *tls.Config
	if a.TLS != nil {
		var err error
		if cfg, err = a.TLS.config(true); err != nil {
			return nil, err
		}
	}

	return tcp.Listen(a.Address, cfg)
}
