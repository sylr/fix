package init

import (
	"github.com/spf13/cobra"

	init_config "sylr.dev/fix/cmd/init/config"
	init_database "sylr.dev/fix/cmd/init/database"
)

// ConfigCmd represents the buy command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize fix configurations",
	Long:  "Initialize fix configurations.",
}

func init() {
	InitCmd.AddCommand(init_config.InitConfigCmd)
	InitCmd.AddCommand(init_database.InitDatabaseCmd)
}
