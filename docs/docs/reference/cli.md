# CLI

## `siren alert`

Manage alerts

### `siren alert list [flags]`

List alerts

```
--end-time uint          end time
--provider-id uint       provider id
--provider-type string   provider type
--resource-name string   resource name
--start-time uint        start time
````

## `siren completion [bash|zsh|fish|powershell]`

Generate shell completion scripts

## `siren config <command>`

Manage siren CLI configuration

### `siren config init`

Initialize CLI configuration

### `siren config list`

List client configuration settings

## `siren environment`

List of supported environment variables

## `siren job <command>`

Manage siren jobs

### `siren job run`

Trigger a job

#### `siren job run cleanup_queue [flags]`

Cleanup stale messages in queue

```
-c, --config string   Config file path (default "config.yaml")
````

## `siren namespace`

Manage namespaces

### `siren namespace create [flags]`

Create a new namespace

```
-f, --file string   path to the namespace config
````

### `siren namespace delete`

Delete a namespace details

### `siren namespace edit [flags]`

Edit a namespace

```
-f, --file string   Path to the namespace config
    --id uint       namespace id
````

### `siren namespace list`

List namespaces

### `siren namespace view [flags]`

View a namespace details

```
--format string   Print output with the selected format (default "yaml")
````

## `siren provider`

Manage providers

### `siren provider create [flags]`

Create a new provider

```
-f, --file string   path to the provider config
````

### `siren provider delete`

Delete a provider details

### `siren provider edit [flags]`

Edit a provider

```
-f, --file string   Path to the provider config
    --id uint       provider id
````

### `siren provider list`

List providers

### `siren provider view [flags]`

View a provider details

```
--format string   Print output with the selected format (default "yaml")
````

## `siren receiver`

Manage receivers

### `siren receiver create [flags]`

Create a new receiver

```
-f, --file string   path to the receiver config
````

### `siren receiver delete`

Delete a receiver details

### `siren receiver edit [flags]`

Edit a receiver

```
-f, --file string   Path to the receiver config
    --id uint       receiver id
````

### `siren receiver list`

List receivers

### `siren receiver send [flags]`

Send a receiver notification

```
-f, --file string   Path to the receiver notification message
    --id uint       receiver id
````

### `siren receiver view [flags]`

View a receiver details

```
--format string   Print output with the selected format (default "yaml")
````

## `siren rule`

Manage rules

### `siren rule edit [flags]`

Edit a rule

```
-f, --file string   Path to the rule config
    --id uint       rule id
````

### `siren rule list [flags]`

List rules

```
--group-name string         rule group name
--name string               rule name
--namespace string          rule namespace
--provider-namespace uint   rule provider namespace id
--template string           rule template
````

### `siren rule upload`

Upload Rules YAML file

## `siren server <command>`

Run siren server

### `siren server init [flags]`

Initialize server

```
-o, --output string   Output config file path (default "./config.yaml")
````

### `siren server migrate [flags]`

Run DB Schema Migrations

```
-c, --config string   Config file path (default "./config.yaml")
````

### `siren server start [flags]`

Start server on default port 8080

```
-c, --config string   Config file path (default "config.yaml")
````

## `siren subscription`

Manage subscriptions

### `siren subscription create [flags]`

Create a new subscription

```
-f, --file string   path to the subscription config
````

### `siren subscription delete`

Delete a subscription details

### `siren subscription edit [flags]`

Edit a subscription

```
-f, --file string   Path to the subscription config
    --id uint       subscription id
````

### `siren subscription list`

List subscriptions

### `siren subscription view [flags]`

View a subscription details

```
--format string   Print output with the selected format (default "yaml")
````

## `siren template`

Manage templates

### `siren template delete`

Delete a template details

### `siren template list [flags]`

List templates

```
--tag string   template tag name
````

### `siren template render [flags]`

Render a template details

```
-f, --file string     path to the template config
    --format string   Print output with the selected format (default "yaml")
    --name string     template name
````

### `siren template upload`

Upload Templates YAML file

### `siren template upsert [flags]`

Create or edit a new template

```
-f, --file string   path to the template config
````

### `siren template view [flags]`

View a template details

```
--format string   Print output with the selected format (default "yaml")
````

## `siren worker <command> <worker_command>`

Start or manage Siren's workers

### `siren worker start <command>`

Start a siren worker

#### `siren worker start notification_dlq_handler [flags]`

A notification dlq handler

```
-c, --config string   Config file path (default "config.yaml")
````

#### `siren worker start notification_handler [flags]`

A notification handler

```
-c, --config string   Config file path (default "config.yaml")
````

