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
	//go:embed alertmanagerde.tmpl
	deTmplateString string
	//go:embed alertmanagervar.tmpl
	varTmplateString string
	//go:embed alertmanagerconfig.goyaml
	configYamlString string
)

type SlackCredential struct {
	Webhook  string
	Channel  string
	Username string
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
	Teams  map[string]TeamCredentials
}

type Client interface {
	SyncConfig(credentials EntityCredentials) error
}

type AlertmanagerClient struct {
	CortextClient client.CortexClient
	vartmplStr    string
	detmplStr     string
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
		CortextClient: *amClient,
		detmplStr:     deTmplateString,
		vartmplStr:    varTmplateString,
	}, nil
}

func (am AlertmanagerClient) SyncConfig(credentials EntityCredentials) error {
	cfg, err := generateAlertmanagerConfig(credentials)
	if err != nil {
		return err
	}
	templates := map[string]string{
		"var.tmpl": am.vartmplStr,
		"de.tmpl":  am.detmplStr,
	}

	ctx := client.NewContextWithTenantId(context.Background(), credentials.Entity)
	err = am.CortextClient.CreateAlertmanagerConfig(ctx, cfg, templates)
	if err != nil {
		return err
	}
	return nil
}

func generateAlertmanagerConfig(credentials EntityCredentials) (string, error) {
	delims := template.New("alertmanagerConfigTemplate").Delims("[[", "]]")
	parse, err := delims.Parse(configYamlString)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = parse.Execute(&tpl, credentials)
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
