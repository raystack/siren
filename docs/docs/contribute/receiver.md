# Add a New Receiver Plugin

More details about the concept of receiver plugin can be found [here](../concepts/plugin.md#receiver). In this part, we will show hot to add a new receiver plugin to write notifications to a `file`.

## Defining Configs

To write a file, we need a `url` of the file. This might be the only config that we needs. We also don't need to define `SubscriptionConfig` since we don't need a specific config for the subscription.
```go
type ReceiverConfig struct {
	URL string `mapstructure:"url"`
}
```
We define a `NotificationConfig` which only embeds `ReceiverConfig`. This is helpful to separate the concern for a specific use-cases in some plugins.
```go
type NotificationConfig struct {
	ReceiverConfig `mapstructure:",squash"`
}
```

For `file` type, we don't need an `AppConfig` as for now. So we don't need to add one in `plugins/receivers/config.go`.

Now that we already have defined all configs needed, we needs to implement all interfaces needed by defining a new `ReceiverService` and `NotificationService`.

## Implement ConfigResolver

We need to create a new `ReceiverService` and implement `ConfigResolver`. For `file` receiver, we don't need to do transformation of configs before and after writing and reading from the DB. Therefore, we only needs to implement two `ConfigResolver` methods: `PreHookTransformConfigs` to validate the config before storing it to the DB and `BuildNotificationConfig` to merge `ReceiverConfig` into a `NotificationConfig`.

```go
type ReceiverService struct {
	base.UnimplementedReceiverService
}

func NewReceiverService() *ReceiverService {
	return &ReceiverService{}
}

func (s *ReceiverService) PreHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(configurations, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	if err := receiverConfig.Validate(); err != nil {
		return nil, errors.ErrInvalid.WithMsgf(err.Error())
	}

	return configurations, nil
}

func (s *ReceiverService) BuildNotificationConfig(subsConfs map[string]interface{}, receiverConfs map[string]interface{}) (map[string]interface{}, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(receiverConfs, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	notificationConfig := NotificationConfig{
		ReceiverConfig: *receiverConfig,
	}

	return notificationConfig.AsMap(), nil
}
```

## Implement Notifier

We need to create a new `NotificationService` and implement `Notifier`. For `file` receiver, we don't need to do transformation of configs before and after enqueue and dequeue. Therefore, we only needs to implement two `Notifier` methods: `PreHookTransformConfigs` to validate the config before enqueuing notification message and `Publish` to send notifications (to write notifications to a file under `url`).

```go
type NotificationService struct {
	base.UnimplementedNotificationService
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) Publish(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationMessage.Configs, notificationConfig); err != nil {
		return false, err
	}

	bodyBytes, err := json.Marshal(notificationMessage.Details)
	if err != nil {
		return false, err
	}
	if err := s.validateFilePath(notificationConfig.URL); err != nil {
		return false, err
	}

	fileInstance, err := os.OpenFile(notificationConfig.URL, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return false, err
	}

	byteNewLine := []byte("\n")
	bodyBytes = append(bodyBytes, byteNewLine...)
	_, err = fileInstance.Write(bodyBytes)
	if err != nil {
		return false, err
	}

	return false, nil
}


func (s *NotificationService) PreHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return nil, err
	}

	if err := notificationConfig.Validate(); err != nil {
		return nil, err
	}

	return notificationConfig.AsMap(), nil
}
```

## Integrate New Plugin with Siren

1. Define and add your new type of plugin called `file` inside `core/receivers/type.go`.
2. Initialize your plugin receiver service and notification service and add to the `ConfigResolver` and `Notifier` registries map.
3. To make sure notification handler and dlq handler process your new type, don't forget to add your new receiver type in notification message & dlq handler config or make it default to support all receiver types.
