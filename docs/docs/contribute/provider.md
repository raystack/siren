# Add a New Provider Plugin

Provider plugin is being used to update rule configurations to the provider, setup some runtime provider configuration config if any, and transform incoming alert data in Siren's Hook API to list of *alert.Alert model.

## Requirements

1. You need to find a way for the provider to send alert to Siren's hook API `{version}/alerts/{provider_type}/{provider_id}`. You need to implement a `TransformToAlerts` method to parse and transform the body of incoming request received by Siren's hook API to the list of Siren's *alert.Alert model.

2. You need to implement `UpsertRules` method to push all rules changes to the provider everytime a rule is upserted.

3. If provider supports setting the configuration dynamically and if for each provider's namespace there is a possibility to have different config (e.g. for multitenancy purpose like CortexMetrics), you could implement `SyncRuntimeConfig` function to push all configurations changes to the provider everytime a new namespace is created or an existing namespace is updated.


## Defining Configs

If there is a need to have a generic config for the provider that is being loaded during start-up, you could add a new `AppConfig` and assign the config to `plugins/providers/config.go` to expose it to the app-level config. Siren will recognize and read the config when starting up.

## Integrate New Plugin with Siren

1. Define and add your new type of plugin inside `core/providers/type.go`.
2. Initialize your plugin receiver service and notification service and add to the `ConfigResolver` and `Notifier` registries map in `cli/deps`.
3. To make sure notification handler and dlq handler process your new type, don't forget to add your new receiver type in notification message & dlq handler config or make it default to support all receiver types.
