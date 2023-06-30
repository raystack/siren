package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/salt/cmdx"
	"github.com/raystack/salt/printer"
	"github.com/raystack/siren/core/provider"
	"github.com/raystack/siren/pkg/errors"
	sirenv1beta1 "github.com/raystack/siren/proto/raystack/siren/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"
)

func providersCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "provider",
		Aliases: []string{"providers"},
		Short:   "Manage providers",
		Long: heredoc.Doc(`
			Work with providers.
			
			Providers are the system for which we intend to mange monitoring and alerting.
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(
		listProvidersCmd(cmdxConfig),
		createProviderCmd(cmdxConfig),
		getProviderCmd(cmdxConfig),
		updateProviderCmd(cmdxConfig),
		deleteProviderCmd(cmdxConfig),
	)

	return cmd
}

func listProvidersCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List providers",
		Long: heredoc.Doc(`
			List all registered providers.
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

			res, err := client.ListProviders(ctx, &sirenv1beta1.ListProvidersRequest{})
			if err != nil {
				return err
			}

			if res.GetProviders() == nil {
				return errors.New("no response from server")
			}

			spinner.Stop()
			providers := res.GetProviders()
			report := [][]string{}

			// TODO unclear log
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

	return cmd
}

func createProviderCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new provider",
		Long: heredoc.Doc(`
			Create a new provider.
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

			var providerConfig provider.Provider
			if err := parseFile(filePath, &providerConfig); err != nil {
				return err
			}

			grpcCredentials, err := structpb.NewStruct(providerConfig.Credentials)
			if err != nil {
				return err
			}

			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.CreateProvider(ctx, &sirenv1beta1.CreateProviderRequest{
				Host:        providerConfig.Host,
				Urn:         providerConfig.URN,
				Name:        providerConfig.Name,
				Type:        providerConfig.Type,
				Credentials: grpcCredentials,
				Labels:      providerConfig.Labels,
			})

			if err != nil {
				return err
			}

			spinner.Stop()
			printer.Successf("Provider created with id: %v", res.GetId())
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to the provider config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func getProviderCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View a provider details",
		Long: heredoc.Doc(`
			View a provider.

			Display the id, name, and other information about a provider.
		`),
		Example: heredoc.Doc(`
			$ siren provider view 1
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
				return fmt.Errorf("invalid provider id: %v", err)
			}

			res, err := client.GetProvider(ctx, &sirenv1beta1.GetProviderRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			if res.GetProvider() == nil {
				return errors.New("no response from server")
			}

			provider := &provider.Provider{
				ID:          res.GetProvider().GetId(),
				Host:        res.GetProvider().GetHost(),
				URN:         res.GetProvider().GetUrn(),
				Name:        res.GetProvider().GetName(),
				Type:        res.GetProvider().GetType(),
				Credentials: res.GetProvider().GetCredentials().AsMap(),
				Labels:      res.GetProvider().GetLabels(),
				CreatedAt:   res.GetProvider().GetCreatedAt().AsTime(),
				UpdatedAt:   res.GetProvider().GetUpdatedAt().AsTime(),
			}

			spinner.Stop()
			if err := printer.File(provider, format); err != nil {
				return fmt.Errorf("failed to format provider: %v", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "yaml", "Print output with the selected format")

	return cmd
}

func updateProviderCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	var id uint64
	var filePath string
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a provider",
		Long: heredoc.Doc(`
			Edit an existing provider.
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

			var providerConfig provider.Provider
			if err := parseFile(filePath, &providerConfig); err != nil {
				return err
			}

			grpcCredentials, err := structpb.NewStruct(providerConfig.Credentials)
			if err != nil {
				return err
			}

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

			spinner.Stop()
			printer.Success("Successfully updated provider")
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "provider id")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the provider config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func deleteProviderCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a provider details",
		Example: heredoc.Doc(`
			$ siren provider delete 1
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
				return fmt.Errorf("invalid provider id: %v", err)
			}

			_, err = client.DeleteProvider(ctx, &sirenv1beta1.DeleteProviderRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			spinner.Stop()
			printer.Success("Successfully deleted provider")
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	return cmd
}
