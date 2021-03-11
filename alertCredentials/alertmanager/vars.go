package alertmanager

var alertmanagerConfigTemplate = `  templates:
    - 'de.tmpl'
    - 'var.tmpl'
  global:
    pagerduty_url: https://events.pagerduty.com/v2/enqueue
    resolve_timeout: 5m
  receivers:
    - name: default[[- /*gotype: github.com/odpf/siren/alertCredentials.EntityCredentials*/ -]]
  [[- range $key, $team := .Teams -]]
  [[if eq $team.Slackcredentials.Critical.Webhook ""]]
  [[else]]
    - name: slack-critical-[[$team.Name]]
      slack_configs:
        - channel: '[[$team.Slackcredentials.Critical.Channel]]'
          api_url: '[[$team.Slackcredentials.Critical.Webhook]]'
          username: '[[$team.Slackcredentials.Critical.Username]]'
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
[[- end -]]
[[- if eq $team.Slackcredentials.Warning.Webhook "" -]][[- else ]]
    - name: slack-warning-[[$team.Name]]
      slack_configs:
        - channel: '[[$team.Slackcredentials.Warning.Channel]]'
          api_url: '[[$team.Slackcredentials.Warning.Webhook]]'
          username: '[[$team.Slackcredentials.Warning.Username]]'
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
[[- end]]
[[- if eq $team.PagerdutyCredential "" -]][[- else ]]
    - name: pagerduty-[[$team.Name]]
      pagerduty_configs:
        - service_key: '[[$team.PagerdutyCredential]]'
[[- end]]
[[- end]]
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
    routes:[[- range $key, $team := .Teams]]
      - match:
          team: '[[$team.Name]]'
        routes:
          - match:
              severity: "CRITICAL"
              environment: "production"
            receiver: [[- if eq $team.PagerdutyCredential "" ]] default[[- else ]] pagerduty-[[$team.Name]] [[- end]]
            continue: true
          - match:
              severity: "CRITICAL"
            receiver: [[if eq $team.Slackcredentials.Critical.Webhook ""]] default [[else]]slack-critical-[[$team.Name]][[end]]
          - match:
              severity: "WARNING"
            receiver: [[if eq $team.Slackcredentials.Warning.Webhook ""]] default [[else]]slack-warning-[[$team.Name]][[end]]
[[end]]
`

