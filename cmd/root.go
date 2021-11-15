package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute runs the command line interface
func Execute() {
	rootCmd := &cobra.Command{
		Use: "siren",
	}

	cliConfig, err := readConfig()
	if err != nil {
		fmt.Println(err)
	}

	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(serveCmd())
	rootCmd.AddCommand(migrateCmd())
	rootCmd.AddCommand(uploadCmd())
	rootCmd.AddCommand(providersCmd(cliConfig))
	rootCmd.AddCommand(namespacesCmd(cliConfig))
	rootCmd.AddCommand(receiversCmd(cliConfig))
	rootCmd.AddCommand(templatesCmd(cliConfig))
	rootCmd.AddCommand(rulesCmd(cliConfig))
	rootCmd.AddCommand(alertsCmd(cliConfig))
	rootCmd.CompletionOptions.DisableDescriptions = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
