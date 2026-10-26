package complete

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/config"
)

func TestContext(t *testing.T) {
	options := config.GetOptions()
	t.Cleanup(func() { options.Config = "" })

	options.Config = filepath.Join(t.TempDir(), "none")
	if names, directive := Context(nil, nil, ""); names != nil || directive != cobra.ShellCompDirectiveError {
		t.Errorf("Context(no config) = %v, %v, want an error", names, directive)
	}

	options.Config = filepath.Join(t.TempDir(), "config")
	conf := "contexts:\n- name: a\n- null\n- name: b\nsessions:\n- name: s\n"
	if err := os.WriteFile(options.Config, []byte(conf), 0o600); err != nil {
		t.Fatal(err)
	}

	if names, directive := Context(nil, nil, ""); !slices.Equal(names, []string{"a", "b"}) || directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("Context() = %v, %v, want [a b]", names, directive)
	}
	if names, _ := Session(nil, nil, ""); !slices.Equal(names, []string{"s"}) {
		t.Errorf("Session() = %v, want [s]", names)
	}
	if names, _ := Initiator(nil, nil, ""); len(names) != 0 {
		t.Errorf("Initiator() = %v, want none", names)
	}
}
