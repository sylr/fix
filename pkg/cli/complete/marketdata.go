package complete

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/utils"
)

func MDEntryTypes(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.MDEntryTypes), cobra.ShellCompDirectiveNoFileComp
}

func SubscriptionRequestTypes(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.SubscriptionRequestTypes), cobra.ShellCompDirectiveNoFileComp
}

func MDUpdateTypes(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.MDUpdateTypes), cobra.ShellCompDirectiveNoFileComp
}
