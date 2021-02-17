# Siren

Alerting on data pipelines with support for multi tenancy

### Installation

#### Compiling from source
Siren requires the following dependencies:

* Docker
* Golang (version 1.14 or above)
* Git

Run the following commands to compile from source
```
$ git clone git@github.com:odpf/siren.git
$ cd siren
$ go build main.go
```

To run tests locally
```
$ go test
```

To run server locally
```
$ go run main.go serve
```

