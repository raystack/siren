import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# 2.3. Register a Receiver

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

Siren supports several types of receiver to send notification to. For this tour, let's pick the simplest receiver: `file`. With `file` receiver, all published notifications will be written to a file. Let's create a receivers `file` using Siren CLI.

Prepare receiver detail and register the receiver with Siren CLI.
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
$ siren receiver create --file receiver_2.yaml
```

Once done, you will get a message.

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

