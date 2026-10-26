package config

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	yage "sylr.dev/yaml/age/v3"
	yaml "sylr.dev/yaml/v3"
)

type ageCache map[[sha256.Size]byte]age.Identity

var (
	sshAgeIdentitiesCache      ageCache     = make(ageCache)
	sshAgeIdentitiesCacheMutex sync.RWMutex = sync.RWMutex{}
)

// ReadYAMLNoAge reads the configuration file path, leaving the values
// encrypted with age as they are.
func ReadYAMLNoAge(path string) (*fixConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fix := fixConfig{}
	decoder := yaml.NewDecoder(bytes.NewReader(file))
	decoder.KnownFields(true)
	if err := decoder.Decode(&fix); err != nil {
		return nil, decodeError(err)
	}

	return &fix, nil
}

// ReadYAML reads the configuration file path, decrypting the values
// encrypted with age for the SSH keys in ~/.ssh.
func ReadYAML(path string, interactive bool) (*fixConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// yage.Wrapper decodes the decrypted document without KnownFields.
	if err := checkKeys(file); err != nil {
		return nil, err
	}

	fix := fixConfig{}
	w := yage.Wrapper{
		Value:      &fix,
		Identities: GetAgeIdentities(interactive),
	}
	if err := yaml.NewDecoder(bytes.NewReader(file)).Decode(&w); err != nil {
		return nil, err
	}

	return &fix, nil
}

// checkKeys returns the errors of a strict decoding of file, the values
// still encrypted, but those of values that do not fit their field, which
// encrypted ones may not until decrypted.
func checkKeys(file []byte) error {
	decoder := yaml.NewDecoder(bytes.NewReader(file))
	decoder.KnownFields(true)

	var typeErr *yaml.TypeError
	err := decoder.Decode(&fixConfig{})
	if !errors.As(err, &typeErr) {
		return decodeError(err)
	}

	var errs []string
	for _, e := range typeErr.Errors {
		if !strings.Contains(e, ": cannot unmarshal ") {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}

	return decodeError(&yaml.TypeError{Errors: errs})
}

// GetAgeIdentities returns the age identities of the SSH keys in ~/.ssh,
// those protected by a passphrase only if interactive.
func GetAgeIdentities(interactive bool) []age.Identity {
	var ids []age.Identity
	var paths []string

	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var globs = []string{
		filepath.Join(home, ".ssh", "id_rsa") + "*",
		filepath.Join(home, ".ssh", "id_ed25519") + "*",
	}

	for _, glob := range globs {
		files, err := filepath.Glob(glob)
		if err == nil {
			for _, file := range files {
				if !strings.HasSuffix(file, ".pub") {
					paths = append(paths, file)
				}
			}
		}
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		sshids, err := parseSSHIdentity(path, content, interactive)
		if err != nil {
			// Keys that cannot be used, as those protected by a passphrase
			// when not interactive, are skipped.
			continue
		}
		ids = append(ids, sshids...)
	}

	return ids
}

func parseSSHIdentity(name string, pemBytes []byte, interactive bool) ([]age.Identity, error) {
	if id := getSSHAgeIdentityFromCache(pemBytes); id != nil {
		return []age.Identity{id}, nil
	}

	id, err := agessh.ParseIdentity(pemBytes)
	if sshErr, ok := err.(*ssh.PassphraseMissingError); ok {
		if !interactive {
			return nil, fmt.Errorf("key %q password protected but we are in non interactive mode", name)
		}
		pubKey := sshErr.PublicKey
		if pubKey == nil {
			pubKey, err = readPubFile(name)
			if err != nil {
				return nil, err
			}
		}
		passphrasePrompt := func() ([]byte, error) {
			fmt.Fprintf(os.Stderr, "Enter passphrase for %q: ", name)
			pass, err := readPassphrase()
			if err != nil {
				return nil, fmt.Errorf("could not read passphrase for %q: %v", name, err)
			}
			return pass, nil
		}
		id, err = agessh.NewEncryptedSSHIdentity(pubKey, pemBytes, passphrasePrompt)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, fmt.Errorf("malformed SSH identity in %q: %v", name, err)
	}

	// An identity protected by a passphrase keeps the key once decrypted:
	// cached, the passphrase is asked once per run. Only interactive runs
	// ask for passphrases.
	if interactive {
		addSSHAgeIdentityToCache(pemBytes, id)
	}

	return []age.Identity{id}, nil
}

func readPubFile(name string) (ssh.PublicKey, error) {
	if name == "-" {
		return nil, fmt.Errorf(`failed to obtain public key for "-" SSH key
Use a file for which the corresponding ".pub" file exists, or convert the private key to a modern format with "ssh-keygen -p -m RFC4716"`)
	}
	f, err := os.Open(name + ".pub")
	if err != nil {
		return nil, fmt.Errorf(`failed to obtain public key for %q SSH key: %v
Ensure %q exists, or convert the private key %q to a modern format with "ssh-keygen -p -m RFC4716"`, name, err, name+".pub", name)
	}
	defer f.Close()
	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %v", name+".pub", err)
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(contents)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %q: %v", name+".pub", err)
	}
	return pubKey, nil
}

func readPassphrase() ([]byte, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		tty, err := os.Open("/dev/tty")
		if err != nil {
			return nil, fmt.Errorf("standard input is not available or not a terminal, and opening /dev/tty failed: %v", err)
		}
		defer tty.Close()
		fd = int(tty.Fd())
	}
	defer fmt.Fprintf(os.Stderr, "\n")
	p, err := term.ReadPassword(fd)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getSSHAgeIdentityFromCache(pemBytes []byte) age.Identity {
	sha256 := sha256.Sum256(pemBytes)

	sshAgeIdentitiesCacheMutex.RLock()
	defer sshAgeIdentitiesCacheMutex.RUnlock()

	if id, ok := sshAgeIdentitiesCache[sha256]; ok {
		return id
	}

	return nil
}

func addSSHAgeIdentityToCache(pemBytes []byte, id age.Identity) {
	sha256 := sha256.Sum256(pemBytes)

	sshAgeIdentitiesCacheMutex.Lock()
	defer sshAgeIdentitiesCacheMutex.Unlock()

	sshAgeIdentitiesCache[sha256] = id
}

// v1Key matches the unknown keys of fix v1's format: v2's start with a
// lowercase letter.
var v1Key = regexp.MustCompile(`field [A-Z]\w* not found in type config\.`)

// decodeError points to the README when err may come from a configuration
// in fix v1's format.
func decodeError(err error) error {
	if err != nil && v1Key.MatchString(err.Error()) {
		return fmt.Errorf("%w (the configuration format changed in fix v2, see https://github.com/sylr/fix#migrating-from-v1)", err)
	}

	return err
}
