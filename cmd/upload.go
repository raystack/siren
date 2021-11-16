package cmd

import (
	"context"
	"errors"
	"fmt"
	sirenv1beta1 "github.com/odpf/siren/api/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/domain"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

type variables struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type rule struct {
	Template  string      `yaml:"template"`
	Enabled   bool        `yaml:"enabled"`
	Variables []variables `yaml:"variables"`
}

type ruleYaml struct {
	ApiVersion        string          `yaml:"apiVersion"`
	Entity            string          `yaml:"entity"`
	Type              string          `yaml:"type"`
	Namespace         string          `yaml:"namespace"`
	ProviderNamespace string          `yaml:"providerNamespace"`
	Rules             map[string]rule `yaml:"rules"`
}

type templatedRule struct {
	Record      string            `yaml:"record,omitempty"`
	Alert       string            `yaml:"alert,omitempty"`
	Expr        string            `yaml:"expr"`
	For         string            `yaml:"for,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type template struct {
	Name       string            `yaml:"name"`
	ApiVersion string            `yaml:"apiVersion"`
	Type       string            `yaml:"type"`
	Body       []templatedRule   `yaml:"body"`
	Tags       []string          `yaml:"tags"`
	Variables  []domain.Variable `yaml:"variables"`
}

type yamlObject struct {
	Type string `yaml:"type"`
}

func uploadCmd(c *configuration) *cobra.Command {
	return &cobra.Command{
		Use:   "upload",
		Short: "Upload Rules or Templates YAML file",
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

			s := UploaderService(client)
			result, err := s.Upload(args[0])

			//print all resources(succeed or failed in upsert)
			if err != nil {
				fmt.Println(err)
				return err
			}
			switch obj := result.(type) {
			case *sirenv1beta1.Template:
				printTemplate(obj)
			case []*sirenv1beta1.Rule:
				printRules(obj)
			default:
				return errors.New("unknown response")
			}
			return nil
		},
	}
}

//Service talks to siren's HTTP Client
type Service struct {
	SirenClient sirenv1beta1.SirenServiceClient
}

func UploaderService(siren sirenv1beta1.SirenServiceClient) *Service {
	return &Service{
		SirenClient: siren,
	}
}

var fileReader = ioutil.ReadFile

func (s Service) Upload(fileName string) (interface{}, error) {
	yamlFile, err := fileReader(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return nil, err
	}
	var y yamlObject
	err = yaml.Unmarshal(yamlFile, &y)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(y.Type) == "template" {
		return s.UploadTemplate(yamlFile)
	} else if strings.ToLower(y.Type) == "rule" {
		return s.UploadRule(yamlFile)
	} else {
		return nil, errors.New("unknown type given")
	}
}

func (s Service) UploadTemplate(yamlFile []byte) (*sirenv1beta1.Template, error) {
	var t template
	err := yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		return nil, err
	}
	body, err := yaml.Marshal(t.Body)
	if err != nil {
		return nil, err
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

	template, err := s.SirenClient.UpsertTemplate(context.Background(), &sirenv1beta1.UpsertTemplateRequest{
		Name:      t.Name,
		Body:      string(body),
		Variables: variables,
		Tags:      t.Tags,
	})
	if err != nil {
		return nil, err
	}

	//update associated rules for this template
	data, err := s.SirenClient.ListRules(context.Background(), &sirenv1beta1.ListRulesRequest{
		Template: t.Name,
	})
	if err != nil {
		return nil, err
	}

	associatedRules := data.Rules
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

		_, err := s.SirenClient.UpdateRule(context.Background(), &sirenv1beta1.UpdateRuleRequest{
			GroupName:         associatedRule.GroupName,
			Namespace:         associatedRule.Namespace,
			Template:          associatedRule.Template,
			Variables:         updatedVariables,
			ProviderNamespace: associatedRule.ProviderNamespace,
			Enabled:           associatedRule.Enabled,
		})

		if err != nil {
			fmt.Println("failed to update rule of ID: ", associatedRule.Id, "\tname: ", associatedRule.Name)
			return nil, err
		}
		fmt.Println("successfully updated rule of ID: ", associatedRule.Id, "\tname: ", associatedRule.Name)
	}
	return template.Template, nil
}

func (s Service) UploadRule(yamlFile []byte) ([]*sirenv1beta1.Rule, error) {
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

		data, err := s.SirenClient.ListProviders(context.Background(), &sirenv1beta1.ListProvidersRequest{
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

		result, err := s.SirenClient.UpdateRule(context.Background(), payload)
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

func printTemplate(template *sirenv1beta1.Template) {
	if template == nil {
		return
	}
	fmt.Println("Upserted Template")
	fmt.Println("ID:", template.Id)
	fmt.Println("Name:", template.Name)
	fmt.Println("Tags:", template.Tags)
	fmt.Println("Variables:", template.Variables)
	fmt.Println("CreatedAt At:", template.CreatedAt)
	fmt.Println("UpdatedAt At:", template.UpdatedAt)

}
