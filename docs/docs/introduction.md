# Introduction

Siren provides alerting on metrics of your applications using [Cortex metrics](https://cortexmetrics.io/) in a simple
DIY configuration. With Siren, you can define templates(using go templates standard), and create/edit/enable/disable
prometheus rules on demand. It also gives flexibility to manage bulk of rules via YAML files. Siren can be integrated
with any client such as CI/CD pipelines, Self-Serve UI, microservices etc.

![Siren Architecture](/img/overview.svg)

## Key Features

- **Rule Templates:** Siren provides a way to define templates over prometheus Rule, which can be reused to create
  multiple instances of same rule with configurable thresholds.
- **Multi-tenancy:** Rules created with Siren are by default multi-tenancy aware.
- **DIY Interface:** Siren can be used to easily create/edit prometheus rules. It also provides soft delete(disable)
  so that you can preserve thresholds in case you need to reuse the same alert.
- **Managing bulk rules:** Siren enables users to manage bulk alerts using YAML files in specified format using simple
  CLI.
- **Slack Notifications:** Siren can be used to send out notifications in slack to user's DM, private and public
  channel.
- **Credentials Management:** Siren can store slack and pagerduty credentials, sync them with Cortex alertmanager to
  deliver alerts on proper channels, in a multi-tenant fashion. It gives a simple interface to rotate the credentials on
  demand via HTTP API.
- **Alert History:** Siren can store alerts triggered via Cortex Alertmanager, which can be used for audit purposes.

## Usage

Explore the following resources to get started with Siren:

- [Guides](./guides/overview.md) provides guidance on usage.
- [Concepts](./concepts/overview.md) describes all important Siren concepts including system architecture.
- [Reference](./reference/configuration.md) contains the details about configurations and other aspects of Siren.
- [Contribute](./contribute/contribution.md) contains resources for anyone who wants to contribute to Siren.
