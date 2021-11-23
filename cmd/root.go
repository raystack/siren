package cmd

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"os"

	"github.com/spf13/cobra"
)

// Execute runs the command line interface
func Execute() {
	rootCmd := &cobra.Command{
		Use:   "siren <command> <subcommand> [flags]",
		Short: "siren",
		Long: heredoc.Doc(`
			Siren.

			Siren provides alerting on metrics of your applications using Cortex metrics
			in a simple DIY configuration. With Siren, you can define templates(using go templates), and 
			create/edit/enable/disable prometheus rules on demand.`),
		SilenceUsage:  true,
		SilenceErrors: true,
		Annotations: map[string]string{
			"group:core": "true",
			"help:learn": heredoc.Doc(`
				Use 'siren <command> <subcommand> --help' for more information about a command.
				Read the manual at https://odpf.gitbook.io/siren/
			`),
			"help:feedback": heredoc.Doc(`
				Open an issue here https://github.com/odpf/siren/issues
			`),
		},
	}

	cliConfig, err := readConfig()
	if err != nil {
		fmt.Println(err)
	}

	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(serveCmd())
	rootCmd.AddCommand(migrateCmd())
	rootCmd.AddCommand(providersCmd(cliConfig))
	rootCmd.AddCommand(namespacesCmd(cliConfig))
	rootCmd.AddCommand(receiversCmd(cliConfig))
	rootCmd.AddCommand(templatesCmd(cliConfig))
	rootCmd.AddCommand(rulesCmd(cliConfig))
	rootCmd.AddCommand(alertsCmd(cliConfig))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
