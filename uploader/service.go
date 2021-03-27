package uploader

import (
	"context"
	"encoding/json"
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
	Upload(string) error
	UploadTemplates([]byte) error
	UploadRules([]byte) error
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

func (s Service) Upload(fileName string) error {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return err
	}
	var y yamlObject
	err = yaml.Unmarshal(yamlFile, &y)
	if err != nil {
		return err
	}
	if strings.ToLower(y.Type) == "template" {
		return s.UploadTemplates(yamlFile)
	} else if strings.ToLower(y.Type) == "rule" {
		return s.UploadRules(yamlFile)
	} else {
		return errors.New("unknown type given")
	}
}

func (s Service) UploadTemplates(yamlFile []byte) error {
	var t template
	err := yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		return err
	}
	if t.Type != "template" {
		return errors.New("object was not of template type")
	}
	body, err := yaml.Marshal(t.Body)
	if err != nil {
		return err
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
		fmt.Println(result)
	}
	response, _ := json.Marshal(result)
	fmt.Println(string(response))
	return nil
}

func (s Service) UploadRules(yamlFile []byte) error {
	var yamlBody ruleYaml
	err := yaml.Unmarshal(yamlFile, &yamlBody)
	if err != nil {
		return err
	}
	if yamlBody.Type != "rule" {
		return errors.New("object was not of rule type")
	}

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
		response, _ := json.Marshal(result)
		fmt.Println(string(response))
		if err != nil {
			fmt.Println(fmt.Sprintf("rule %s/%s/%s/%s upload error",
				payload.Namespace, payload.Entity, payload.GroupName, payload.Template), err)
		} else {
			fmt.Println(fmt.Sprintf("successfully uploaded %s/%s/%s/%s",
				payload.Namespace, payload.Entity, payload.GroupName, payload.Template))
		}
	}
	return nil
}
