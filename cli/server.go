package cli

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/log"
	"github.com/odpf/salt/printer"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/internal/server"
	"github.com/odpf/siren/pkg/pgc"
	"github.com/odpf/siren/pkg/secret"
	"github.com/odpf/siren/pkg/telemetry"
	"github.com/odpf/siren/pkg/worker"
	"github.com/odpf/siren/pkg/zaputil"
	"github.com/odpf/siren/plugins/queues"
	"github.com/odpf/siren/plugins/queues/inmemory"
	"github.com/odpf/siren/plugins/queues/postgresq"
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
			return StartServer(cmd.Context(), cfg)
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

			if err := pgc.Migrate(cfg.DB); err != nil {
				return err
			}
			printer.Success("Migration done")
			printer.Space()
			printer.SuccessIcon()
			return nil
		},
	}

	c.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return c
}

func StartServer(ctx context.Context, cfg config.Config) error {
	logger := initLogger(cfg.Log)

	telemetry.Init(ctx, cfg.Telemetry, logger)
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.Telemetry.ServiceName),
		newrelic.ConfigLicense(cfg.Telemetry.NewRelicAPIKey),
	)
	if err != nil {
		return err
	}

	dbClient, err := db.New(cfg.DB)
	if err != nil {
		return err
	}

	pgClient, err := pgc.NewClient(logger, dbClient)
	if err != nil {
		return err
	}

	encryptor, err := secret.New(cfg.Service.EncryptionKey)
	if err != nil {
		return fmt.Errorf("cannot initialize encryptor: %w", err)
	}

	var queue, dlq notification.Queuer
	switch cfg.Notification.Queue.Kind {
	case queues.KindPostgres:
		queue, err = postgresq.New(logger, cfg.DB)
		if err != nil {
			return err
		}
		dlq, err = postgresq.New(logger, cfg.DB, postgresq.WithStrategy(postgresq.StrategyDLQ))
		if err != nil {
			return err
		}
	default:
		queue = inmemory.New(logger, 50)
		dlq = inmemory.New(logger, 10)
	}

	apiDeps, notifierRegistry, err := InitAPIDeps(logger, cfg, pgClient, encryptor, queue)
	if err != nil {
		return err
	}

	// notification
	// run worker
	cancelWorkerChan := make(chan struct{})
	wg := &sync.WaitGroup{}

	if cfg.Notification.MessageHandler.Enabled {
		workerTicker := worker.NewTicker(logger, worker.WithTickerDuration(cfg.Notification.MessageHandler.PollDuration), worker.WithID("message-handler"))
		notificationHandler := notification.NewHandler(cfg.Notification.MessageHandler, logger, queue, notifierRegistry,
			notification.HandlerWithIdentifier(workerTicker.GetID()))
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerTicker.Run(ctx, cancelWorkerChan, func(ctx context.Context, runningAt time.Time) error {
				return notificationHandler.Process(ctx, runningAt)
			})
		}()
	}
	if cfg.Notification.DLQHandler.Enabled {
		workerDLQTicker := worker.NewTicker(logger, worker.WithTickerDuration(cfg.Notification.DLQHandler.PollDuration), worker.WithID("dlq-handler"))
		notificationDLQHandler := notification.NewHandler(cfg.Notification.DLQHandler, logger, dlq, notifierRegistry,
			notification.HandlerWithIdentifier(workerDLQTicker.GetID()))
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerDLQTicker.Run(ctx, cancelWorkerChan, func(ctx context.Context, runningAt time.Time) error {
				return notificationDLQHandler.Process(ctx, runningAt)
			})
		}()
	}

	err = server.RunServer(
		ctx,
		cfg.Service,
		logger,
		nrApp,
		apiDeps,
	)

	logger.Info("server stopped", "error", err)

	// stopping server first before cancelling worker
	close(cancelWorkerChan)
	// wait for all workers to stop before stopping the queue
	wg.Wait()

	const gracefulStopQueueWaitPeriod = 5 * time.Second
	timeoutCtx, cancel := context.WithTimeout(context.Background(), gracefulStopQueueWaitPeriod)
	defer cancel()

	if err := queue.Stop(timeoutCtx); err != nil {
		logger.Error("error stopping queue", "error", err)
	}
	if err := dlq.Stop(timeoutCtx); err != nil {
		logger.Error("error stopping dlq", "error", err)
	}

	if err := pgClient.Close(); err != nil {
		return err
	}

	return err
}

func initLogger(cfg config.Log) log.Logger {
	defaultConfig := zap.NewProductionConfig()
	defaultConfig.Level = zap.NewAtomicLevelAt(getZapLogLevelFromString(cfg.Level))

	if cfg.GCPCompatible {
		defaultConfig = zap.Config{
			Level:       zap.NewAtomicLevelAt(getZapLogLevelFromString(cfg.Level)),
			Encoding:    "json",
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			EncoderConfig:    zaputil.EncoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	return log.NewZap(log.ZapWithConfig(
		defaultConfig,
		zap.Fields(zaputil.ServiceContext(serviceName)),
		zap.AddCaller(),
		zap.AddStacktrace(zap.DPanicLevel),
	))

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
