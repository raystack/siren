# Providers

Siren providers represent a monitoring server. We define alerts and alerts routing configuration inside providers. For
example the below provider, describes how a cortex monitoring server info is stored inside Siren.

```json
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

# Namespaces

Monitoring providers usually have tenants, a sharded way of storing and querying telemetry data. Siren calls them  
**namespaces**. You should create a namespaces with the same name as the tenant of your monitoring providers. Example:
the below namespace represent the tenant "odpf".

```json
{
  "id": "10",
  "urn": "odpf",
  "name": "odpf",
  "provider": "3",
  "credentials": {},
  "created_at": "2021-10-26T06:00:52.627917Z",
  "updated_at": "2021-10-26T06:00:52.627917Z"
}
```

# API Interface

## Providers API interface

### Providers creation

```text
POST /v1beta1/providers HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 154

{
  "host": "http://localhost:9009",
  "urn": "localhost_cortex",
  "name": "localhost_cortex",
  "type": "cortex",
  "credentials": {},
  "labels": {}
}
```

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

### Providers update

```text
PUT /v1beta1/providers/4 HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 155

{
  "host": "http://localhost:9009",
  "urn": "localhost_cortex",
  "name": "localhost_cortex_1",
  "type": "cortex",
  "credentials": {},
  "labels": {}
}
```

### Getting a provider

```text
GET /v1beta1/providers/4 HTTP/1.1
Host: localhost:3000
```

### Getting all providers

```text
GET /v1beta1/providers HTTP/1.1
Host: localhost:3000
```

### Deleting a provider

```text
DELETE /v1beta1/providers/4 HTTP/1.1
Host: localhost:3000
```

**Note:**

1. Before deleting the provider, you will need to delete dependant resources(namespaces).

## Namespace API interface

**Note:** These operations on namespaces inside Siren doesn't affect the actual tenant in the monitoring provider. The
user should manage the tenants themselves.

### Creating a namespace

```text
POST /v1beta1/namespaces HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 103

{
    "name": "test",
    "urn": "test",
    "provider": "5",
    "credentials": {},
    "labels": {}
}
```

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

```text
PUT /v1beta1/namespaces/13 HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 104

{
    "name": "test2",
    "urn": "test",
    "provider": "5",
    "credentials": {},
    "labels": {}
}
```

### Getting a namespace

```text
GET /v1beta1/namespaces/13 HTTP/1.1
Host: localhost:3000
```

### Getting all namespace

```text
GET /v1beta1/namespaces HTTP/1.1
Host: localhost:3000
```

### Deleting a namespace

```text
GET /v1beta1/namespaces/13 HTTP/1.1
Host: localhost:3000
```

# CLI Interface

## Provider CLI interface

With CLI, you will need a YAML file in the below specified format to create/update providers.

**Example usage**

```text
Work with providers.

Providers are the system for which we intend to mange monitoring and alerting.

Usage:
  siren provider [command]

Aliases:
  provider, providers

Available Commands:
  create      Create a new provider
  delete      Delete a provider details
  edit        Edit a provider
  list        List providers
  view        View a provider details

Flags:
  -h, --help   help for provider

Use "siren provider [command] --help" for more information about a command.
```

### Provider create

```yaml
# input.yaml
host: http://localhost:9009
urn: localhost-dev-cortex
name: dev-cortex
type: cortex
```

```shell
go run main.go providers create --file test.yam
```

## Namespace CLI interface

```text
Work with namespaces.

namespaces are used for multi-tenancy for a given provider.

Usage:
  siren namespace [command]

Aliases:
  namespace, namespaces

Available Commands:Ò
  create      Create a new namespace
  delete      Delete a namespace details
  edit        Edit a namespace
  list        List namespaces
  view        View a namespace details

Flags:
  -h, --help   help for namespace

Use "siren namespace [command] --help" for more information about a command.
```

The usage is same as providers.
