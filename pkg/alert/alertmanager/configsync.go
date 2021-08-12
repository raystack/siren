package alertmanager

import (
	"bytes"
	"context"
	"text/template"

	_ "embed"

	"github.com/grafana/cortex-tools/pkg/client"
	"github.com/odpf/siren/domain"
	"github.com/prometheus/alertmanager/config"
)

var (
	//go:embed helper.tmpl
	helperTemplateString string
	//go:embed config.goyaml
	configYamlString string
)

type SlackCredential struct {
	Channel  string
}

type SlackConfig struct {
	Critical SlackCredential
	Warning  SlackCredential
}

type TeamCredentials struct {
	PagerdutyCredential string
	Slackcredentials    SlackConfig
	Name                string
}

type EntityCredentials struct {
	Entity string
	Token  string
	Teams  map[string]TeamCredentials
}

type AlertManagerConfig struct {
	EntityCredentials EntityCredentials
	AlertHistoryHost  string
}

type Client interface {
	SyncConfig(credentials AlertManagerConfig) error
}

type AlertmanagerClient struct {
	CortextClient  client.CortexClient
	helperTemplate string
}

func NewClient(c domain.CortexConfig) (AlertmanagerClient, error) {
	config := client.Config{
		Address: c.Address,
	}
	amClient, err := client.New(config)
	if err != nil {
		return AlertmanagerClient{}, err
	}

	if err != nil {
		return AlertmanagerClient{}, err
	}
	return AlertmanagerClient{
		CortextClient:  *amClient,
		helperTemplate: helperTemplateString,
	}, nil
}

func (am AlertmanagerClient) SyncConfig(config AlertManagerConfig) error {
	cfg, err := generateAlertmanagerConfig(config)
	if err != nil {
		return err
	}
	templates := map[string]string{
		"helper.tmpl": am.helperTemplate,
	}

	ctx := client.NewContextWithTenantId(context.Background(), config.EntityCredentials.Entity)
	err = am.CortextClient.CreateAlertmanagerConfig(ctx, cfg, templates)
	if err != nil {
		return err
	}
	return nil
}

func generateAlertmanagerConfig(alertManagerConfig AlertManagerConfig) (string, error) {
	delims := template.New("alertmanagerConfigTemplate").Delims("[[", "]]")
	parse, err := delims.Parse(configYamlString)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = parse.Execute(&tpl, alertManagerConfig)
	if err != nil {
		return "", err
	}
	configStr := tpl.String()
	_, err = config.Load(configStr)
	if err != nil {
		return "", err
	}
	return configStr, nil
}
