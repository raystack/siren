# Alert credentials management

You can use Siren to configure Alertmanager to route alerts to get notified when alerts are triggered from the rules you
have defined. Siren stores slack and pagerduty credentials per team which gets set in Alertmanager Configuration. This
can be configured from the HTTP APIs.

**Creating/Updating Slack and Pagerduty Credentials**

Here is an example to create slack and pagerduty credentials routing configuration in Alertmanager.

A valid request payload has slack webhook, slack channel name and slack username for warning and critical severity
levels along with pagerduty service key. You also need to mention the entity (tenant id) which is tenant id in cortex
alertmanager.

```http request
PUT /alertingCredentials/teams/siren_devs HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 381

{
  "entity": "odpf",
  "pagerduty_credentials": "x67foo9bar",
  "slack_config": {
    "critical": {
      "channel": "siren_devs_critical",
      "username": "friday-bot",
      "webhook": "https://hooks.slack.com/services/abcd/efgh"
    },
    "warning": {
      "channel": "siren_devs_warning",
      "username": "friday-bot",
      "webhook": "https://hooks.slack.com/services/abcd/ikjl"
    }
  }
}
```

This will update the alertmanager configs of `odpf` entity with the credentials as given.

**How routing works**

Once alert credential has been created via Siren, there will be an entry created inside `receivers` inside Cortex
Alertmanager config of `odpf` tenant which would look like:

```yaml
      - name: slack-critical-siren_devs
        slack_configs:
          - channel: 'siren_devs-critical'
            api_url: 'https://hooks.slack.com/services/abcd/efgh'
            username: 'friday-bot'
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
      - name: slack-warning-siren_devs
        slack_configs:
          - channel: 'siren_devs-warning'
            api_url: 'https://hooks.slack.com/services/abcd/ikjl'
            username: 'friday-bot'
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
      - name: pagerduty-EIM
        pagerduty_configs:
          - service_key: 'x67foo9bar'
```

So if the rule has a label of `team` as `siren_devs` the alert will be processed by these receivers.

**Fetching alerting credentials**

```http request
GET /alertingCredentials/teams/siren_devs HTTP/1.1
Host: localhost:3000
```

This will return the stored credentials for team `siren_devs`
