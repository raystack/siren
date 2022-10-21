import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# Notification

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

Notification is one of main features in Siren. Siren capables to send notification to various receivers (e.g. Slack, PagerDuty). Notification in Siren could be sent directly to a receiver or user could subscribe notifications by providing key-value label matchers. For the latter, Siren routes notification to specific receivers by matching notification key-value labels with the provided label matchers.

## Sending a message/notification

We could send a notification to a specific receiver by passing a `receiver_id` in the path params and correct payload format in the body. The payload format needs to follow receiver type contract. 

### Example: Sending Notification to Slack

If receiver is slack, the `payload.data` should be within the expected [slack](#slack) payload format.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren receiver create --file receiver.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers/51/send
  --header 'content-type: application/json'
  --data-raw '{
    "payload": {
        "data": {
            "channel": "siren-devs",
            "text": "an alert or notification",
            "icon_emoji": ":smile:"
            "attachments": [
                "blocks": [
                    {
                        "type": "section",
                        "text": {
                            "type": "mrkdwn",
                            "text": "New Paid Time Off request from <example.com|Fred Enriquez>\n\n<https://example.com|View request>"
                        }
                    }
                ]
            ]
        }
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

Above end the message to channel name `#siren-devs` with `payload.data` in [slack](#slack) payload format.


## Alerts Notification

For all incoming alerts via Siren hook API, notifications are also generated and published via subscriptions. Siren will match labels from the alerts with label matchers in subscriptions. The assigned receivers for all matched subscriptions will get the notifications. More details are explained [here](./alert_history.md). Sending notification message requires notification message payload to be in the same format as what receiver expected. The format can be found in the detail in [reference](../reference/receiver.md).