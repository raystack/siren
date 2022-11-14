# Add a New Provider Plugin

Provider plugin is being used to update rule configurations to the provider, setup some runtime provider configuration config if any, and transform incoming alert data in Siren's Hook API to a list of `*alert.Alert` model. More details about the concept of provider plugin can be found [here](../concepts/plugin.md#provider-plugin). 

## Steps to Add a New Plugin

1. **Add a new package**

    A new plugin could be added as a new package under `plugins/providers/{new_plugin}`. Package name should comply with odpf golang standard naming. See the [handbook](https://odpf.github.io/handbook/programming/go#packages). Ideally you might want to name the package with lower-case and without any `-` or `_` signs.

2. **Defining Configurations** (Optional)

    If there is a need to have a generic config for the provider that is being loaded during start-up, you could add a new `AppConfig` and assign the config to `plugins/providers/config.go` to expose it to the server-level config. Siren will recognize and read the config when starting up.

3. **Implement Interfaces**

    - **Implement `TransformToAlerts`**

        You need to find a way for the provider to send alert to Siren's hook API `{api_version}/alerts/{provider_type}/{provider_id}`. You need to implement a `TransformToAlerts` method to parse and transform the body of incoming request received by Siren's hook API to a list of Siren's `*alert.Alert` model.

    - **Implement `UpsertRules`**
        
        You need to implement `UpsertRules` method to push all rules changes to the provider everytime a rule is upserted.

    - **Implement `SyncRuntimeConfig`** (Optional)

        If a provider supports setting the configuration dynamically and if for each provider's namespace there is a possibility to have different config (e.g. for multitenancy purpose like CortexMetrics), you could implement `SyncRuntimeConfig` function to push all configurations changes to the provider everytime a new namespace is created or an existing namespace is updated.

4. **Integrate the new plugin with Siren main flow**
    - Define and add your new type of plugin inside `core/providers/type.go`.
    - Initialize your plugin provider service and add to the `AlertTransformer`, `ConfigSyncer`, and `RuleUploader` registries map in `cli/deps`.
