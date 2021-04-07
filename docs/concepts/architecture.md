# Architecture

Siren exposes HTTP API to allow rule, template and slack & Pagerduty credentials configuration. It talks to upstream
cortex ruler to configure rules(alerting and recording rules). It talks to Cortex Alertmanager to configure destination
where alerts should go. It stores data around credentials, templates and current state of configured alerts in
PostgresDB.

![Siren Architecture](../assets/siren.jpg)

## System Design

### Components

_**HTTP Server**_

* HTTP Server exposes Restful APIs to allow configuration of rules and templates and alerting credentials.

_**PostgresDB**_

* Used for storing the templates in a predefined schema enabling reuse of same rule body.
* Stores the rules configured via HTTP APIs and used for preserving thresholds when rule is deleted
* Stores Slack and Pagerduty credentials to enable DIY interface for configuring destinations for alerting.

_**Command Line Interface**_

* Provides a way to manage rules and templates using YAML files in format described below.
* Run a web-server from CLI.
* Runs DB Migrations

### Managing Templates via YAML File

Siren gives flexibility to templatise prometheus rules for re-usability purpose. Template can be managed via HTTP APIs  
at`/templates` route in a restful manner. Apart from that, there is a command line interface as well which parses a YAML
file in a specified format (as described below) and upload to Siren using an HTTP Client of Siren Service.
Refer [here](templates.md) for more details around usage and terminology.

### Managing Rules via YAML File

To manage rules in bulk, Siren gives a way to manage rules using YAML files, which you can manage in a version
controlled repository. Using the `upload` command one can upload a rule YAML file in a specified format (as described
below) and upload to Siren using an HTTP Client(comes inbuilt) of Siren Service. Refer [here](rules.md) for more details
around usage and terminology.

**Note:** Updating a template also updates the associated rules.

## Siren Integration

The section details all integrating systems for Siren deployment. These are external systems that Siren connects to.

### Cortex Ruler

* The upstream Cortex ruler deployment which is used for rule creation in proper namespace/group.
  The [`cortex_host`](../reference/configuration.md#-cortex.address) config can be set in Siren.

### Cortex Alertmanager

* The upstream Cortex alertmanager deployment where slack and pagerduty credentials are stored in the proper format.
  The [`cortex_host`](../reference/configuration.md#-cortex.address) config can be set in Siren.
