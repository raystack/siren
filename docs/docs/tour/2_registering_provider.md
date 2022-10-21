import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# 2 - Registering Provider and Namespaces


export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

## 1. Register the provider

The first things we need to set up before we add receivers and testing alerts and notifications are we need to register our [Cortexmetrics](https://cortexmetrics.io/) as provider and its namespaces.

Siren provides HTTP API where we need to send a request to `POST /v1beta1/providers` with a json body to create a provider. Beside that, Siren also has a CLI that interacts to Siren server and we could use it.

To create a new provider with CLI, we need to create a `yaml` file.
```yaml
# input.yaml
host: http://localhost:9009
urn: localhost-dev-cortex
name: dev-cortex
type: cortex
```
If you are in unix system, you could do this
```bash
cat <<EOT >> input.yaml
host: http://localhost:9009
urn: localhost-dev-cortex
name: dev-cortex
type: cortex
EOT
```

Once the file is ready, we can start creating the provider.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
siren provider create --file input.yaml
```
If succeed, you will got this message.
```shell
Provider created with id: 1 ✓
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/providers
  --header 'content-type: application/json'
  --data-raw '{
    "host": "http://localhost:9009",
    "urn": "localhost-dev-cortex",
    "name": "dev-cortex",
    "type": "cortex"
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>



The `id` we got from the provider creation is important to create a namespace later.

## 2. Register namespaces

For multi-tenancy scenario, which Cortex supports, we need to define namespaces in Siren. Assuming there are 2 tenants in Cortex, `odpf` and `non-odpf`, we need to create 2 namespaces. This could be done similar with how we created provider.
```bash
cat <<EOT >> ns1.yaml
urn: odpf-ns
name: odpf-ns
provider:
    id: 1
EOT
```

```bash
cat <<EOT >> ns2.yaml
urn: non-odpf-ns
name: non-odpf-ns
provider:
    id: 1
EOT
```

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
./siren namespace create -f ns1.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces
  --header 'content-type: application/json'
  --data-raw '{
    "urn": "odpf-ns",
    "name": "odpf-ns",
    "provider": {
        "id": 1
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
./siren namespace create -f ns2.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces
  --header 'content-type: application/json'
  --data-raw '{
    "urn": "non-odpf-ns",
    "name": "non-odpf-ns",
    "provider": {
        "id": 2
    }
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>


## 3. Verify Created Providers and Namespaces

To make sure all providers and namespaces are properly created, we could try query Siren with Siren CLI.

See what providers exist in Siren.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
./siren provider list
```
```shell
Showing 2 of 2 providers
 
ID      TYPE    URN                     NAME       
1       cortex  localhost-dev-cortex    dev-cortex

For details on a provider, try: siren provider view <id>
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/providers'`}
    </CodeBlock>
  </TabItem>
</Tabs>


See what namespaces exist in Siren.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```shell
./siren namespace list
```
```shell
Showing 2 of 2 namespaces
 
ID      URN             NAME       
1       odpf-ns         odpf-ns    
2       non-odpf-ns     non-odpf-ns

For details on a namespace, try: siren namespace view <id>
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/namespaces'`}
    </CodeBlock>
  </TabItem>
</Tabs>
