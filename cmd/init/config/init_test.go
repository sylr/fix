package init_config

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"sylr.dev/fix/v2/config"
)

func TestTemplate(t *testing.T) {
	tpl, err := template.ParseFS(templates, "templates/config")
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "config")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := tpl.ExecuteTemplate(file, "config", nil); err != nil {
		t.Fatal(err)
	}
	file.Close()

	conf, err := config.ReadYAMLNoAge(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := conf.Validate(); err != nil {
		t.Fatal(err)
	}

	if len(conf.Contexts) != 2 || len(conf.Sessions) != 2 || len(conf.Initiators) != 1 || len(conf.Acceptors) != 1 {
		t.Fatalf("unexpected config: %+v", conf)
	}

	for _, s := range conf.Sessions {
		if !s.IsSupported() {
			t.Errorf("session %s is not a FIX version fix supports", s.Name)
		}
	}
}
