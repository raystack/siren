import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# 2.5 Subscribing to Alert Notifications

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

Notifications can be subscribed and routed to the defined receivers by adding a subscription. In this part, we will trigger an alert to CortexMetrics manually by calling CortexMetrics `POST /alerts` API and expect CortexMetrics to trigger webhook-notification and calling Siren alerts hook API. On Siren side, we expect a notification is published everytime the hook API is being called.

> The way CortexMetrics monitor a specific metric and auto-trigger an alert are out of this `tour` scope.


The first thing that we should do is knowing what would be the labels sent by CortexMetrics. The labels should be defined when we were defining [rules](./4configuring_provider_alerting_rules.md). Assuming the labels sent by CortexMetrics are these:

```yaml
severity: WARNING
team: odpf
service: some-service
environment: integration
resource_name: some-resource
```

We want to subscribe all notifications owned by `odpf` team and has severity `WARNING` regardless the service name and route the notification to `file` with receiver id `2` (the one that we created in the [previous](./3registering_receivers.md) part).

Prepare a subscription detail and create a new subscription with Siren CLI.
```bash
cat <<EOT >> cpu_subs.yaml
urn: subscribe-cpu-odpf-warning
namespace: 1
receivers:
  - id: 1
  - id: 2
match
  team: odpf
  severity: WARNING
EOT
```

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

  ```shell
  $ siren subscription create --file cpu_subs.yaml
  ```
  
  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/subscriptions'
--header 'Content-Type: application/json'
--header 'Accept: application/json'
--data-raw '{
  "urn": "subscribe-cpu-odpf-warning",
  "namespace": 1,
  "receivers": [
    {
      "id": 1
    },
    {
      "id": 2
    }
  ],
  "match": {
    "team": "odpf",
    "severity": "WARNING"
  }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

Once a subscription is created, let's manually trigger alert in CortexMetrics with this cURL.

<Tabs groupId="api">
  <TabItem value="http" label="HTTP">

  ```bash
  curl --location --request POST 'http://localhost:9009/api/prom/alertmanager/api/v1/alerts'
  --header 'X-Scope-OrgId: odpf-ns'
  --header 'Content-Type: application/json' \
  --data-raw '[
      {
          "state": "firing",
          "value": 1,
          "labels": {
              "severity": "WARNING",
              "team": "odpf",
              "service": "some-service",
              "environment": "integration"
          },
          "annotations": {
              "resource": "test_alert",
              "metricName": "test_alert",
              "metricValue": "1",
              "template": "alert_test"
          }
      }
  ]'
  ```

  </TabItem>
</Tabs>



If succeed, the response should be like this.
```json
{"status":"success"}
```


Now, we need to expect CortexMetrics to send alerts notification to our Siren API `/alerts/cortex/:providerId`. If that is the case, the alert should also be stored and published to the receivers in the matching subscriptions. You might want to wait for a CortexMetrics `group_wait` (usually 30s) until alerts are triggered by Cortex Alertmanager.

Let's verify the alert is stored inside our DB.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren alert list --provider-id 1 --provider-type cortex --resource-name test_alert
```
The result would be something like this.
```shell
Showing 1 of 1 alerts
 
ID      PROVIDER_ID     RESOURCE_NAME   METRIC_NAME     METRIC_VALUE    SEVERITY
1       1               test_alert      test_alert      1               WARNING 

For details on a alert, try: siren alert view <id>
```
  
  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/alerts?providerId=1&providerType=cortex&resourceName=test_alert`}
    </CodeBlock>
  </TabItem>
</Tabs>


We also expect notifications have been published to the receiver id `2`. You can check a new notification is already added in `./out-file-sink2.json` with this value.

```json
{"environment":"integration","generatorUrl":"","groupKey":"{}:{severity=\"WARNING\"}","metricName":"test_alert","metricValue":"1","numAlertsFiring":1,"resource":"test_alert","routing_method":"subscribers","service":"some-service","severity":"WARNING","status":"firing","team":"odpf","template":"alert_test"}
```

## What Next?

This is the end of `alerting rules and subscription` tour. If you want to know how to send on-demand notification to a receiver, you could check the [first tour](../notifications/1sending_notifications_overview.md).

Apart from the tour, we recommend completing the [guides](../../guides/overview.md). You could also check out the remainder of the documentation in the [reference](../../reference/server_configuration.md) and [concepts](../../concepts/overview.md) sections for your specific areas of interest. We've aimed to provide as much documentation as we can for the various components of Siren to give you a full understanding of Siren's surface area. If you are interested to contribute, check out the [contribution](../../contribute/contribution.md) page.
