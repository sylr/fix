package utils

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ValidateRequiredFlags is a public version of cobra's Command.validateRequiredFlags()
func ValidateRequiredFlags(c *cobra.Command) error {
	if c.DisableFlagParsing {
		return nil
	}

	flags := c.Flags()
	missingFlagNames := []string{}
	flags.VisitAll(func(flag *pflag.Flag) {
		requiredAnnotation, found := flag.Annotations[cobra.BashCompOneRequiredFlag]
		if !found {
			return
		}
		if (requiredAnnotation[0] == "true") && !flag.Changed {
			missingFlagNames = append(missingFlagNames, flag.Name)
		}
	})

	if len(missingFlagNames) > 0 {
		return fmt.Errorf(`required flag(s) "%s" not set`, strings.Join(missingFlagNames, `", "`))
	}
	return nil
}
