# Glossary

## Provider

Monitoring and Alerting Provider. E.g. Cortexmetrics, Influx 2.0, Prometheus etc.

## Namespace

Used to represent multi-tenancy in a given provider. Cortex calls it a tenant, influx 2.0 calls it organization.

## Receiver

Receivers are alert routing and notification mediums. Examples: Slack, Pagerduty, HTTP POST Webhooks etc. They contain global level configs which enables clients to use this medium for alert routing or sending notifications.

## Rule

Alerting rules that are set within a provider

## Template

An abstraction of data in Siren that could make data definition reusable. Currently template can be used to define [rules](../guides/rule.md) and [notification's](../guides/notification.md) message body.

## Subscription

Using subscriptions one can get notified when a set of conditions are true on a triggered alert.

## Notification

A message to be sent to the specific receivers. Notification could be sent directly to receivers or sent by matching subscription's labels.

## Alert History

Triggered Alert History. Siren provides simple endpoints to accept alert trigger event from various alerting providers e.g. Prometheus Alertmanager, Kapacitor, Influx 2.0 etc.

## Notification Vendor

External parties that has capability to communicates to the end-user with their own medium e.g. Slack, PagerDuty.

## Notification Message Payload

Notification Message Payload is the data that are being sent to the notification vendor in the format that meets notification vendor's contract.