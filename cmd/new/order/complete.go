package neworder

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/utils"
)

func completeOrderSide(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.OrderSidesReversed), cobra.ShellCompDirectiveNoFileComp
}

func completeOrderType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.OrderTypesReversed), cobra.ShellCompDirectiveNoFileComp
}

func completeOrderTimeInForce(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return utils.PrettyOptionValues(dict.OrderTimeInForcesReversed), cobra.ShellCompDirectiveNoFileComp
}
