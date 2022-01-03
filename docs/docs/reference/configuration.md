# Configuration

| Go struct                     | YAML path        | ENV var          | default               | Valid values                                                                                                     |
| ----------------------------- | ---------------- | ---------------- | --------------------- | ---------------------------------------------------------------------------------------------------------------- |
| Config.LogConfig.Level        | log.level        | LOG_LEVEL        | info                  | debug,info,warn,error,dpanic,panic,fatal                                                                         |
| Config.Port                   | port             | PORT             | 8080                  | 0-65535                                                                                                          |
| Config.DBConfig.User          | db.user          | DB_USER          | postgres              | [PostgreSQL identifiers](https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS) |
| Config.DBConfig.Name          | db.name          | DB_NAME          | postgres              | [PostgreSQL identifiers](https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS) |
| Config.DBConfig.Port          | db.port          | DB_PORT          | 5432                  | 0-65535                                                                                                          |
| Config.DBConfig.Password      | db.password      | DB_PASSWORD      |                       | valid PostgreSQL password                                                                                        |
| Config.DBConfig.SslMode       | db.sslmode       | DB_SSLMODE       | disable               | [libpq sslmode](https://www.postgresql.org/docs/9.1/libpq-ssl.html)                                              |
| Config.DBConfig.LogLevel      | db.log_level     | DB_LOG_LEVEL     | info                  | silent,error,warn,info                                                                                           |
| Config.DBConfig.Host          | db.host          | DB_HOST          | localhost             | valid hostname name or IP address                                                                                |
| Config.NewRelicConfig.Enabled | newrelic.enabled | NEWRELIC_ENABLED | false                 | bool                                                                                                             |
| Config.NewRelicConfig.License | newrelic.license | NEWRELIC_LICENSE |                       | 40 char NewRelic license key                                                                                     |
| Config.NewRelicConfig.AppName | newrelic.appname | NEWRELIC_APPNAME | siren                 | string                                                                                                           |

## How to configure

There are 3 ways to configure siren:

- Using env variables
- Using a yaml file
- or using a combination of both

### Using env variables

Example:

```sh
export PORT=9999
go run main.go serve
```

This will run the service on port 9999 instead of the default 8080

### Using a yaml file

For default values and the structure of the yaml file please check file - `config.yaml.example`

Usage example:

```sh
cp config.yaml.example config.yaml
# make any modifications to the configs as required
go run main.go serve
```

### Using a combinnation of both

If any key that is set via both env vars and yaml the value set in env vars will take effect.
