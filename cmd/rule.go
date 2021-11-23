package cmd

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	sirenv1beta1 "github.com/odpf/siren/api/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/domain"
	"github.com/spf13/cobra"
)

func rulesCmd(c *configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rule",
		Aliases: []string{"rules"},
		Short:   "Manage rules",
		Long: heredoc.Doc(`
			Work with rules.
			
			rules are used for alerting within a provider.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
	}

	cmd.AddCommand(listRulesCmd(c))
	cmd.AddCommand(updateRuleCmd(c))
	cmd.AddCommand(uploadRuleCmd(c))

	return cmd
}

func listRulesCmd(c *configuration) *cobra.Command {
	var name string
	var namespace string
	var groupName string
	var template string
	var providerNamespace uint64
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List rules",
		Long: heredoc.Doc(`
			List all rules.
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

			res, err := client.ListRules(ctx, &sirenv1beta1.ListRulesRequest{
				Name:              name,
				GroupName:         groupName,
				Namespace:         namespace,
				ProviderNamespace: providerNamespace,
			})
			if err != nil {
				return err
			}

			rules := res.Rules
			report := [][]string{}

			fmt.Printf(" \nShowing %d of %d rules\n \n", len(rules), len(rules))
			report = append(report, []string{"ID", "NAME", "GROUP_NAME", "TEMPLATE", "ENABLED"})

			for _, p := range rules {
				report = append(report, []string{
					fmt.Sprintf("%v", p.GetId()),
					p.GetName(),
					p.GetGroupName(),
					p.GetTemplate(),
					strconv.FormatBool(p.GetEnabled()),
				})
			}
			printer.Table(os.Stdout, report)

			fmt.Println("\nFor details on a rule, try: siren rule view <id>")
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "rule name")
	cmd.Flags().StringVar(&namespace, "namespace", "", "rule namespace")
	cmd.Flags().StringVar(&groupName, "group-name", "", "rule group name")
	cmd.Flags().StringVar(&template, "template", "", "rule template")
	cmd.Flags().Uint64Var(&providerNamespace, "provider-namespace", 0, "rule provider namespace id")

	return cmd
}

func updateRuleCmd(c *configuration) *cobra.Command {
	var id uint64
	var filePath string
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a rule",
		Long: heredoc.Doc(`
			Edit an existing rule.
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var ruleConfig domain.Rule
			if err := parseFile(filePath, &ruleConfig); err != nil {
				return err
			}

			ctx := context.Background()
			client, cancel, err := createClient(ctx, c.Host)
			if err != nil {
				return err
			}
			defer cancel()

			variables := make([]*sirenv1beta1.Variables, 0)
			for _, variable := range ruleConfig.Variables {
				variables = append(variables, &sirenv1beta1.Variables{
					Name:        variable.Name,
					Type:        variable.Type,
					Value:       variable.Value,
					Description: variable.Description,
				})
			}

			_, err = client.UpdateRule(ctx, &sirenv1beta1.UpdateRuleRequest{
				Enabled:           ruleConfig.Enabled,
				GroupName:         ruleConfig.GroupName,
				Namespace:         ruleConfig.Namespace,
				Template:          ruleConfig.Template,
				Variables:         variables,
				ProviderNamespace: ruleConfig.ProviderNamespace,
			})
			if err != nil {
				return err
			}

			fmt.Println("Successfully updated rule")

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "rule id")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the rule config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func uploadRuleCmd(c *configuration) *cobra.Command {
	var fileReader = ioutil.ReadFile
	return &cobra.Command{
		Use:   "upload",
		Short: "Upload Rules YAML file",
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

			if strings.ToLower(yamlObject.Type) == "rule" {
				result, err := uploadRule(client, yamlFile)
				if err != nil {
					return err
				}
				printRules(result)
			} else {
				return errors.New("yaml is not rule type")
			}
			return nil
		},
	}
}

func uploadRule(client sirenv1beta1.SirenServiceClient, yamlFile []byte) ([]*sirenv1beta1.Rule, error) {
	var yamlBody ruleYaml
	err := yaml.Unmarshal(yamlFile, &yamlBody)
	if err != nil {
		return nil, err
	}
	var successfullyUpsertedRules []*sirenv1beta1.Rule

	for groupName, v := range yamlBody.Rules {
		var ruleVariables []*sirenv1beta1.Variables
		for i := 0; i < len(v.Variables); i++ {
			v := &sirenv1beta1.Variables{
				Name:  v.Variables[i].Name,
				Value: v.Variables[i].Value,
			}
			ruleVariables = append(ruleVariables, v)
		}

		if yamlBody.ProviderNamespace == "" {
			return nil, errors.New("provider namespace is required")
		}

		data, err := client.ListProviders(context.Background(), &sirenv1beta1.ListProvidersRequest{
			Urn: yamlBody.ProviderNamespace,
		})
		if err != nil {
			return nil, err
		}

		provideres := data.Providers
		if len(provideres) == 0 {
			return nil, errors.New(fmt.Sprintf("no provider found with urn: %s", yamlBody.ProviderNamespace))
		}

		payload := &sirenv1beta1.UpdateRuleRequest{
			GroupName:         groupName,
			Namespace:         yamlBody.Namespace,
			Template:          v.Template,
			Variables:         ruleVariables,
			ProviderNamespace: provideres[0].Id,
			Enabled:           v.Enabled,
		}

		result, err := client.UpdateRule(context.Background(), payload)
		if err != nil {
			fmt.Println(fmt.Sprintf("rule %s/%s/%s upload error",
				payload.Namespace, payload.GroupName, payload.Template), err)
			return successfullyUpsertedRules, err
		} else {
			successfullyUpsertedRules = append(successfullyUpsertedRules, result.Rule)
			fmt.Println(fmt.Sprintf("successfully uploaded %s/%s/%s",
				payload.Namespace, payload.GroupName, payload.Template))
		}
	}
	return successfullyUpsertedRules, nil
}

func printRules(rules []*sirenv1beta1.Rule) {
	for i := 0; i < len(rules); i++ {
		fmt.Println("Upserted Rule")
		fmt.Println("ID:", rules[i].Id)
		fmt.Println("Name:", rules[i].Name)
		fmt.Println("Enabled:", rules[i].Enabled)
		fmt.Println("Group Name:", rules[i].GroupName)
		fmt.Println("Namespace:", rules[i].Namespace)
		fmt.Println("Template:", rules[i].Template)
		fmt.Println("CreatedAt At:", rules[i].CreatedAt)
		fmt.Println("UpdatedAt At:", rules[i].UpdatedAt)
		fmt.Println()
	}
}