import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# 6 - Subscribing Notifications

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

Notifications can be subscribed and routed to the defined receivers by adding a subscription. In this part, we will simulate how Cortex Ruler trigger an alert to Cortex Alertmanager, and Cortex Alertmanager trigger webhook-notification and calling Siren alerts hook API. On Siren side, we expect a notification is published everytime the hook API is being called.

In this part we will create alerting rules for our Cortex monitoring provider. Rules in Siren relies on [template](../guides/template.md) for its abstraction. We need to create a rule's template first before uploading a rule.

The first thing that we should do is knowing what would be the labels sent by Cortex Alertmanager. The labels should be defined when we were defining [rules](./5_configuring_provider_alerting_rules.md#creating-a-rule). Assuming the labels sent by Cortex Alertmanager are these:

```yaml
severity: WARNING
team: odpf
service: some-service
environment: integration
resource_name: some-resource
```

Later we will try to simulate triggering alert by calling Cortex Alertmanager `POST /alerts` API directly. The way Cortex Ruler monitor a specific metric and trigger an alert to Cortex Alertmanager are out of this `tour` scope.

We want to subscribe all notifications owned by `odpf` team and has severity `WARNING` regardless the service name related with the alerts and route the notification to `file` with receiver id `1` and `2`. Currently there is no CLI to create a subscription (this would need to be added in the future) so we could call Siren HTTP API direclty to create one.

Prepare a subscription detail:
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
  ./siren subscription create --file cpu_subs.yaml
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

Once a subscription is created, let's simulate on how Cortex Ruler trigger an alert by calling Cortex Alertmanager API directly with this cURL.

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


Now, we need to expect Cortex Alertmanager send alerts to our Siren API `/alerts/cortex/:providerId`. If that is the case, the alert should also be stored and published to the receivers in the matching subscriptions. You might want to wait for a Cortex Alertmanager `group_wait` (usually 30s) until alerts are triggered by Cortex Alertmanager.

Let's verify the alert is stored inside our DB.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
./siren alert list --provider-id 1 --provider-type cortex --resource-name test_alert
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



We also expect notifications have been published to the receiver id `1` and `2` similar with the [previous part](./4_sending_notifications_to_receiver.md). You can check a new notification is already added in `./out-file-sink1.json` and `./out-file-sink2.json` with this value.

```json
{"environment":"integration","generatorUrl":"","groupKey":"{}:{severity=\"WARNING\"}","metricName":"test_alert","metricValue":"1","numAlertsFiring":1,"resource":"test_alert","routing_method":"subscribers","service":"some-service","severity":"WARNING","status":"firing","team":"odpf","template":"alert_test"}
```