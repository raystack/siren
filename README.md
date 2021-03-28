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

To view swagger docs of HTTP APIs visit `/documentation` route on the server.
e.g. [http://localhost:3000/documentation](http://localhost:3000/documentation)

#### Config

The config file used by application is `config.yaml` which should be present at the root of this directory.

For any variable the order of precedence is:

1. Env variable
2. Config file
3. Default in Struct defined in the application code

### HTTP Client

The `client` directory holds the HTTP Client for siren service. It's generated using
project [swagger-codegen](https://github.com/swagger-api/swagger-codegen)

Ideally we should generate the client on any changes in the swagger spec of siren service.

The config used for client generation is `client_config.json`

To regenerate the client, run

```
$ swagger-codegen generate -i swagger.yaml -l go -o client -c client_config.json
```

Sample usage of the client:

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antihax/optional"
	"github.com/odpf/siren/client"
)

func main() {
	cfg := &client.Configuration{
		BasePath: "http://localhost:3000",
	}
	x := client.NewAPIClient(cfg)
	options := &client.RulesApiListRulesRequestOpts{
		Namespace: optional.NewString("n1"),
	}
	result, _, err := x.RulesApi.ListRulesRequest(context.Background(), options)
	if err != nil {
		panic(err)
	}
	response, _ := json.Marshal(result)
	fmt.Println(string(response))
}
```

### List of available commands

1. Migrate
    - Runs the DB Migrations `$ go run main.go serve`

2. Upload
    - Parses a YAML File in specified format to upsert templates and rules(alerts) `$ go run main.go upload fileName.yaml`

   **Note:** Updating a template will update all associated rules(alerts). If some new variable is introduced in tempalte, 
   one should always give default values. 
   
#### Sample template YAML File

```yaml
apiVersion: v2
type: template
name: qa
body:
  - alert: QAWarning
    expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.warning]]
    for: "[[.for]]"
    labels:
      severity: WARNING
    annotations:
      description: CPU has been above [[.warning]] for last [[.for]] {{ $labels.host }}
  - alert: QACritical
    expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.critical]]
    for: "[[.for]]"
    labels:
      severity: CRITICAL
    annotations:
      description: CPU has been above [[.critical]] for last [[.for]] {{ $labels.host }}
variables:
  - name: for
    type: string
    default: 10m
    description: For eg 5m, 2h; Golang duration format
  - name: warning
    type: int
    default: 10
  - name: critical
    type: int
    default: 20
tags:
  - systems
```

#### Sample rules YAML file

```yaml
apiVersion: v2
type: rule
namespace: testNamespace
entity: real
rules:
  TestGroup:
    template: test
    status: enabled
    variables:
      - name: for
        value: 15m
      - name: warning
        value: 185
      - name: critical
        value: 195
```
