package cli

import (
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/internal/store/postgres"
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

			logger := initLogger(cfg.Log)
			if err := runPostgresMigrations(logger, cfg); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return cmd
}

func runPostgresMigrations(logger log.Logger, cfg config.Config) error {
	return postgres.Migrate(cfg.DB)
}
