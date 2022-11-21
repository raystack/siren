# Overview

The following topics will describe how to use Siren.

## Managing providers and multi-tenancy

Siren can be used to define alerting rules inside monitoring `provider`. List of supported providers are [here](../concepts/plugin.md#provider-plugin). Support for other providers is on the roadmap, feel free to [contribute](../extend/adding_new_provider.md). Siren also respects the multi-tenancy provided by various monitoring providers using [namespace](./provider_and_namespace.md#namespace). A `namespace` represents a tenant inside a provider. Learn in more detail [here](./provider_and_namespace.md).

## Managing Templates

Siren templates are abstraction to make data definition reusable (e.g. Prometheus rules to reuse same rule body to create multiple rules). Template could be used to define alerting rule body with go templates. Learn in more detail [here](./template.md).

## Managing Rules

Siren rules are defined using a template by providing value for the variables defined inside that template. Learn in more details [here](./rule.md).

## Managing bulk rules and templates

For org wide use cases, where teams need to manage multiple templates and rules Siren CLI can be highly useful. Think GitOps but for alerting. Learn in more detail [here](./rule.md#bulk-rule-management).

## Alerts Subscription

Siren capables to subscribe to an alert and route notifications to the registered receivers in a subscription. Siren uses key-value label matching for routing. For each alerts coming to Siren's webhook, a notification will be generated and routed to specific [receivers](./receiver.md) based on the [subscriptions](./subscription.md).

## Sending On-demand Notifications

Siren provides a way to the user to send direct notification to a supported [receiver](./receiver.md) by calling an API `/receivers/{receiver_id}/send` with a custom payload message.

## Receivers

Receivers represent a notification medium e.g. Slack, PagerDuty, HTTP, which can be used in a [subscription](./subscription.md) to define notification routing configuration in Siren. You can use receivers to send notifications on demand as well as on certain matching conditions. Learn in more detail [here](./receiver.md).

## Subscriptions

Siren can be used to route your notifications (non-alerts or alerts notification) to proper channels (receivers) based on the labels match conditions. You define your own set of label matchers and subscribe to alerts matching these labels in the notification mediums of your choice. Learn in more detail [here](./subscription.md).

## Alert History Subscription

Siren expect [provider](./provider_and_namespace.md) to call Siren back when an alert is triggered, allowing storage of triggered alerts and sending susbcribed notification via Siren. Storing triggered alerts is beneficial to be used for auditing and analytics purposes. Learn in more detail [here](./alert_history.md).

## Deployment

Refer [here](./deployment.md) to learn how to deploy Siren in production.
