package cortex_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/plugins/providers/cortex"
	"github.com/odpf/siren/plugins/providers/cortex/mocks"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v3"
)

func TestMergeRuleNodes(t *testing.T) {
	var dummyRuleNodes = []rulefmt.RuleNode{
		{
			Alert: yaml.Node{Value: "alert-1"},
			Labels: map[string]string{
				"key": "alert-1",
			},
		},
		{
			Alert: yaml.Node{Value: "alert-2"},
			Labels: map[string]string{
				"key": "alert-2",
			},
		},
		{
			Alert: yaml.Node{Value: "alert-3"},
			Labels: map[string]string{
				"key": "alert-3",
			},
		},
		{
			Alert: yaml.Node{Value: "alert-4"},
			Labels: map[string]string{
				"key": "alert-4",
			},
		},
		{
			Alert: yaml.Node{Value: "alert-5"},
			Labels: map[string]string{
				"key": "alert-5",
			},
		},
	}
	var GetOriginRuleNodes = func() []rulefmt.RuleNode {
		temp := make([]rulefmt.RuleNode, len(dummyRuleNodes))
		copy(temp, dummyRuleNodes)
		return temp
	}
	type args struct {
		newRuleNodes []rulefmt.RuleNode
		ruleNodes    []rulefmt.RuleNode
		enabled      bool
	}
	tests := []struct {
		name    string
		args    args
		want    []rulefmt.RuleNode
		wantErr bool
	}{
		{
			name: "should remove rn if there exist rn and enabled is false",
			args: args{
				newRuleNodes: []rulefmt.RuleNode{
					{
						Alert: yaml.Node{
							Value: "alert-3",
						},
					},
					{
						Alert: yaml.Node{
							Value: "alert-4",
						},
					},
				},
				ruleNodes: GetOriginRuleNodes(),
				enabled:   false,
			},
			want: []rulefmt.RuleNode{
				{
					Alert: yaml.Node{Value: "alert-1"},
					Labels: map[string]string{
						"key": "alert-1",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-2"},
					Labels: map[string]string{
						"key": "alert-2",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-5"},
					Labels: map[string]string{
						"key": "alert-5",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should update rn if there exist rn and enabled is true",
			args: args{
				newRuleNodes: []rulefmt.RuleNode{
					{
						Alert: yaml.Node{
							Value: "alert-3",
						},
						Labels: map[string]string{
							"key": "alert-3-new",
						},
					},
					{
						Alert: yaml.Node{
							Value: "alert-4",
						},
						Labels: map[string]string{
							"key": "alert-4-new",
						},
					},
				},
				ruleNodes: GetOriginRuleNodes(),
				enabled:   true,
			},
			want: []rulefmt.RuleNode{
				{
					Alert: yaml.Node{Value: "alert-1"},
					Labels: map[string]string{
						"key": "alert-1",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-2"},
					Labels: map[string]string{
						"key": "alert-2",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-3"},
					Labels: map[string]string{
						"key": "alert-3-new",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-4"},
					Labels: map[string]string{
						"key": "alert-4-new",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-5"},
					Labels: map[string]string{
						"key": "alert-5",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should insert rn if rn is not exist and enabled is true",
			args: args{
				newRuleNodes: []rulefmt.RuleNode{
					{
						Alert: yaml.Node{
							Value: "alert-6",
						},
						Labels: map[string]string{
							"key": "alert-6",
						},
					},
					{
						Alert: yaml.Node{
							Value: "alert-7",
						},
						Labels: map[string]string{
							"key": "alert-7",
						},
					},
				},
				ruleNodes: GetOriginRuleNodes(),
				enabled:   true,
			},
			want: []rulefmt.RuleNode{
				{
					Alert: yaml.Node{Value: "alert-1"},
					Labels: map[string]string{
						"key": "alert-1",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-2"},
					Labels: map[string]string{
						"key": "alert-2",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-3"},
					Labels: map[string]string{
						"key": "alert-3",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-4"},
					Labels: map[string]string{
						"key": "alert-4",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-5"},
					Labels: map[string]string{
						"key": "alert-5",
					},
				},
				{
					Alert: yaml.Node{Value: "alert-6"},
					Labels: map[string]string{
						"key": "alert-6",
					},
				},
				{
					Alert: yaml.Node{
						Value: "alert-7",
					},
					Labels: map[string]string{
						"key": "alert-7",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should return as is rn if rn is not exist and enabled is false",
			args: args{
				newRuleNodes: []rulefmt.RuleNode{
					{
						Alert: yaml.Node{
							Value: "alert-6",
						},
						Labels: map[string]string{
							"key": "alert-6",
						},
					},
					{
						Alert: yaml.Node{
							Value: "alert-7",
						},
						Labels: map[string]string{
							"key": "alert-7",
						},
					},
				},
				ruleNodes: GetOriginRuleNodes(),
				enabled:   false,
			},
			want:    dummyRuleNodes,
			wantErr: false,
		},
		{
			name: "should return empty if rule nodes empty and enabled false",
			args: args{
				newRuleNodes: []rulefmt.RuleNode{
					{
						Alert: yaml.Node{
							Value: "alert-6",
						},
						Labels: map[string]string{
							"key": "alert-6",
						},
					},
					{
						Alert: yaml.Node{
							Value: "alert-7",
						},
						Labels: map[string]string{
							"key": "alert-7",
						},
					},
				},
				enabled: false,
			},
			wantErr: false,
		},
		{
			name: "should return rn if rule nodes empty and enabled true",
			args: args{
				newRuleNodes: []rulefmt.RuleNode{
					{
						Alert: yaml.Node{
							Value: "alert-6",
						},
						Labels: map[string]string{
							"key": "alert-6",
						},
					},
					{
						Alert: yaml.Node{
							Value: "alert-7",
						},
						Labels: map[string]string{
							"key": "alert-7",
						},
					},
				},
				enabled: true,
			},
			want: []rulefmt.RuleNode{
				{
					Alert: yaml.Node{Value: "alert-6"},
					Labels: map[string]string{
						"key": "alert-6",
					},
				},
				{
					Alert: yaml.Node{
						Value: "alert-7",
					},
					Labels: map[string]string{
						"key": "alert-7",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := cortex.MergeRuleNodes(tc.args.ruleNodes, tc.args.newRuleNodes, tc.args.enabled)
			if (err != nil) != tc.wantErr {
				t.Errorf("CompareRuleNode() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !cmp.Equal(got, tc.want) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.want)
			}

		})
	}
}

func TestCortexService_UpsertRule(t *testing.T) {
	var sampleTemplate = template.Template{
		Name: "my-template",
		Body: heredoc.Doc(`
- alert: cpu high warning
  expr: avg by (host, environment) (cpu_usage_user{cpu="cpu-total"}) > [[.warning]]
  for: '[[.for]]'
  labels:
    alertname: CPU usage has been above [[.warning]] for last [[.for]] {{ $labels.host }}
    environment: '{{ $labels.environment }}'
    severity: WARNING
    team: '[[.team]]'
  annotations:
    metric_name: cpu_usage_user
    metric_value: '{{ printf "%0.2f" $value }}'
    resource: '{{ $labels.host }}'
    summary: CPU usage has been {{ printf "%0.2f" $value }} for last [[.for]] on host {{ $labels.host }}
    template: cpu-usage
- alert: cpu high critical
  expr: avg by (host, environment) (cpu_usage_user{cpu="cpu-total"}) > [[.critical]]
  for: '[[.for]]'
  labels:
    alertname: CPU usage has been above [[.warning]] for last [[.for]] {{ $labels.host }}
    environment: '{{ $labels.environment }}'
    severity: CRITICAL
    team: '[[.team]]'
  annotations:
    metric_name: cpu_usage_user
    metric_value: '{{ printf "%0.2f" $value }}'
    resource: '{{ $labels.host }}'
    summary: CPU usage has been {{ printf "%0.2f" $value }} for last [[.for]] on host {{ $labels.host }}
    template: cpu-usage`),
		Variables: []template.Variable{
			{
				Name:        "for",
				Type:        "string",
				Description: "For eg 5m, 2h; Golang duration format",
				Default:     "5m",
			},
			{
				Name:    "warning",
				Type:    "int",
				Default: "85",
			},
			{
				Name:    "critical",
				Type:    "int",
				Default: "90",
			},
			{
				Name:        "team",
				Type:        "string",
				Description: "For eg team name which the alert should go to",
				Default:     "odpf-infra",
			},
		},
		Tags: []string{"system"},
	}

	var sampleRule = rule.Rule{
		Name:      "siren_api_provider-urn_namespace-urn_system_cpu-usage_cpu-usage",
		Namespace: "system",
		GroupName: "cpu-usage",
		Template:  "cpu-usage",
		Enabled:   true,
		Variables: []rule.RuleVariable{
			{
				Name:  "for",
				Value: "5m",
			},
			{
				Name:  "warning",
				Value: "85",
			},
			{
				Name:  "critical",
				Value: "90",
			},
			{
				Name:  "team",
				Value: "odpf-infra",
			},
		},
		ProviderNamespace: 1,
	}

	type args struct {
		rl               *rule.Rule
		templateToUpdate *template.Template
		namespaceURN     string
	}
	tests := []struct {
		name  string
		setup func(*mocks.CortexClient)
		args  args
		err   error
	}{
		{
			name:  "should return error if cannot render the rule and template",
			setup: func(cc *mocks.CortexClient) {},
			args: args{
				rl: &rule.Rule{},
				templateToUpdate: func() *template.Template {
					copiedTemplate := template.Template{}
					copiedTemplate.Body = "[[x"
					return &copiedTemplate
				}(),
				namespaceURN: "odpf",
			},
			err: errors.New("failed to parse template body"),
		},
		{
			name:  "should return error if cannot cannot parse rendered rule to RuleNode",
			setup: func(cc *mocks.CortexClient) {},
			args: args{
				rl: &sampleRule,
				templateToUpdate: func() *template.Template {
					copiedTemplate := sampleTemplate
					copiedTemplate.Body = "name: a"
					return &copiedTemplate
				}(),
				namespaceURN: "odpf",
			},
			err: errors.New("cannot parse upserted rule"),
		},
		{
			name: "should return error if getting rule group from cortex return error",
			setup: func(cc *mocks.CortexClient) {
				cc.EXPECT().GetRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("some error"))
			},
			args: args{
				rl:               &sampleRule,
				templateToUpdate: &sampleTemplate,
				namespaceURN:     "odpf",
			},
			err: errors.New("cannot get rule group from cortex when upserting rules"),
		},
		// {
		// 	name:    "should return error if merge rule nodes return error",
		// 	wantErr: true,
		// },
		{
			name: "should return error if merge rule nodes return empty and delete rule group return error",
			setup: func(cc *mocks.CortexClient) {
				cc.EXPECT().GetRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&rwrulefmt.RuleGroup{}, nil)
				cc.EXPECT().DeleteRuleGroup(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("some error"))
			},
			args: args{
				rl: func() *rule.Rule {
					copiedRule := sampleRule
					copiedRule.Enabled = false
					return &copiedRule
				}(),
				templateToUpdate: &sampleTemplate,
				namespaceURN:     "odpf",
			},
			err: errors.New("error calling cortex: some error"),
		},
		{
			name: "should return nil if create rule group return error",
			setup: func(cc *mocks.CortexClient) {
				cc.EXPECT().GetRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&rwrulefmt.RuleGroup{}, nil)
				cc.EXPECT().CreateRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("rwrulefmt.RuleGroup")).Return(errors.New("some error"))
			},
			args: args{
				rl:               &sampleRule,
				templateToUpdate: &sampleTemplate,
				namespaceURN:     "odpf",
			},
			err: errors.New("error calling cortex: some error"),
		},
		{
			name: "should return nil if create rule group return no error",
			setup: func(cc *mocks.CortexClient) {
				cc.EXPECT().GetRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&rwrulefmt.RuleGroup{}, nil)
				cc.EXPECT().CreateRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("rwrulefmt.RuleGroup")).Return(nil)
			},
			args: args{
				rl:               &sampleRule,
				templateToUpdate: &sampleTemplate,
				namespaceURN:     "odpf",
			},
			err: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCortexClient := new(mocks.CortexClient)
			tc.setup(mockCortexClient)
			s := cortex.NewProviderService(mockCortexClient)
			err := s.UpsertRule(context.Background(), tc.args.rl, tc.args.templateToUpdate, tc.args.namespaceURN)
			if err != nil && tc.err.Error() != err.Error() {
				t.Fatalf("got error %s, expected was %s", err.Error(), tc.err)
			}
		})
	}
}
