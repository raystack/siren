# Templates

Siren templates are an abstraction
over [Prometheus rules](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/). It
utilises [gotemplates](https://golang.org/pkg/text/template/) to provide implements data-driven templates for generating
textual output. The template delimeter used is `[[` and `]]`.

One can create template using either HTTP APIs or CLI. You can check out the swagger file which describes the schema of
request payload for HTTP APIs.

With CLI, you will need a YAML file in the below specified format to create/edit templates. The CLI calls Siren service
templates APIs in turn.

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
are have made 3 templates variables `for`, `warning` and `critical` which denote the appropriate alerting thresholds.
They will be given a value while actual rule(alert) creating.

**Example upload command**

```shell
go run main.go upload cpu_template.yaml
```

### Terminology

| Term        | Description                                                                                                | Example/Default  |
|-------------|------------------------------------------------------------------------------------------------------------|------------------|
| API Version | Which API to use to parse the YAML file                                                                    | v2               |
| Type        | Describes the type of object represented by YAML file                                                      | template         |
| Name        | Name of the template                                                                                       | CPUHigh          |
| Body        | Array of rule body. The body can be templatized in go template format.                                     | See example file |
| Variables   | Array of variables that were templatized in the body with their data type, default value and description.  | See example file |
| Tags        | Array of resources/applications that can utilize this template                                             | VM               |

**Note:**

1. It's suggested to always provide default value for the templated variables. 
