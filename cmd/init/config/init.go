package initconfig

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
)

// InitConfigCmd writes the default configuration.
var InitConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Initialize fix configuration",
	Long:  "Initialize fix configuration.",
	RunE:  Execute,
}

//go:embed templates/*
var templates embed.FS

// templateData is what the configuration template is executed with.
type templateData struct {
	// Store is the store type of the initiators and acceptors: sqlite3 if
	// fix is built with it.
	Store string
}

func newTemplateData() templateData {
	if config.SQLiteStore {
		return templateData{Store: "sqlite3"}
	}

	return templateData{Store: "memory"}
}

func Execute(cmd *cobra.Command, args []string) error {
	path := config.GetOptions().Config

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("%w: %s: %w", errors.ErrConfigCanNotBeCreated, path, err)
	}

	t, err := template.ParseFS(templates, "templates/config")
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("%w: %s", errors.ErrConfigAlreadyExists, path)
	} else if err != nil {
		return fmt.Errorf("%w: %s: %w", errors.ErrConfigCanNotBeCreated, path, err)
	}

	err = t.ExecuteTemplate(file, "config", newTemplateData())
	if cerr := file.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("%w: %s: %w", errors.ErrConfigCanNotBeCreated, path, err)
	}

	fmt.Printf("Configuration written to %s.\n", path)

	return nil
}
