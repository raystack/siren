package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/cmdx"
	"github.com/goto/salt/printer"
	"github.com/goto/siren/core/subscription"
	"github.com/goto/siren/pkg/errors"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"
)

func subscriptionsCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subscription",
		Aliases: []string{"subscriptions"},
		Short:   "Manage subscriptions",
		Long: heredoc.Doc(`
			Work with subscriptions.
			
			Subscribe to a notification with matching label.
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(
		listSubscriptionsCmd(cmdxConfig),
		viewSubscriptionCmd(cmdxConfig),
		createSubscriptionCmd(cmdxConfig),
		updateSubscriptionCmd(cmdxConfig),
		deleteSubscriptionCmd(cmdxConfig),
	)

	return cmd
}

func listSubscriptionsCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List subscriptions",
		Long: heredoc.Doc(`
			List all registered subscriptions.
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			ctx := cmd.Context()

			c, err := loadClientConfig(cmd, cmdxConfig)
			if err != nil {
				return err
			}

			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			resNamespaces, err := client.ListNamespaces(ctx, &sirenv1beta1.ListNamespacesRequest{})
			if err != nil {
				return err
			}

			var namespaceMaps = make(map[uint64]string)
			for _, n := range resNamespaces.GetNamespaces() {
				namespaceMaps[n.GetId()] = n.GetUrn()
			}

			res, err := client.ListSubscriptions(ctx, &sirenv1beta1.ListSubscriptionsRequest{})
			if err != nil {
				return err
			}

			if res.GetSubscriptions() == nil {
				return errors.New("no response from server")
			}

			spinner.Stop()
			subscriptions := res.GetSubscriptions()
			report := [][]string{}

			fmt.Printf(" \nShowing %d of %d subscriptions\n \n", len(subscriptions), len(subscriptions))
			report = append(report, []string{"ID", "URN", "NAMESPACE", "RECEIVERS", "MATCH LABELS"})

			for _, p := range subscriptions {
				var (
					ns string
					ok bool
				)

				if ns, ok = namespaceMaps[p.GetNamespace()]; !ok {
					return fmt.Errorf("unrecognized namespace %d", p.GetNamespace())
				}

				matchStr, err := json.Marshal(p.GetMatch())
				if err != nil {
					return errors.New("cannot marshal match labels")
				}

				receiversStr, err := json.Marshal(p.GetReceivers())
				if err != nil {
					return errors.New("cannot marshal receivers metadata")
				}

				report = append(report, []string{
					fmt.Sprintf("%v", p.GetId()),
					p.GetUrn(),
					ns,
					string(receiversStr),
					string(matchStr),
				})
			}
			printer.Table(os.Stdout, report)

			fmt.Println("\nFor details on a subscription, try: siren subscription view <id>")
			return nil
		},
	}

	return cmd
}

func viewSubscriptionCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View a subscription details",
		Long: heredoc.Doc(`
			View a subscription.

			Display the id, urn, namespace, receivers, and label matchers of a subscription.
		`),
		Example: heredoc.Doc(`
			$ siren subscription view 1
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			ctx := cmd.Context()

			c, err := loadClientConfig(cmd, cmdxConfig)
			if err != nil {
				return err
			}

			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid receiver id: %v", err)
			}

			res, err := client.GetSubscription(ctx, &sirenv1beta1.GetSubscriptionRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			if res.GetSubscription() == nil {
				return errors.New("no response from server")
			}

			var subscriptionReceivers []subscription.Receiver
			for _, sr := range res.GetSubscription().GetReceivers() {
				subscriptionReceivers = append(subscriptionReceivers, subscription.Receiver{
					ID:            sr.GetId(),
					Configuration: sr.GetConfiguration().AsMap(),
				})
			}

			sub := &subscription.Subscription{
				ID:        res.GetSubscription().GetId(),
				URN:       res.GetSubscription().GetUrn(),
				Namespace: res.GetSubscription().GetNamespace(),
				Receivers: subscriptionReceivers,
				Match:     res.GetSubscription().GetMatch(),
				CreatedAt: res.GetSubscription().GetCreatedAt().AsTime(),
				UpdatedAt: res.GetSubscription().GetUpdatedAt().AsTime(),
			}

			spinner.Stop()
			if err := printer.File(sub, format); err != nil {
				return fmt.Errorf("failed to format subscription: %v", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "yaml", "Print output with the selected format")

	return cmd
}

func createSubscriptionCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new subscription",
		Long: heredoc.Doc(`
			Create a new subscription.
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			ctx := cmd.Context()

			c, err := loadClientConfig(cmd, cmdxConfig)
			if err != nil {
				return err
			}

			var subscriptionDetail subscription.Subscription
			if err := parseFile(filePath, &subscriptionDetail); err != nil {
				return err
			}

			var receiverMetadatasPB []*sirenv1beta1.ReceiverMetadata
			for _, rcv := range subscriptionDetail.Receivers {
				grpcConfigurations, err := structpb.NewStruct(rcv.Configuration)
				if err != nil {
					return err
				}
				receiverMetadatasPB = append(receiverMetadatasPB, &sirenv1beta1.ReceiverMetadata{
					Id:            rcv.ID,
					Configuration: grpcConfigurations,
				})
			}

			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.CreateSubscription(ctx, &sirenv1beta1.CreateSubscriptionRequest{
				Urn:       subscriptionDetail.URN,
				Namespace: subscriptionDetail.Namespace,
				Match:     subscriptionDetail.Match,
				Receivers: receiverMetadatasPB,
			})

			if err != nil {
				return err
			}

			spinner.Stop()
			printer.Successf("Subscription created with id: %v", res.GetId())
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to the subscription config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func updateSubscriptionCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	var id uint64
	var filePath string
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a subscription",
		Long: heredoc.Doc(`
			Edit an existing subscription.

			Note: receiver type is immutable.
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			ctx := cmd.Context()

			c, err := loadClientConfig(cmd, cmdxConfig)
			if err != nil {
				return err
			}

			var subscriptionDetail subscription.Subscription
			if err := parseFile(filePath, &subscriptionDetail); err != nil {
				return err
			}

			var receiverMetadatasPB []*sirenv1beta1.ReceiverMetadata
			for _, rcv := range subscriptionDetail.Receivers {
				grpcConfigurations, err := structpb.NewStruct(rcv.Configuration)
				if err != nil {
					return err
				}
				receiverMetadatasPB = append(receiverMetadatasPB, &sirenv1beta1.ReceiverMetadata{
					Id:            rcv.ID,
					Configuration: grpcConfigurations,
				})
			}

			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			_, err = client.UpdateSubscription(ctx, &sirenv1beta1.UpdateSubscriptionRequest{
				Urn:       subscriptionDetail.URN,
				Namespace: subscriptionDetail.Namespace,
				Match:     subscriptionDetail.Match,
				Receivers: receiverMetadatasPB,
			})
			if err != nil {
				return err
			}

			spinner.Stop()
			printer.Successf("Successfully updated subscription with id %d", id)
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "subscription id")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the subscription config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func deleteSubscriptionCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a subscription details",
		Example: heredoc.Doc(`
			$ siren subscription delete 1
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			ctx := cmd.Context()

			c, err := loadClientConfig(cmd, cmdxConfig)
			if err != nil {
				return err
			}

			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid subscription id: %v", err)
			}

			_, err = client.DeleteSubscription(ctx, &sirenv1beta1.DeleteSubscriptionRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			spinner.Stop()
			printer.Success("Successfully deleted subscription")
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	return cmd
}
