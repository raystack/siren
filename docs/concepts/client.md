# HTTP Client

The `client` directory holds the HTTP Client for siren service. It's generated using
project [swagger-codegen](https://github.com/swagger-api/swagger-codegen) with `swagger.yaml` file which can be found at
the root.

This client is used by command line interface to parse YAML files and call the HTTP APIs of Siren service to create or
update templates and rules.

Ideally we should generate the client on any changes in the swagger spec of siren service.

The config used for client generation is `client_config.json`

To generate the client, run

```shell
$ swagger-codegen generate -i api/handlers/swagger.yaml -l go -o client -c client_config.json
```

**Sample usage**

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
