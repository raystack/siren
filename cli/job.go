package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/cmdx"
	"github.com/odpf/salt/log"
	"github.com/odpf/salt/printer"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/plugins/queues"
	"github.com/odpf/siren/plugins/queues/postgresq"
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
		Args: cobra.ExactValidArgs(1),
		ValidArgs: []string{
			"cleanup_queue",
		},
		Example: heredoc.Doc(`
			$ siren job run cleanup_queue
		`),
	}

	cmd.AddCommand(
		jobRunCleanupQueueCommand(),
	)

	return cmd
}

func jobRunCleanupQueueCommand() *cobra.Command {
	var configFile string

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
			printer.Info("Running job cleanup_queue(%s)", cfg.Notification.Queue.Kind.String())
			if err := queue.Cleanup(cmd.Context(), queues.FilterCleanup{}); err != nil {
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
	return cmd
}
