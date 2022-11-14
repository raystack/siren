import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# Notification

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

To understand more concepts of notification in Siren, you can visit this [page](../concepts/notification.md).

## Sending a message/notification

We could send a notification to a specific receiver by passing a `receiver_id` in the path params and correct payload format in the body. The payload format needs to follow receiver type contract. 

### Example: Sending Notification to Slack

Assuming there is a slack receiver registered in Siren with ID `51`. Sending to that receiver would require us to have a `payload.data` that have the same format as the expected [slack](../receivers/slack.md#message-payload) payload format.

```yaml title=payload.yaml
payload:
  data:
    channel: siren-devs
    text: an alert or notification
    icon_emoji: ":smile:"
    attachments:
      - blocks: 
        - type: section
          text:
            type: mrkdwn
            text: |-
              New Paid Time Off request from <example.com|Fred Enriquez>

              <https://example.com|View request>
```

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren receiver send --id 51 --file payload.yaml
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

For all incoming alerts via Siren hook API, notifications are also generated and published via subscriptions. Siren will match labels from the alerts with label matchers in subscriptions. The assigned receivers for all matched subscriptions will get the notifications. More details are explained [here](./alert_history.md). 

Siren has a default template for alerts notification for each receiver. Go to the [Receivers](../receivers/slack.md#default-alert-template) section to explore the default template defined by Siren.