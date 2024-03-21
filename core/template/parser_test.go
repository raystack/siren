package template_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/goto/siren/core/receiver"
	"github.com/goto/siren/core/template"
	testdatatemplate_test "github.com/goto/siren/test/e2e_test/testdata/templates"
	"github.com/stretchr/testify/assert"
)

func Test_Parser(t *testing.T) {
	t.Run("successfully parse rule body", func(t *testing.T) {
		expectedRules := []template.Rule{
			{
				Alert: "cpu high warning",
				Expr:  "avg by (host, environment) (cpu_usage_user{cpu=\"cpu-total\"}) > [[.warning]]",
				For:   "[[.for]]",
				Labels: map[string]string{
					"severity":    "WARNING",
					"alert_name":  "CPU usage has been above [[.warning]] for last [[.for]] {{ $labels.host }}",
					"environment": "{{ $labels.environment }}",
					"team":        "[[.team]]",
				},
				Annotations: map[string]string{
					"dashboard":    "https://dashboard.gotocompany.com/xxx",
					"summary":      "CPU usage has been {{ printf \"%0.2f\" $value }} for last [[.for]] on host {{ $labels.host }}",
					"resource":     "{{ $labels.host }}",
					"template":     "cpu-usage",
					"metric_name":  "cpu_usage_user",
					"metric_value": "{{ printf \"%0.2f\" $value }}",
				},
			},
			{
				Alert: "cpu high critical",
				Expr:  "avg by (host, environment) (cpu_usage_user{cpu=\"cpu-total\"}) > [[.critical]]",
				For:   "[[.for]]",
				Labels: map[string]string{
					"severity":    "CRITICAL",
					"alert_name":  "CPU usage has been above [[.critical]] for last [[.for]] {{ $labels.host }}",
					"environment": "{{ $labels.environment }}",
					"team":        "[[.team]]",
				},
				Annotations: map[string]string{
					"dashboard":    "https://dashboard.gotocompany.com/xxx",
					"summary":      "CPU usage has been {{ printf \"%0.2f\" $value }} for last [[.for]] on host {{ $labels.host }}",
					"resource":     "{{ $labels.host }}",
					"template":     "cpu-usage",
					"metric_name":  "cpu_usage_user",
					"metric_value": "{{ printf \"%0.2f\" $value }}",
				},
			},
		}
		fl, err := template.YamlStringToFile(testdatatemplate_test.SampleRuleTemplate)
		assert.NoError(t, err)

		tpl, err := template.ParseFile(fl)
		assert.NoError(t, err)

		rules, err := template.RulesBody(tpl)
		assert.NoError(t, err)

		if diff := cmp.Diff(rules, expectedRules); diff != "" {
			t.Fatalf("got diff %v", diff)
		}
	})

	t.Run("successfully parse message body", func(t *testing.T) {
		expectedMessages := []template.Message{
			{
				ReceiverType: receiver.TypeSlack,
				Content: `title: |
  Several lines of text.
  with some "quotes" of various 'types'.
  Escapes (like \n) don't do anything.

  Newlines can be added by leaving a blank line.
  Additional leading whitespace is ignored.`,
			},
			{
				ReceiverType: receiver.TypePagerDuty,
				Content: `title: |
  Plain flow scalars are picky about the (:) and (#) characters. 
  They can be in the string, but (:) cannot appear before a space or newline.
  And (#) cannot appear after a space or newline; doing this will cause a syntax error. 
  If you need to use these characters you are probably better off using one of the quoted styles instead.`,
			},
			{
				ReceiverType: receiver.TypeHTTP,
				Content: `title: {{.Data.title}}
description: |
  Plain flow scalars are picky about the (:) and (#) characters. 
  They can be in the string, but (:) cannot appear before a space or newline.
  And (#) cannot appear after a space or newline; doing this will cause a syntax error. 
  If you need to use these characters you are probably better off using one of the quoted styles instead.
category: {{.Labels.category}}`,
			},
		}
		fl, err := template.YamlStringToFile(testdatatemplate_test.SampleMessageTemplate)
		assert.NoError(t, err)

		tpl, err := template.ParseFile(fl)
		assert.NoError(t, err)

		messages, err := template.MessagesFromBody(tpl)
		assert.NoError(t, err)

		if diff := cmp.Diff(messages, expectedMessages); diff != "" {
			t.Fatalf("got diff %v", diff)
		}
	})
}
