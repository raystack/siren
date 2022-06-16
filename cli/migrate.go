package cli

import (
	"fmt"
	"net/http"

	"github.com/odpf/siren/config"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/server"
	"github.com/odpf/siren/internal/store"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/odpf/siren/pkg/cortex"
	"github.com/odpf/siren/pkg/secret"
	"github.com/odpf/siren/pkg/slack"
	"github.com/pkg/errors"
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

			// TODO to be refactored later: https://github.com/odpf/siren/issues/105
			gormDB, err := postgres.New(cfg.DB)
			if err != nil {
				return err
			}

			httpClient := &http.Client{}
			repositories := store.NewRepositoryContainer(gormDB)
			encryptor, err := secret.New(cfg.EncryptionKey)
			if err != nil {
				return fmt.Errorf("cannot initialize encryptor: %w", err)
			}

			templateService := template.NewService(repositories.TemplatesRepository)
			alertHistoryService := alert.NewService(repositories.AlertRepository)

			providerService := provider.NewService(repositories.ProviderRepository)

			namespaceService := namespace.NewService(encryptor, repositories.NamespaceRepository)

			if cfg.Cortex.PrometheusAlertManagerConfigYaml == "" || cfg.Cortex.PrometheusAlertManagerHelperTemplate == "" {
				return errors.New("empty prometheus alert manager config template")
			}

			cortexClient, err := cortex.NewClient(cortex.Config{Address: cfg.Cortex.Address},
				cortex.WithHelperTemplate(cfg.Cortex.PrometheusAlertManagerConfigYaml, cfg.Cortex.PrometheusAlertManagerHelperTemplate),
			)
			if err != nil {
				return errors.Wrap(err, "failed to init cortex client")
			}

			ruleService := rule.NewService(
				repositories.RuleRepository,
				templateService,
				namespaceService,
				providerService,
				cortexClient,
			)

			slackClient := slack.NewClient(slack.ClientWithHTTPClient(httpClient))
			receiverService := receiver.NewService(repositories.ReceiverRepository, slackClient, encryptor)

			subscriptionService := subscription.NewService(repositories.SubscriptionRepository, providerService, namespaceService, receiverService, cortexClient)

			return server.RunMigrations(
				templateService,
				ruleService,
				alertHistoryService,
				providerService,
				namespaceService,
				receiverService,
				subscriptionService)
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return cmd
}
