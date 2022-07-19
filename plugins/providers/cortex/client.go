package cortex

import (
	"bytes"
	"context"
	"html/template"

	"github.com/grafana/cortex-tools/pkg/client"
	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	promconfig "github.com/prometheus/alertmanager/config"
)

//go:generate mockery --name=CortexCaller -r --case underscore --with-expecter --structname CortexCaller --filename cortex_caller.go --output=./mocks
type CortexCaller interface {
	CreateAlertmanagerConfig(ctx context.Context, cfg string, templates map[string]string) error
	CreateRuleGroup(ctx context.Context, namespace string, rg rwrulefmt.RuleGroup) error
	DeleteRuleGroup(ctx context.Context, namespace, groupName string) error
	GetRuleGroup(ctx context.Context, namespace, groupName string) (*rwrulefmt.RuleGroup, error)
	ListRules(ctx context.Context, namespace string) (map[string][]rwrulefmt.RuleGroup, error)
	GetAlertmanagerConfig(ctx context.Context) (string, map[string]string, error)
}

type Client struct {
	cortexClient   CortexCaller
	helperTemplate string
	configYaml     string
}

func NewClient(cfg Config, opts ...ClientOption) (*Client, error) {
	c := &Client{
		helperTemplate: HelperTemplateString,
		configYaml:     ConfigYamlString,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.cortexClient == nil {
		cortexCfg := client.Config{
			Address: cfg.Address,
		}
		cortexClient, err := client.New(cortexCfg)
		if err != nil {
			return nil, err
		}
		c.cortexClient = cortexClient
	}

	return c, nil
}

// CreateAlertmanagerConfig uploads an alert manager config to cortex
func (c *Client) CreateAlertmanagerConfig(amConfigs AlertManagerConfig, tenantID string) error {
	cfg, err := c.generateAlertmanagerConfig(amConfigs)
	if err != nil {
		return err
	}
	templates := map[string]string{
		"helper.tmpl": c.helperTemplate,
	}

	newCtx := NewContext(context.Background(), tenantID)
	err = c.cortexClient.CreateAlertmanagerConfig(newCtx, cfg, templates)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) generateAlertmanagerConfig(amConfigs AlertManagerConfig) (string, error) {
	delims := template.New("alertmanagerConfigTemplate").Delims("[[", "]]")
	parse, err := delims.Parse(c.configYaml)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = parse.Execute(&tpl, amConfigs)
	if err != nil {
		// it is unlikely that the code returns error here
		return "", err
	}
	configStr := tpl.String()
	_, err = promconfig.Load(configStr)
	if err != nil {
		return "", err
	}
	return configStr, nil
}

func (c *Client) CreateRuleGroup(ctx context.Context, namespace string, rg rwrulefmt.RuleGroup) error {
	return c.cortexClient.CreateRuleGroup(ctx, namespace, rg)
}

func (c *Client) DeleteRuleGroup(ctx context.Context, namespace, groupName string) error {
	return c.cortexClient.DeleteRuleGroup(ctx, namespace, groupName)
}

func (c *Client) GetRuleGroup(ctx context.Context, namespace, groupName string) (*rwrulefmt.RuleGroup, error) {
	return c.cortexClient.GetRuleGroup(ctx, namespace, groupName)
}

func (c *Client) ListRules(ctx context.Context, namespace string) (map[string][]rwrulefmt.RuleGroup, error) {
	return c.cortexClient.ListRules(ctx, namespace)
}
