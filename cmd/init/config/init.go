package init_config

import (
	"compress/bzip2"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/errors"
)

// ConfigCmd represents the buy command
var InitConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Initialize fix configuration",
	Long:  "Initialize fix configuration.",
	RunE:  Execute,
}

//go:embed templates/*
var templates embed.FS

func Execute(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()

	configDir := path.Dir(options.Config)
	err := os.MkdirAll(configDir, 0700)
	if err != nil {
		return err
	}

	// Dictionnaries
	dicts, err := fs.Glob(templates, "templates/FIX*.xml.bz2")
	if err != nil {
		return err
	}

	for _, dict := range dicts {
		bzf, err := templates.Open(dict)
		if err != nil {
			return err
		}
		defer bzf.Close()

		bunzf := bzip2.NewReader(bzf)
		b, err := io.ReadAll(bunzf)
		if err != nil {
			return err
		}

		basenamebz2 := path.Base(dict)
		filename := strings.TrimSuffix(basenamebz2, filepath.Ext(basenamebz2))
		err = os.WriteFile(strings.Join([]string{configDir, filename}, string(os.PathSeparator)), b, 0600)
		if err != nil {
			return err
		}
	}

	// Config
	_, err = os.Stat(options.Config)
	if err == nil {
		return fmt.Errorf("%w: %s", errors.ConfigAlreadyExists, options.Config)
	}

	config, err := os.OpenFile(options.Config, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ConfigCanNotBeCreated, options.Config)
	}
	defer config.Close()

	data := struct{ ConfigDir string }{
		ConfigDir: configDir,
	}

	t, err := template.ParseFS(templates, "templates/config")
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(config, "config", &data)
	if err != nil {
		return err
	}

	return err
}
