# File
|||
|---|---|
|**type**|`file`|

File receiver write down outgoing notifications to a file located within a `url` in the configurations.


## Configurations in API

The `url` below should be file path url e.g. `./folder_a/folder_b/file.json`.

```json
"configurations": {
  "url": <string>
}
```

## Configurations Stored in DB

Same like [Configurations in API](#configurations-in-api)

## Subscription

File receiver does not have `SubscriptionConfig`.

## Message Payload

### Contract
No specific message payload contract for File receiver. Payload will be written to file as-is in `ndjson` format.
