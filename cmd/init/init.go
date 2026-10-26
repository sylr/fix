package initcmd

import (
	"github.com/spf13/cobra"

	initconfig "sylr.dev/fix/v2/cmd/init/config"
	initdatabase "sylr.dev/fix/v2/cmd/init/database"
)

// InitCmd groups the commands creating fix's configuration and databases.
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize fix configurations",
	Long:  "Initialize fix configurations.",
}

func init() {
	InitCmd.AddCommand(initconfig.InitConfigCmd)
	InitCmd.AddCommand(initdatabase.InitDatabaseCmd)
}
