# Plugin

Siren decouples `provider` and `receiver` as a plugin. The purpose is to ease the extension of new plugin. We welcome all contributions to add new plugin.

## Provider

Provider responsibility is to accept incoming rules configuration from Siren and send alerts to the designated Siren's Hook API. Provider plugin needs to fulfill some interfaces. More detail about interfaces can be found in [contribution](../contribute/provider.md) page. Supported providers are:

- [Cortexmetrics](https://cortexmetrics.io/)

### Configurations

Siren provider plugin could have an `AppConfig`: 
- A config of provider plugin that is being loaded when the Siren app is started. `AppConfig` can be set up via environment variable or config file. Usually this is a generic config of a specific receiver regardless which namespace the provider is attached to. If your plugin requires `AppConfig`, you can set the config inside `plugins/providers/config.go`.
### Interface

```go
type ProviderPlugin interface {
	// AlertTransformer
	TransformToAlerts(ctx context.Context, providerID uint64, body map[string]interface{}) ([]*alert.Alert, int, error) 

	// ConfigSyncer
	SyncRuntimeConfig(ctx context.Context, namespaceURN string, prov provider.Provider) error
	
	// RuleUploader
	UpsertRule(ctx context.Context, namespaceURN string, prov provider.Provider, rl *rule.Rule, templateToUpdate *template.Template) error
}
```
**AlertTransformer** interface is being used by alert service to transform incoming alert body in Siren's Hook API to list of *alert.Alert

**ConfigSyncer** interface is being used by namespace service to synchronize runtime-provider configs for a specific namespace. In cortex, it is being used to sync alertmanager config.

**RuleUploader** interface is being used to upsert rules to the provider. It support templatization of rules.

## Receiver

Receiver defines where the notification Siren sends to. Receiver plugin needs to fulfill some interfaces. More detail about interfaces can be found in [contribution](../contribute/receiver.md) page. Supported providers are:

- [Slack](https://api.slack.com/methods/chat.postMessage)
- [PagerDuty Events API v1](https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTc3-events-api-v1)
- HTTP
- File

Receiver plugin is being used by two different services: receiver service and notification service. Receiver service handles the way the receiver is being stored, updated, fetched, and removed. Notification service uses receiver plugin to send notification. Each service has its own interface that needs to be implemented.

### Configurations

Siren receiver plugin could have several configs: 
1. `ReceiverConfig`
	- A config that will be part of `receiver.Receiver` struct and will be stored inside the DB's receivers table.
2. `SubscriptionConfig` (optional)
	- Only if needed
	- Defined and used if the receivers inside subscription requires another additional configs rather than `ReceiverConfig`. For example, Slack stores encrypted `token` when storing receiver information inside the DB but has another config `channel_name` on subscription level.
3. `NotificationConfig`
	- Embeds `ReceiverConfig` and `SubscriptionConfig` (if exists).
4. `AppConfig`
	- A config of receiver plugin that is being loaded when the Siren app is started. `AppConfig` can be set up via environment variable or config file. Usually this is a generic config of a specific receiver regardless where the notification is being sent to (e.g. http config, receiver host, etc...). If your plugin requires `AppConfig`, you can set the config inside `plugins/receivers/config.go`.

> In Siren receiver plugins, all configs will be transform back and forth from `map[string]interface{}` to struct using [mitchellh/mapstructure](https://github.com/mitchellh/mapstructure). You might also need to add more functions to validate and transform configs to `map[string]interface{}`.

### Interface

```go
type ReceiverPlugin interface {
	// Config Resolver
	PreHookDBTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error)
	PostHookDBTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error)
	BuildData(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error)
	
	// Notifier
	PreHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error)
	PostHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error)
	GetSystemDefaultTemplate() string
	Send(ctx context.Context, notificationMessage notification.Message) (bool, error)
}
```

ConfigResolver interface is being used by receiver service to manage receivers. It is an interface for the receiver to resolve all configs and functions.

- **BuildData** is optional. It is being used in Get and List Receiver where `data` field in Receiver is being populated to send back custom information to the users. Slack receiver uses this to show what channels does the slack receiver support in a workspace.
- **PreHookDBTransformConfigs** is being used to transform configs (e.g. encryption) before the config is being stored in the DB.
- **PostHookDBTransformConfigs** is being used to transform configs (e.g. decryption) after the config is being fetched from the DB.

Notifier interface is being used by notification service and consists of all functionalities to publish notifications.

- **PreHookQueueTransformConfigs** is being used to transform configs (e.g. encryption) before the config is being enqueued.
- **PostHookQueueTransformConfigs** is being used to transform configs (e.g. decryption) after the config is being dequeued.
- **GetSystemDefaultTemplate** assigns default template for alert notifications. It is expected for hook API to transform the data into the [system-default template variables](#notification-system-default-template).
- **Send** handles how message is being sent. The first return argument is `retryable` boolean to indicate whether an error is a `retryable` error or not. If it is API call, usually response status code 429 or 5xx is retriable. You can use `pkg/retrier` to retry the call.

## Base Plugin

Siren provide base plugin in `plugins/receivers/base` and `plugins/providers/base` which needs to be embedded in all plugins service struct. By doing so, you just need to implement all interface's methods that you only need. The unimplemented methods one will already be handled by the `base` plugin.

## Notification System Default Template

Each receiver might want to have a specific template to transform alerts into a receiver's contract. In that case one needs to define a template in a `.goyaml` file inside `plugins/receivers/{type}/config/*.goyaml`. It is expected for each Hook API to [populate the Data and Labels of notification](#hook-api) with the same variable keys. Below are the supported variables that could be used inside the template.
```
Data
- id
- status ("firing"/"resolved")
- resource
- template
- metricValue
- metricName
- generatorUrl
- numAlertsFiring
- dashboard
- playbook
- summary

Labels
- severity ("WARNING"/"CRITICAL")
- alertname
- others specific rules
```