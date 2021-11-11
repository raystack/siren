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

	rootCmd.AddCommand(configCommand())
	rootCmd.AddCommand(serveCommand())
	rootCmd.AddCommand(migrateCommand())
	rootCmd.AddCommand(uploadCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
