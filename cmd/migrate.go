package cmd

import (
	"github.com/odpf/siren/app"
	"github.com/odpf/siren/config"
	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := config.LoadConfig()
			return app.RunMigrations(c)
		},
	}
}
