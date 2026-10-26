package complete

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/config"
)

// Context completes with the names of the contexts of the configuration
// file.
func Context(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	conf, err := config.ReadYAMLNoAge(config.GetOptions().Config)
	if err != nil {
		return failed(err)
	}

	return names(conf.Contexts), cobra.ShellCompDirectiveNoFileComp
}

// Acceptor completes with the names of the acceptors of the configuration
// file.
func Acceptor(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	conf, err := config.ReadYAMLNoAge(config.GetOptions().Config)
	if err != nil {
		return failed(err)
	}

	return names(conf.Acceptors), cobra.ShellCompDirectiveNoFileComp
}

// Initiator completes with the names of the initiators of the
// configuration file.
func Initiator(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	conf, err := config.ReadYAMLNoAge(config.GetOptions().Config)
	if err != nil {
		return failed(err)
	}

	return names(conf.Initiators), cobra.ShellCompDirectiveNoFileComp
}

// Session completes with the names of the sessions of the configuration
// file.
func Session(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	conf, err := config.ReadYAMLNoAge(config.GetOptions().Config)
	if err != nil {
		return failed(err)
	}

	return names(conf.Sessions), cobra.ShellCompDirectiveNoFileComp
}

// failed reports err on stderr, where the shell completion shows it.
func failed(err error) ([]string, cobra.ShellCompDirective) {
	cobra.CompErrorln(err.Error())
	return nil, cobra.ShellCompDirectiveError
}

// names returns the names of items, skipping the null ones.
func names[T any, P interface {
	*T
	config.HasName
}](items []P) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		if item != nil {
			names = append(names, item.GetName())
		}
	}

	return names
}
