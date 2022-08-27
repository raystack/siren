package cli

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	"github.com/odpf/siren/pkg/errors"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"github.com/spf13/cobra"
)

func alertsCmd(c *configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "alert",
		Aliases: []string{"alerts"},
		Short:   "Manage alerts",
		Long: heredoc.Doc(`
			Work with alerts.
			
			alerts are historical events triggered by alerting providers like cortex, influx etc.
		`),
		Example: heredoc.Doc(`
			$ siren alert list --provider-name=cortex --provider-id=1 --resource-name=demo
			$ siren alert list --provider-name=cortex --provider-id=1 --resource-name=demo --start-time=1636959300000 --end-time=1636959369716
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
	}

	cmd.AddCommand(listAlertsCmd(c))

	return cmd
}

func listAlertsCmd(c *configuration) *cobra.Command {
	var providerName string
	var providerId uint64
	var resouceName string
	var startTime uint64
	var endTime uint64
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List alerts",
		Long: heredoc.Doc(`
			List all alerts.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.ListAlerts(ctx, &sirenv1beta1.ListAlertsRequest{
				ProviderName: providerName,
				ProviderId:   providerId,
				ResourceName: resouceName,
				StartTime:    startTime,
				EndTime:      endTime,
			})
			if err != nil {
				return err
			}

			if res.GetAlerts() == nil {
				return errors.New("no response from server")
			}

			alerts := res.GetAlerts()
			report := [][]string{}

			// TODO unclear log
			fmt.Printf(" \nShowing %d of %d alerts\n \n", len(alerts), len(alerts))
			report = append(report, []string{"ID", "PROVIDER_ID", "RESOURCE_NAME", "METRIC_NAME", "METRIC_VALUE", "SEVERITY"})

			for _, p := range alerts {
				report = append(report, []string{
					fmt.Sprintf("%v", p.GetId()),
					strconv.FormatUint(p.GetProviderId(), 10),
					p.GetResourceName(),
					p.GetMetricName(),
					p.GetMetricValue(),
					p.GetSeverity(),
				})
			}
			printer.Table(os.Stdout, report)

			fmt.Println("\nFor details on a alert, try: siren alert view <id>")
			return nil
		},
	}

	cmd.Flags().StringVar(&providerName, "provider-name", "", "provider name")
	cmd.MarkFlagRequired("provider-name")
	cmd.Flags().Uint64Var(&providerId, "provider-id", 0, "provider id")
	cmd.MarkFlagRequired("provider-id")
	cmd.Flags().StringVar(&resouceName, "resource-name", "", "resource name")
	cmd.Flags().Uint64Var(&startTime, "start-time", 0, "start time")
	cmd.Flags().Uint64Var(&endTime, "end-time", 0, "end time")
	return cmd
}
