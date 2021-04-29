# Usage

The following topics will describe how to use Siren.

## CLI Interface

1. Serve
    - Runs the Server  `$ go run main.go serve`

2. Migrate
    - Runs the DB Migrations `$ go run main.go migrate`

3. Upload
    - Parses a YAML File in specified format to upsert templates and rules(
      alerts) `$ go run main.go upload fileName.yaml`. Read more about the Rules and Templates [here](../concepts).

## Managing Templates

Siren templates are abstraction over prometeheus rules to reuse same rule body to create multiple rules. The rule body
is templated using go templates. Learn in more detail [here](./templates.md).

## Managing Rules

Siren rules are defined using a template by providing value for the variables defined inside that template. Learn in
more details [here](./rules.md)

## Managing bulk rules and templates

For org wide use cases, where teams need to manage multiple templates and rules Siren CLI can be highly useful. Think
GitOps but for alerting. Learn in More detail [here](./bulk_rules.md)

## Alerting Credentials Management

Siren stores slack and pagerduty credentials which can be used to configure Cortex Alertmanager to route alerts to Slack
and Pagerduty. Learn in more detail [here](./alert_credential.md).

## Alert History Subscription

Siren can configure CortexAlertmanager to call Siren back allowing storage of triggered alerts. This can be used for
auditing and analytics purposes. Learn in more detail [here](./alert_history.md).

## Deployment

Refer [here](./deployment.md) to learn how to deploy siren in production.

## Monitoring

Refer [here](./monitoring.md) to for more details on monitoring siren.

## Troubleshooting

Troubleshooting [guide](./troubleshooting.md). 
