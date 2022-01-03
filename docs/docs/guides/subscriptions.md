# Subscriptions

Siren lets you subscribe to the rules when they are triggered. You can define custom matching conditions and use
[receivers](./receivers.md) to describe which medium you want to use for getting the notifications when those rules are
triggered. Siren syncs this configuration in the respective monitoring provider.

**Example Subscription:**

```json
{
  "id": "385",
  "urn": "siren-dev-prod-critical",
  "namespace": "10",
  "receivers": [
    {
      "id": "2"
    },
    {
      "id": "1",
      "configuration": {
        "channel_name": "siren-dev-critical"
      }
    }
  ],
  "match": {
    "environment": "production",
    "severity": "CRITICAL"
  },
  "created_at": "2021-12-10T10:38:22.364353Z",
  "updated_at": "2021-12-10T10:38:22.364353Z"
}
```

The above means whenever any alert which has labels matching the `match`viz:
`"environment": "production", "severity": "CRITICAL"`, send this alert to two medium defined by receivers with id: `2`
and `1`. Assuming the receivers id `2` to be of Pagerduty type, a PD call will be invoked and assuming the receiver with
id `1` to be slack type, a message will be sent to the channel #siren-dev-critical.

**Upstream sync example**

The logical equivalence of this routing configuration is put in the respective monitoring provider by Siren. For ex: if
the provider is cortex, an alertmanager configuration will be created depicting the above routing logic.

```yaml
templates:
  - "helper.tmpl"
global:
  pagerduty_url: https://events.pagerduty.com/v2/enqueue
  resolve_timeout: 5m
  slack_api_url: https://slack.com/api/chat.postMessage
receivers:
  - name: default
  - name: slack_siren-dev-prod-critical_receiverId_1_idx_0
    slack_configs:
      - channel: "siren-dev-critical"
        http_config:
          bearer_token: "secret-taken-from-receiver-config"
        icon_emoji: ":eagle:"
        link_names: false
        send_resolved: true
        color: '{{ template "slack.color" . }}'
        title: ""
        pretext: '{{template "slack.pretext" . }}'
        text: '{{ template "slack.body" . }}'
        actions:
          - type: button
            text: "Runbook :books:"
            url: '{{template "slack.runbook" . }}'
          - type: button
            text: "Dashboard :bar_chart:"
            url: '{{template "slack.dashboard" . }}'
  - name: pagerduty_siren-dev-prod-critical_receiverId_2_idx_1
    pagerduty_configs:
      - service_key: "secret-taken-from-receiver-config"
route:
  group_by:
    - alertname
    - severity
    - owner
    - service_name
    - time_stamp
    - identifier
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  receiver: default
  routes:
    - receiver: slack_siren-dev-prod-critical_receiverId_1_idx_0
      match:
        environment: production
        severity: CRITICAL
      continue: true
    - receiver: pagerduty_siren-dev-prod-critical_receiverId_2_idx_1
      match:
        environment: production
        severity: CRITICAL
      continue: true

```

As you can see, Siren dynamically defined two receivers: `slack_siren-dev-prod-critical_receiverId_1_idx_0`
and `pagerduty_siren-dev-prod-critical_receiverId_2_idx_1` and used them in the routing tree as per the match
conditions.

This alertmanager config is for the tenant defined by namespace with id `10` as mentioned in the example. This is an
example config, the actual config will contain all subscriptions that belong to namespace with id `10`

## API Interface

### Create a subscription

```text
POST /v1beta1/subscriptions HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 363

{
    "urn": "siren-dev-prod-critical",
    "receivers": [
        {
            "id": "1",
            "configuration": {
                "channel_name": "siren-dev-critical"
            }
        },
        {
            "id": "2"
        }
    ],
    "match": {
        "severity": "CRITICAL",
        "environment": "production"
    },
    "namespace": "10"
}
```

### Update a subscription

```text
POST /v1beta1/subscriptions HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 392

{
    "urn": "siren-dev-prod-critical",
    "receivers": [
        {
            "id": "1",
            "configuration": {
                "channel_name": "siren-dev-critical"
            }
        },
        {
            "id": "2"
        }
    ],
    "match": {
        "severity": "CRITICAL",
        "environment": "production",
        "team": "siren-dev"
    },
    "namespace": "10"
}
```

### Get all subscriptions

```text
GET /v1beta1/subscriptions HTTP/1.1
Host: localhost:3000
```

### Get a subscriptions

```text
GET /v1beta1/subscriptions/10 HTTP/1.1
Host: localhost:3000
```

### Delete subscriptions

```text
DELETE /v1beta1/subscriptions/10 HTTP/1.1
Host: localhost:3000
```
