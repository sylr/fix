package database

import (
	"fmt"

	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
)

func ValidateOptions(cmd *cobra.Command, args []string) error {
	options := config.GetOptions()
	conf, err := config.ReadYAML(options.Config, options.Interactive)

	if err != nil {
		return fmt.Errorf("unable to read configuration: %w", err)
	}

	err = conf.Validate()
	if err != nil {
		return err
	}

	// Retrieve the global config pointer
	fixConfig := config.GetConfig()

	// Set the config retrieved in the config file into the global config pointer
	*fixConfig = *conf

	return nil
}
