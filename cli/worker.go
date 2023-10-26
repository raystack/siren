package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/siren/config"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/pkg/worker"
	"github.com/goto/siren/pkg/zaputil"
	"github.com/goto/siren/plugins/queues"
	"github.com/goto/siren/plugins/queues/postgresq"
	"github.com/spf13/cobra"
)

func workerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "worker <command> <worker_command>",
		Aliases: []string{"w"},
		Short:   "Start or manage Siren's workers",
		Long: heredoc.Doc(`
			A command to start or manage Siren's workers.

			A worker is an instance in Siren that run detached from the server.
		`),
		Example: heredoc.Doc(`
			$ siren worker start notification_handler
			$ siren worker start notification_handler -c ./config.yaml
		`),
	}

	cmd.AddCommand(
		workerStartCommand(),
	)

	return cmd
}

func workerStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start <command>",
		Aliases: []string{"w"},
		Short:   "Start a siren worker",
		Long:    "Command to start a siren worker.",
		Example: heredoc.Doc(`
			$ siren worker start notification_handler
			$ siren server start notification_handler -c ./config.yaml
		`),
	}

	cmd.AddCommand(
		workerStartNotificationHandlerCommand(),
		workerStartNotificationDLQHandlerCommand(),
	)

	return cmd
}

func workerStartNotificationHandlerCommand() *cobra.Command {
	var configFile string

	c := &cobra.Command{
		Use:   "notification_handler",
		Short: "A notification handler",
		Long:  "Start a handler to dequeue and publish notification messages.",
		Example: heredoc.Doc(`
			$ siren worker start notification_handler
			$ siren worker start notification_handler -c ./config.yaml
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
			}

			cancelWorkerChan := make(chan struct{})

			if err := StartNotificationHandlerWorker(ctx, cfg, cancelWorkerChan); err != nil {
				return err
			}

			<-ctx.Done()
			close(cancelWorkerChan)

			return nil
		},
	}

	c.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file path")
	return c
}

func workerStartNotificationDLQHandlerCommand() *cobra.Command {
	var configFile string

	c := &cobra.Command{
		Use:   "notification_dlq_handler",
		Short: "A notification dlq handler",
		Long:  "Start a handler to dequeue dlq and publish notification messages.",
		Example: heredoc.Doc(`
			$ siren worker start notification_dlq_handler
			$ siren worker start notification_dlq_handler -c ./config.yaml
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
			}

			cancelWorkerChan := make(chan struct{})

			if err := StartNotificationDLQHandlerWorker(ctx, cfg, cancelWorkerChan); err != nil {
				return err
			}

			<-ctx.Done()
			close(cancelWorkerChan)

			return nil
		},
	}

	c.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file path")
	return c
}

func StartNotificationHandlerWorker(ctx context.Context, cfg config.Config, cancelWorkerChan chan struct{}) error {
	logger := zaputil.InitLogger(serviceName, cfg.Log.Level, cfg.Log.GCPCompatible)

	_, nrApp, pgClient, notifierRegistry, _, err := InitDeps(ctx, logger, cfg, nil, false)
	if err != nil {
		return err
	}

	var queue notification.Queuer
	switch cfg.Notification.Queue.Kind {
	case queues.KindPostgres:
		queue, err = postgresq.New(logger, cfg.DB)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf(heredoc.Docf(`
						unsupported kind of queue for worker: %s
						supported queue kind are:
						- postgres
						`, cfg.Notification.Queue.Kind.String()))
	}
	workerTicker := worker.NewTicker(logger, worker.WithTickerDuration(cfg.Notification.MessageHandler.PollDuration), worker.WithID("message-worker"))
	notificationHandler := notification.NewHandler(cfg.Notification.MessageHandler, logger, nrApp, queue, notifierRegistry,
		notification.HandlerWithIdentifier(workerTicker.GetID()))

	go func() {
		workerTicker.Run(ctx, cancelWorkerChan, func(ctx context.Context, runningAt time.Time) error {
			return notificationHandler.Process(ctx, runningAt)
		})

		logger.Info("closing all clients")
		if err := pgClient.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()

	return nil
}

func StartNotificationDLQHandlerWorker(ctx context.Context, cfg config.Config, cancelWorkerChan chan struct{}) error {
	logger := zaputil.InitLogger(serviceName, cfg.Log.Level, cfg.Log.GCPCompatible)

	_, nrApp, pgClient, notifierRegistry, _, err := InitDeps(ctx, logger, cfg, nil, false)
	if err != nil {
		return err
	}

	var queue notification.Queuer
	switch cfg.Notification.Queue.Kind {
	case queues.KindPostgres:
		queue, err = postgresq.New(logger, cfg.DB, postgresq.WithStrategy(postgresq.StrategyDLQ))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf(heredoc.Docf(`
				unsupported kind of queue for worker: %s
				supported queue kind are:
				- postgres
				`, string(cfg.Notification.Queue.Kind)))
	}

	workerTicker := worker.NewTicker(logger, worker.WithTickerDuration(cfg.Notification.DLQHandler.PollDuration), worker.WithID("dlq-worker"))
	notificationHandler := notification.NewHandler(cfg.Notification.DLQHandler, logger, nrApp, queue, notifierRegistry,
		notification.HandlerWithIdentifier("dlq-"+workerTicker.GetID()))
	go func() {
		workerTicker.Run(ctx, cancelWorkerChan, func(ctx context.Context, runningAt time.Time) error {
			return notificationHandler.Process(ctx, runningAt)
		})

		logger.Info("closing all clients")
		if err := pgClient.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()

	return nil
}
