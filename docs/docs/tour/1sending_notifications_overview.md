import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# 1 Sending On-demand Notification

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost


This tour shows you how to send a notification to a receiver. You need to pick to which receiver you want send the notification to. If the receiver is not registered yet, you could add a new one. There are several [receivers](../receivers) supported in Siren. In this tour we want to send a notification to a `file` receiver.

## 1.1 Register a Receiver

With `file` receiver, all published notifications will be written to a file. Let's create a `file` receiver using Siren CLI.

Prepare receiver detail.

```bash  title=receiver_1.yaml
name: file-sink-1
type: file
labels:
    key1: value1
    key2: value2
configurations:
    url: ./out-file-sink1.json
```

Register the receiver with this command.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren receiver create --file receiver_1.yaml
```

Once done, you will get a message.

```bash
Receiver created with id: 1 ✓
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers
  --header 'content-type: application/json'
  --data-raw '{
    "name": "file-sink-1",
    "type": "file",
    "labels": {
        "key1": "value1",
        "key2": "value2"
    },
    "configurations": {
        "url": "./out-file-sink1.json"
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

You could verify the registered receiver by getting all receivers or get the new registered receiver by passing the ID. This command is to get all receivers in Siren.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren receiver list
```
  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers
  `}
    </CodeBlock>
  </TabItem>
</Tabs>

Or view a specific receiver with its ID with this command. For example the ID is `1`.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren receiver view 1
```
  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers/1
  `}
    </CodeBlock>
  </TabItem>
</Tabs>

## 1.2 Sending Notification to a Receiver

In the previous [part](#11-register-a-receiver), we have already registered a receiver and got back the receiver ID. We need to prepare the message payload as yaml to be sent by Siren CLI. The message is expected to be in a key value format and placed under `payload.data`.

Prepare a message to send to receiver 1.
```bash title=message_file_1.yaml
payload:
    data:
        text: this is notification to file 1
        a_field: a_value
        another_field: another_value
```
Then we can run `receiver send` command and target the receiver id `1` with flag `--id`.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren receiver send --id 1 --file message_file_1.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers/1/send
  --header 'content-type: application/json'
  --data-raw '{
    "payload": {
        "data": {
            "text": "this is notification to file 1",
            "a_field": "a_value",
            "another_field": "another_value"
        }
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

If succeed, onw new file have been created: `out-file-sink1.json` and the file will have this content.
```json
// out-file-sink1.json
{"a_field":"a_value","another_field":"another_value","routing_method":"receiver","text":"this is notification to file 1"}
```

## What Next?

Well done, you have just completed a tour to send and on-demand notification. The [second tour](./2alerting_rules_subscriptions_overview.md) will be around how to create alerting rules and subscribe a notification if an alert is triggered.

Apart from the tour, we recommend completing the [guides](../guides/overview.md). You could also check out the remainder of the documentation in the [reference](../reference/server_configuration.md) and [concepts](../concepts/overview.md) sections for your specific areas of interest. We've aimed to provide as much documentation as we can for the various components of Siren to give you a full understanding of Siren's surface area. If you are interested to contribute, check out the [contribution](../contribute/contribution.md) page.
