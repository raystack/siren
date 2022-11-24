# Server Configuration

Server configuration in siren is required to configure server, workers, and jobs. We can generate the default configuration with Siren CLI.
```bash
siren server init
```
Above command will generate a `./config.yaml` file in the same folder. When starting the server, Siren server will auto detect the `./config.yaml` and read all configs inside it to be used when starting up the server. Below is the Siren server configuration.

```yaml
db:
  driver: <string>

  url: <string>

  max_idle_conns: <int>

  max_open_conns: <int>

  # db connection max life time config e.g. 10ms
  conn_max_life_time: <string duration> | default="10ms"
  
  # db connection max query timeout config e.g. 100ms
  max_query_timeout: <string duration> | default="100ms"

# instrumentation/metrics related configurations.
telemetry:
  # debug_addr is used for exposing the pprof, zpages & `/metrics` endpoints. if
  # not set, all of the above are disabled.
  debug_addr: <string> | default="localhost:8081"

  # enable_cpu enables collection of runtime CPU metrics. available in `/metrics`.
  enable_cpu: <bool> | default=true

  # enable_memory enables collection of runtime memory metrics. available via `/metrics`.
  enable_memory: <bool> | default=true

  # sampling_fraction indicates the sampling rate for tracing. 1 indicates all traces
  # are collected and 0 means no traces.
  sampling_fraction: <bool> | default=1

  # service_name is the identifier used in trace exports, NewRelic, etc for the
  # dex instance.
  service_name: <string> | default="siren"

  # enable_newrelic enables exporting NewRelic instrumentation in addition to the
  # OpenCensus.
  enable_newrelic: <bool> | default=false

  # new relic app name, if left empty, app name will be service_name
  newrelic_app_name: <string> | default=""

  # newrelic_api_key must be a valid NewRelic License key.
  newrelic_api_key: <string> | default="____LICENSE_STRING_OF_40_CHARACTERS_____"

  # enable_otel_agent enables the OpenTelemetry Exporter for both traces and views.
  enable_otel_agent: <bool> | default=false

  # otel_agent_addr is the addr of OpenTelemetry Collector/Agent. This is where the
  # opene-telemetry exporter will publish the collected traces/views to.
  otel_agent_addr: <string> | default="localhost:8088"

service:
  host: <string> | default="localhost"

  port: <int> | default=8080
  
  encryption_key: <string> | default="_ENCRYPTIONKEY_OF_32_CHARACTERS_"
  
log:
  level: <string> | default="info"

  # log format will be compatible with gcp logging if this is set to true
  gcp_compatible: <bool> | default=true

providers:
  cortex:
    group_wait: <string> | default="30s"

    webhook_base_api: <string> | default="http://localhost:8080/v1beta1/alerts/cortex"

    http_client:
      <httpclient>

receivers:
  slack:
    # host of slack api, default value is hardcoded as `https://slack.com/api`
    apihost: <string> | default=""
    
    retry:
      <retry>
      
    httpclient:
      <httpclient>

  pagerduty:
    # host of pagerduty api, default value is hardcoded as `https://events.pagerduty.com`
    api_host: <string> | default=""

    retry:
      <retry>
      
    httpclient:
      <httpclient>
      
  http:
    retry:
      <retry>
      
    httpclient:
      <httpclient>

notification:
  queue:
    # queue to use (supported are: inmemory, postgres)
    kind: <string> | default="inmemory"

  message_handler:
    <message_handler>

  dlq_handler:
    <message_handler>
```

The `<retry>` block above could be represented like below.
```yaml
retry:
    # duration to wait before retrying a call to api
    wait_duration: <string duration> | default="20ms"
    
    enable_backoff: <bool> | default=false

    # number of trial the client does the work (e.g. api call)
    max_tries: 3

    # won't retry the call if there is a failure if enable is false
    enable: <bool> | default=true
```

The `<httpclient>` block above could be represented like below.
```yaml
httpclient:
    # if set to 0, will use the default value from net/http library DefaultTransport: 30000
    timeout_ms: <int> | default=0

    # if set to 0, will use the default value from net/http library: 0 means no limit
    max_conns_per_host: <int> | default=0

    # if set to 0, will use the default value from net/http library DefaultTransport: 100
    max_idle_conns: <int> | default=0

    # if set to 0, will use the default value from net/http library: 2
    max_idle_conns_per_host: <int> | default=0

    # if set to 0, will use the default value from net/http library DefaultTransport: 90000
    idle_conn_timeout_ms: <int> | default=0
```

The `<message_handler>` block above could be represented like below.
```yaml
message_handler:
    # disable message handler worker if `enabled` is false
    enabled: <bool> | default=true

    # duration to dequeue and publish messages
    poll_duration: <string duration> | default="5s"

    # types of receiver that need to be supported by the handler (e.g. slack, http, pagerduty, file)
    receiver_types: <list of string> | default="[slack, http, pagerduty, file]"\

    # number of messages to dequeue and publish at once
    batch_size: <int> | default=1
```

**Convert YAML to Environment Variable**
If you prefer to use env variable instead of a yaml file. You could also represent the config in the env variable. Each alphanumeric character in config need to be uppercased and the nested config is merged into a single word separated by an underscore `_`. This is similar like what [viper](https://github.com/spf13/viper) does.

__Example__

```yaml
# yaml config
db:
  driver: postgres
  url: postgres://postgres:@localhost:5432/siren_development?sslmode=disable
newrelic:
  license: ____LICENSE_STRING_OF_40_CHARACTERS_____
service:
  port: 8080
  encryption_key: _ENCRYPTIONKEY_OF_32_CHARACTERS_
```
The environment variable will be like this.
```bash
DB_DRIVER=postgres
DB_URL=postgres://postgres:@localhost:5432/siren_development?sslmode=disable
NEWRELIC_LICENSE=____LICENSE_STRING_OF_40_CHARACTERS_____
SERVICE_PORT=8080
SERVICE_ENCRYPTION_KEY=_ENCRYPTIONKEY_OF_32_CHARACTERS_
```

## How to configure

There are 3 ways to configure siren:

- Using env variables
- Using a yaml file
- or using a combination of both

### Using env variables

Example:

```sh
export PORT=9999
siren server start
```

This will run the service on port 9999 instead of the default 8080

### Using a yaml file

For default values and the structure of the yaml file, generate yaml config file with:
```bash
siren server init
```
This will generate a `./config.yaml` file. Now you can make modification to the config yaml as you wish and then start Siren server.

```sh
siren server start
```

### Using a combination of both

If any key that is set via both env vars and yaml the value set in env vars will take effect.
