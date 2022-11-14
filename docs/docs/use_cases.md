# Use Cases

Siren is an Incident Management Platform. Siren integrates with several monitoring and alerting providers (CortexMetrics, Prometheus, InfluxDB, etc...) and orchestrates alerting rules in a simple DIY configuration. Apart from that, siren could also send notifications based on the triggered alerts or sending on-demand notifications to the supported receivers (slack, pagerduty, etc...).

This page describes some of Siren use cases and provides some related resources for better understanding. There might be some other use cases not mentioned in this page that are also suitable with Siren.

## Alerting Rules Orchestration

A rule is an expression that should be met, given the metrics, to trigger an alert. Each monitoring & alerting provider has its own way to define alerting rules and it is relatively easy to do so. However, the basic solution of this does not give good flexibility when the users and teams are getting bigger and there is a need to do self-serve alerting rules creation. Siren's role is to provide an abstraction on top of that so each user could create alerting rules in self-serve basis (via API or CLI or a UI).

## Alerting Rules Templating

We noticed there are several times when multiple users or teams using the same rules with just different threshold numbers or labels. Creating multiple similar rules for different purpose is not necessary and would give overhead to maintain. Siren provides [templating](./guides/template.md) feature to templatize rules given some variables so users could reuse the existing templates one if that is suitable for them.

## Alert Notifications Subscription

Most monitoring and alerting providers have their own feature to notify a specific channel when an alert is triggered. If an organization uses different monitoring and alerting providers, the responsibility to send notification would be passed on to the respective providers. With Siren, notification responsibility will be unified in Siren. This approach will be more maintainable and easier for auditing the notifications. Siren handles all alert notification subscriptions where user could define subscriptions and Siren publishes notifications if the labels in subscriptions match with the labels in the triggered alerts.

## Sending On-demand Notification

There is also a case when a non-alert event needs to be sent as notification with a custom payload. Siren could be used to send on-demand [notifications](./guides/notification.md) too. You just need to pick to which receiver that is registered in Siren you want to send notifications too or create a new one if it does not exist yet and send a notification to it.
