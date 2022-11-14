package cli

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/cmdx"

	"github.com/spf13/cobra"
)

const serviceName = "siren"

func New() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "siren <command> <subcommand> [flags]",
		Short: "siren",
		Long: heredoc.Doc(`
			Work seamlessly with your observability stack.`),
		SilenceUsage:  true,
		SilenceErrors: true,
		Annotations: map[string]string{
			"group": "core",
			"help:learn": heredoc.Doc(`
				Use 'siren <command> <subcommand> --help' for more information about a command.
				Read the manual at https://odpf.gitbook.io/siren/
			`),
			"help:feedback": heredoc.Doc(`
				Open an issue here https://github.com/odpf/siren/issues
			`),
			"help:environment": heredoc.Doc(`
				See 'siren help environment' for the list of supported environment variables.
			`),
		},
	}

	cmdxConfig := cmdx.SetConfig("siren")

	rootCmd.AddCommand(serverCmd())
	rootCmd.AddCommand(configCmd(cmdxConfig))
	rootCmd.AddCommand(providersCmd(cmdxConfig))
	rootCmd.AddCommand(namespacesCmd(cmdxConfig))
	rootCmd.AddCommand(receiversCmd(cmdxConfig))
	rootCmd.AddCommand(templatesCmd(cmdxConfig))
	rootCmd.AddCommand(rulesCmd(cmdxConfig))
	rootCmd.AddCommand(subscriptionsCmd(cmdxConfig))
	rootCmd.AddCommand(alertsCmd(cmdxConfig))
	rootCmd.AddCommand(jobCmd(cmdxConfig))
	rootCmd.AddCommand(workerCmd())

	// Help topics
	cmdx.SetHelp(rootCmd)
	rootCmd.AddCommand(cmdx.SetCompletionCmd("siren"))
	rootCmd.AddCommand(cmdx.SetHelpTopicCmd("environment", envHelp))
	rootCmd.AddCommand(cmdx.SetRefCmd(rootCmd))

	cmdx.SetClientHook(rootCmd, func(cmd *cobra.Command) {
		// client config
		cmd.PersistentFlags().StringP("host", "h", "", "Siren API service to connect to")
	})

	return rootCmd
}
