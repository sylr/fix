package init_config

import (
	"embed"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
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

	t, err := template.ParseFS(templates, "templates/config")
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(config, "config", nil)
	if err != nil {
		return err
	}

	return err
}
