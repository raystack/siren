# Use Cases

As an Incident Management Platform, Siren integrates with several monitoring and alerting providers (CortexMetrics, Prometheus, InfluxDB, etc) and orchestrates alerting rules in a simple DIY configuration. Siren capables to subscribe to alerts and send notifications based on the triggered alerts or sending on-demand notifications to the supported receivers (slack, pagerduty, etc).

This page describes some of Siren use cases and provides some related resources for better understanding. There might be some other use cases not mentioned in this page that are also suitable with Siren.

## Alerting

### Alerting Rules Orchestration

A rule is an expression that should be met, given the metrics, to trigger an alert. Each monitoring & alerting provider has its own way to define alerting rules and it is relatively easy to do so. However it does not give that much flexibility when the users and teams are getting bigger and there is a need to do self-serve alerting rules creation. Siren provides an abstraction on top of that to give more flexibility in creating alerting rules (via API, CLI, or a UI).

### Alerting Rules Templating

We noticed there are several times when multiple users or teams using the same rules with just different threshold numbers or labels. Creating multiple similar rules for different purposes is not necessary and would give more overhead to maintain. Siren provides [templating](./guides/template.md) feature to templatize rules given some variables so users could reuse the existing templates to define rules.

## Notification

### Alert Notifications Subscription

Most monitoring and alerting providers have their own feature to notify a specific channel when an alert is triggered. If an organization uses different monitoring and alerting providers, the responsibility to send notification would be passed on to the respective providers. The number of supported notification channels is also vary depends on the provider. This would give limitation if one needs to send a notification to the unsupported channel in a provider. 

With Siren, notification responsibility will be unified in Siren. This approach will be more maintainable and easier to audit. Siren handles all alert notification subscriptions where user could define subscriptions and Siren publishes notifications if the labels in subscriptions match with the labels in the triggered alerts.

Siren is also designed to be easily extended with a new notification channel as a [new receiver plugin](./extend/adding_new_receiver.md) to support more use cases.


### Sending On-demand Notification

There is also a case when a non-alert event needs to be sent as notification with a custom payload. Siren could be used to send on-demand [notifications](./guides/notification.md) too. One just need to pick to which receiver to send notifications too or create a new one if it does not exist yet and send a notification to it.
