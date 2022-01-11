# Architecture

Siren exposes HTTP API to allow rule, template and slack & Pagerduty credentials configuration. It talks to upstream
cortex ruler to configure rules(alerting and recording rules). It talks to Cortex Alertmanager to configure the
destination where alerts should go. It stores data around credentials, templates and current state of configured alerts
in PostgresDB. It also stores alerts triggered via Cortex Alertmanager.

![Siren Architecture](../../static/img/siren_architecture.svg)

## System Design

### Technologies

Siren is developed with

- Golang - Programming language
- Docker - container engine to start postgres and cortex to aid development
- Cortex - multi-tenant prometheus based monitoring stack
- Postgres - a relational database

### Components

_**GRPC Server and HTTP Gateway**_

* GRPC Server exposes RPC APIs and RESTfull APIs (via GRPC gateway) to allow configuration of rules, templates, alerting
  credentials and storing triggered alert history.

_**PostgresDB**_

* Used for storing the templates in a predefined schema enabling reuse of same rule body.
* Stores the rules configured via HTTP APIs and used for preserving thresholds when rule is deleted
* Stores Slack and Pagerduty credentials to enable DIY interface for configuring destinations for alerting.

_**Command Line Interface**_

* Provides a way to manage rules and templates using YAML files in a format described below.
* Run a web-server from CLI.
* Runs DB Migrations
* Manage templates, rules, providers, namespaces and receivers

### Managing Templates via YAML File

Siren gives flexibility to templatize prometheus rules for re-usability purpose. Template can be managed via APIs(REST
and GRPC). Apart from that, there is a command line interface as well which parses a YAML file in a specified format (as
described below) and upload to Siren using an HTTP Client of Siren Service. Refer [here](../guides/templates.md) for
more details around usage and terminology.

### Managing Rules via YAML File

To manage rules in bulk, Siren gives a way to manage rules using YAML files, which you can manage in a version
controlled repository. Using the `upload` command one can upload a rule YAML file in a specified format (as described
below) and upload to Siren using the GRPC Client(comes inbuilt) of Siren Service. Refer [here](../guides/rules.md) for
more details around usage and terminology.

**Note:** Updating a template also updates the associated rules.

## Siren Integration

The monitoring providers supported are:

1. Cortex metrics.

The section details all integrating systems for Siren deployment.

### Cortex Ruler

* The upstream Cortex ruler deployment which is used for rule creation in the proper namespace/group. You can create
  a [provider](../guides/providers.md) for that purpose and provide appropriate hostname.

### Cortex Alertmanager

* The upstream Cortex alertmanager deployment where routing configurations are stored in the proper format. Sirenstores
  subscriptions which gets synced in the alertmanager. Cortex Alertmanger hostname is fetched
  from [provider's](../guides/providers.md) host key. 
  