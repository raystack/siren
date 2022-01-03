# Alert History Subscription

Siren can store the alerts triggered via Cortex Alertmanager. Cortex alertmanager is configured to call Siren API, using
a webhook receiver. This is done by adding a subscription using an HTTP Receiver on empty match condition, which will
result in calling thus HTTP Receiver on all alerts

Example Receiver

```json
{
  "id": "1",
  "name": "alert-history-receiver",
  "type": "http",
  "labels": {
    "team": "siren-devs-alert-history"
  },
  "configurations": {
    "url": "http://localhost:3000/v1beta1/alerts/cortex/3"
  }
}
```

Note that the url has `cortex/3` at the end, which means this will be able to parse alert history payloads from cortex
type and store in DB by making it belong to provider id `3`.

We will need the subscription as well, example:

```json
{
  "id": "384",
  "urn": "alert-history-subscription",
  "namespace": "10",
  "receivers": [
    {
      "id": "1"
    }
  ],
  "match": {}
}
```

After this, as soon as any alert is sent by Alertmanager to slack or pagerduty, it will be sent to Siren for storage
purpose.

![Siren Alert History](../assets/alerthistory.jpg)

The parsing of payload from alert manager depends on a particular syntax. you can configure your templates to follow
this syntax, with proper annotations to identify:

- Which alert was triggered
- Which resource this alert refers to
- On Which metric, this alert was triggered
- What was the metric value for alert trigger
- What was the severity of alert(CRITICAL, WARNING or RESOLVED)

An Example template:

```yaml
apiVersion: v2
type: template
name: CPU
body:
  - alert: CPUWarning
    expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.warning]]
    for: "[[.for]]"
    labels:
      severity: WARNING
    annotations:
      description: CPU has been above [[.warning]] for last [[.for]] {{ $labels.host }}
      resource: { { $labels.instance } }
      template: CPU
      metricName: cpu_usage_user
      metricValue: { { $labels.cpu_usage_user } }
  - alert: CPUCritical
    expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.critical]]
    for: "[[.for]]"
    labels:
      severity: CRITICAL
    annotations:
      description: CPU has been above [[.critical]] for last [[.for]] {{ $labels.host }}
      resource: { { $labels.instance } }
      template: CPU
      metricName: cpu_usage_user
      metricValue: { { $labels.cpu_usage_user } }
variables:
  - name: for
    type: string
    default: 10m
    description: For eg 5m, 2h; Golang duration format
  - name: warning
    type: int
    default: 80
  - name: critical
    type: int
    default: 90
tags:
  - systems
```

Please note that, the mandatory keys, in order to successfully store Alert History is,

```yaml
labels:
  severity: CRITICAL
annotations:
  resource: { { $labels.instance } }
  template: CPU
  metricName: cpu_usage_user
  metricValue: { { $labels.cpu_usage_user } }

```

The keys are pretty obvious to match with what was described in bullets points in the introduction above.

In the above example we can see, the alert annotation has sufficient values for alert history storage. We can set up
cortex alertmanager, to call Siren AlertHistory APIs as a webhook receiver. The above annotations and labels will be
parsed by Siren APIs, to be stored in the database.

**Alert History Creation via API**

```text
POST /v1beta1/alerts/cortex/1 HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 357

{
    "alerts": [
        {
            "status": "firing",
            "labels": {
                "severity": "CRITICAL"
            },
            "annotations": {
                "resource": "apolloVM",
                "template": "CPU",
                "metricName": "cpu_usage_user",
                "metricValue": "90"
            }
        }
    ]
}
```

The request body of Alertmanager POST call to configured webhook looks something like (after you have followed the
labels and annotations c in the templates) above snippet.

The alerts API will parse the above payload and store in the database, which you can fetch via the GET APIs with proper
filters of startTime, endTime. See the [swagger](../../api/handlers/swagger.yaml) file for more details.
