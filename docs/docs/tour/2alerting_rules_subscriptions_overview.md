import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# 2 Alerting Rules and Subscription

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

This tour shows you how could we create alerting rules and we want to subscribe to a notification triggered by an alert. If you want to know how to send on-demand notification to a receiver, you could go to the [first tour](./1sending_notifications_overview.md).

As mentioned previously, we will be using [CortexMetrics](https://cortexmetrics.io/docs/getting-started/) as a provider. We need to register the provider and create a provider namespace in Siren first before creating any rule and subscription. 

> Provider is implemented as a plugin in Siren. You can learn more about Siren Plugin concepts [here](/docs/docs/concepts/plugin.md). We also welcome all contributions to add new provider plugins. Learn more about how to add a new provider plugin [here](/docs/docs/extend/adding_new_provider.md).

Once an alert triggered, the subscription labels will be matched with alert's labels. If all subscription labels matched, receiver's subscripton will get the alert notification.

## 2.1 Register a Provider and Namespaces

### Register a Provider

To create a new provider with CLI, we need to create a `yaml` file that contains provider detail.

```yaml title=provider.yaml
host: http://localhost:9009
urn: localhost-dev-cortex
name: dev-cortex
type: cortex
```

Once the file is ready, we can create the provider with Siren CLI.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren provider create --file provider.yaml
```

If succeed, you will got this message.

```shell
Provider created with id: 1 ✓
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/providers
  --header 'content-type: application/json'
  --data-raw '{
    "host": "http://localhost:9009",
    "urn": "localhost-dev-cortex",
    "name": "dev-cortex",
    "type": "cortex"
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

The `id` we got from the provider creation is important to create a namespace later.

### Register Namespaces

For multi-tenant scenario, which [CortexMetrics](https://cortexmetrics.io/) supports, we need to define namespaces in Siren. Assuming there are 2 tenants in Cortex, `odpf` and `non-odpf`, we need to create 2 namespaces. This could be done in similar way with how we created provider.

```bash  title=ns1.yaml
urn: odpf-ns
name: odpf-ns
provider:
    id: 1
```

```bash  title=ns2.yaml
urn: non-odpf-ns
name: non-odpf-ns
provider:
    id: 1
```

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren namespace create --file ns1.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces
  --header 'content-type: application/json'
  --data-raw '{
    "urn": "odpf-ns",
    "name": "odpf-ns",
    "provider": {
        "id": 1
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren namespace create --file ns2.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces
  --header 'content-type: application/json'
  --data-raw '{
    "urn": "non-odpf-ns",
    "name": "non-odpf-ns",
    "provider": {
        "id": 2
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Verify Created Provider and Namespaces

To make sure all provider and namespaces are properly created, we could try query Siren with Siren CLI.

See what providers exist in Siren.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren provider list
```

```shell
Showing 1 of 1 providers

ID      TYPE    URN                     NAME
1       cortex  localhost-dev-cortex    dev-cortex

For details on a provider, try: siren provider view <id>
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/providers'`}
    </CodeBlock>
  </TabItem>
</Tabs>

See what namespaces exist in Siren.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren namespace list
```

```shell
Showing 2 of 2 namespaces

ID      URN             NAME
1       odpf-ns         odpf-ns
2       non-odpf-ns     non-odpf-ns

For details on a namespace, try: siren namespace view <id>
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces'`}
    </CodeBlock>
  </TabItem>
</Tabs>

## 2.2 Register a Receiver

Siren supports several types of receiver to send notification to. For this tour, let's pick the simplest receiver: `file`. If the receiver is not added in Siren yet, you could add one using `siren receiver create`. See [receiver guide](/docs/docs/guides/receiver.md) to explore more on how to work with `siren receiver` command.

With `file` receiver, all published notifications will be written to a file. Let's create a receivers `file` using Siren CLI.

> We welcome all contributions to add new type of receiver plugins. See [Extend](/docs/docs/extend/adding_new_receiver.md) section to explore how to add a new type of receiver plugin to Siren

Prepare receiver detail and register the receiver with Siren CLI.
```bash  title=receiver_2.yaml
name: file-sink-2
type: file
labels:
    key1: value1
    key2: value2
configurations:
    url: ./out-file-sink2.json
```

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren receiver create --file receiver_2.yaml
```

Once done, you will get a message.

```bash
Receiver created with id: 2 ✓
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers
  --header 'content-type: application/json'
  --data-raw '{
    "name": "file-sink-2",
    "type": "file",
    "labels": {
        "key1": "value1",
        "key2": "value2"
    },
    "configurations": {
        "url": "./out-file-sink2.json"
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>


## 2.3 Configuring Provider Alerting Rules

In this part, we will create alerting rules for our CortexMetrics monitoring provider. Rules in Siren relies on [template](../guides/template.md) for its abstraction. To create a rule, we need to create a template first.

### Creating a Rule's Template

We will create a rule's template to monitor CPU usage. 
```yaml  title=cpu_template.yaml
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
    - alert: CPUCritical
      expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.critical]]
      for: "[[.for]]"
      labels:
          severity: CRITICAL
      annotations:
          description: CPU has been above [[.critical]] for last [[.for]] {{ $labels.host }}
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

We named the template above as `CPU`, the body in the template is the data that will be interpolated with variables. Notice that template body format is similar with [Prometheus alerting rules](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/). This is because Cortex uses the same rules as prometheus and Siren will translate the rendered rules to the Cortex alerting rules. Let's save the template above into a file called `cpu_template.yaml` and upload our template to Siren using 

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren template upload cpu_template.yaml
```

  </TabItem>
</Tabs>

You could verify the newly created template using this command.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren template list
```
```shell
Showing 1 of 1 templates
 
ID      NAME    TAGS   
1       CPU     systems

For details on a template, try: siren template view <name>
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/templates`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Creating a Rule

Now we already have a `CPU` template, we can create a rule based on that template. Let's prepare a rule and save it in a file called `cpu_test.yaml`.
```yaml  title=cpu_test.yaml
apiVersion: v2
type: rule
namespace: odpf
provider: localhost-dev-cortex
providerNamespace: odpf-ns
rules:
    cpuGroup:
        template: CPU
        enabled: true
        variables:
            - name: for
              value: 15m
            - name: warning
              value: 185
            - name: critical
              value: 195
```
We defined a rule based on `CPU` template for namespace urn `odpf-ns` and provider urn `localhost-dev-cortex`. The rule group name is `cpuGroup` and there are also some variables to be assign to the template when the template is rendered. Let's upload the rule with Siren CLI.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren rule upload cpu_test.yaml
```
If succeed, you will get this message.
```shell
Upserted Rule
ID: 4
```
  </TabItem>
</Tabs>

You could verify the created rules by getting all registered rules in CortexMetrics with cURL.
<Tabs groupId="api">
  <TabItem value="http" label="HTTP">

  ```bash
  curl --location --request GET 'http://localhost:9009/api/v1/rules' \
  --header 'X-Scope-OrgId: odpf-ns'
  ```

  </TabItem>
</Tabs>

The response body should be in `yaml` format and like this.
```yaml
odpf:
    - name: cpuGroup
      rules:
        - alert: CPUWarning
          expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > 185
          for: 15m
          labels:
            severity: WARNING
          annotations:
            description: CPU has been above 185 for last 15m {{ $labels.host }}
        - alert: CPUCritical
          expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > 195
          for: 15m
          labels:
            severity: CRITICAL
          annotations:
            description: CPU has been above 195 for last 15m {{ $labels.host }}
```

If there is a response like above, that means the rule that we created in Siren was already synchronized to the provider. Next, we can add a subscription to the alert and try to trigger an alert to verify whether we got a notification alert or not.

## 2.4 Subscribing to Alert Notifications

Notifications can be subscribed and routed to the defined receivers by adding a subscription. In this part, we will trigger an alert to CortexMetrics manually by calling CortexMetrics `POST /alerts` API and expect CortexMetrics to trigger webhook-notification and calling Siren alerts hook API. On Siren side, we expect a notification is published everytime the hook API is being called.


> If you are curious about how notification in Siren works, you can read the concepts [here](/docs/docs/concepts/notification.md).

The first thing that we should do is knowing what would be the labels sent by CortexMetrics. The labels should be defined when we were defining [rules](#23-configuring-provider-alerting-rules). Assuming the labels sent by CortexMetrics are these:

```yaml
severity: WARNING
team: odpf
service: some-service
environment: integration
resource_name: some-resource
```

We want to subscribe all notifications owned by `odpf` team and has severity `WARNING` regardless the service name and route the notification to `file` with receiver id `2` (the one that we created in the [previous](#22-register-a-receiver) part).

Prepare a subscription detail and create a new subscription with Siren CLI.
```bash  title=cpu_subs.yaml
urn: subscribe-cpu-odpf-warning
namespace: 1
receivers:
  - id: 1
  - id: 2
match
  team: odpf
  severity: WARNING
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

Once a subscription is created, let's manually trigger alert in CortexMetrics with this cURL. The way CortexMetrics monitor a specific metric and auto-trigger an alert are out of this `tour` scope.

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

This is the end of `alerting rules and subscription` tour. If you want to know how to send on-demand notification to a receiver, you could check the [first tour](./1sending_notifications_overview.md).

Apart from the tour, we recommend completing the [guides](../guides/overview.md). You could also check out the remainder of the documentation in the [reference](../reference/server_configuration.md) and [concepts](../concepts/overview.md) sections for your specific areas of interest. We've aimed to provide as much documentation as we can for the various components of Siren to give you a full understanding of Siren's surface area. If you are interested to contribute, check out the [contribution](../contribute/contribution.md) page.
