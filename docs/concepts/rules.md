# Rules

Siren rules are creating from predefined [templates](templates.md) by providing values of the variables of the template.

One can create rules using either HTTP APIs or CLI. You can check out the swagger file which describes the schema of
request payload for HTTP APIs.

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

### Terminology

| Term        | Description                                                              | Example/Default  |
|-------------|--------------------------------------------------------------------------|------------------|
| API Version | Which API to use to parse the YAML file                                  | v2               |
| Type        | Describes the type of object represented by YAML file                    | rule             |
| Namespace   | Corresponds to Cortex namespace in which rule will be created            | kafka            |
| Entity      | Corresponds to tenant name in cortex                                     | odpf             |
| Rules       | Map of GroupNames describing what template is used in a particular group | See example file |
