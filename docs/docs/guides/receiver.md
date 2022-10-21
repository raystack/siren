import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# Receiver

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

You can use receivers to send notifications on demand as well as on certain matching conditions (API for this is in the roadmap). Subscriptions use receivers to define routing configuration in Siren. With Siren subscriptions, incoming alerts via webhook will be routed to the pre-registered receivers by matching the subscriptions label. More info about notification concept is [here](../concepts/notification.md). The how-to sending notification can be found [here](../guides/notification.md).

### Create a Receiver

Each receiver type might require different kind of configurations. A `configurations` field is a dynamic field that has to be filled depend on the receiver type. Below is the example to add a [PagerDuty](#pagerduty) receiver. A `labels` field is a KV-string value to label each receiver.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren receiver create --file receiver.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers
  --header 'content-type: application/json'
  --data-raw '{
    "name": "doc-pagerduty-receiver",
    "type": "pagerduty",
    "labels": {
        "team": "siren-devs"
    },
    "configurations": {
        "service_key": "eq23r23rfewf3qwf3wf3w2f23wf32qwfw3fw3"
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Update a Receiver

**Note:** While updating a receiver, you will have to make sure all subscriptions that are using this receivers get refreshed(updated), since subscriptions use receivers to create routing configuration dynamically. Receiver type is immutable.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren receiver edit --id 61 --file receiver.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request PUT
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers/61
  --header 'content-type: application/json'
  --data-raw '{
    "name": "doc-http-receiver",
    "type": "http",
    "labels": {
        "team": "siren-devs"
    },
    "configurations": {
        "url": "http://localhost:4001"
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Getting a receiver

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren receiver view 61
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET --url `}{defaultHost}{`/`}{apiVersion}{`/receivers/61`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Getting all receivers

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren receiver list
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET --url `}{defaultHost}{`/`}{apiVersion}{`/receivers`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Deleting a receiver

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren receiver delete 61
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request DELETE --url `}{defaultHost}{`/`}{apiVersion}{`/receivers/61`}
    </CodeBlock>
  </TabItem>
</Tabs>

## CLI Interface

```text
Receivers are the medium to send notification for which we intend to mange configuration.


Usage
  siren receiver [flags]

Core commands
  create         Create a new receiver
  delete         Delete a receiver details
  edit           Edit a receiver
  list           List receivers
  send           Send a receiver notification
  view           View a receiver details

Flags
  -h, --host string   Siren API service to connect to
Use "siren receiver [command] --help" for more information about a command.
```

