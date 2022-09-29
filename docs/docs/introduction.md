# Introduction

Siren orchestrates alerting rules of your applications using a monitoring and alerting provider e.g. [Cortex metrics](https://cortexmetrics.io/) and sending notifications in a simple DIY configuration. With Siren, you can define templates (using go templates standard), create/edit/enable/disable alerting rules on demand, and sending out notifications. It also gives flexibility to manage bulk of rules via YAML files. Siren can be integrated with any client such as CI/CD pipelines, Self-Serve UI, microservices etc.

![Siren Overview](/img/overview.svg)

## Key Features

- **Rule Templates:** Siren provides a way to define templates over alerting rule which can be reused to create multiple instances of the same rule with configurable thresholds.
- **Multi-tenancy:** Rules created with Siren are by default multi-tenancy aware.
- **DIY Interface:** Siren can be used to easily create/edit alerting rules. It also provides soft delete(disable) so that you can preserve thresholds in case you need to reuse the same alert.
- **Managing bulk rules:** Siren enables users to manage bulk alerting rules using YAML files in specified format with simple CLI.
- **Receivers:** Siren can be used to send out notifications to several channels (e.g. slack, pagerduty, email etc).
- **Subscriptions** Siren can be used to subscribe to notifications (with desired matching conditions) via the channel of your choice.
- **Alert History:** Siren can store alerts triggered by monitoring & alerting provider e.g. Cortex Alertmanager, which can be used for audit purposes.

## Usage

Explore the following resources to get started with Siren:

- [Tour](tour/introduction.md) allows you to explore Siren features quickly.
- [Concepts](concepts/overview.md) describes all important Siren concepts including system architecture.
- [Guides](guides/overview.md) provides guidance on usage.
- [Reference](reference/server_configuration.md) contains the details about configurations and other aspects of Siren.
- [Contribute](contribute/contribution.md) contains resources for anyone who wants to contribute to Siren.
