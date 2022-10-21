import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# Alert History

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

Siren can store the alerts triggered by provider e.g. Cortex Alertmanager. Provider needs to be configured to call Siren API using a webhook.

## Cortex Alertmanager

For Cortex Alertmanager, everytime a provider is added, a default [webhook_config](https://prometheus.io/docs/alerting/latest/configuration/#webhook_config) receiver with empty route condition is set, which will result in calling the defined Siren API's on all alerts. This sync also happens everytime the server started.

Example Cortex Alertmanager config

```yaml
receivers:
  - name: default
    webhook_configs:
      - send_resolved: true
        http_config:
          follow_redirects: true
        url: http://localhost:8080/v1beta1/alerts/cortex/3 # siren API
        max_alerts: 0
```

Note that the url has `cortex/3` at the end, which means this will be able to parse alert history payloads from cortex type and store in DB by making it belong to provider id `3`. After this, as soon as any alert is triggered, it will be sent to Siren for history and a notification will also be published.

All information on triggered alerts depend on the alerting rule configured in Siren (synced to Cortex Ruler). Main information that should exist in the rule are:

- Which alert was triggered
- Which resource this alert refers to
- On which metric, this alert was triggered
- What was the metric value for alert trigger
- What was the severity of alert (CRITICAL, WARNING or RESOLVED)

For reusability, rule in siren need to be defined based on a [template](./template.md). Template's body describes what data that is rendered once the variable is applied.

An Example of rule's template:

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

Please note that, the mandatory keys, in order to successfully store Alert History of Cortex Alertmanager provider is,

```yaml
labels:
  severity: CRITICAL
annotations:
  resource: { { $labels.instance } }
  template: CPU
  metricName: cpu_usage_user
  metricValue: { { $labels.cpu_usage_user } }
```

The keys are pretty obvious to match with what was described in bullets points in the introduction above. The above annotations and labels will be parsed by Siren APIs to be stored in the database and would affect the content of notification message.

### Alert History Creation via API

<Tabs groupId="api">
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/alerts/cortex/1
  --header 'content-type: application/json'
  --data-raw '{
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
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

The request body of Alertmanager POST call to configured webhook looks something like (after you have followed the labels and annotations in the templates) the above snippet. The contract complies with Cortex Alertmanager [webhook_config](https://prometheus.io/docs/alerting/latest/configuration/#webhook_config) body. Siren's alerts API will parse the above payload and also store in the database, which you can fetch via the GET APIs with proper filters of startTime, endTime. See the swagger file for more details.


**Alert Notification Payload Template**

For each receiver, Siren has a default notification payload template to render Cortex alert notification. See [notification](./notification.md#message-payload-format).