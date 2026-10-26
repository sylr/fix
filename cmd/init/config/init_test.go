package initconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
)

func TestExecute(t *testing.T) {
	options := config.GetOptions()
	old := options.Config
	t.Cleanup(func() { options.Config = old })
	options.Config = filepath.Join(t.TempDir(), "sub", "config")

	if err := Execute(nil, nil); err != nil {
		t.Fatal(err)
	}

	conf, err := config.ReadYAMLNoAge(options.Config)
	if err != nil {
		t.Fatal(err)
	}

	if err := conf.Validate(); err != nil {
		t.Fatal(err)
	}

	if len(conf.Contexts) != 2 || len(conf.Sessions) != 2 || len(conf.Initiators) != 1 || len(conf.Acceptors) != 1 {
		t.Fatalf("unexpected config: %+v", conf)
	}

	if store := conf.Initiators[0].Store; store.IsMemory() == config.SQLiteStore {
		t.Errorf("store %q, built with sqlite3: %v", store.Type, config.SQLiteStore)
	}

	if err := os.WriteFile(options.Config, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Execute(nil, nil); !errors.Is(err, errors.ErrConfigAlreadyExists) {
		t.Errorf("existing configuration: %v", err)
	}
	if b, _ := os.ReadFile(options.Config); string(b) != "keep" {
		t.Errorf("existing configuration overwritten: %q", b)
	}
}

// TestREADME checks that the README's example is the configuration written.
func TestREADME(t *testing.T) {
	readme, err := os.ReadFile("../../../README.md")
	if err != nil {
		t.Fatal(err)
	}
	tpl, err := templates.ReadFile("templates/config")
	if err != nil {
		t.Fatal(err)
	}

	example := strings.ReplaceAll(string(tpl), "{{ .Store }}", "sqlite3")
	if !strings.Contains(string(readme), "```yaml\n"+example+"```\n") {
		t.Error("the README's example is not templates/config")
	}
}
