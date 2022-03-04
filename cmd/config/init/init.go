package configinit

import (
	"embed"
	"fmt"
	"io"
	"os"

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

//go:embed config.template
var configTemplate embed.FS

func Execute(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()

	if _, err := os.Stat(options.Config); err == nil {
		return fmt.Errorf("%w: %s", errors.ConfigAlreadyExists, options.Config)
	}

	fd, _ := configTemplate.Open("config.template")
	b, _ := io.ReadAll(fd)
	err := os.WriteFile(options.Config, b, 0600)

	return err
}
