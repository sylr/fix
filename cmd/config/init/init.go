package configinit

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/errors"
)

// ConfigCmd represents the buy command
var ConfigInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize fix configuration",
	Long:  "Initialize fix configuration.",
	RunE:  Execute,
}

//go:embed templates/*
var templates embed.FS

func Execute(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()

	_, err := os.Stat(options.Config)
	if err == nil {
		return fmt.Errorf("%w: %s", errors.ConfigAlreadyExists, options.Config)
	}

	config, err := os.OpenFile(options.Config, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ConfigCanNotBeCreated, options.Config)
	}

	defer config.Close()

	configDir := path.Dir(options.Config)
	data := struct{ ConfigDir string }{
		ConfigDir: configDir,
	}

	t, err := template.ParseFS(templates, "templates/**")
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(config, "config", &data)

	// Dictionnaries
	dicts, err := fs.Glob(templates, "templates/FIX*.xml")
	fmt.Println(dicts)
	if err != nil {
		return err
	}

	for _, dict := range dicts {
		b, err := templates.ReadFile(dict)
		if err != nil {
			return err
		}

		err = os.WriteFile(strings.Join([]string{configDir, path.Base(dict)}, string(os.PathSeparator)), b, 0600)
		if err != nil {
			return err
		}
	}

	return err
}
