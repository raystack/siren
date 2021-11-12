package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"github.com/spf13/cobra"
)

func templatesCmd(c *configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template",
		Aliases: []string{"templates"},
		Short:   "Manage templates",
		Long: heredoc.Doc(`
			Work with templates.
			
			templates are used for alert abstraction.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
	}

	cmd.AddCommand(listTemplatesCmd(c))
	cmd.AddCommand(upsertTemplateCmd(c))
	cmd.AddCommand(getTemplateCmd(c))
	cmd.AddCommand(deleteTemplateCmd(c))
	cmd.AddCommand(renderTemplateCmd(c))

	return cmd
}

func listTemplatesCmd(c *configuration) *cobra.Command {
	var tag string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List templates",
		Long: heredoc.Doc(`
			List all registered templates.
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

			res, err := client.ListTemplates(ctx, &sirenv1.ListTemplatesRequest{
				Tag: tag,
			})
			if err != nil {
				return err
			}

			templates := res.Templates
			report := [][]string{}

			fmt.Printf(" \nShowing %d of %d templates\n \n", len(templates), len(templates))
			report = append(report, []string{"ID", "NAME", "TAGS"})

			for _, p := range templates {
				report = append(report, []string{
					fmt.Sprintf("%v", p.GetId()),
					p.GetName(),
					strings.Join(p.GetTags(), ","),
				})
			}
			printer.Table(os.Stdout, report)

			fmt.Println("\nFor details on a template, try: siren template view <name>")
			return nil
		},
	}

	cmd.Flags().StringVar(&tag, "tag", "", "template tag name")

	return cmd
}

func upsertTemplateCmd(c *configuration) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "Create or edit a new template",
		Long: heredoc.Doc(`
			Create or edit a new template.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var templateConfig domain.Template
			if err := parseFile(filePath, &templateConfig); err != nil {
				return err
			}

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			variables := make([]*sirenv1.TemplateVariables, 0)
			for _, variable := range templateConfig.Variables {
				variables = append(variables, &sirenv1.TemplateVariables{
					Name:        variable.Name,
					Type:        variable.Type,
					Default:     variable.Default,
					Description: variable.Description,
				})
			}

			res, err := client.UpsertTemplate(ctx, &sirenv1.UpsertTemplateRequest{
				Name:      templateConfig.Name,
				Body:      templateConfig.Body,
				Tags:      templateConfig.Tags,
				Variables: variables,
			})

			if err != nil {
				return err
			}

			fmt.Printf("template created with id: %v\n", res.GetTemplate().GetId())

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to the template config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func getTemplateCmd(c *configuration) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View a template details",
		Long: heredoc.Doc(`
			View a template.

			Display the Id, name, and other information about a template.
		`),
		Example: heredoc.Doc(`
			$ siren template view <template_name>
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

			name := args[0]
			res, err := client.GetTemplateByName(ctx, &sirenv1.GetTemplateByNameRequest{
				Name: name,
			})
			if err != nil {
				return err
			}

			templateData := res.Template

			variables := make([]domain.Variable, 0)
			for _, variable := range templateData.GetVariables() {
				variables = append(variables, domain.Variable{
					Name:        variable.Name,
					Type:        variable.Type,
					Default:     variable.Default,
					Description: variable.Description,
				})
			}

			template := &domain.Template{
				ID:        uint(templateData.GetId()),
				Name:      templateData.GetName(),
				Body:      templateData.GetBody(),
				Tags:      templateData.GetTags(),
				Variables: variables,
				CreatedAt: templateData.CreatedAt.AsTime(),
				UpdatedAt: templateData.UpdatedAt.AsTime(),
			}

			if err := printer.Text(template, format); err != nil {
				return fmt.Errorf("failed to format template: %v", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "yaml", "Print output with the selected format")

	return cmd
}

func deleteTemplateCmd(c *configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a template details",
		Example: heredoc.Doc(`
			$ siren template delete 1
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

			name := args[0]
			_, err = client.DeleteTemplate(ctx, &sirenv1.DeleteTemplateRequest{
				Name: name,
			})
			if err != nil {
				return err
			}

			fmt.Println("Successfully deleted template")
			return nil
		},
	}

	return cmd
}

func renderTemplateCmd(c *configuration) *cobra.Command {
	var name string
	var filePath string
	var format string
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render a template details",

		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var variableConfig struct {
				Variables map[string]string
			}
			if err := parseFile(filePath, &variableConfig); err != nil {
				return err
			}

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			template, err := client.RenderTemplate(ctx, &sirenv1.RenderTemplateRequest{
				Name:      name,
				Variables: variableConfig.Variables,
			})
			if err != nil {
				return err
			}

			if err := printer.Text(template, format); err != nil {
				return fmt.Errorf("failed to format template: %v", err)
			}
			return nil

		},
	}

	cmd.Flags().StringVar(&name, "name", "", "template name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to the template config")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVar(&format, "format", "yaml", "Print output with the selected format")

	return cmd
}
