package cortex

import (
	"bytes"
	"context"
	"fmt"
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
}

// Client is a wrapper of cortex-tools client
type Client struct {
	cortexClient   CortexCaller
	helperTemplate string
	configYaml     string
}

// NewClient creates a new Client
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
// this function merges alertmanager template defined in config/helper.tmpl
// and a rendered alertmanager config template usually used in
// subscription flow
func (c *Client) CreateAlertmanagerConfig(ctx context.Context, amConfigs AlertManagerConfig, tenantID string) error {
	cfg, err := c.generateAlertmanagerConfig(amConfigs)
	if err != nil {
		return err
	}
	templates := map[string]string{
		"helper.tmpl": c.helperTemplate,
	}

	err = c.cortexClient.CreateAlertmanagerConfig(NewContextWithTenantID(ctx, tenantID), cfg, templates)
	if err != nil {
		return fmt.Errorf("cortex client: %w", err)
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

// CreateRuleGroup creates a rule group in a namespace in a tenant
// in cortex ruler. this will replace the existing rule group if exist
func (c *Client) CreateRuleGroup(ctx context.Context, tenantID string, namespace string, rg rwrulefmt.RuleGroup) error {
	err := c.cortexClient.CreateRuleGroup(NewContextWithTenantID(ctx, tenantID), namespace, rg)
	if err != nil {
		return fmt.Errorf("cortex client: %w", err)
	}
	return nil
}

// DeleteRuleGroup removes a rule group in a namespace in a tenant
// in cortex ruler
func (c *Client) DeleteRuleGroup(ctx context.Context, tenantID, namespace, groupName string) error {
	err := c.cortexClient.DeleteRuleGroup(NewContextWithTenantID(ctx, tenantID), namespace, groupName)
	if err != nil {
		return fmt.Errorf("cortex client: %w", err)
	}
	return nil
}

// GetRuleGroup fetchs a rule group in a namespace in a tenant
// in cortex ruler
func (c *Client) GetRuleGroup(ctx context.Context, tenantID, namespace, groupName string) (*rwrulefmt.RuleGroup, error) {
	results, err := c.cortexClient.GetRuleGroup(ctx, namespace, groupName)
	if err != nil {
		return nil, fmt.Errorf("cortex client: %w", err)
	}

	return results, nil
}
