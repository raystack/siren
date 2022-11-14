# PagerDuty
|||
|---|---|
|**type**|`pagerduty`|

Siren's PagerDuty receiver tied to a PagerDuty Service. Siren requires a `v1` integration key/service key of a PagerDuty service to communicate and the `Events API v1` of the PagerDuty Service needs to be enabled. [Here](https://support.pagerduty.com/docs/services-and-integrations) is more information on how to create a new service.

## Configurations in API

```json
"configurations": {
    "service_key": <string>
}
```
## Configurations Stored in DB

Same like [Configurations in API](#configurations-in-api)

## Subscription

PagerDuty receiver does not have `SubscriptionConfig`.

## Message Payload

### Contract

Pagerduty has `v1` and `v2` events API. What Siren's support currently is sending event to PagerDuty events `v1` API with this [contract](https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTc3-events-api-v1).

```yaml
# v1
service_key: <string>
event_type: <string>
incident_key: <string>
description: <string>
client: <string>
client_url: <string>
details:
  - <key1>: <any>
    <key2>: <any>
  - <key3>: <any>
    <key4>: <any>
    .
    .
contexts:
  - type:  <string>
    src:  <string>
    href:  <string>
    text:  <string>
    alt:  <string>
  - type:  <string>
    src:  <string>
    href:  <string>
    text:  <string>
    alt:  <string>
    .
    .
```
### Default Alert Template

Siren has a PagerDuty default notification [template](../../../plugins/receivers/pagerduty/config/default_alert_template_body_v1.goyaml) used by all alert notifications.
