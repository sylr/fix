package sendorder

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"sylr.dev/fix/pkg/dict"
)

func prettyTags[T any](input map[string]T) []string {
	output := make([]string, 0, len(input))
	for k := range input {
		output = append(output, strings.ToLower(k))
	}
	sort.Strings(output)
	return output
}

func completeOrderSide(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return prettyTags(dict.OrderSidesReversed), cobra.ShellCompDirectiveNoFileComp
}

func completeOrderType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return prettyTags(dict.OrderTypesReversed), cobra.ShellCompDirectiveNoFileComp
}

func completeOrderTimeInForce(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return prettyTags(dict.OrderTimeInForcesReversed), cobra.ShellCompDirectiveNoFileComp
}
