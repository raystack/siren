package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/salt/cmdx"
	"github.com/raystack/salt/printer"
	"github.com/raystack/siren/core/namespace"
	"github.com/raystack/siren/pkg/errors"
	sirenv1beta1 "github.com/raystack/siren/proto/raystack/siren/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"
)

func namespacesCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "namespace",
		Aliases: []string{"namespaces"},
		Short:   "Manage namespaces",
		Long: heredoc.Doc(`
			Work with namespaces.
			
			Namespaces are used for multi-tenancy for a given provider.
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(
		listNamespacesCmd(cmdxConfig),
		createNamespaceCmd(cmdxConfig),
		getNamespaceCmd(cmdxConfig),
		updateNamespaceCmd(cmdxConfig),
		deleteNamespaceCmd(cmdxConfig),
	)

	return cmd
}

func listNamespacesCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List namespaces",
		Long: heredoc.Doc(`
			List all registered namespaces.
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

			res, err := client.ListNamespaces(ctx, &sirenv1beta1.ListNamespacesRequest{})
			if err != nil {
				return err
			}

			if res.GetNamespaces() == nil {
				return errors.New("no response from server")
			}

			spinner.Stop()
			namespaces := res.GetNamespaces()
			report := [][]string{}

			fmt.Printf(" \nShowing %d of %d namespaces\n \n", len(namespaces), len(namespaces))
			report = append(report, []string{"ID", "URN", "NAME"})

			for _, p := range namespaces {
				report = append(report, []string{
					fmt.Sprintf("%v", p.GetId()),
					p.GetUrn(),
					p.GetName(),
				})
			}
			printer.Table(os.Stdout, report)

			fmt.Println("\nFor details on a namespace, try: siren namespace view <id>")
			return nil
		},
	}

	return cmd
}

func createNamespaceCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new namespace",
		Long: heredoc.Doc(`
			Create a new namespace.
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

			var namespaceConfig namespace.Namespace
			if err := parseFile(filePath, &namespaceConfig); err != nil {
				return err
			}

			grpcCredentials, err := structpb.NewStruct(namespaceConfig.Credentials)
			if err != nil {
				return err
			}

			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.CreateNamespace(ctx, &sirenv1beta1.CreateNamespaceRequest{
				Provider:    namespaceConfig.Provider.ID,
				Urn:         namespaceConfig.URN,
				Name:        namespaceConfig.Name,
				Credentials: grpcCredentials,
				Labels:      namespaceConfig.Labels,
			})

			if err != nil {
				return err
			}

			spinner.Stop()
			printer.Successf("Namespace created with id: %v", res.GetId())
			printer.Space()
			printer.SuccessIcon()
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to the namespace config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func getNamespaceCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View a namespace details",
		Long: heredoc.Doc(`
			View a namespace.

			Display the id, name, and other information about a namespace.
		`),
		Example: heredoc.Doc(`
			$ siren namespace view 1
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
				return fmt.Errorf("invalid namespace id: %v", err)
			}

			res, err := client.GetNamespace(ctx, &sirenv1beta1.GetNamespaceRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			if res.GetNamespace() == nil {
				return errors.New("no response from server")
			}

			nspace := &namespace.Namespace{
				ID:          res.GetNamespace().GetId(),
				URN:         res.GetNamespace().GetUrn(),
				Name:        res.GetNamespace().GetName(),
				Credentials: res.GetNamespace().GetCredentials().AsMap(),
				Labels:      res.GetNamespace().GetLabels(),
				CreatedAt:   res.GetNamespace().GetCreatedAt().AsTime(),
				UpdatedAt:   res.GetNamespace().GetUpdatedAt().AsTime(),
			}

			spinner.Stop()
			if err := printer.File(nspace, format); err != nil {
				return fmt.Errorf("failed to format namespace: %v", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "yaml", "Print output with the selected format")

	return cmd
}

func updateNamespaceCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	var id uint64
	var filePath string
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a namespace",
		Long: heredoc.Doc(`
			Edit an existing namespace.
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

			var namespaceConfig namespace.Namespace
			if err := parseFile(filePath, &namespaceConfig); err != nil {
				return err
			}

			grpcCredentials, err := structpb.NewStruct(namespaceConfig.Credentials)
			if err != nil {
				return err
			}

			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.UpdateNamespace(ctx, &sirenv1beta1.UpdateNamespaceRequest{
				Id:          id,
				Provider:    namespaceConfig.Provider.ID,
				Name:        namespaceConfig.Name,
				Credentials: grpcCredentials,
				Labels:      namespaceConfig.Labels,
			})
			if err != nil {
				return err
			}

			spinner.Stop()
			printer.Successf("Successfully updated namespace with id %d", res.GetId())
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "namespace id")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the namespace config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func deleteNamespaceCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a namespace details",
		Example: heredoc.Doc(`
			$ siren namespace delete 1
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
				return fmt.Errorf("invalid namespace id: %v", err)
			}

			_, err = client.DeleteNamespace(ctx, &sirenv1beta1.DeleteNamespaceRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			spinner.Stop()
			printer.Success("Successfully deleted namespace")
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	return cmd
}
