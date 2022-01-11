# Usage

The following topics will describe how to use Siren.

## CLI Interface

```text
Siren provides alerting on metrics of your applications using Cortex metrics
in a simple DIY configuration. With Siren, you can define templates(using go templates), and
create/edit/enable/disable prometheus rules on demand.

Available Commands:
  alert       Manage alerts
  completion  generate the autocompletion script for the specified shell
  config      manage siren CLI configuration
  help        Help about any command
  migrate     Migrate database schema
  namespace   Manage namespaces
  provider    Manage providers
  receiver    Manage receivers
  rule        Manage rules
  serve       Run server
  template    Manage templates
```

## Managing providers and multi-tenancy

Siren can be used define alerts and their routing configurations inside monitoring "providers". List of supported
providers:

- CortexMetrics.

Support for other providers is also planned, feel free to contribute. Siren also respects the multi-tenancy provided by
various monitoring providers using "namespaces". Namespace simply represents a tenant inside your provider. Learn in
more detail [here](./providers.md).

## Managing Templates

Siren templates are abstraction over Prometheus rules to reuse same rule body to create multiple rules. The rule body is
templated using go templates. Learn in more detail [here](./templates.md).

## Managing Rules

Siren rules are defined using a template by providing value for the variables defined inside that template. Learn in
more details [here](./rules.md)

## Managing bulk rules and templates

For org wide use cases, where teams need to manage multiple templates and rules Siren CLI can be highly useful. Think
GitOps but for alerting. Learn in More detail [here](./bulk_rules.md)

## Receivers

Receivers represent a notification medium, which can be used to define routing configuration in the monitoring
providers, to control the behaviour of how your alerts are notified. Few examples: Slack receiver, HTTP receiver,
Pagerduty receivers. You can use receivers to send notifications on demand as well as on certain matching conditions.
Learn in more detail [here](./receivers.md).

## Subscriptions

Siren can be used to configure various monitoring providers to route your alerts to proper channels based on your match
conditions. You define your own set of selectors and subscribe to alerts matching these selectors in the notification
mediums of your choice. Learn in more detail [here](./subscriptions.md).

## Alert History Subscription

Siren can configure Cortex Alertmanager to call Siren back, allowing storage of triggered alerts. This can be used for
auditing and analytics purposes. Alert History is simply a "subscription" defined using an "HTTP receiver" on all
alerts.

Learn in more detail [here](./alert_history.md).

## Deployment

Refer [here](./deployment.md) to learn how to deploy Siren in production.

## Monitoring

Refer [here](./monitoring.md) to for more details on monitoring siren.

## Troubleshooting

Troubleshooting [guide](./troubleshooting.md).
