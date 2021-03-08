# Siren

Alerting on data pipelines with support for multi tenancy

### Installation

#### Compiling from source

Siren requires the following dependencies:

* Docker
* Golang (version 1.15 or above)
* Git

Run the application dependecies using Docker:

```
$ docker-compose up
```

Update the configs(db credentials etc.) as per your dev machine and docker configs.

Run the following commands to compile from source

```
$ git clone git@github.com:odpf/siren.git
$ cd siren
$ go build main.go
```

To run tests locally

```
$ make test
```

To run tests locally with coverage

```
$ make test-coverage
```

To run server locally

```
$ go run main.go serve
```

To view swagger docs of HTTP APIs visit `/docs` route on the server.
e.g. [http://localhost:3000/docs](http://localhost:3000/docs)

#### Config

The config file used by application is `config.yaml` which should be present at the root of this directory.

For any variable the order of precedence is:

1. Env variable
2. Config file
3. Default in Struct defined in the application code
