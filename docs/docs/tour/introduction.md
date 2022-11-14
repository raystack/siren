# Introduction

This tour introduces you to Siren. Along the way you will learn how to manage alerting rules, notification receivers, and subscribing to notifications.

## Prerequisites

This tour requires you to have Siren CLI tool installed on your local machine. You can run `siren version` to verify the installation. Please follow [installation](../installation.md) and [configuration](../reference/server_configuration.md) guides if you do not have it installed already.

Siren client CLI talks to Siren server to configure and fetch rules, subscriptions, and notifications. Please make sure you also have a Siren server running. You can also run server locally with `siren server start` command. For more details check the [deployment](../guides/deployment.md) guide.

## Help
At any time you can run the following commands.

```
# Check the installed version for Siren cli tool
$ siren version

# See the help for a command
$ siren --help
```

The list of all available commands are as follows:

```
CORE COMMANDS
  alert           Manage alerts
  namespace       Manage namespaces
  provider        Manage providers
  receiver        Manage receivers
  rule            Manage rules
  subscription    Manage subscriptions
  template        Manage templates

ADDITIONAL COMMANDS
  completion      Generate shell completion scripts
  config          Manage siren CLI configuration
  help            Help about any command
  job             Manage siren jobs
  environment     List of supported environment variables
  reference       Comprehensive reference of all commands
  server          Run siren server
  worker          Start or manage Siren's workers
```

Help command can also be run on any sub command with syntax `siren <command> <subcommand> --help`. Here is an example for the same.

```
$ siren rule --help
```
Check the reference for Siren cli commands.

```
$ siren reference
```
## Background for this tutorial

This tour introduces you to two different scenarios
1. [Sending on-demand notification to a receiver](./1sending_notifications_overview.md)
    - Register a receiver
    - Send notification to the receiver
2. [Setting up alerting rules and subscribing to the alerts](2alerting_rules_subscriptions_overview.md)
    - Register a CortexMetrics provider
    - Create a new namespace
    - Register a receiver that we want to send the notification to
    - Create a subscription to define the routing so alert notification will be routed to the registered receivers

The tour takes approximately 20 minutes to complete.