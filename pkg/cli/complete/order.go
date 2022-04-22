package complete

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/utils"
)

func OrderSide(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.OrderSides), cobra.ShellCompDirectiveNoFileComp
}

func OrderType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.OrderTypes), cobra.ShellCompDirectiveNoFileComp
}

func OrderTimeInForce(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.OrderTimeInForces), cobra.ShellCompDirectiveNoFileComp
}

func OrderPartyIDSource(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.PartyIDSources), cobra.ShellCompDirectiveNoFileComp
}

func OrderPartySubIDTypes(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.PartySubIDTypes), cobra.ShellCompDirectiveNoFileComp
}

func OrderPartyIDRole(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.PartyRoles), cobra.ShellCompDirectiveNoFileComp
}

func OrderOriginationRole(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.OrderOriginations), cobra.ShellCompDirectiveNoFileComp
}
