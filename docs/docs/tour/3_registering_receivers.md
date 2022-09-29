import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# 3 - Registering Receivers

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost


## 1. Register a receiver

Siren supports several types of receiver to send notification to. For this tour, let's pick the simplest receiver: `file`. With `file` receiver, all published notifications will be written to a file. Let's create two receivers `file` with different filename using Siren CLI.

Prepare receivers detail:
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
Register the receiver with this command

Prepare a receiver detail:
```bash
cat <<EOT >> receiver_2.yaml
name: file-sink-2
type: file
labels:
    key1: value1
    key2: value2
configurations:
    url: ./out-file-sink2.json
EOT
```

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
./siren receiver create -f receiver_1.yaml
```
Once done, you will get messages
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

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
./siren receiver create -f receiver_2.yaml
```
Once done, you will get messages
```bash
Receiver created with id: 2 ✓
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers
  --header 'content-type: application/json'
  --data-raw '{
    "name": "file-sink-2",
    "type": "file",
    "labels": {
        "key1": "value1",
        "key2": "value2"
    },
    "configurations": {
        "url": "./out-file-sink2.json"
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>
