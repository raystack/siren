package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/odpf/siren/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	"github.com/odpf/siren/core/template"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"github.com/spf13/cobra"
)

type templatedRule struct {
	Record      string            `yaml:"record,omitempty"`
	Alert       string            `yaml:"alert,omitempty"`
	Expr        string            `yaml:"expr"`
	For         string            `yaml:"for,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type templateStruct struct {
	Name       string              `yaml:"name"`
	ApiVersion string              `yaml:"apiVersion"`
	Type       string              `yaml:"type"`
	Body       []templatedRule     `yaml:"body"`
	Tags       []string            `yaml:"tags"`
	Variables  []template.Variable `yaml:"variables"`
}

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
	cmd.AddCommand(uploadTemplateCmd(c))

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

			res, err := client.ListTemplates(ctx, &sirenv1beta1.ListTemplatesRequest{
				Tag: tag,
			})
			if err != nil {
				return err
			}

			if res.GetTemplates() == nil {
				return errors.New("no response from server")
			}

			templates := res.GetTemplates()
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
			var templateConfig template.Template
			if err := parseFile(filePath, &templateConfig); err != nil {
				return err
			}

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			variables := make([]*sirenv1beta1.TemplateVariables, 0)
			for _, variable := range templateConfig.Variables {
				variables = append(variables, &sirenv1beta1.TemplateVariables{
					Name:        variable.Name,
					Type:        variable.Type,
					Default:     variable.Default,
					Description: variable.Description,
				})
			}

			res, err := client.UpsertTemplate(ctx, &sirenv1beta1.UpsertTemplateRequest{
				Name:      templateConfig.Name,
				Body:      templateConfig.Body,
				Tags:      templateConfig.Tags,
				Variables: variables,
			})

			if err != nil {
				return err
			}

			fmt.Printf("template created with id: %v\n", res.GetId())

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
			res, err := client.GetTemplate(ctx, &sirenv1beta1.GetTemplateRequest{
				Name: name,
			})
			if err != nil {
				return err
			}

			if res.GetTemplate() == nil {
				return errors.New("no response from server")
			}

			templateData := res.GetTemplate()

			variables := make([]template.Variable, 0)
			for _, variable := range templateData.GetVariables() {
				variables = append(variables, template.Variable{
					Name:        variable.Name,
					Type:        variable.Type,
					Default:     variable.Default,
					Description: variable.Description,
				})
			}

			template := &template.Template{
				ID:        templateData.GetId(),
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
			_, err = client.DeleteTemplate(ctx, &sirenv1beta1.DeleteTemplateRequest{
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

			template, err := client.RenderTemplate(ctx, &sirenv1beta1.RenderTemplateRequest{
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

func uploadTemplateCmd(c *configuration) *cobra.Command {
	var fileReader = ioutil.ReadFile
	return &cobra.Command{
		Use:   "upload",
		Short: "Upload Templates YAML file",
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

			yamlFile, err := fileReader(args[0])
			if err != nil {
				fmt.Printf("Error reading YAML file: %s\n", err)
				return err
			}

			var yamlObject struct {
				Type string `yaml:"type"`
			}
			err = yaml.Unmarshal(yamlFile, &yamlObject)
			if err != nil {
				return err
			}

			if strings.ToLower(yamlObject.Type) == "template" {
				templateID, err := uploadTemplate(client, yamlFile)
				if err != nil {
					return err
				}
				//TODO might need to log the actual template or log error here
				printTemplateID(templateID)
				return nil
			}
			return errors.New("yaml is not rule type")
		},
	}
}

func uploadTemplate(client sirenv1beta1.SirenServiceClient, yamlFile []byte) (uint64, error) {
	var t templateStruct
	err := yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		return 0, err
	}
	body, err := yaml.Marshal(t.Body)
	if err != nil {
		return 0, err
	}

	variables := make([]*sirenv1beta1.TemplateVariables, 0)
	for _, variable := range t.Variables {
		variables = append(variables, &sirenv1beta1.TemplateVariables{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}

	template, err := client.UpsertTemplate(context.Background(), &sirenv1beta1.UpsertTemplateRequest{
		Name:      t.Name,
		Body:      string(body),
		Variables: variables,
		Tags:      t.Tags,
	})
	if err != nil {
		return 0, err
	}

	//update associated rules for this template
	data, err := client.ListRules(context.Background(), &sirenv1beta1.ListRulesRequest{
		Template: t.Name,
	})
	if err != nil {
		return 0, err
	}

	if data.GetRules() == nil {
		return 0, errors.New("no response from server")
	}

	associatedRules := data.GetRules()
	for i := 0; i < len(associatedRules); i++ {
		associatedRule := associatedRules[i]

		var updatedVariables []*sirenv1beta1.Variables
		for j := 0; j < len(associatedRules[i].Variables); j++ {
			ruleVar := &sirenv1beta1.Variables{
				Name:        associatedRules[i].Variables[j].Name,
				Value:       associatedRules[i].Variables[j].Value,
				Type:        associatedRules[i].Variables[j].Type,
				Description: associatedRules[i].Variables[j].Description,
			}
			updatedVariables = append(updatedVariables, ruleVar)
		}

		_, err := client.UpdateRule(context.Background(), &sirenv1beta1.UpdateRuleRequest{
			GroupName:         associatedRule.GroupName,
			Namespace:         associatedRule.Namespace,
			Template:          associatedRule.Template,
			Variables:         updatedVariables,
			ProviderNamespace: associatedRule.ProviderNamespace,
			Enabled:           associatedRule.Enabled,
		})

		if err != nil {
			fmt.Println("failed to update rule of ID: ", associatedRule.Id, "\tname: ", associatedRule.Name)
			return 0, err
		}
		fmt.Println("successfully updated rule of ID: ", associatedRule.Id, "\tname: ", associatedRule.Name)
	}
	return template.GetId(), nil
}

func printTemplateID(templateID uint64) {
	fmt.Println("Upserted Template")
	fmt.Println("ID:", templateID)
}
