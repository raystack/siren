# Rules

Siren rules are creating from predefined [templates](templates.md) by providing values of the variables of the template.

One can create templates using either HTTP APIs or CLI.

## API interface

A rule is uniquely identified with the combination of provider's namespace(uniquely identifies which provider and
namespace), template name, optional namespace, optional group name.

One can choose any namespace and group name. In cortex terminology, namespace is a collection of groups. Groups can have
one or more rules.

**Creating/Updating a rule**

The below snippet describes an example of rule creation/update. Same API can be used to enable or disable alerts.

```text
PUT /v1beta1/rules HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 298

{
  "namespace": "odpf",
  "group_name": "CPUHigh",
  "template": "CPU",
  "providerNamespace": "3"
  "variables": [
    {
      "name": "for",
      "value": "15m",
      "type": "string"
    },
     {
      "name": "team",
      "value": "odpf",
      "type": "string"
    }
  ],
  "enabled": true,
}
```

Here we are using CPU template and providing value for few variables("for", "team"). In case some variables value is not
provided default will be picked from the template's definition.

### Terminology of the request body

| Term              | Description                                                    | Example/Default   |
|-------------------|----------------------------------------------------------------|-------------------|
| Namespace         | Corresponds to Cortex namespace in which rule will be created  | kafka             |
| Group Name        | Corresponds to Cortex group name in which rule will be created | CPUHigh           |
| providerNamespace | Corresponds to a tenant in a provider                          | 4                 |
| Template          | what template is used to create the rule                       | CPU               |
| Variables         | Value of variables defined inside the template                 | See example above |
| Enabled           | boolean describing if the rule is enabled or not               | true              |

**Fetching rules**

Rules can be fetched and filtered with multiple parameters. An example of all filters is described below.

```text
GET /v1beta1/rules?namespace=foo&providerNamespace=4&group_name=CPUHigh&template=CPU HTTP/1.1
Host: localhost:3000
```

## CLI Interface

```text
Work with rules.

rules are used for alerting within a provider.

Usage:
  siren rule [command]

Aliases:
  rule, rules

Available Commands:
  edit        Edit a rule
  list        List rules
  upload      Upload Rules YAML file

Flags:
  -h, --help   help for rule

Use "siren rule [command] --help" for more information about a command.
```

With CLI, you will need a YAML file in the below specified format to create/edit rules.
**Example rule file**

```yaml
apiVersion: v2
type: rule
namespace: demo
provider: localhost-cortex
providerNamespace: test
rules:
  TestGroup:
    template: CPU
    status: enabled
    variables:
      - name: for
        value: 15m
      - name: warning
        value: 185
      - name: critical
        value: 195
```

In the above example, we are creating rules
inside `demo` [namespace](https://cortexmetrics.io/docs/api/#get-rule-groups-by-namespace)
under `test` [tenant](https://cortexmetrics.io/docs/architecture/#the-role-of-prometheus) of `localhost-cortex`
provider.

The rules array defines actual rules defined over the templates. Here `TestGroup` is the name of the group which will be
created/updated with the rule defined by `CPU` template. The example shows the value of variables provided in creating
rules(alert).

**Example upload command**

```shell
go run main.go rule upload cpu_rule.yaml
```

The yaml file can be edited and re-uploaded to edit the rule thresholds.

### Terminology

| Term              | Description                                                              | Example/Default   |
|-------------------|--------------------------------------------------------------------------|-------------------|
| API Version       | Which API to use to parse the YAML file                                  | v2                |
| Type              | Describes the type of object represented by YAML file                    | rule              |
| Namespace         | Corresponds to Cortex namespace in which rule will be created            | kafka             |
| Entity            | Corresponds to tenant name in cortex                                     | odpf              |
| Rules             | Map of GroupNames describing what template is used in a particular group | See example file  |
| Variables         | Value of variables defined inside the template                           | See example above |
| provider          | URN of monitoring provider to be used                                    | localhost-cortex  |
| providerNamespace | URN of tenant to choose inside the monitoring provider                   | test              |
