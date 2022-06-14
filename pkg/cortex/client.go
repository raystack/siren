package cortex

import (
	"bytes"
	"context"
	"html/template"

	"github.com/grafana/cortex-tools/pkg/client"
	promconfig "github.com/prometheus/alertmanager/config"
)

type Client struct {
	cortexClient   CortexCaller
	helperTemplate string
	configYaml     string
}

func NewClient(cfg Config, opts ...ClientOption) (*Client, error) {
	c := &Client{}

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

func (c *Client) CreateRuleGroup(ctx context.Context, namespace string, rg RuleGroup) error {
	return c.cortexClient.CreateRuleGroup(ctx, namespace, rg.RuleGroup)
}

func (c *Client) DeleteRuleGroup(ctx context.Context, namespace, groupName string) error {
	return c.cortexClient.DeleteRuleGroup(ctx, namespace, groupName)
}

func (c *Client) GetRuleGroup(ctx context.Context, namespace, groupName string) (*RuleGroup, error) {
	rg, err := c.cortexClient.GetRuleGroup(ctx, namespace, groupName)
	if err != nil {
		return nil, err
	}
	return &RuleGroup{RuleGroup: *rg}, nil
}

func (c *Client) ListRules(ctx context.Context, namespace string) (map[string][]RuleGroup, error) {
	crgsMap, err := c.cortexClient.ListRules(ctx, namespace)
	if err != nil {
		return nil, err
	}

	rgsMap := make(map[string][]RuleGroup)
	for k, crgs := range crgsMap {
		rgs := []RuleGroup{}
		for _, crg := range crgs {
			rgs = append(rgs, RuleGroup{crg})
		}
		rgsMap[k] = rgs
	}

	return rgsMap, nil
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
