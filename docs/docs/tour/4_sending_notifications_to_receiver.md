import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# 4 - Sending Notification to Receiver

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

In previous [part](./3_registering_receivers.md), we have already registered several receivers and got back the receiver IDs. We could send a notification to the receivers with `/receivers/:receiverId/send` API. We can use Siren CLI to do this.

We need to prepare the message payload as yaml to be sent by Siren CLI. The message is expected to be placed under `payload.data`.

Prepare a message to send to receiver 1.
```bash
cat <<EOT >> message_file_1.yaml
payload:
    data:
        text: this is notification to file 1
        a_field: a_value
        another_field: another_value
EOT
```
Then we can run `receiver send` command and target the receiver id `1` with flag `--id`.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
./siren receiver send --id 1 -f message_file_1.yaml
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

Prepare a message to send to receiver 2.
```bash
cat <<EOT >> message_file_2.yaml
payload:
    data:
        text: this is notification to file 2
        a_field: a_value
        another_field: another_value
EOT
```
Then we can run `receiver send` command and target the receiver id `2` with flag `--id`.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
./siren receiver send --id 2 -f message_file_2.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers/2/send
  --header 'content-type: application/json'
  --data-raw '{
    "payload": {
        "data": {
            "text": "this is notification to file 2",
            "a_field": "a_value",
            "another_field": "another_value"
        }
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

If succeed, two new files have been created: `out-file-sink1.json` and `out-file-sink2.json`. Each file will have this content:
```json
// out-file-sink1.json
{"a_field":"a_value","another_field":"another_value","routing_method":"receiver","text":"this is notification to file 1"}
```
```json
// out-file-sink2.json
{"a_field":"a_value","another_field":"another_value","routing_method":"receiver","text":"this is notification to file 2"}
```