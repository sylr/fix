package config

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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

func ReadYAMLNoAge(path string) (*fixConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	in := bytes.NewBuffer(file)
	fix := fixConfig{}
	decoder := yaml.NewDecoder(in)
	decoder.KnownFields(true)
	err = decoder.Decode(&fix)

	if err != nil {
		return nil, err
	}

	return &fix, nil
}

func ReadYAML(path string, interactive bool) (*fixConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	in := bytes.NewBuffer(file)
	fix := fixConfig{}

	ids := GetAgeIdentities(interactive)
	w := yage.Wrapper{
		Value:      &fix,
		Identities: ids,
	}
	decoder := yaml.NewDecoder(in)
	decoder.KnownFields(true)
	err = decoder.Decode(&w)

	if err != nil {
		return nil, err
	}

	return &fix, nil
}

func GetAgeIdentities(interactive bool) []age.Identity {
	var ids []age.Identity
	var paths []string

	var globs = []string{
		path.Join(os.Getenv("HOME"), ".ssh", "id_rsa") + "*",
		path.Join(os.Getenv("HOME"), ".ssh", "id_ed25519") + "*",
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
		content, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}

		sshids, err := parseSSHIdentity(path, content, interactive)
		if err != nil {
			// If the key is explicitly requested, this error will be caught
			// below, otherwise ignore it silently.
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

	// Only cache keys in interactive mode because we probably don't want to keep
	// in memory unlocked ssh private keys.
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
	contents, err := ioutil.ReadAll(f)
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
