import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# Template

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

Templates concept in Siren is used for abstraction. The usage is versatile enough to be used to abstract out rules and notification format. It utilises [go-templates](https://golang.org/pkg/text/template/) to provide data-driven templates for generating textual output. The template delimiter used is `[[` and `]]`.

One can create templates using either HTTP APIs or CLI.

## API interface

### Template creation/update

Templates can be created using Siren APIs. The below snippet describes an example.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren template upsert --file template.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request PUT
  --url `}{defaultHost}{`/`}{apiVersion}{`/templates
  --header 'content-type: application/json'
  --data-raw '{
    "name": "CPU",
    "body": "- alert: CPUHighWarning\n  expr: avg by (host) (cpu_usage_user{cpu=\"cpu-total\"}) > [[.warning]]\n  for: '[[.for]]'\n  labels:\n    severity: WARNING\n    team: '[[ .team ]]'\n  annotations:\n    dashboard: https://example.com\n    description: CPU has been above [[.warning]] for last [[.for]] {{ $labels.host }}\n- alert: CPUHighCritical\n  expr: avg by (host) (cpu_usage_user{cpu=\"cpu-total\"}) > [[.critical]]\n  for: '[[.for]]'\n  labels:\n    severity: CRITICAL\n    team: '[[ .team ]]'\n  annotations:\n    dashboard: example.com\n    description: CPU has been above [[.critical]] for last [[.for]] {{ $labels.host }}\n",
    "tags": [
        "firehose",
        "dagger"
    ],
    "variables": [
        {
            "name": "team",
            "type": "string",
            "default": "odpf",
            "description": "Name of the team that owns the deployment"
        },
        {
            "name": "for",
            "type": "string",
            "default": "10m",
            "description": "For eg 5m, 2h; Golang duration format"
        },
        {
            "name": "warning",
            "type": "int",
            "default": "85",
            "description": ""
        },
        {
            "name": "critical",
            "type": "int",
            "default": "95",
            "description": ""
        }
    ]
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Terminology of the request body

| Term      | Description                                                                                               | Example/Default   |
| --------- | --------------------------------------------------------------------------------------------------------- | ----------------- |
| Name      | Name of the template                                                                                      | CPUHigh           |
| Body      | Array of rule body. The body can be templatized in go template format.                                    | See example above |
| Variables | Array of variables that were templatized in the body with their data type, default value and description. | See example above |
| Tags      | Array of resources/applications that can utilize this template                                            | VM                |

The response body will look like this:

```json
{
  "id": 38,
  "CreatedAt": "2021-04-29T16:20:48.061862+05:30",
  "UpdatedAt": "2021-04-29T16:22:19.978837+05:30",
  "name": "CPU",
  "body": "- alert: CPUHighWarning\n  expr: avg by (host) (cpu_usage_user{cpu=\"cpu-total\"}) > [[.warning]]\n  for: '[[.for]]'\n  labels:\n    severity: WARNING\n    team: '[[ .team ]]'\n  annotations:\n    dashboard: https://example.com\n    description: CPU has been above [[.warning]] for last [[.for]] {{ $labels.host }}\n- alert: CPUHighCritical\n  expr: avg by (host) (cpu_usage_user{cpu=\"cpu-total\"}) > [[.critical]]\n  for: '[[.for]]'\n  labels:\n    severity: CRITICAL\n    team: '[[ .team ]]'\n  annotations:\n    dashboard: example.com\n    description: CPU has been above [[.critical]] for last [[.for]] {{ $labels.host }}\n",
  "tags": ["firehose", "dagger"],
  "variables": [
    {
      "name": "team",
      "type": "string",
      "default": "odpf",
      "description": "Name of the team that owns the deployment"
    },
    {
      "name": "for",
      "type": "string",
      "default": "10m",
      "description": "For eg 5m, 2h; Golang duration format"
    },
    {
      "name": "warning",
      "type": "int",
      "default": "85",
      "description": ""
    },
    {
      "name": "critical",
      "type": "int",
      "default": "95",
      "description": ""
    }
  ]
}
```

### Fetching a template

**Fetching by Name**

Here is an example to fetch a template using name.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren template view cpu
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/templates/cpu`}
    </CodeBlock>
  </TabItem>
</Tabs>

**Fetching by Tags**

Here is an example to fetch a templates matching the tag.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren template list --tag firehose
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/templates?tag=firehose`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Deleting a template

```text
DELETE /v1beta1/templates/cpu HTTP/1.1
Host: localhost:3000
```

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren template delete cpu
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request DELETE
  --url `}{defaultHost}{`/`}{apiVersion}{`/templates/cpu`}
    </CodeBlock>
  </TabItem>
</Tabs>


**Note:**

1. Updating a template via API will not upload the associated rules.

## CLI interface

With CLI, you will need a YAML file in the below specified format to create/update templates. The CLI calls Siren
service templates APIs in turn.

**Example template file**

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

In the above example, we are using one template to define rules of two severity labels viz WARNING and CRITICAL. Here we
have made 3 templates variables `for`, `warning` and `critical` which denote the appropriate alerting thresholds. They
will be given a value while actual rule(alert) creating.

```bash
siren template upload cpu_template.yaml
```

### Terminology

| Term        | Description                                                                                               | Example/Default  |
| ----------- | --------------------------------------------------------------------------------------------------------- | ---------------- |
| API Version | Which API to use to parse the YAML file                                                                   | v2               |
| Type        | Describes the type of object represented by YAML file                                                     | template         |
| Name        | Name of the template                                                                                      | CPUHigh          |
| Body        | Array of rule body. The body can be templatized in go template format.                                    | See example file |
| Variables   | Array of variables that were templatized in the body with their data type, default value and description. | See example file |
| Tags        | Array of resources/applications that can utilize this template                                            | VM               |

**Note:**

1. It's suggested to always provide default value for the templated variables.
2. Updating a template used by rules via CLI will update all associated rules.
