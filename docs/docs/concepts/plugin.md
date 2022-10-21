# Plugin

Siren decouples various `provider`, `receiver`, and `queue` as a plugin. The purpose is to ease the extension of new plugin. We welcome all contributions to add new plugin.

## Provider

Provider responsibility is to accept incoming rules configuration from Siren and send alerts to the designated Siren Hook API. Provider plugin needs to fulfill some interfaces. More detail about interfaces can be found in [contribution](../contribute/provider.md) page. Supported providers are:
- [Cortexmetrics](https://cortexmetrics.io/)

## Receiver

Receiver defines where the notification Siren sends to. Receiver plugin needs to fulfill some interfaces. More detail about interfaces can be found in [contribution](../contribute/receiver.md) page. Supported providers are:
- [Slack](https://api.slack.com/methods/chat.postMessage)
- [PagerDuty Events API v1](https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTc3-events-api-v1)
- HTTP
- File


Receiver plugin is being used by two different services: receiver service and notification service. Receiver service handles the way the receiver is being stored, updated, fetched, and removed. Notification service uses receiver plugin to send notification. Each service has its own interface that needs to be implemented.

### Configurations

Siren receiver plugins have several configs: `ReceiverConfig`, `SubscriptionConfig` (if needed), and `NotificationConfig`. 

- **ReceiverConfig** is a config that will be part of `receiver.Receiver` struct and will be stored inside the DB's receivers table.
- **SubscriptionConfig** is optional. Subscription config is defined and used if the receivers inside subscription requires another additional configs rather than `ReceiverConfig`. For example, Slack stores encrypted `token` when storing receiver information inside the DB but has another config `channel_name` on subscription level.
- **NotificationConfig** embeds `ReceiverConfig` and `SubscriptionConfig` (if needed).
- **AppConfig** is a config of receiver plugins that is being loaded when the Siren app is started. `AppConfig` can be set up via environment variable or config file. Usually this is a generic config of a specific receiver regardless where the notification is being sent to (e.g. http config, receiver host, etc...). If your plugin requires `AppConfig`, you can set the config inside `plugins/receivers/config.go`.

In Siren receiver plugins, all configs will be transform back and forth from `map[string]interface{}` to struct using [mitchellh/mapstructure](https://github.com/mitchellh/mapstructure). You might also need to add more functions to validate and transform configs to `map[string]interface{}`.


### Interface

#### ConfigResolver

ConfigResolver is being used by receiver service to manage receivers. It is an interface for the receiver to resolve all configs and functions. 

```go
type ConfigResolver interface {
    // TODO might be removed
	BuildData(ctx context.Context, configs map[string]interface{}) (map[string]interface{}, error)
    // TODO might be removed
	BuildNotificationConfig(subscriptionConfigMap map[string]interface{}, receiverConfigMap map[string]interface{}) (map[string]interface{}, error)
	PreHookTransformConfigs(ctx context.Context, configs map[string]interface{}) (map[string]interface{}, error)
	PostHookTransformConfigs(ctx context.Context, configs map[string]interface{}) (map[string]interface{}, error)
}
```

- **BuildData** is being used in GetReceiver where `data` field in Receiver is being populated. This might not relevant anymore for our current use case and might be deprecated later.
- **BuildNotificationConfig** is being used for subscription. This might not relevant anymore for our current use case and might be deprecated later.
- **PreHookTransformConfigs** is being used to transform configs (e.g. encryption) before the config is being stored in the DB.
- **PostHookTransformConfigs** is being used to transform configs (e.g. decryption) after the config is being fetched from the DB.

#### Notifier

Notifier interface is being used by notification service and consists of all functionalities to publish notifications.

```go
type Notifier interface {
	PreHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error)
	PostHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error)
	DefaultTemplateOfProvider(templateName string) string
	Publish(ctx context.Context, message Message) (bool, error)
}
```

- **PreHookTransformConfigs** is being used to transform configs (e.g. encryption) before the config is being enqueued.
- **PostHookTransformConfigs** is being used to transform configs (e.g. decryption) after the config is being dequeued.
- **DefaultTemplateOfProvider** assigns default provider template for alert notifications of as specific provider. Each provider might send alerts with different format, the template needs to build notification specific message out of the alerts for each provider. Each provider has to have a reserved template name (e.g. `template.ReservedName_xxx`) and all alerts coming from the provider needs to use the template with the reserved name.
- **Publish** handles how message is being sent. The first return argument is `retryable` boolean to indicate whether an error is a `retryable` error or not. If it is API call, usually response status code 429 or 5xx is retriable. You can use `pkg/retrier` to retry the call.


### Base Plugin

Siren provide base plugin in `plugins/receivers/base` which can be embedded in all plugins service struct. By doing so, you just need to implement all interfaces' method that you only need. The unimplemented methods one will already be handled by the `base` plugin.


## Queue
Queue is used as a buffer for the outbound notifications. Siren has a pluggable queue where user could choose which Queue to use in the [config](../reference/server_configuration.md). Supported Queues are:
- In-Memory
- PostgreSQl