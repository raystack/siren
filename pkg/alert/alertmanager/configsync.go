package alertmanager

import (
	"bytes"
	"context"
	"github.com/grafana/cortex-tools/pkg/client"
	"github.com/markbates/pkger"
	"github.com/odpf/siren/domain"
	"github.com/prometheus/alertmanager/config"
	"text/template"
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

	deTemplatePath := "/pkg/alert/alertmanagerde.tmpl"
	deTmplateString, err := readTemplateString(err, deTemplatePath)
	if err != nil {
		return AlertmanagerClient{}, err
	}
	varTmplPath := "/pkg/alert/alertmanagervar.tmpl"
	varTmplateString, err := readTemplateString(err, varTmplPath)
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
	configYaml, err := pkger.Open("/pkg/alert/alertmanagerconfig.goyaml")
	if err != nil {
		return "", err
	}
	defer configYaml.Close()
	configYamlBuf := new(bytes.Buffer)
	configYamlBuf.ReadFrom(configYaml)
	delims := template.New("alertmanagerConfigTemplate").Delims("[[", "]]")
	parse, err := delims.Parse(configYamlBuf.String())
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

func readTemplateString(err error, templatePath string) (string, error) {
	tmpl, err := pkger.Open(templatePath)
	if err != nil {
		return "", err
	}
	defer tmpl.Close()
	detmplBuf := new(bytes.Buffer)
	_, err = detmplBuf.ReadFrom(tmpl)
	if err != nil {
		return "", err
	}
	return detmplBuf.String(), nil
}
