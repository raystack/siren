# CortexMetrics

|||
|---|---|
|**type**|`cortex`|

[CortexMetrics](https://cortexmetrics.io/) is a Horizontally scalable, highly available, multi-tenant, long term storage for Prometheus. It could run in multi-tenant scenario where it isolates data and queries from multiple different independent Prometheus sources in a single cluster, allowing untrusted parties to share the same cluster.

Similar with how prometheus works, CortexMetrics consumes metrics sent from other services, store, and evaluate the metrics with the configured rules. Alerts will be triggered and processed by CortexMetrics' alert manager and notifications will be sent to the designated sinks (webhook, slack, pagerduty, etc..). 

Since Siren handles all subscriptions and notifications routing, Siren configures CortexMetrics to send all alerts only to Siren webhook API.

## Multi-tenancy

Tenants in CortexMetrics are mapped to [Namespaces](../guides/provider_and_namespace.md#namespace) in Siren. To integrate multiple tenants, you need to create multiple namespaces for each tenant. Each tenant will have different configuration.


## Server Configuration
There is a generic CortexMetrics configuration in Siren server configuration that could be used to tune the CortexMetrics. The configuration will always be synchronized everytime a namespace in Siren is created or updated.

Here is a config that is part of the server configuration. Please note that the config will always be applied to all CortexMetrics registered in Siren and only synchronized when a namespace in Siren is created or updated. Siren server restart is also required to get the latest value update of these configs.

```yaml
...
providers:
  cortex:
    group_wait: 30s
    webhook_base_api: http://localhost:8080/v1beta1/alerts/cortex
...
```
- The `group_wait` config usage is similar with the one in CortexMetrics alert manager [configuration](https://prometheus.io/docs/alerting/latest/configuration/#example).
- The `webhook_base_api` defined the base API that will be appended with `provider_id` for each specific provider. If a namespace of provider with id `3` is updated, Siren will configure the webhook receiver in CortexMetrics with this URL: `http://localhost:8080/v1beta1/alerts/cortex/3`.