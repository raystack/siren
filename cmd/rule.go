package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
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

			res, err := client.ListRules(ctx, &sirenv1.ListRulesRequest{
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

			variables := make([]*sirenv1.Variables, 0)
			for _, variable := range ruleConfig.Variables {
				variables = append(variables, &sirenv1.Variables{
					Name:        variable.Name,
					Type:        variable.Type,
					Value:       variable.Value,
					Description: variable.Description,
				})
			}

			_, err = client.UpdateRule(ctx, &sirenv1.UpdateRuleRequest{
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
