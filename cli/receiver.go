package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	"github.com/odpf/siren/core/receiver"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
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
	cmd.AddCommand(notifyReceiverCmd(c))
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

			if res.GetData() == nil {
				return errors.New("no response from server")
			}

			receivers := res.GetData()
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
			var receiverConfig receiver.Receiver
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
	cmd.MarkFlagRequired("file")

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

			if res.GetData() == nil {
				return errors.New("no response from server")
			}

			receiver := &receiver.Receiver{
				ID:             res.GetData().GetId(),
				Name:           res.GetData().GetName(),
				Type:           res.GetData().GetType(),
				Configurations: res.GetData().GetConfigurations().AsMap(),
				Labels:         res.GetData().GetLabels(),
				Data:           res.GetData().GetData().AsMap(),
				CreatedAt:      res.GetData().GetCreatedAt().AsTime(),
				UpdatedAt:      res.GetData().GetUpdatedAt().AsTime(),
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
			var receiverConfig receiver.Receiver
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

			fmt.Printf("Successfully updated receiver with id %q", id)

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "receiver id")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the receiver config")
	cmd.MarkFlagRequired("file")

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

func notifyReceiverCmd(c *configuration) *cobra.Command {
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
			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			rcv, err := client.GetReceiver(ctx, &sirenv1beta1.GetReceiverRequest{
				Id: id,
			})
			if err != nil {
				return err
			}

			if rcv.GetData() == nil {
				return errors.New("no response from server")
			}

			notifyReceiverReq := &sirenv1beta1.NotifyReceiverRequest{}
			notifyReceiverReq.Id = rcv.GetData().GetId()
			switch rcv.GetData().GetType() {
			case "slack":
				var slackConfig *structpb.Struct
				if err := parseFile(filePath, &slackConfig); err != nil {
					return err
				}

				notifyReceiverReq.Payload = slackConfig
			}

			_, err = client.NotifyReceiver(ctx, notifyReceiverReq)
			if err != nil {
				return err
			}

			fmt.Println("Successfully send receiver notification")

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "receiver id")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the receiver notification config")
	cmd.MarkFlagRequired("file")

	return cmd
}