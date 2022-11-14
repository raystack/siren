# HTTP

|||
|---|---|
|**type**|`http`|

HTTP receiver submit notification with `HTTP POST` to a url.

## Configurations in API

```json
"configurations": {
  "url": <string>
}
```
## Configurations Stored in DB

Same like [Configurations in API](#configurations-in-api)

## Subscription

HTTP receiver does not have `SubscriptionConfig`.

## Message Payload

### Contract

No specific message payload contract for HTTP receiver. Payload will be sent as-is. If defined by [templates](../guides/template.md), payload will be the same with the generated payload by template.
