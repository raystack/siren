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

type AMReceiverConfig struct {
	Receiver      string
	Type          string
	Match         map[string]string
	Configuration map[string]string
}

type AMConfig struct {
	Receivers []AMReceiverConfig
}

type Client interface {
	SyncConfig(AMConfig, string) error
}

type AlertmanagerClient struct {
	CortexClient   client.CortexClient
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
		CortexClient:   *amClient,
		helperTemplate: helperTemplateString,
	}, nil
}

func (am AlertmanagerClient) SyncConfig(config AMConfig, tenant string) error {
	cfg, err := generateAlertmanagerConfig(config)
	if err != nil {
		return err
	}
	templates := map[string]string{
		"helper.tmpl": am.helperTemplate,
	}

	ctx := client.NewContextWithTenantId(context.Background(), tenant)
	err = am.CortexClient.CreateAlertmanagerConfig(ctx, cfg, templates)
	if err != nil {
		return err
	}
	return nil
}

func generateAlertmanagerConfig(alertManagerConfig AMConfig) (string, error) {
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
