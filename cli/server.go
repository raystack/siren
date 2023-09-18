package cli

import (
	"context"
	"sync"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/printer"
	"github.com/goto/siren/config"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/internal/server"
	"github.com/goto/siren/internal/store/postgres"
	"github.com/goto/siren/pkg/worker"
	"github.com/goto/siren/pkg/zaputil"
	"github.com/goto/siren/plugins/queues"
	"github.com/goto/siren/plugins/queues/inmemory"
	"github.com/goto/siren/plugins/queues/postgresq"
	"github.com/spf13/cobra"
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

			if err := postgres.Migrate(cfg.DB); err != nil {
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
	logger := zaputil.InitLogger(serviceName, cfg.Log.Level, cfg.Log.GCPCompatible)

	var queue, dlq notification.Queuer
	var err error
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

	apiDeps, nrApp, pgClient, notifierRegistry, providersPluginManager, err := InitDeps(ctx, logger, cfg, queue)
	if err != nil {
		return err
	}

	// notification
	// run worker
	cancelWorkerChan := make(chan struct{})
	wg := &sync.WaitGroup{}

	if cfg.Notification.MessageHandler.Enabled {
		workerTicker := worker.NewTicker(logger, worker.WithTickerDuration(cfg.Notification.MessageHandler.PollDuration), worker.WithID("message-handler"))
		notificationHandler := notification.NewHandler(cfg.Notification.MessageHandler, logger, nrApp, queue, notifierRegistry,
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
		notificationDLQHandler := notification.NewHandler(cfg.Notification.DLQHandler, logger, nrApp, dlq, notifierRegistry,
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

	logger.Warn("stopping plugins")
	providersPluginManager.Stop()
	logger.Warn("all plugins stopped")

	logger.Warn("stopping queue...")
	if err := queue.Stop(timeoutCtx); err != nil {
		logger.Error("error stopping queue", "error", err)
	}
	logger.Warn("queue stopped...")

	logger.Warn("stopping dlq...")
	if err := dlq.Stop(timeoutCtx); err != nil {
		logger.Error("error stopping dlq", "error", err)
	}
	logger.Warn("dlq stopped...")

	logger.Warn("closing db...")
	if err := pgClient.Close(); err != nil {
		logger.Error("error when closing db", "err", err)
	}
	logger.Warn("db closed...")

	logger.Warn("exiting plugins...")
	if err := pgClient.Close(); err != nil {
		logger.Error("error when closing db", "err", err)
	}
	logger.Warn("plugins exited...")

	return err
}
