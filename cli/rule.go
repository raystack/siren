package cli

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/odpf/siren/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/cmdx"
	"github.com/odpf/salt/printer"
	"github.com/odpf/siren/core/rule"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"github.com/spf13/cobra"
)

func rulesCmd(cmdxConfig *cmdx.Config) *cobra.Command {
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
			"client":     "true",
		},
	}

	cmd.AddCommand(
		listRulesCmd(cmdxConfig),
		updateRuleCmd(cmdxConfig),
		uploadRuleCmd(cmdxConfig),
	)

	return cmd
}

func listRulesCmd(cmdxConfig *cmdx.Config) *cobra.Command {
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

			res, err := client.ListRules(ctx, &sirenv1beta1.ListRulesRequest{
				Name:              name,
				GroupName:         groupName,
				Namespace:         namespace,
				ProviderNamespace: providerNamespace,
			})
			if err != nil {
				return err
			}

			if res.GetRules() == nil {
				return errors.New("no response from server")
			}

			spinner.Stop()
			rules := res.GetRules()
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

func updateRuleCmd(cmdxConfig *cmdx.Config) *cobra.Command {
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
			spinner := printer.Spin("")
			defer spinner.Stop()

			ctx := cmd.Context()

			c, err := loadClientConfig(cmd, cmdxConfig)
			if err != nil {
				return err
			}

			var ruleConfig rule.Rule
			if err := parseFile(filePath, &ruleConfig); err != nil {
				return err
			}

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

			spinner.Stop()
			printer.Success("Successfully updated rule")
			printer.Space()
			printer.SuccessIcon()

			return nil
		},
	}

	cmd.Flags().Uint64Var(&id, "id", 0, "rule id")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the rule config")
	cmd.MarkFlagRequired("file")

	return cmd
}

func uploadRuleCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	var fileReader = os.ReadFile
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload Rules YAML file",
		Annotations: map[string]string{
			"group:core": "true",
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
				rulesID, err := UploadRules(client, yamlFile)
				if err != nil {
					return err
				}

				spinner.Stop()
				//TODO might need to print the actual rule here or log error rules
				printRulesID(rulesID)
				return nil
			}
			return errors.New("yaml is not rule type")
		},
	}

	return cmd
}

func UploadRules(client sirenv1beta1.SirenServiceClient, yamlFile []byte) ([]uint64, error) {
	var yamlBody rule.RuleFile
	err := yaml.Unmarshal(yamlFile, &yamlBody)
	if err != nil {
		return nil, err
	}
	var successfullyUpsertedRulesID []uint64

	for groupName, v := range yamlBody.Rules {
		var ruleVariables []*sirenv1beta1.Variables
		for i := 0; i < len(v.Variables); i++ {
			v := &sirenv1beta1.Variables{
				Name:  v.Variables[i].Name,
				Value: v.Variables[i].Value,
			}
			ruleVariables = append(ruleVariables, v)
		}

		if yamlBody.Provider == "" {
			return nil, errors.New("provider is required")
		}

		if yamlBody.ProviderNamespace == "" {
			return nil, errors.New("provider namespace is required")
		}

		providersData, err := client.ListProviders(context.Background(), &sirenv1beta1.ListProvidersRequest{
			Urn: yamlBody.Provider,
		})
		if err != nil {
			return nil, err
		}

		if providersData.GetProviders() == nil {
			return nil, errors.New("provider not found")
		}

		var provider *sirenv1beta1.Provider
		providers := providersData.GetProviders()
		if len(providers) != 0 {
			provider = providers[0]
		} else {
			return nil, errors.New("provider not found")
		}

		res, err := client.ListNamespaces(context.Background(), &sirenv1beta1.ListNamespacesRequest{})
		if err != nil {
			return nil, err
		}

		if res.GetNamespaces() == nil {
			return nil, errors.New("no response of getting list of namespaces from server")
		}

		var providerNamespace *sirenv1beta1.Namespace
		for _, ns := range res.GetNamespaces() {
			if ns.GetUrn() == yamlBody.ProviderNamespace && ns.Provider == provider.Id {
				providerNamespace = ns
				break
			}
		}

		if providerNamespace == nil {
			return nil, fmt.Errorf("no namespace found with urn: %s under provider %s", yamlBody.ProviderNamespace, provider.Name)
		}

		payload := &sirenv1beta1.UpdateRuleRequest{
			GroupName:         groupName,
			Namespace:         yamlBody.Namespace,
			Template:          v.Template,
			Variables:         ruleVariables,
			ProviderNamespace: providerNamespace.Id,
			Enabled:           v.Enabled,
		}

		result, err := client.UpdateRule(context.Background(), payload)
		if err != nil {
			fmt.Println(fmt.Sprintf("rule %s/%s/%s upload error",
				payload.Namespace, payload.GroupName, payload.Template), err)
			return successfullyUpsertedRulesID, err
		}
		successfullyUpsertedRulesID = append(successfullyUpsertedRulesID, result.GetId())
		fmt.Printf("successfully uploaded %s/%s/%s",
			payload.Namespace, payload.GroupName, payload.Template)

	}
	return successfullyUpsertedRulesID, nil
}

func printRulesID(rulesID []uint64) {
	for i := 0; i < len(rulesID); i++ {
		fmt.Println("Upserted Rule")
		fmt.Println("ID:", rulesID[i])
		fmt.Println()
	}
}
