package alertmanager

import (
	"bytes"
	"github.com/odpf/siren/domain"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAlertmanagerConfig(t *testing.T) {
	slackTestConfig := AMReceiverConfig{
		Receiver: "config1",
		Type:     "slack",
		Match: map[string]string{
			"foo": "bar"},
		Configuration: map[string]string{
			"token":        "xoxb",
			"channel_name": "test",
		},
	}
	pagerdutyTestConfig := AMReceiverConfig{
		Receiver: "config2",
		Type:     "pagerduty",
		Match: map[string]string{
			"bar": "baz",
		},
		Configuration: map[string]string{
			"service_key": "1234",
		},
	}
	httpTestConfig := AMReceiverConfig{
		Receiver: "config3",
		Type:     "http",
		Match:    map[string]string{},
		Configuration: map[string]string{
			"url": "http://localhost:3000",
		},
	}
	receivers := []AMReceiverConfig{slackTestConfig, pagerdutyTestConfig, httpTestConfig}
	config := AMConfig{
		Receivers: receivers,
	}

	expectedConfigStr :=
		`  templates:
    - 'helper.tmpl'
  global:
    pagerduty_url: https://events.pagerduty.com/v2/enqueue
    resolve_timeout: 5m
    slack_api_url: https://slack.com/api/chat.postMessage
  receivers:
    - name: default
    - name: slack_config1
      slack_configs:
        - channel: 'test'
          http_config:
            bearer_token: 'xoxb'
          icon_emoji: ':eagle:'
          link_names: false
          send_resolved: true
          color: '{{ template "slack.color" . }}'
          title: ''
          pretext: '{{template "slack.pretext" . }}'
          text: '{{ template "slack.body" . }}'
          actions:
            - type: button
              text: 'Runbook :books:'
              url: '{{template "slack.runbook" . }}'
            - type: button
              text: 'Dashboard :bar_chart:'
              url: '{{template "slack.dashboard" . }}'
    - name: pagerduty_config2
      pagerduty_configs:
        - service_key: '1234'
    - name: http_config3
      webhook_configs:
        - url: 'http://localhost:3000'
  route:
    group_by:
      - alertname
      - severity
      - owner
      - service_name
      - time_stamp
	  - identifier
    group_wait: 30s
    group_interval: 30m
    repeat_interval: 4h
    receiver: default
    routes:
      - receiver: slack_config1
        match:
          foo: bar
        continue: true
      - receiver: pagerduty_config2
        match:
          bar: baz
        continue: true
      - receiver: http_config3
        continue: true
`
	configStr, err := generateAlertmanagerConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, strings.Fields(expectedConfigStr), strings.Fields(configStr))
}

type ConfigCompat struct {
	TemplateFiles      map[string]string `yaml:"template_files"`
	AlertmanagerConfig string            `yaml:"alertmanager_config"`
}

func TestSyncConfig(t *testing.T) {
	config := AMReceiverConfig{
		Receiver: "config1",
		Type:     "slack",
		Match: map[string]string{
			"foo": "bar"},
		Configuration: map[string]string{
			"token":        "xoxb",
			"channel_name": "test",
		},
	}

	receiverConfig := AMConfig{Receivers: []AMReceiverConfig{config}}

	t.Run("should return error if alertmanager response code is non-2xx", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()
		client, err := NewClient(domain.CortexConfig{
			Address: ts.URL,
		})
		if err != nil {
			t.Fatal(err)
		}
		err = client.SyncConfig(receiverConfig, "fake")
		assert.Error(t, err)

	})

	t.Run("should return nil on successful sync", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenant := r.Header.Get("X-Scope-Orgid")
			assert.Equal(t, "fake", tenant)
			requestBody := ConfigCompat{}
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			err := yaml.Unmarshal(buf.Bytes(), &requestBody)
			if err != nil {
				t.Fatal(err)
			}
			assert.NotEmpty(t, requestBody.AlertmanagerConfig)
			helperTemplate := requestBody.TemplateFiles["helper.tmpl"]
			assert.NotEmpty(t, helperTemplate)
		}))
		defer ts.Close()
		client, err := NewClient(domain.CortexConfig{
			Address: ts.URL,
		})
		if err != nil {
			t.Fatal(err)
		}
		err = client.SyncConfig(receiverConfig, "fake")
		if err != nil {
			t.Fatal(err)
		}
	})
}
