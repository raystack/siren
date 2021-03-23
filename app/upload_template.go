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
