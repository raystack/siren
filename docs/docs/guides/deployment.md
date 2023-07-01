# Server Installation

There are several approaches to setup Siren Server

1. [Using the CLI](#using-the-cli)
2. [Using the Docker](#use-the-docker-image)
3. [Using the Helm Chart](#use-the-helm-chart)

## General pre-requisites

- PostgreSQL (version 13 or above)
- Monitoring Providers
  - Ex: CortexMetrics

## Using the CLI

### Pre-requisites for CLI

- [Create siren config file](../tour/setup_server.md#initialization)

To run the Siren server use command:

```sh
$ siren server start -c <path-to-config>
```

## Use the Docker

To run the Siren server using Docker, you need to have Docker installed on your system. You can find the installation instructions [here](https://docs.docker.com/get-docker/).

You can choose to set the configuration using environment variables or a config file. The environment variables will override the config file.

### Using environment variables

All the configs can be passed as environment variables using underscore `_` as the delimiter between nested keys. See the following examples

See [configuration reference](../reference/server_configuration.md) for the list of all the configuration keys.

```sh title=".env"
DB_DRIVER=postgres
DB_URL=postgres://postgres:@localhost:5432/siren_development?sslmode=disable
SERVICE_PORT=8080
SERVICE_ENCRYPTION_KEY=<32 characters encryption key>
```

Run the following command to start the server

```sh
$ docker run -d \
    --restart=always \
    -p 8080:8080 \
    --env-file .env \
    --name siren-server \
    raystack/siren:<version> \
    server start
```

### Using config file

```yaml title="config.yaml"
db:
  driver: postgres
  url: postgres://postgres:@localhost:5432/siren_integration?sslmode=disable
service:
  port: 8080
  encryption_key: <32 characters encryption key>
```

Run the following command to start the server

```sh
$ docker run -d \
    --restart=always \
    -p 8080:8080 \
    -v $(pwd)/config.yaml:/config.yaml \
    --name siren-server \
    raystack/siren:<version> \
    server start -c /config.yaml
```

## Use the Helm chart

### Pre-requisites for Helm chart

Siren can be installed in Kubernetes using the Helm chart from https://github.com/raystack/charts.

Ensure that the following requirements are met:

- Kubernetes 1.14+
- Helm version 3.x is [installed](https://helm.sh/docs/intro/install/)

### Add Raystack Helm repository

Add Raystack chart repository to Helm:

```
helm repo add raystack https://raystack.github.io/charts/
```

You can update the chart repository by running:

```
helm repo update
```

### Setup helm values

The following table lists the configurable parameters of the Siren chart and their default values.

See full helm values guide [here](https://github.com/raystack/charts/tree/main/stable/siren#values).

```yaml title="values.yaml"
app:
  ## Value to fully override guardian.name template
  nameOverride: ""
  ## Value to fully override guardian.fullname template
  fullnameOverride: ""

  image:
    repository: raystack/siren
    pullPolicy: Always
    tag: latest
  container:
    args:
      - server
      - start
    livenessProbe:
      httpGet:
        path: /ping
        port: tcp
    readinessProbe:
      httpGet:
        path: /ping
        port: tcp

  migration:
    enabled: true
    args:
      - server
      - migrate

  service:
    annotations:
      projectcontour.io/upstream-protocol.h2c: tcp

  ingress:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: contour
    hosts:
      - host: siren.example.com
        paths:
          - path: /
            pathType: ImplementationSpecific
            backend:
              service:
                # name: backend_01
                port:
                  number: 8080

  config:
    LOG_LEVEL: info
    SERVICE_PORT: 8080

  secretConfig:
    ENCRYPTION_SECRET_KEY:
    NOTIFIER_ACCESS_TOKEN:
    DB_URL: postgres://siren:<password>@localhost:5432/siren_integration?sslmode=disable
```

And install it with the helm command line along with the values file:

```sh
$ helm install my-release -f values.yaml raystack/siren
```
