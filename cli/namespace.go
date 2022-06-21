package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	"github.com/odpf/siren/core/namespace"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"
)

func namespacesCmd(c *configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "namespace",
		Aliases: []string{"namespaces"},
		Short:   "Manage namespaces",
		Long: heredoc.Doc(`
			Work with namespaces.
			
			namespaces are used for multi-tenancy for a given provider.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
	}

	cmd.AddCommand(listNamespacesCmd(c))
	cmd.AddCommand(createNamespaceCmd(c))
	cmd.AddCommand(getNamespaceCmd(c))
	cmd.AddCommand(updateNamespaceCmd(c))
	cmd.AddCommand(deleteNamespaceCmd(c))
	return cmd
}

func listNamespacesCmd(c *configuration) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List namespaces",
		Long: heredoc.Doc(`
			List all registered namespaces.
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

			res, err := client.ListNamespaces(ctx, &sirenv1beta1.ListNamespacesRequest{})
			if err != nil {
				return err
			}

			if res.GetData() == nil {
				return errors.New("no response from server")
			}

			namespaces := res.GetData()
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
}

func createNamespaceCmd(c *configuration) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new namespace",
		Long: heredoc.Doc(`
			Create a new namespace.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var namespaceConfig namespace.Namespace
			if err := parseFile(filePath, &namespaceConfig); err != nil {
				return err
			}

			grpcCredentials, err := structpb.NewStruct(namespaceConfig.Credentials)
			if err != nil {
				return err
			}

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.CreateNamespace(ctx, &sirenv1beta1.CreateNamespaceRequest{
				Provider:    namespaceConfig.Provider,
				Urn:         namespaceConfig.URN,
				Name:        namespaceConfig.Name,
				Credentials: grpcCredentials,
				Labels:      namespaceConfig.Labels,
			})

			if err != nil {
				return err
			}

			fmt.Printf("namespace created with id: %v\n", res.GetId())

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to the namespace config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func getNamespaceCmd(c *configuration) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View a namespace details",
		Long: heredoc.Doc(`
			View a namespace.

			Display the Id, name, and other information about a namespace.
		`),
		Example: heredoc.Doc(`
			$ siren namespace view 1
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
				return fmt.Errorf("invalid namespace id: %v", err)
			}

			res, err := client.GetNamespace(ctx, &sirenv1beta1.GetNamespaceRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			if res.GetData() == nil {
				return errors.New("no response from server")
			}

			nspace := &namespace.Namespace{
				ID:          res.GetData().GetId(),
				URN:         res.GetData().GetUrn(),
				Name:        res.GetData().GetName(),
				Credentials: res.GetData().GetCredentials().AsMap(),
				Labels:      res.GetData().GetLabels(),
				CreatedAt:   res.GetData().GetCreatedAt().AsTime(),
				UpdatedAt:   res.GetData().GetUpdatedAt().AsTime(),
			}

			if err := printer.Text(nspace, format); err != nil {
				return fmt.Errorf("failed to format namespace: %v", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "yaml", "Print output with the selected format")

	return cmd
}

func updateNamespaceCmd(c *configuration) *cobra.Command {
	var id uint64
	var filePath string
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a namespace",
		Long: heredoc.Doc(`
			Edit an existing namespace.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var namespaceConfig namespace.Namespace
			if err := parseFile(filePath, &namespaceConfig); err != nil {
				return err
			}

			grpcCredentials, err := structpb.NewStruct(namespaceConfig.Credentials)
			if err != nil {
				return err
			}

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.UpdateNamespace(ctx, &sirenv1beta1.UpdateNamespaceRequest{
				Id:          id,
				Provider:    namespaceConfig.Provider,
				Name:        namespaceConfig.Name,
				Credentials: grpcCredentials,
				Labels:      namespaceConfig.Labels,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Successfully updated namespace with id %q", res.GetId())

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "namespace id")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the namespace config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func deleteNamespaceCmd(c *configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a namespace details",
		Example: heredoc.Doc(`
			$ siren namespace delete 1
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
				return fmt.Errorf("invalid namespace id: %v", err)
			}

			_, err = client.DeleteNamespace(ctx, &sirenv1beta1.DeleteNamespaceRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			fmt.Println("Successfully deleted namespace")
			return nil
		},
	}

	return cmd
}
