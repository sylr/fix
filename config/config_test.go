package config

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"filippo.io/age/armor"
	"golang.org/x/crypto/ssh"
)

// validConfig returns a configuration Validate accepts, for tests to break.
func validConfig() fixConfig {
	return fixConfig{
		Contexts: []*Context{
			{Name: "client", Initiator: "i", Sessions: []string{"client"}},
			{Name: "server", Acceptor: "a", Sessions: []string{"server"}},
		},
		Initiators: []*Initiator{{Name: "i", Addresses: []string{"127.0.0.1:5005"}}},
		Acceptors:  []*Acceptor{{Name: "a", Address: ":5005"}},
		Sessions: []*Session{
			{Name: "client", BeginString: "FIXT.1.1", DefaultApplVerID: "FIX.5.0SP2", SenderCompID: "A", TargetCompID: "B"},
			{Name: "server", BeginString: "FIX.4.4", SenderCompID: "B", TargetCompID: "A", HeartBtInt: 30 * time.Second},
		},
	}
}

func TestConfigValidate(t *testing.T) {
	if conf := validConfig(); conf.Validate() != nil {
		t.Fatalf("Validate() = %v", conf.Validate())
	}

	for name, breakIt := range map[string]func(*fixConfig){
		"nil context":             func(c *fixConfig) { c.Contexts = append(c.Contexts, nil) },
		"nil session":             func(c *fixConfig) { c.Sessions = append(c.Sessions, nil) },
		"nil initiator":           func(c *fixConfig) { c.Initiators = append(c.Initiators, nil) },
		"nil acceptor":            func(c *fixConfig) { c.Acceptors = append(c.Acceptors, nil) },
		"unnamed context":         func(c *fixConfig) { c.Contexts[0].Name = "" },
		"unnamed session":         func(c *fixConfig) { c.Sessions[0].Name = "" },
		"duplicate acceptor":      func(c *fixConfig) { c.Acceptors = append(c.Acceptors, &Acceptor{Name: "a", Address: ":1"}) },
		"no addresses":            func(c *fixConfig) { c.Initiators[0].Addresses = nil },
		"address without port":    func(c *fixConfig) { c.Initiators[0].Addresses = []string{"127.0.0.1"} },
		"address with empty port": func(c *fixConfig) { c.Initiators[0].Addresses = []string{"localhost:"} },
		"no acceptor address":     func(c *fixConfig) { c.Acceptors[0].Address = "" },
		"bad acceptor address":    func(c *fixConfig) { c.Acceptors[0].Address = "5005" },
		"negative dialTimeout":    func(c *fixConfig) { c.Initiators[0].DialTimeout = -time.Second },
		"negative reconnectMin":   func(c *fixConfig) { c.Initiators[0].ReconnectMin = -time.Second },
		"reconnectMin > max": func(c *fixConfig) {
			c.Initiators[0].ReconnectMin, c.Initiators[0].ReconnectMax = 10*time.Second, 5*time.Second
		},
		"reconnectMin > default max": func(c *fixConfig) { c.Initiators[0].ReconnectMin = time.Minute },
		"unknown store":              func(c *fixConfig) { c.Acceptors[0].Store.Type = "postgres" },
		"initiator clientAuth":       func(c *fixConfig) { c.Initiators[0].TLS = &TLS{ClientAuth: true} },
		"no senderCompID":            func(c *fixConfig) { c.Sessions[0].SenderCompID = "" },
		"no targetCompID":            func(c *fixConfig) { c.Sessions[1].TargetCompID = "" },
		"no beginString":             func(c *fixConfig) { c.Sessions[0].BeginString = "" },
		"unsupported beginString":    func(c *fixConfig) { c.Sessions[1].BeginString = "FIX.4.2" },
		"typo in beginString":        func(c *fixConfig) { c.Sessions[1].BeginString = "FIX44" },
		"no defaultApplVerID":        func(c *fixConfig) { c.Sessions[0].DefaultApplVerID = "" },
		"unknown defaultApplVerID":   func(c *fixConfig) { c.Sessions[0].DefaultApplVerID = "FIX.5.0.SP2" },
		"unsupported ApplVerID":      func(c *fixConfig) { c.Sessions[0].DefaultApplVerID = "FIX.5.0" },
		"FIX.4.4 defaultApplVerID":   func(c *fixConfig) { c.Sessions[1].DefaultApplVerID = "FIX.4.4" },
		"FIX.4.4 two dictionaries":   func(c *fixConfig) { c.Sessions[1].Dictionaries = Dictionaries{Session: "a.xml", App: "b.xml"} },
		"heartBtInt fraction":        func(c *fixConfig) { c.Sessions[0].HeartBtInt = 1500 * time.Millisecond },
		"negative heartBtInt":        func(c *fixConfig) { c.Sessions[0].HeartBtInt = -time.Second },
		"negative logonTimeout":      func(c *fixConfig) { c.Sessions[0].LogonTimeout = -time.Second },
		"negative logoutTimeout":     func(c *fixConfig) { c.Sessions[1].LogoutTimeout = -time.Second },
	} {
		conf := validConfig()
		breakIt(&conf)
		if err := conf.Validate(); err == nil {
			t.Errorf("%s: no error", name)
		}
	}

	for name, fixIt := range map[string]func(*fixConfig){
		"ApplVerID code":   func(c *fixConfig) { c.Sessions[0].DefaultApplVerID = "9" },
		"FIX.4.4 app dict": func(c *fixConfig) { c.Sessions[1].Dictionaries.App = "FIX44.xml" },
		"negative latency": func(c *fixConfig) { c.Sessions[0].MaxLatency = -1 },
		"reconnect in range": func(c *fixConfig) {
			c.Initiators[0].ReconnectMin, c.Initiators[0].ReconnectMax = time.Minute, time.Hour
		},
		"several addresses": func(c *fixConfig) { c.Initiators[0].Addresses = []string{"[::1]:5005", "example.com:fix"} },
	} {
		conf := validConfig()
		fixIt(&conf)
		if err := conf.Validate(); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

func TestContextGetInitiatorAcceptor(t *testing.T) {
	old := config
	t.Cleanup(func() { config = old })
	config = validConfig()

	if _, err := config.Contexts[1].GetInitiator(); err == nil {
		t.Error("GetInitiator() of an acceptor context: no error")
	}
	if _, err := config.Contexts[0].GetAcceptor(); err == nil {
		t.Error("GetAcceptor() of an initiator context: no error")
	}
	if i, err := config.Contexts[0].GetInitiator(); err != nil || i.Name != "i" {
		t.Errorf("GetInitiator() = %v, %v", i, err)
	}
}

// writeFile writes content to a file in a temporary directory.
func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

const testYAML = `
contexts:
- name: client
  initiator: i
  sessions: [client]
initiators:
- name: i
  addresses: [127.0.0.1:5005]
  %s
sessions:
- name: client
  beginString: FIXT.1.1
  defaultApplVerID: FIX.5.0SP2
  senderCompID: A
  targetCompID: B
`

func TestReadYAMLUnknownKeys(t *testing.T) {
	// No SSH key: ReadYAML must not read the user's.
	t.Setenv("HOME", t.TempDir())

	for name, tc := range map[string]struct {
		extra, want string
		v1          bool
	}{
		"tls typo":        {extra: "tsl: {caFile: ca.pem}", want: "field tsl not found"},
		"nested typo":     {extra: "tls: {serverName: x, insecureSkipVerfy: true}", want: "field insecureSkipVerfy not found"},
		"store typo":      {extra: "store: {typ: sqlite3}", want: "field typ not found"},
		"v1 key":          {extra: "SocketConnectHost: 127.0.0.1", want: "field SocketConnectHost not found", v1: true},
		"duplicate key":   {extra: "addresses: [127.0.0.1:5006]", want: "already"},
		"wrong type only": {extra: "dialTimeout: [5s]", want: "cannot unmarshal"},
	} {
		path := writeFile(t, "config", strings.Replace(testYAML, "%s", tc.extra, 1))

		for reader, read := range map[string]func(string) (*fixConfig, error){
			"ReadYAML":      func(p string) (*fixConfig, error) { return ReadYAML(p, false) },
			"ReadYAMLNoAge": ReadYAMLNoAge,
		} {
			_, err := read(path)
			switch {
			case err == nil:
				t.Errorf("%s: %s: no error", name, reader)
			case !strings.Contains(err.Error(), tc.want):
				t.Errorf("%s: %s: %v, want %q", name, reader, err, tc.want)
			case strings.Contains(err.Error(), "migrating-from-v1") != tc.v1:
				t.Errorf("%s: %s: %v: v1 hint should be %v", name, reader, err, tc.v1)
			}
		}
	}

	path := writeFile(t, "config", strings.Replace(testYAML, "%s", "dialTimeout: 3s", 1))
	conf, err := ReadYAML(path, false)
	if err != nil {
		t.Fatal(err)
	}
	if conf.Initiators[0].DialTimeout != 3*time.Second {
		t.Errorf("DialTimeout = %v", conf.Initiators[0].DialTimeout)
	}
}

// encrypt returns value encrypted with age for recipient, armored.
func encrypt(t *testing.T, recipient age.Recipient, value string) string {
	t.Helper()
	var buf bytes.Buffer
	a := armor.NewWriter(&buf)
	w, err := age.Encrypt(a, recipient)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte(value)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := a.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// indent indents s for a YAML literal block.
func indent(s string) string {
	return "    " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n    ")
}

func TestReadYAMLAge(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	block, err := ssh.MarshalPrivateKey(priv, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(home, ".ssh"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".ssh", "id_ed25519"), pem.EncodeToMemory(block), 0o600); err != nil {
		t.Fatal(err)
	}
	sshPub, err := ssh.NewPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	recipient, err := agessh.NewEd25519Recipient(sshPub)
	if err != nil {
		t.Fatal(err)
	}

	yaml := testYAML + "  password: !crypto/age |\n" + indent(encrypt(t, recipient, "secret")) + "\n" +
		"  heartBtInt: !crypto/age |\n" + indent(encrypt(t, recipient, "5s")) + "\n"

	conf, err := ReadYAML(writeFile(t, "config", strings.Replace(yaml, "%s", "", 1)), false)
	if err != nil {
		t.Fatal(err)
	}
	if s := conf.Sessions[0]; s.Password != "secret" || s.HeartBtInt != 5*time.Second {
		t.Errorf("Password = %q, HeartBtInt = %v", s.Password, s.HeartBtInt)
	}

	_, err = ReadYAML(writeFile(t, "config", strings.Replace(yaml, "%s", "tsl: {}", 1)), false)
	if err == nil || !strings.Contains(err.Error(), "field tsl not found") {
		t.Errorf("unknown key next to encrypted values: %v", err)
	}
}

// writeCA writes a self-signed CA certificate and its key, and returns their
// paths.
func writeCA(t *testing.T) (certFile, keyFile string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		DNSNames:              []string{"localhost"},
	}
	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	certFile = writeFile(t, "cert.pem", string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})))
	keyFile = writeFile(t, "key.pem", string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})))
	return certFile, keyFile
}

func TestTLS(t *testing.T) {
	for name, tc := range map[string]struct {
		tls      TLS
		acceptor bool
	}{
		"acceptor without cert":           {TLS{}, true},
		"cert without key":                {TLS{CertFile: "c"}, false},
		"acceptor clientAuth without CA":  {TLS{CertFile: "c", KeyFile: "k", ClientAuth: true}, true},
		"acceptor CA without clientAuth":  {TLS{CertFile: "c", KeyFile: "k", CAFile: "ca"}, true},
		"acceptor serverName":             {TLS{CertFile: "c", KeyFile: "k", ServerName: "x"}, true},
		"acceptor insecureSkipVerify":     {TLS{CertFile: "c", KeyFile: "k", InsecureSkipVerify: true}, true},
		"initiator clientAuth":            {TLS{ClientAuth: true}, false},
		"initiator clientAuth and caFile": {TLS{CAFile: "ca", ClientAuth: true}, false},
	} {
		if err := tc.tls.validate(tc.acceptor); err == nil {
			t.Errorf("%s: no error", name)
		}
	}

	cert, key := writeCA(t)

	acceptor, err := (&TLS{CertFile: cert, KeyFile: key, CAFile: cert, ClientAuth: true}).config(true)
	if err != nil {
		t.Fatal(err)
	}
	if acceptor.ClientAuth != tls.RequireAndVerifyClientCert || acceptor.ClientCAs == nil || acceptor.RootCAs != nil || len(acceptor.Certificates) != 1 {
		t.Errorf("acceptor TLS config = %+v", acceptor)
	}

	initiator, err := (&TLS{ServerName: "localhost", CAFile: cert}).config(false)
	if err != nil {
		t.Fatal(err)
	}
	if initiator.ClientAuth != tls.NoClientCert || initiator.RootCAs == nil || initiator.ClientCAs != nil || initiator.ServerName != "localhost" {
		t.Errorf("initiator TLS config = %+v", initiator)
	}

	if _, err := (&Acceptor{Name: "a", Address: "127.0.0.1:0", TLS: &TLS{CAFile: cert, ClientAuth: true}}).Listen(); err == nil {
		t.Error("Listen() without a certificate: no error")
	}
}

func TestExpandPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("FIX_TEST_DIR", "/var/fix")

	for in, want := range map[string]string{
		"~/.fix/x.db":           filepath.Join(home, ".fix", "x.db"),
		"$FIX_TEST_DIR/x.db":    "/var/fix/x.db",
		"/tmp/~/x.db":           "/tmp/~/x.db",
		"~user/x.db":            "~user/x.db",
		"relative/path/x.db":    "relative/path/x.db",
		"${FIX_TEST_DIR}/f.xml": "/var/fix/f.xml",
		"":                      "",
	} {
		if got := expandPath(in); got != want {
			t.Errorf("expandPath(%q) = %q, want %q", in, got, want)
		}
	}

	path, err := Store{Type: "sqlite3"}.SQLitePath("initiator.db")
	if err != nil || path != filepath.Join(home, ".fix", "initiator.db") {
		t.Errorf("SQLitePath() = %q, %v", path, err)
	}
}
