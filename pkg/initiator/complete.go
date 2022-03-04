package initiator

import (
	"github.com/spf13/cobra"
	"sylr.dev/fix/config"
)

func completeContext(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	options := config.GetOptions()
	fixConfig := config.GetConfig()

	if conf, err := config.ReadYAMLNoAge(options.Config); err == nil {
		*fixConfig = *conf
	}

	contexts := make([]string, 0, len(fixConfig.Contexts))
	for _, context := range fixConfig.Contexts {
		contexts = append(contexts, context.Name)
	}

	return contexts, cobra.ShellCompDirectiveNoFileComp
}

func completeInitiator(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	options := config.GetOptions()
	fixConfig := config.GetConfig()

	if conf, err := config.ReadYAMLNoAge(options.Config); err == nil {
		*fixConfig = *conf
	}

	acceptors := make([]string, 0, len(fixConfig.Initiators))
	for _, acceptor := range fixConfig.Initiators {
		acceptors = append(acceptors, acceptor.Name)
	}

	return acceptors, cobra.ShellCompDirectiveNoFileComp
}

func completeSession(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	options := config.GetOptions()
	fixConfig := config.GetConfig()

	if conf, err := config.ReadYAMLNoAge(options.Config); err == nil {
		*fixConfig = *conf
	}

	sessions := make([]string, 0, len(fixConfig.Sessions))
	for _, session := range fixConfig.Sessions {
		sessions = append(sessions, session.Name)
	}

	return sessions, cobra.ShellCompDirectiveNoFileComp
}
