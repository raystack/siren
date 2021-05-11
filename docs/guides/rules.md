# Rules

Siren rules are creating from predefined [templates](templates.md) by providing values of the variables of the template.

One can create templates using either HTTP APIs or CLI.

## API interface

A rule is uniquely identified with the combination of namespace, entity(tenant name), group name and template name.

One can choose any namespace and group name. A cortex ruler namespace is a collection of groups. Groups can have many
one or more rules.

**Creating/Updating a rule**

The below snippet describes an example of template creation/update. Same API can be used to enable or disable alerts.

```text
PUT /rules HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 298

{
  "namespace": "odpf",
  "group_name": "CPUHigh",
  "template": "CPU",
  "entity": "odpf",
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
  "status": "enabled"
}
```

Here we are using CPU template and providing value for few variables. In case some variables value is not provided
default will be picked.

### Terminology of the request body

| Term        | Description                                                              | Example/Default  |
|-------------|--------------------------------------------------------------------------|------------------|
| Namespace   | Corresponds to Cortex namespace in which rule will be created            | kafka            |
| Group Name  | Corresponds to Cortex group name in which rule will be created           | CPUHigh            |
| Entity      | Corresponds to tenant name in cortex                                     | odpf             |
| Rules       | Map of GroupNames describing what template is used in a particular group | See example above |
| Variables   | Value of variables defined inside the template                           | See example above |

**Fetching rules**

Rules can be fetched and filtered with multiple parameters. An example of all filters is described below.

```text
GET /rules?namespace=foo&entity=odpf&group_name=CPUHigh&status=enabled&template=CPU HTTP/1.1
Host: localhost:3000
```

## CLI Interface

With CLI, you will need a YAML file in the below specified format to create/edit rules. The CLI calls Siren service
rules APIs in turn.

**Example rule file**

```yaml
apiVersion: v2
type: rule
namespace: demo
entity: odpf
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
under `odpf` [tenant](https://cortexmetrics.io/docs/architecture/#the-role-of-prometheus) as denoted by namespace and
entity keys.

The rules array defines actual rules defined over the templates. Here `TestGroup` is the name of the group which will be
created/updated with the rule defined by `CPU` template. The example shows the value of variables provided in creating
rules(alert).

**Example upload command**

```shell
go run main.go upload cpu_rule.yaml
```

The yaml file can be edited and re-uploaded to edit the rule thresholds.

### Terminology

| Term        | Description                                                              | Example/Default  |
|-------------|--------------------------------------------------------------------------|------------------|
| API Version | Which API to use to parse the YAML file                                  | v2               |
| Type        | Describes the type of object represented by YAML file                    | rule             |
| Namespace   | Corresponds to Cortex namespace in which rule will be created            | kafka            |
| Entity      | Corresponds to tenant name in cortex                                     | odpf             |
| Rules       | Map of GroupNames describing what template is used in a particular group | See example file |
