package complete

import (
	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
)

func Context(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	options := config.GetOptions()
	fixConfig := config.GetConfig()

	if conf, err := config.ReadYAMLNoAge(options.Config); err == nil {
		*fixConfig = *conf
	} else {
		return nil, cobra.ShellCompDirectiveError
	}

	contexts := make([]string, 0, len(fixConfig.Contexts))
	for _, context := range fixConfig.Contexts {
		contexts = append(contexts, context.Name)
	}

	return contexts, cobra.ShellCompDirectiveNoFileComp
}

func Acceptor(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	options := config.GetOptions()
	fixConfig := config.GetConfig()

	if conf, err := config.ReadYAMLNoAge(options.Config); err == nil {
		*fixConfig = *conf
	} else {
		return nil, cobra.ShellCompDirectiveError
	}

	acceptors := make([]string, 0, len(fixConfig.Acceptors))
	for _, acceptor := range fixConfig.Initiators {
		acceptors = append(acceptors, acceptor.Name)
	}

	return acceptors, cobra.ShellCompDirectiveNoFileComp
}

func Initiator(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	options := config.GetOptions()
	fixConfig := config.GetConfig()

	if conf, err := config.ReadYAMLNoAge(options.Config); err == nil {
		*fixConfig = *conf
	} else {
		return nil, cobra.ShellCompDirectiveError
	}

	initiators := make([]string, 0, len(fixConfig.Initiators))
	for _, initiator := range fixConfig.Initiators {
		initiators = append(initiators, initiator.Name)
	}

	return initiators, cobra.ShellCompDirectiveNoFileComp
}

func Session(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	options := config.GetOptions()
	fixConfig := config.GetConfig()

	if conf, err := config.ReadYAMLNoAge(options.Config); err == nil {
		*fixConfig = *conf
	} else {
		return nil, cobra.ShellCompDirectiveError
	}

	sessions := make([]string, 0, len(fixConfig.Sessions))
	for _, session := range fixConfig.Sessions {
		sessions = append(sessions, session.Name)
	}

	return sessions, cobra.ShellCompDirectiveNoFileComp
}
