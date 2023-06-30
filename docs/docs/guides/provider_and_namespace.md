import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# Provider and Namespace

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

Siren provider represents a monitoring server. We define alerts and alerts routing configuration inside providers. For example the below provider, describes how a cortex monitoring server info is stored inside Siren.

```json
# provider
{
  "id": "3",
  "host": "http://localhost:9009",
  "urn": "localhost_cortex",
  "name": "localhost_cortex",
  "type": "cortex",
  "credentials": {},
  "created_at": "2021-10-26T06:00:52.627917Z",
  "updated_at": "2021-10-26T06:00:52.627917Z"
}
```

Monitoring providers usually have tenants, a sharded way of storing and querying telemetry data. Siren calls them **namespaces**. You should create a namespaces with the same name as the tenant of your monitoring providers. Example: the below namespace represent the tenant "raystack".

```json
# namespace
{
  "id": "10",
  "urn": "raystack",
  "name": "raystack",
  "provider": {
    "id": 3
  },
  "credentials": {},
  "created_at": "2021-10-26T06:00:52.627917Z",
  "updated_at": "2021-10-26T06:00:52.627917Z"
}
```

## Provider

### Provider creation

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren provider create --file provider.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/providers
  --header 'content-type: application/json'
  --data-raw '{
  "host": "http://localhost:9009",
  "urn": "localhost_cortex",
  "name": "localhost_cortex",
  "type": "cortex",
  "credentials": {},
  "labels": {}
}'`}
    </CodeBlock>

  </TabItem>
</Tabs>

**Terminology of the request body**

| Term        | Description                                                | Example                    |
| ----------- | ---------------------------------------------------------- | -------------------------- |
| host        | Fully qualified path for the provider                      | http://localhost:9009      |
| urn         | Unique name for this provider (uneditable)                 | localhost_cortex           |
| name        | Name of the proider (editable)                             | localhost_cortex           |
| type        | type of the provider(cortex/influx etc.)                   | cortex                     |
| credentials | key value pair to be used for authentication with the host | {"bearer_token":"x2y4rd5"} |
| labels      | key value pair that can be used as label selector          | {"environment":"dev"}      |

The response body will look like this:

```json
{
  "id": "1",
  "host": "http://localhost:9009",
  "urn": "localhost_cortex",
  "name": "localhost_cortex",
  "type": "cortex",
  "credentials": {},
  "created_at": "2022-01-03T05:10:47.880209Z",
  "updated_at": "2022-01-03T05:10:47.880209Z"
}
```

### Provider update

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren provider edit --id 4 --file provider.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request PUT
  --url `}{defaultHost}{`/`}{apiVersion}{`/providers
  --header 'content-type: application/json'
  --data-raw '{
  "host": "http://localhost:9009",
  "urn": "localhost_cortex",
  "name": "localhost_cortex_1",
  "type": "cortex",
  "credentials": {},
  "labels": {}
}'`}
    </CodeBlock>

  </TabItem>
</Tabs>

### Getting a provider

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren provider view 1
```

  </TabItem>

  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET --url `}{defaultHost}{`/`}{apiVersion}{`/providers/1`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Getting all providers

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren provider list
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET --url `}{defaultHost}{`/`}{apiVersion}{`/providers`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Deleting a provider

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren provider delete 4
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request DELETE --url `}{defaultHost}{`/`}{apiVersion}{`/providers/4`}
    </CodeBlock>
  </TabItem>
</Tabs>

**Note:**

1. Before deleting the provider, you will need to delete dependant resources (namespaces).

## Namespace

**Note:** These operations on namespaces inside Siren doesn't affect the actual tenant in the monitoring provider. The user should manage the tenants themselves.

### Creating a namespace

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren namespace create --file namespace.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces
  --header 'content-type: application/json'
  --data-raw '{
    "name": "test",
    "urn": "test",
    "provider": "5",
    "credentials": {},
    "labels": {}
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

**Terminology of the request body**

| Term        | Description                                                | Example                    |
| ----------- | ---------------------------------------------------------- | -------------------------- |
| urn         | Unique name for this namespace (uneditable)                | test-tenant                |
| name        | Name of the tenant (editable)                              | test-tenant                |
| provider    | id of the provider to which this tenant belongs            | 1                          |
| labels      | key value pair that can be used as label selector          | {"environment":"dev"}      |
| credentials | key value pair to be used for authentication with the host | {"bearer_token":"x2y4rd5"} |

The response body will look like this:

```json
{
  "id": "13",
  "urn": "test",
  "name": "test",
  "provider": "5",
  "credentials": {},
  "created_at": "2022-01-03T07:06:30.884113Z",
  "updated_at": "2022-01-03T07:06:30.884113Z"
}
```

### Updating a namespace

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren namespace edit --id 13 --file namespace.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request PUT
  --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces/13
  --header 'content-type: application/json'
  --data-raw '{
    "name": "test2",
    "urn": "test",
    "provider": "5",
    "credentials": {},
    "labels": {}
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Getting a namespace

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren namespace view 13
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces/13`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Getting all namespaces

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren namespace list
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces`}
    </CodeBlock>
  </TabItem>
</Tabs>

### Deleting a namespace

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren namespace delete 13
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request DELETE --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces/13`}
    </CodeBlock>
  </TabItem>
</Tabs>
