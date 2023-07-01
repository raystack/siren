# Setup Server

Siren binary contains both the CLI client and the server itself. Each has it's own configuration in order to run. Server configuration contains information such as database credentials, log severity, etc. while CLI client configuration only has configuration about which server to connect.

## Server

### Pre-requisites

Dependencies

- PostgreSQL
- [CortexMetrics](https://cortexmetrics.io/docs/getting-started/)

You need to prepare and run above dependencies first before running Siren. Siren also has a [`docker-compose.yaml`](https://github.com/raystack/siren/blob/main/docker-compose.yaml) file in its repo that has all required dependencies. If you are interested to use it, you just need to `git clone` the [repo](https://github.com/raystack/siren) and run `docker-compose up` in the root project.

### Initialization

This steps assumes all dependencies already up and running. Create a server config `config.yaml` file (`siren server init`) in the root folder of siren project or use `--config` flag to customize to a certain config file location or you can also use environment variables to provide the server config.

Setup up a database in postgres and provide the details in the DB field as given in the example below. For the purpose of this tutorial, we'll assume that the username is `postgres`, database name is `siren_development`, host and port are `localhost` and `5432`.

```yaml
db:
  driver: postgres
  url: postgres://postgres:@localhost:5432/siren_development?sslmode=disable
  ...
service:
  host: localhost
  port: 8080
  encryption_key: _ENCRYPTIONKEY_OF_32_CHARACTERS_
log:
  level: info
  gcp_compatible: true
providers:
  cortex:
    group_wait: 30s
    webhook_base_api: http://host.docker.internal:8080/v1beta1/alerts/cortex
receivers:
  slack:
    ...
```

We are using CortexMetrics as a provider. We need to set the `webhook_base_api` in `providers.cortex` config. If you are using the dockerized CortexMetrics with Docker Desktop, you could put the value as `http://host.docker.internal:8080/v1beta1/alerts/cortex`. This will enable CortexMetrics inside container to talk with Siren in our host machine later. If you are running CortexMetrics inside your host, you could put the value as `http://localhost:8080/v1beta1/alerts/cortex`.

You also might want to set the `provider.cortex.group_wait` value to `1s` so alert will be sent to the webhook immediately once it was triggered.

> If you're new to YAML and want to learn more, see Learn [YAML in Y minutes](https://learnxinyminutes.com/docs/yaml/).

### Starting the server

Database migration is required during the first server initialization. In addition, re-running the migration command might be needed in a new release to apply the new schema changes (if any). It's safer to always re-run the migration script before deploying/starting a new release.

To initialize the database schema, Run Migrations with the following command:

```sh
$ siren server migrate
```

To run the Siren server use command:

```sh
$ siren server start
```

Using `--config` flag

```sh
$ siren server migrate --config=<path-to-file>
```

```sh
$ siren server start --config=<path-to-file>
```

### Using environment variables

All the configs can be passed as environment variables using underscore `_` as the delimiter between nested keys. See the following examples

```yaml
db:
  driver: postgres
  url: postgres://postgres:@localhost:5432/siren_development?sslmode=disable
service:
  host: localhost
  port: 8080
  encryption_key: _ENCRYPTIONKEY_OF_32_CHARACTERS_
```

Here is the corresponding environment variable for the above

| Configuration key      | Environment variable   |
| ---------------------- | ---------------------- |
| db.driver              | DB_DRIVER              |
| db.url                 | DB_URL                 |
| service.host           | SERVICE_HOST           |
| service.port           | SERVICE_PORT           |
| service.encryption_key | SERVICE_ENCRYPTION_KEY |

Set the env variable using export

```
$ export SERVICE_PORT=8080
```

---

## CLI Client

### Initialization

Siren CLI supports CLI client to communicate with a Siren server. To initialize the client configuration, run the following command:

```sh
$ siren config init
```

A yaml file will be created in the `~/.config/raystack/siren.yaml` directory. Open this file to configure the host for Siren server as in the example below:

```yaml
host: "localhost:8080"
```
