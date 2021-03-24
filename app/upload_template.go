package app

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
)

type variables struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type rule struct {
	Template  string      `json:"template"`
	Status    string      `json:"enabled"`
	Variables []variables `json:"variables"`
}

type ruleYaml struct {
	ApiVersion string          `json:"apiVersion"`
	Type       string          `json:"type"`
	Namespace  string          `json:"namespace"`
	Rules      map[string]rule `json:"rules"`
}

type templatedRule struct {
	Alert       string            `json:"alert"`
	Expr        string            `json:"expr"`
	For         string            `json:"for"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type template struct {
	Name       string            `json:"name"`
	ApiVersion string            `json:"apiVersion"`
	Type       string            `json:"type"`
	Body       []templatedRule   `json:"body"`
	Tags       []string          `json:"tags"`
	Variables  []domain.Variable `json:"variables"`
}

func UploadTemplates(c *domain.Config, fileName string) error {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return err
	}
	var t template
	err = yaml.Unmarshal(yamlFile, &t)
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
	cfg := &client.Configuration{
		BasePath: c.SirenService.Host,
	}
	payload := domain.Template{
		Body:      string(body),
		Name:      t.Name,
		Variables: t.Variables,
		Tags:      t.Tags,
	}
	sirenClient := client.NewAPIClient(cfg)
	options := &client.TemplatesApiCreateTemplateRequestOpts{
		Body: optional.NewInterface(payload),
	}
	result, _, err := sirenClient.TemplatesApi.CreateTemplateRequest(context.Background(), options)
	if err != nil {
		fmt.Println(result)
	}
	response, _ := json.Marshal(result)
	fmt.Println(string(response))
	return nil
}

func UploadRules(c *domain.Config, fileName string, entity string) error {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return err
	}
	var yamlBody ruleYaml
	err = yaml.Unmarshal(yamlFile, &yamlBody)
	if err != nil {
		return err
	}
	if yamlBody.Type != "rule" {
		return errors.New("object was not of rule type")
	}
	cfg := &client.Configuration{
		BasePath: c.SirenService.Host,
	}
	sirenClient := client.NewAPIClient(cfg)

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
			Entity:    entity,
			GroupName: k,
			Template:  v.Template,
			Status:    v.Status,
			Variables: vars,
		}
		options := &client.RulesApiCreateRuleRequestOpts{
			Body: optional.NewInterface(payload),
		}
		result, _, err := sirenClient.RulesApi.CreateRuleRequest(context.Background(), options)
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
