package uploader

import (
	"context"
	"errors"
	"fmt"
	"github.com/antihax/optional"
	"github.com/odpf/siren/client"
	"github.com/odpf/siren/domain"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"strings"
)

type RulesAPICaller interface {
	CreateRuleRequest(ctx context.Context, localVarOptionals *client.RulesApiCreateRuleRequestOpts) (client.Rule, *http.Response, error)
	ListRulesRequest(ctx context.Context, localVarOptionals *client.RulesApiListRulesRequestOpts) ([]client.Rule, *http.Response, error)
}

type TemplatesAPICaller interface {
	CreateTemplateRequest(ctx context.Context, localVarOptionals *client.TemplatesApiCreateTemplateRequestOpts) (client.Template, *http.Response, error)
}

type SirenClient struct {
	client       *client.APIClient
	RulesAPI     RulesAPICaller
	TemplatesAPI TemplatesAPICaller
}

func NewSirenClient(host string) *SirenClient {
	cfg := &client.Configuration{
		BasePath: host,
	}
	apiClient := client.NewAPIClient(cfg)
	return &SirenClient{
		client:       apiClient,
		RulesAPI:     apiClient.RulesApi,
		TemplatesAPI: apiClient.TemplatesApi,
	}
}

type Uploader interface {
	Upload(string) (interface{}, error)
	UploadTemplates([]byte) (client.Template, error)
	UploadRules([]byte) ([]*client.Rule, error)
}

//Service talks to siren's HTTP Client
type Service struct {
	SirenClient *SirenClient
}

func NewService(c *domain.SirenServiceConfig) *Service {
	return &Service{
		SirenClient: NewSirenClient(c.Host),
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
		return s.UploadTemplates(yamlFile)
	} else if strings.ToLower(y.Type) == "rule" {
		return s.UploadRules(yamlFile)
	} else {
		return nil, errors.New("unknown type given")
	}
}

func (s Service) UploadTemplates(yamlFile []byte) (*client.Template, error) {
	var t template
	err := yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		return nil, err
	}
	body, err := yaml.Marshal(t.Body)
	if err != nil {
		return nil, err
	}
	payload := domain.Template{
		Body:      string(body),
		Name:      t.Name,
		Variables: t.Variables,
		Tags:      t.Tags,
	}
	options := &client.TemplatesApiCreateTemplateRequestOpts{
		Body: optional.NewInterface(payload),
	}
	result, _, err := s.SirenClient.TemplatesAPI.CreateTemplateRequest(context.Background(), options)
	if err != nil {
		return nil, err
	}
	//update associated rules for this template
	associatedRules, _, err := s.SirenClient.RulesAPI.ListRulesRequest(context.Background(), &client.RulesApiListRulesRequestOpts{
		Template: optional.NewString(t.Name),
	})
	if err != nil {
		return &result, err
	}
	for i := 0; i < len(associatedRules); i++ {
		associatedRule := associatedRules[i]
		var updatedVariables []domain.RuleVariable
		for j := 0; j < len(associatedRules[i].Variables); j++ {
			ruleVar := domain.RuleVariable{
				Name:        associatedRules[i].Variables[j].Name,
				Value:       associatedRules[i].Variables[j].Value,
				Type:        associatedRules[i].Variables[j].Type_,
				Description: associatedRules[i].Variables[j].Description,
			}
			updatedVariables = append(updatedVariables, ruleVar)
		}
		updateRulePayload := domain.Rule{
			Namespace: associatedRule.Namespace,
			Entity:    associatedRule.Entity,
			GroupName: associatedRule.GroupName,
			Template:  associatedRule.Template,
			Status:    associatedRule.Status,
			Variables: updatedVariables,
		}
		updateOptions := &client.RulesApiCreateRuleRequestOpts{
			Body: optional.NewInterface(updateRulePayload),
		}
		_, _, err := s.SirenClient.RulesAPI.CreateRuleRequest(context.Background(), updateOptions)
		if err != nil {
			fmt.Println("failed to update rule of ID: ", associatedRule.Id, "\tname: ", associatedRule.Name)
			return &result, err
		} else {
			fmt.Println("successfully updated rule of ID: ", associatedRule.Id, "\tname: ", associatedRule.Name)
		}
	}
	return &result, nil
}

func (s Service) UploadRules(yamlFile []byte) ([]*client.Rule, error) {
	var yamlBody ruleYaml
	err := yaml.Unmarshal(yamlFile, &yamlBody)
	if err != nil {
		return nil, err
	}
	var successfullyUpsertedRules []*client.Rule

	for k, v := range yamlBody.Rules {
		var vars []domain.RuleVariable
		for i := 0; i < len(v.Variables); i++ {
			v := domain.RuleVariable{
				Name:  v.Variables[i].Name,
				Value: v.Variables[i].Value,
			}
			vars = append(vars, v)
		}
		payload := domain.Rule{
			Namespace: yamlBody.Namespace,
			Entity:    yamlBody.Entity,
			GroupName: k,
			Template:  v.Template,
			Status:    v.Status,
			Variables: vars,
		}
		options := &client.RulesApiCreateRuleRequestOpts{
			Body: optional.NewInterface(payload),
		}
		result, _, err := s.SirenClient.RulesAPI.CreateRuleRequest(context.Background(), options)
		if err != nil {
			fmt.Println(fmt.Sprintf("rule %s/%s/%s/%s upload error",
				payload.Namespace, payload.Entity, payload.GroupName, payload.Template), err)
			return successfullyUpsertedRules, err
		} else {
			successfullyUpsertedRules = append(successfullyUpsertedRules, &result)
			fmt.Println(fmt.Sprintf("successfully uploaded %s/%s/%s/%s",
				payload.Namespace, payload.Entity, payload.GroupName, payload.Template))
		}
	}
	return successfullyUpsertedRules, nil
}
