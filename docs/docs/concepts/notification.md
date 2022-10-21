# Notification

Notification is one of main features in Siren. Siren capables to send notification to various receivers (e.g. Slack, PagerDuty). Notification in Siren could be sent directly to a receiver or user could subscribe notifications by providing key-value label matchers. For the latter, Siren routes notification to specific receivers by matching notification key-value labels with the provided label matchers.

Below is how the notification is implemented in Siren

![Notification in Siren](../../static/img/siren_notification.svg)

**Notification Source** is a point where a notification is generated. In Siren, there are two points of notification source: `/receivers/{id}/notify` API and alerts hook API. The first one will generate and publish notification when the API is invoked with some message in its payload and the latter one will always generate and publish notification everytime the hook API is called by the provider.

## Notification Message Payload

There is currently no abstraction on notification payload so user needs to pass notification message payload in the same format as what receiver (notification vendor) expected. See [reference](../reference/receiver.md) on how the format for each receiver is expected and how to send notification.

## Templating Notification Message Payload

Message payload in notification could also be reused by defining template and passing some variables needed. See [template](../guides/template.md) for further details on how to use the template feature.

## Queue

Queue is used as a buffer to avoid pressure when notifications are being sent. Siren implements Queue as a plugin. Currently there are two kind of queue plugin supported: in-memory (not for production usage) and postgres. User could choose the which queue to use by mentioning it in the [config](../reference/server_configuration.md).

### In-memory Queue

In-memory queue simulates a queue with Go channel. The usage is intended to be used in development only.