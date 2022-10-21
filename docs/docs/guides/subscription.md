import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# Subscription

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

Siren lets you subscribe to a notification when they are triggered. You can define custom matching conditions and use
[receivers](./receiver.md) to describe which medium you want to use for getting the notifications when a notification is triggered. A notification could be triggered on-demand via API or by the incoming alerts via webhook.

**Example Subscription:**

```json
{
  "id": "385",
  "urn": "siren-dev-prod-critical",
  "namespace": "10",
  "receivers": [
    {
      "id": "2"
    },
    {
      "id": "1",
      "configuration": {
        "channel_name": "siren-dev-critical"
      }
    }
  ],
  "match": {
    "environment": "production",
    "severity": "CRITICAL"
  },
  "created_at": "2021-12-10T10:38:22.364353Z",
  "updated_at": "2021-12-10T10:38:22.364353Z"
}
```

The above means whenever any alert which has labels matching the labels
`"environment": "production", "severity": "CRITICAL"`, send this alert to two medium defined by receivers with id: `2` and `1`. Assuming the receivers id `2` to be of Pagerduty type, a PD call will be invoked and assuming the receiver with id `1` to be slack type, a message will be sent to the channel #siren-dev-critical.

## API Interface

### Create a subscription

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren subscription create --file subscription.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/subscriptions
  --header 'content-type: application/json'
  --data-raw '{
    "urn": "siren-dev-prod-critical",
    "receivers": [
        {
            "id": "1",
            "configuration": {
                "channel_name": "siren-dev-critical"
            }
        },
        {
            "id": "2"
        }
    ],
    "match": {
        "severity": "CRITICAL",
        "environment": "production"
    },
    "namespace": "10"
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Update a subscription

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren subscription edit --id 10 --file subscription.yaml
```


  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request PUT
  --url `}{defaultHost}{`/`}{apiVersion}{`/subscriptions/10
  --header 'content-type: application/json'
  --data-raw '{
    "urn": "siren-dev-prod-critical",
    "receivers": [
        {
            "id": "1",
            "configuration": {
                "channel_name": "siren-dev-critical"
            }
        },
        {
            "id": "2"
        }
    ],
    "match": {
        "severity": "CRITICAL",
        "environment": "production"
    },
    "namespace": "10"
}'`}
    </CodeBlock>

  </TabItem>
</Tabs>

### Get all subscriptions
<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>


```bash
$ siren subscription list
```


  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/subscriptions`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Get a subscriptions

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren subscription view 10
```


  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/subscriptions/10`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Delete subscriptions

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren subscription delete 10
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request DELETE
  --url `}{defaultHost}{`/`}{apiVersion}{`/subscriptions/10`}
    </CodeBlock>
  </TabItem>
</Tabs>
