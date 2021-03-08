package cmd

import (
	"github.com/odpf/siren/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Run DB migrations",
		RunE:  migrate,
	})
}

func migrate(cmd *cobra.Command, args []string) error {
	c := app.LoadConfig()
	return app.RunMigrations(c)
}
