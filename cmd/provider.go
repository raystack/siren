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
	"google.golang.org/protobuf/types/known/structpb"
)

func providersCmd(c *configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "provider",
		Aliases: []string{"providers"},
		Short:   "Manage providers",
		Long: heredoc.Doc(`
			Work with providers.
			
			Providers are the system for which we intend to mange monitoring and alerting.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
	}

	cmd.AddCommand(listProvidersCmd(c))
	cmd.AddCommand(createProviderCmd(c))
	cmd.AddCommand(getProviderCmd(c))
	cmd.AddCommand(updateProviderCmd(c))
	cmd.AddCommand(deleteProviderCmd(c))
	return cmd
}

func listProvidersCmd(c *configuration) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List providers",
		Long: heredoc.Doc(`
			List all registered providers.
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

			res, err := client.ListProviders(ctx, &sirenv1beta1.ListProvidersRequest{})
			if err != nil {
				return err
			}

			providers := res.Providers
			report := [][]string{}

			fmt.Printf(" \nShowing %d of %d providers\n \n", len(providers), len(providers))
			report = append(report, []string{"ID", "TYPE", "URN", "NAME"})

			for _, p := range providers {
				report = append(report, []string{
					fmt.Sprintf("%v", p.GetId()),
					p.GetType(),
					p.GetUrn(),
					p.GetName(),
				})
			}
			printer.Table(os.Stdout, report)

			fmt.Println("\nFor details on a provider, try: siren provider view <id>")
			return nil
		},
	}
}

func createProviderCmd(c *configuration) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new provider",
		Long: heredoc.Doc(`
			Create a new provider.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var providerConfig domain.Provider
			if err := parseFile(filePath, &providerConfig); err != nil {
				return err
			}

			grpcCredentials, err := structpb.NewStruct(providerConfig.Credentials)
			if err != nil {
				return err
			}

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.CreateProvider(ctx, &sirenv1beta1.CreateProviderRequest{
				Host:        providerConfig.Host,
				Urn:         providerConfig.Urn,
				Name:        providerConfig.Name,
				Type:        providerConfig.Type,
				Credentials: grpcCredentials,
				Labels:      providerConfig.Labels,
			})

			if err != nil {
				return err
			}

			fmt.Printf("Provider created with id: %v\n", res.GetId())

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to the provider config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func getProviderCmd(c *configuration) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View a provider details",
		Long: heredoc.Doc(`
			View a provider.

			Display the Id, name, and other information about a provider.
		`),
		Example: heredoc.Doc(`
			$ siren provider view 1
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
				return fmt.Errorf("invalid provider id: %v", err)
			}

			res, err := client.GetProvider(ctx, &sirenv1beta1.GetProviderRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			provider := &domain.Provider{
				Id:          res.GetId(),
				Host:        res.GetHost(),
				Urn:         res.GetUrn(),
				Name:        res.GetName(),
				Type:        res.GetType(),
				Credentials: res.GetCredentials().AsMap(),
				Labels:      res.GetLabels(),
				CreatedAt:   res.CreatedAt.AsTime(),
				UpdatedAt:   res.UpdatedAt.AsTime(),
			}

			if err := printer.Text(provider, format); err != nil {
				return fmt.Errorf("failed to format provider: %v", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "yaml", "Print output with the selected format")

	return cmd
}

func updateProviderCmd(c *configuration) *cobra.Command {
	var id uint64
	var filePath string
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a provider",
		Long: heredoc.Doc(`
			Edit an existing provider.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var providerConfig domain.Provider
			if err := parseFile(filePath, &providerConfig); err != nil {
				return err
			}

			grpcCredentials, err := structpb.NewStruct(providerConfig.Credentials)
			if err != nil {
				return err
			}

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			_, err = client.UpdateProvider(ctx, &sirenv1beta1.UpdateProviderRequest{
				Id:          id,
				Host:        providerConfig.Host,
				Name:        providerConfig.Name,
				Type:        providerConfig.Type,
				Credentials: grpcCredentials,
				Labels:      providerConfig.Labels,
			})
			if err != nil {
				return err
			}

			fmt.Println("Successfully updated provider")

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "provider id")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the provider config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func deleteProviderCmd(c *configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a provider details",
		Example: heredoc.Doc(`
			$ siren provider delete 1
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
				return fmt.Errorf("invalid provider id: %v", err)
			}

			_, err = client.DeleteProvider(ctx, &sirenv1beta1.DeleteProviderRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			fmt.Println("Successfully deleted provider")
			return nil
		},
	}

	return cmd
}
