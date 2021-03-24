package alertmanager

import (
	"bytes"
	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateAlertmanagerConfig(t *testing.T) {
	credentials := EntityCredentials{
		Entity: "de-infra",
		Teams: map[string]TeamCredentials{
			"eureka": {
				Name: "eureka",
				Slackcredentials: SlackConfig{
					Critical: SlackCredential{
						Webhook:  "http://eurcritical.com",
						Channel:  "dss",
						Username: "ss",
					},
				},
			},
			"wonder-woman": {
				Name: "wonder-woman",
				Slackcredentials: SlackConfig{
					Warning: SlackCredential{
						Webhook:  "http://eurcritical.com",
						Channel:  "dss",
						Username: "ss",
					},
				},
				PagerdutyCredential: "abc",
			},
		},
	}
	expectedConfigStr :=
		`  templates:
    - 'de.tmpl'
    - 'var.tmpl'
  global:
    pagerduty_url: https://events.pagerduty.com/v2/enqueue
    resolve_timeout: 5m
  receivers:
    - name: default
    - name: slack-critical-eureka
      slack_configs:
        - channel: 'dss'
          api_url: 'http://eurcritical.com'
          username: 'ss'
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
  
    - name: slack-warning-wonder-woman
      slack_configs:
        - channel: 'dss'
          api_url: 'http://eurcritical.com'
          username: 'ss'
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
    - name: pagerduty-wonder-woman
      pagerduty_configs:
        - service_key: 'abc'
  route:
    group_by:
      - alertname
      - severity
      - owner
      - service_name
      - time_stamp
    group_wait: 30s
    group_interval: 5m
    repeat_interval: 4h
    receiver: default
    routes:
      - match:
          team: 'eureka'
        routes:
          - match:
              severity: "CRITICAL"
              environment: "production"
            receiver: default
            continue: true
          - match:
              severity: "CRITICAL"
            receiver: slack-critical-eureka
          - match:
              severity: "WARNING"
            receiver:  default 

      - match:
          team: 'wonder-woman'
        routes:
          - match:
              severity: "CRITICAL"
              environment: "production"
            receiver: pagerduty-wonder-woman
            continue: true
          - match:
              severity: "CRITICAL"
            receiver:  default 
          - match:
              severity: "WARNING"
            receiver: slack-warning-wonder-woman

`
	configStr, err := generateAlertmanagerConfig(credentials)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedConfigStr, configStr)

}

type ConfigCompat struct {
	TemplateFiles      map[string]string `yaml:"template_files"`
	AlertmanagerConfig string            `yaml:"alertmanager_config"`
}

func TestSyncConfig(t *testing.T) {
	credentials := EntityCredentials{
		Entity: "greek",
		Teams: map[string]TeamCredentials{
			"eureka": {
				Name: "eureka",
				Slackcredentials: SlackConfig{
					Critical: SlackCredential{
						Webhook:  "http://eurcritical.com",
						Channel:  "dss",
						Username: "ss",
					},
				},
			},
			"wonder": {
				Name: "wonder",
				Slackcredentials: SlackConfig{
					Warning: SlackCredential{
						Webhook:  "http://eurcritical.com",
						Channel:  "dss",
						Username: "ss",
					},
				},
				PagerdutyCredential: "abc",
			},
		},
	}
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
		err = client.SyncConfig(credentials)
		assert.Error(t, err)

	})
	t.Run("should return nil on successful sync", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenant := r.Header.Get("X-Scope-Orgid")
			assert.Equal(t, "greek", tenant)
			requestBody := ConfigCompat{}
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			err := yaml.Unmarshal(buf.Bytes(), &requestBody)
			if err != nil {
				t.Fatal(err)
			}
			assert.NotEmpty(t, requestBody.AlertmanagerConfig)
			vartmpl := requestBody.TemplateFiles["var.tmpl"]
			detmpl := requestBody.TemplateFiles["de.tmpl"]
			assert.NotEmpty(t, vartmpl)
			assert.NotEmpty(t, detmpl)
		}))
		defer ts.Close()
		client, err := NewClient(domain.CortexConfig{
			Address: ts.URL,
		})
		if err != nil {
			t.Fatal(err)
		}
		err = client.SyncConfig(credentials)
		if err != nil {
			t.Fatal(err)
		}
	})
}
