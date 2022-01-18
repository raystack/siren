package cmd

import (
	"github.com/odpf/siren/app"
	"github.com/odpf/siren/config"
	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database schema",
		Annotations: map[string]string{
			"group:other": "dev",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				return err
			}
			return app.RunMigrations(cfg)
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return cmd
}
