package cortexv1plugin_test

import (
	"context"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/siren/core/namespace"
	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/core/rule"
	"github.com/goto/siren/core/template"
	"github.com/goto/siren/pkg/errors"
	cortexv1plugin "github.com/goto/siren/plugins/providers/cortex/v1"
	"github.com/goto/siren/plugins/providers/cortex/v1/mocks"
	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/mock"
)

func TestService_UpsertRule(t *testing.T) {
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
				Default:     "gotocompany-infra",
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
				Value: "gotocompany-infra",
			},
		},
		ProviderNamespace: 1,
	}

	type args struct {
		rl               *rule.Rule
		templateToUpdate *template.Template
		ns               namespace.Namespace
	}
	tests := []struct {
		name  string
		setup func(*mocks.CortexCaller)
		args  args
		err   error
	}{
		{
			name:  "should return error if cannot render the rule and template",
			setup: func(cc *mocks.CortexCaller) {},
			args: args{
				rl: &rule.Rule{},
				templateToUpdate: func() *template.Template {
					copiedTemplate := template.Template{}
					copiedTemplate.Body = "[[x"
					return &copiedTemplate
				}(),
				ns: namespace.Namespace{
					URN: "gotocompany",
				},
			},
			err: errors.New("template: parser:1: function \"x\" not defined"),
		},
		{
			name:  "should return error if cannot cannot parse rendered rule to RuleNode",
			setup: func(cc *mocks.CortexCaller) {},
			args: args{
				rl: &sampleRule,
				templateToUpdate: func() *template.Template {
					copiedTemplate := sampleTemplate
					copiedTemplate.Body = "name: a"
					return &copiedTemplate
				}(),
				ns: namespace.Namespace{
					URN: "gotocompany",
				},
			},
			err: errors.New("cannot parse upserted rule"),
		},
		{
			name: "should return error if getting rule group from cortexv1 return error",
			setup: func(cc *mocks.CortexCaller) {
				cc.EXPECT().GetRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("some error"))
			},
			args: args{
				rl:               &sampleRule,
				templateToUpdate: &sampleTemplate,
				ns: namespace.Namespace{
					URN: "gotocompany",
				},
			},
			err: errors.New("cannot get rule group from cortexv1 when upserting rules"),
		},
		{
			name: "should return error if merge rule nodes return empty and delete rule group return error",
			setup: func(cc *mocks.CortexCaller) {
				cc.EXPECT().GetRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&rwrulefmt.RuleGroup{}, nil)
				cc.EXPECT().DeleteRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("some error"))
			},
			args: args{
				rl: func() *rule.Rule {
					copiedRule := sampleRule
					copiedRule.Enabled = false
					return &copiedRule
				}(),
				templateToUpdate: &sampleTemplate,
				ns: namespace.Namespace{
					URN: "gotocompany",
				},
			},
			err: errors.New("error calling cortexv1: some error"),
		},
		{
			name: "should return nil if create rule group return error",
			setup: func(cc *mocks.CortexCaller) {
				cc.EXPECT().GetRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&rwrulefmt.RuleGroup{}, nil)
				cc.EXPECT().CreateRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("rwrulefmt.RuleGroup")).Return(errors.New("some error"))
			},
			args: args{
				rl:               &sampleRule,
				templateToUpdate: &sampleTemplate,
				ns: namespace.Namespace{
					URN: "gotocompany",
				},
			},
			err: errors.New("error calling cortexv1: some error"),
		},
		{
			name: "should return nil if create rule group return no error",
			setup: func(cc *mocks.CortexCaller) {
				cc.EXPECT().GetRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&rwrulefmt.RuleGroup{}, nil)
				cc.EXPECT().CreateRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("rwrulefmt.RuleGroup")).Return(nil)
			},
			args: args{
				rl:               &sampleRule,
				templateToUpdate: &sampleTemplate,
				ns: namespace.Namespace{
					URN: "gotocompany",
				},
			},
			err: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockcortexv1Client := new(mocks.CortexCaller)
			tc.setup(mockcortexv1Client)
			s := cortexv1plugin.NewPluginService(hclog.NewNullLogger(), cortexv1plugin.WithCortexClient(mockcortexv1Client))
			err := s.UpsertRule(context.Background(), tc.args.ns, provider.Provider{}, tc.args.rl, tc.args.templateToUpdate)
			if err != nil && tc.err.Error() != err.Error() {
				t.Fatalf("got error %s, expected was %s", err.Error(), tc.err)
			}
		})
	}
}
