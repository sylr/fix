package complete

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/utils"
)

func SecurityListRequestType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.SecurityListRequestTypes), cobra.ShellCompDirectiveNoFileComp
}
