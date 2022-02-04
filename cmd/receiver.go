package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	sirenv1beta1 "github.com/odpf/siren/api/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/domain"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func receiversCmd(c *configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "receiver",
		Aliases: []string{"receivers"},
		Short:   "Manage receivers",
		Long: heredoc.Doc(`
			Work with receivers.
			
			Receivers are the medium to send notification for which we intend to mange configuration.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
	}

	cmd.AddCommand(listReceiversCmd(c))
	cmd.AddCommand(createReceiverCmd(c))
	cmd.AddCommand(getReceiverCmd(c))
	cmd.AddCommand(updateReceiverCmd(c))
	cmd.AddCommand(deleteReceiverCmd(c))
	cmd.AddCommand(sendReceiverNotificationCmd(c))
	return cmd
}

func listReceiversCmd(c *configuration) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List receivers",
		Long: heredoc.Doc(`
			List all registered receivers.
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

			res, err := client.ListReceivers(ctx, &emptypb.Empty{})
			if err != nil {
				return err
			}

			receivers := res.Receivers
			report := [][]string{}

			fmt.Printf(" \nShowing %d of %d receivers\n \n", len(receivers), len(receivers))
			report = append(report, []string{"ID", "TYPE", "NAME"})

			for _, p := range receivers {
				report = append(report, []string{
					fmt.Sprintf("%v", p.GetId()),
					p.GetType(),
					p.GetName(),
				})
			}
			printer.Table(os.Stdout, report)

			fmt.Println("\nFor details on a receiver, try: siren receiver view <id>")
			return nil
		},
	}
}

func createReceiverCmd(c *configuration) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new receiver",
		Long: heredoc.Doc(`
			Create a new receiver.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var receiverConfig domain.Receiver
			if err := parseFile(filePath, &receiverConfig); err != nil {
				return err
			}

			grpcConfigurations, err := structpb.NewStruct(receiverConfig.Configurations)
			if err != nil {
				return err
			}

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.CreateReceiver(ctx, &sirenv1beta1.CreateReceiverRequest{
				Name:           receiverConfig.Name,
				Type:           receiverConfig.Type,
				Configurations: grpcConfigurations,
				Labels:         receiverConfig.Labels,
			})

			if err != nil {
				return err
			}

			fmt.Printf("Receiver created with id: %v\n", res.GetId())

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to the receiver config")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func getReceiverCmd(c *configuration) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View a receiver details",
		Long: heredoc.Doc(`
			View a receiver.

			Display the Id, name, and other information about a receiver.
		`),
		Example: heredoc.Doc(`
			$ siren receiver view 1
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid receiver id: %v", err)
			}

			res, err := client.GetReceiver(ctx, &sirenv1beta1.GetReceiverRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			receiver := &domain.Receiver{
				Id:             res.GetId(),
				Name:           res.GetName(),
				Type:           res.GetType(),
				Configurations: res.GetConfigurations().AsMap(),
				Labels:         res.GetLabels(),
				Data:           res.GetData().AsMap(),
				CreatedAt:      res.CreatedAt.AsTime(),
				UpdatedAt:      res.UpdatedAt.AsTime(),
			}

			if err := printer.Text(receiver, format); err != nil {
				return fmt.Errorf("failed to format receiver: %v", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "yaml", "Print output with the selected format")

	return cmd
}

func updateReceiverCmd(c *configuration) *cobra.Command {
	var id uint64
	var filePath string
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a receiver",
		Long: heredoc.Doc(`
			Edit an existing receiver.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var receiverConfig domain.Receiver
			if err := parseFile(filePath, &receiverConfig); err != nil {
				return err
			}

			grpcConfigurations, err := structpb.NewStruct(receiverConfig.Configurations)
			if err != nil {
				return err
			}

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			_, err = client.UpdateReceiver(ctx, &sirenv1beta1.UpdateReceiverRequest{
				Id:             id,
				Name:           receiverConfig.Name,
				Type:           receiverConfig.Type,
				Configurations: grpcConfigurations,
				Labels:         receiverConfig.Labels,
			})
			if err != nil {
				return err
			}

			fmt.Println("Successfully updated receiver")

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "receiver id")
	_ = cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the receiver config")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func deleteReceiverCmd(c *configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a receiver details",
		Example: heredoc.Doc(`
			$ siren receiver delete 1
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid receiver id: %v", err)
			}

			_, err = client.DeleteReceiver(ctx, &sirenv1beta1.DeleteReceiverRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			fmt.Println("Successfully deleted receiver")
			return nil
		},
	}

	return cmd
}

func sendReceiverNotificationCmd(c *configuration) *cobra.Command {
	var id uint64
	var filePath string
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send a receiver notification",
		Long: heredoc.Doc(`
			Send a notification to receiver.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var notificationConfig sirenv1beta1.SendReceiverNotificationRequest

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			receiver, err := client.GetReceiver(ctx, &sirenv1beta1.GetReceiverRequest{
				Id: id,
			})
			if err != nil {
				return err
			}

			notificationConfig.Id = id
			switch receiver.Type {
			case "slack":
				var slackConfig *sirenv1beta1.SendReceiverNotificationRequest_Slack
				if err := parseFile(filePath, &slackConfig); err != nil {
					return err
				}

				notificationConfig.Data = slackConfig
			}

			_, err = client.SendReceiverNotification(ctx, &notificationConfig)
			if err != nil {
				return err
			}

			fmt.Println("Successfully send receiver notification")

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "receiver id")
	_ = cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the receiver notification config")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
