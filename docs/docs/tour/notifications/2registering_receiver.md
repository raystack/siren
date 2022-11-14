import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# 1.2 Register a Receiver

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost


As mentioned in the [overview](./1sending_notifications_overview.md) before, we will send a notification to a `file` receiver. With `file` receiver, all published notifications will be written to a file. Let's create a `file` receiver using Siren CLI.

Prepare receiver detail.

```bash
cat <<EOT >> receiver_1.yaml
name: file-sink-1
type: file
labels:
    key1: value1
    key2: value2
configurations:
    url: ./out-file-sink1.json
EOT
```

Register the receiver with this command.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
$ siren receiver create -f receiver_1.yaml
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
