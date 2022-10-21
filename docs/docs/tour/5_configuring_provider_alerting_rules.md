import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# 5 - Configuring Provider Alerting Rules

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

In this part we will create alerting rules for our Cortex monitoring provider. Rules in Siren relies on [template](../guides/template.md) for its abstraction. We need to create a rule's template first before uploading a rule.

## Creating a Rule's Template

For now, we will create a rule's template to monitor CPU usage. More details about template is [here](../guides/template.md).

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

We named the template above as `CPU`, the body in the template is the data that will be interpolated with variables and rendered. Notice that template body format is similar with [Prometheus alerting rules](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/). This is because Cortex uses the same rules as prometheus and Siren will translate the rendered rules to the Cortex alerting rules. Let's save the template above into a file called `cpu_template.yaml`.

Let's upload our template to Siren using Siren CLI.
```shell
./siren template upload cpu_template.yaml
```

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
./siren template upload cpu_template.yaml
```

  </TabItem>
</Tabs>

You could verify the newly created template.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
./siren template list
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

## Creating a Rule

Now we already have a `CPU` template, we can create a rule based on that template. Let's prepare a rule and save it in a file called `cpu_test.yaml`.
```yaml
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
Above we defined a rule based on `CPU` template for namespace urn `odpf-ns` and provider urn `localhost-dev-cortex`. The rule group name is `cpuGroup` and there are also some variables to be assign to the template when the template is rendered.

Upload the rule with Siren CLI.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
./siren rule upload cpu_test.yaml
```
If succeed, you will get this message.
```shell
Upserted Rule
ID: 4
```
  </TabItem>
</Tabs>



You could verify the created rules by getting the added rules from Cortex provider with cURL.

<Tabs groupId="api">
  <TabItem value="http" label="HTTP">

  ```bash
  curl --location --request GET 'http://localhost:9009/api/v1/rules' \
  --header 'X-Scope-OrgId: odpf-ns'
  ```

  </TabItem>
</Tabs>

The response body should be in `yaml` format and like this
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