package cli

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/cmdx"
	"github.com/goto/salt/log"
	"github.com/goto/salt/printer"
	"github.com/goto/siren/config"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/internal/jobs"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/plugins/queues"
	"github.com/goto/siren/plugins/queues/postgresq"
	"github.com/spf13/cobra"
)

func jobCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "job <command>",
		Aliases: []string{"jobs"},
		Short:   "Manage siren jobs",
		Long: heredoc.Doc(`
			Execute a siren's job.
			
			A job is a task in Siren that could be executed and stopped once the task is done. 
			The Job is usually run as a CronJob to be executed on a specified time.
		`),
		Example: heredoc.Doc(`
			$ siren job run cleanup_idempotency
			$ siren job run cleanup_queue
		`),
	}

	cmd.AddCommand(
		jobRunCommand(cmdxConfig),
	)

	return cmd
}

func jobRunCommand(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Trigger a job",
		Long: heredoc.Doc(`
			Trigger a job
		`),
		Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: []string{
			"cleanup_queue",
			"cleanup_idempotency",
		},
		Example: heredoc.Doc(`
			$ siren job run cleanup_idempotency
			$ siren job run cleanup_queue
		`),
	}

	cmd.AddCommand(
		jobRunCleanupQueueCommand(),
		jobRunCleanupIdempotencyCommand(),
	)

	return cmd
}

func jobRunCleanupQueueCommand() *cobra.Command {
	var (
		configFile        string
		publishedDuration string
		pendingDuration   string
	)

	cmd := &cobra.Command{
		Use:   "cleanup_queue",
		Short: "Cleanup stale messages in queue",
		Long: heredoc.Doc(`
			Cleaning up all published messages in queue with last updated 
			more than specific threshold (default 7 days) from now() and
			
			(Optional) cleaning up all pending messages in queue with last updated 
			more than specific threshold (default 7 days) from now().
		`),
		Example: heredoc.Doc(`
			$ siren job run cleanup_queue
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
			}

			var queue notification.Queuer
			switch cfg.Notification.Queue.Kind {
			case queues.KindPostgres:
				queue, err = postgresq.New(log.NewZap(), cfg.DB)
				if err != nil {
					return err
				}
			default:
				printer.Info("Cleanup queue job only works for postgres queue")
				return nil
			}

			spinner := printer.Spin("")
			defer spinner.Stop()
			printer.Infof("Running job cleanup_queue(%s)\n", cfg.Notification.Queue.Kind.String())
			if err := queue.Cleanup(cmd.Context(), queues.FilterCleanup{
				MessagePendingTimeThreshold:   pendingDuration,
				MessagePublishedTimeThreshold: publishedDuration,
			}); err != nil {
				return err
			}
			spinner.Stop()
			printer.Success(fmt.Sprintf("Job cleanup_queue(%s) finished", cfg.Notification.Queue.Kind.String()))
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file path")
	cmd.Flags().StringVarP(&publishedDuration, "published", "p", "168h", "Cleanup treshold for published messages in string (e.g. 10h, 30m)")
	cmd.Flags().StringVarP(&publishedDuration, "pending", "s", "", "Cleanup treshold for pending messages in string (e.g. 10h, 30m)")

	return cmd
}

func jobRunCleanupIdempotencyCommand() *cobra.Command {
	var (
		configFile  string
		ttlDuration string
	)

	cmd := &cobra.Command{
		Use:   "cleanup_idempotency",
		Short: "Cleanup idempotencies outside TTL",
		Long: heredoc.Doc(`
			Cleaning up all idempotencies data outside TTL.

			Default value is 24 hours.
		`),
		Example: heredoc.Doc(`
			$ siren job run cleanup_idempotency
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configFile)
			if err != nil {
				return err
			}

			logger := initLogger(cfg.Log)

			apiDeps, _, pgClient, _, err := InitDeps(cmd.Context(), logger, cfg, nil)
			if err != nil {
				return err
			}

			spinner := printer.Spin("")
			defer spinner.Stop()
			printer.Infof("Running job cleanup_idempotency with ttl %s\n", ttlDuration)

			ttlInTimeDuration, err := time.ParseDuration(ttlDuration)
			if err != nil {
				return err
			}

			jobHandler := jobs.NewHandler(logger, apiDeps.NotificationService)
			if err := jobHandler.CleanupIdempotencies(cmd.Context(), ttlInTimeDuration); err != nil {
				spinner.Stop()
				if errors.Is(err, errors.ErrNotFound) {
					printer.Success(fmt.Sprintf("No data outside TTL %s", ttlDuration))
				} else {
					logger.Error(err.Error())
				}
			} else {
				spinner.Stop()
				printer.Success("Job cleanup_idempotency finished")
				printer.Space()
				printer.SuccessIcon()
			}

			if err := pgClient.Close(); err != nil {
				logger.Error(err.Error())
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file path")
	cmd.Flags().StringVarP(&ttlDuration, "ttl", "t", "24h", "TTL duration of idempotency data in golang duration format (e.g. 10h, 30m)")

	return cmd
}
