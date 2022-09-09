package cli

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/log"
	"github.com/odpf/salt/printer"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/server"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/odpf/siren/pkg/cortex"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/secret"
	"github.com/odpf/siren/pkg/slack"
	"github.com/odpf/siren/pkg/telemetry"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server <command>",
		Aliases: []string{"s"},
		Short:   "Run siren server",
		Long:    "Server management commands.",
		Example: heredoc.Doc(`
			$ siren server init
			$ siren server start
			$ siren server start -c ./config.yaml
			$ siren server migrate
			$ siren server migrate -c ./config.yaml
		`),
	}

	cmd.AddCommand(
		serverInitCommand(),
		serverStartCommand(),
		serverMigrateCommand(),
	)

	return cmd
}

func serverInitCommand() *cobra.Command {
	var configFile string
	// var resourcesURL string
	// var rulesURL string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize server",
		Long: heredoc.Doc(`
			Initializing server. Creating a sample of siren server config.
			Default: ./config.yaml
		`),
		Example: "siren server init",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Init(configFile); err != nil {
				return err
			}

			printer.Successf("Server config created: %s", configFile)
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "output", "o", "./config.yaml", "Output config file path")

	return cmd
}

func serverStartCommand() *cobra.Command {
	var configFile string

	c := &cobra.Command{
		Use:     "start",
		Short:   "Start server on default port 8080",
		Example: "siren server start",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
			}
			return runServer(cfg)
		},
	}

	c.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file path")
	return c
}

func serverMigrateCommand() *cobra.Command {
	var configFile string

	c := &cobra.Command{
		Use:     "migrate",
		Short:   "Run DB Schema Migrations",
		Example: "siren migrate",
		RunE: func(c *cobra.Command, args []string) error {
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
			}

			if err := postgres.Migrate(cfg.DB); err != nil {
				return err
			}
			return nil
		},
	}

	c.Flags().StringVarP(&configFile, "config", "c", "", "Config file path")
	return c
}

func runServer(cfg config.Config) error {
	nr, err := telemetry.New(cfg.NewRelic)
	if err != nil {
		return err
	}

	logger := initLogger(cfg.Log.Level)

	dbClient, err := db.New(cfg.DB)
	if err != nil {
		return err
	}

	pgClient, err := postgres.NewClient(logger, dbClient)
	if err != nil {
		return err
	}

	httpClient := &http.Client{}
	encryptor, err := secret.New(cfg.EncryptionKey)
	if err != nil {
		return fmt.Errorf("cannot initialize encryptor: %w", err)
	}

	templateRepository := postgres.NewTemplateRepository(pgClient)
	templateService := template.NewService(templateRepository)

	alertRepository := postgres.NewAlertRepository(pgClient)
	alertHistoryService := alert.NewService(alertRepository)

	providerRepository := postgres.NewProviderRepository(pgClient)
	providerService := provider.NewService(providerRepository)

	namespaceRepository := postgres.NewNamespaceRepository(pgClient)
	namespaceService := namespace.NewService(encryptor, namespaceRepository)

	if cfg.Cortex.PrometheusAlertManagerConfigYaml == "" || cfg.Cortex.PrometheusAlertManagerHelperTemplate == "" {
		return errors.New("empty prometheus alert manager config template")
	}

	cortexClient, err := cortex.NewClient(cortex.Config{Address: cfg.Cortex.Address},
		cortex.WithHelperTemplate(cfg.Cortex.PrometheusAlertManagerConfigYaml, cfg.Cortex.PrometheusAlertManagerHelperTemplate),
	)
	if err != nil {
		return fmt.Errorf("failed to init cortex client: %w", err)
	}

	ruleRepository := postgres.NewRuleRepository(pgClient)
	ruleService := rule.NewService(
		ruleRepository,
		templateService,
		namespaceService,
		providerService,
		cortexClient,
	)

	slackClient := slack.NewClient(slack.ClientWithHTTPClient(httpClient))
	slackReceiverService := receiver.NewSlackService(slackClient, encryptor)
	httpReceiverService := receiver.NewHTTPService()
	pagerDutyReceiverService := receiver.NewPagerDutyService()
	receiverRepository := postgres.NewReceiverRepository(pgClient)
	receiverService := receiver.NewService(
		receiverRepository,
		map[string]receiver.TypeService{
			receiver.TypeSlack:     slackReceiverService,
			receiver.TypeHTTP:      httpReceiverService,
			receiver.TypePagerDuty: pagerDutyReceiverService,
		},
	)

	subscriptionRepository := postgres.NewSubscriptionRepository(pgClient)
	subscriptionService := subscription.NewService(subscriptionRepository, providerService, namespaceService, receiverService, cortexClient)

	return server.RunServer(
		cfg.Service,
		logger,
		nr,
		templateService,
		ruleService,
		alertHistoryService,
		providerService,
		namespaceService,
		receiverService,
		subscriptionService)
}

func initLogger(logLevel string) log.Logger {
	defaultConfig := zap.NewProductionConfig()
	defaultConfig.Level = zap.NewAtomicLevelAt(getZapLogLevelFromString(logLevel))
	return log.NewZap(log.ZapWithConfig(defaultConfig, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel)))
}

// getZapLogLevelFromString helps to set logLevel from string
func getZapLogLevelFromString(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
