# Usage

The following topics will describe how to use Siren.

## CLI Interface

```text
Siren provides alerting on metrics of your applications using Cortex metrics
in a simple DIY configuration. With Siren, you can define templates(using go templates), and
create/edit/enable/disable prometheus rules on demand.

Available Commands:
  alert          Manage alerts
  completion     Generate shell completion scripts
  config         Manage siren CLI configuration
  help           Help about any command
  job            Manage siren jobs
  namespace      Manage namespaces
  provider       Manage providers
  receiver       Manage receivers
  rule           Manage rules
  server         Run siren server
  template       Manage templates
  worker         Start or manage Siren's workers
```

## Managing providers and multi-tenancy

Siren can be used to define alerting rules inside monitoring `providers`. List of supported
providers:

- [CortexMetrics](http://cortexmetrics.io).

Support for other providers is on the roadmap, feel free to contribute. Siren also respects the multi-tenancy provided by various monitoring providers using `namespaces`. A `namespace` represents a tenant inside your provider. Learn in more detail [here](./provider_and_namespace.md).

## Managing Templates

Siren templates are abstraction to make data definition reusable (e.g. Prometheus rules to reuse same rule body to create multiple rules). Template could be used to define alerting rule body with go templates. Learn in more detail [here](./template.md).

## Managing Rules

Siren rules are defined using a template by providing value for the variables defined inside that template. Learn in more details [here](./rule.md)

## Managing bulk rules and templates

For org wide use cases, where teams need to manage multiple templates and rules Siren CLI can be highly useful. Think GitOps but for alerting. Learn in More detail [here](./rule.md#bulk-rule-management)

## Notifications

Siren capables to send notifications which could route a notification into a specific channel defined by a [receiver](./receiver.md). Siren uses key-value label matching for routing. There are two kind of notification route method, `direct notification to receiver` and `subscription-based notification`.

- **Direct Notification to Receiver:** Siren provides a way to the user to send direct notification to a supported receiver by calling an API `/receivers/{receiver_id}/notify` with a specific payload message.

- **Subscription-based Notification:** The subscription-based notification is currently only works for triggered alerts. For each alerts coming to Siren's webhook, a notification will be generated and routed to specific [receivers](./receiver.md) based on the [subscriptions](./subscription.md).

## Receivers

Receivers represent a notification medium e.g. Slack, PagerDuty, HTTP, which can be used to define routing configuration in Siren to control the behaviour of how the notifications are notified. You can use receivers to send notifications on demand as well as on certain matching conditions. Learn in more detail [here](./receiver.md).

## Subscriptions

Siren can be used to route your notifications (non-alerts or alerts notification) to proper channels (receivers) based on the labels match conditions. You define your own set of label matchers and subscribe to alerts matching these labels in the notification mediums of your choice. Learn in more detail [here](./subscription.md).

## Alert History Subscription

Siren expect `provider` to call Siren back when an alert is triggered, allowing storage of triggered alerts and sending notification via Siren. Storing triggered alerts is beneficial to be used for auditing and analytics purposes. Alert History is simply a `subscription` defined using an `HTTP receiver` on all alerts. Learn in more detail [here](./alert_history.md).

## Deployment, Monitoring, & Troubleshooting

Refer [here](./deployment.md) to learn how to deploy Siren in production.
