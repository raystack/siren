# Schema Design

Siren stores templates, rules and triggered alerts history in PostgresDB.

We use GORM to handle database interactions and running migrations. GORM make it easier to create tables from Golang
Struct declaration.

There are the tables as of now as described below:

- **alerts:** Stores the triggered alert history.

- **templates:** Stores the templates uploaded via HTTP APIs.

- **rules:** Stores the rules configured and their state and thresholds defined.

- **slack_credentials:** Stores the slack webhook credentials which is put in alertmanager configs of a tenant

- **pagerduty_credentials:** Stores the pagerduty credentials which is put in alertmanager configs of a tenant

**Templates table:**

| Column     | Type                     | Description                                                                                          | Example                                                                                |
| ---------- | ------------------------ | ---------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| id         | bigint                   | Primary key                                                                                          | 1                                                                                      |
| created_at | timestamp with time zone | Creation timestamp                                                                                   | `2021-03-05 12:37:56.905618+05:30`                                                     |
| updated_at | timestamp with time zone | Last update timestamp                                                                                | `2021-03-05 12:37:56.905618+05:30`                                                     |
| name       | text                     | name of the template, should be unique                                                               | `cpuHigh`                                                                              |
| tags       | text[]                   | Tags array represented which resource types can use this template                                    | `{kafka, airflow}`                                                                     |
| body       | text                     | Alert or recording rule body                                                                         | See examples body in [here](../guides/templates.md)                                    |
| variables  | jsonb                    | JSON variable listing all variables in the body with their data type, description and default value. | `[{"name": "for", "type": "string", "default": "bar", "description": "group period"}]` |

**Rules Table:**

| Column     | Type                     | Description                                                                                          | Example                                  |
| ---------- | ------------------------ | ---------------------------------------------------------------------------------------------------- | ---------------------------------------- |
| id         | bigint                   | Primary key                                                                                          | 1                                        |
| created_at | timestamp with time zone | Creation timestamp                                                                                   | `2021-03-05 12:37:56.905618+05:30`       |
| updated_at | timestamp with time zone | Last update timestamp                                                                                | `2021-03-05 12:37:56.905618+05:30`       |
| namespace  | text[]                   | the ruler namespace in which this rule should be created                                             | `kafka`                                  |
| entity     | text                     | tenant name in which rule should be created                                                          | `odpf`                                   |
| group_name | text                     | the ruler namespace in which this rule should be created                                             | `testGroup`                              |
| status     | text                     | running status of alert (enabled or disabled)                                                        | `enabled`                                |
| template   | text                     | the template which should be used for rule body                                                      | `CPUHigh`                                |
| name       | text                     | name of the rule, must be unique, constructed as per `siren_api_entity_namespace_groupName_template` | `siren_api_odpf_kafka_testGroup_cpuHigh` |

**Alerts table:**

| Column       | Type                     | Description                                            | Example                            |
| ------------ | ------------------------ | ------------------------------------------------------ | ---------------------------------- |
| id           | bigint                   | Primary key                                            | 1                                  |
| created_at   | timestamp with time zone | Creation timestamp                                     | `2021-03-05 12:37:56.905618+05:30` |
| updated_at   | timestamp with time zone | Last update timestamp                                  | `2021-03-05 12:37:56.905618+05:30` |
| resource     | text                     | resource on which the alert was triggered              | `kafkaMachine1`                    |
| template     | text                     | name of template which used for this rule              | `cpuHigh`                          |
| metric_name  | text                     | the metric on which alert was triggered                | `cpu usgae %`                      |
| metric_value | text                     | value of above metric on which the alert was triggered | `95%`                              |
| level        | text                     | severity level of alert (CRITICAL, WARNING, RESOLVED)  | `CRITICAL`                         |

**Slack Credentials:**

| Column       | Type                     | Description                                                                    | Example                                 |
| ------------ | ------------------------ | ------------------------------------------------------------------------------ | --------------------------------------- |
| id           | bigint                   | Primary key                                                                    | 1                                       |
| created_at   | timestamp with time zone | Creation timestamp                                                             | `2021-03-05 12:37:56.905618+05:30`      |
| updated_at   | timestamp with time zone | Last update timestamp                                                          | `2021-03-05 12:37:56.905618+05:30`      |
| deleted_at   | timestamp with time zone | Deletion time stamp                                                            | `2021-03-05 12:37:56.905618+05:30`      |
| channel_name | text                     | name of slack channel                                                          | `siren-devs`                            |
| username     | text                     | username which will send the alert message in the channel                      | `siren-bot`                             |
| webhook      | text                     | preconfigured slack webhook                                                    | `https://hooks.slack.com/services/abcd` |
| level        | text                     | which severity levels alerts should be sent to this channel(WARNING, CRITICAL) | `WARNING`                               |
| team_name    | text                     | name of the team who owns this alert credential                                | `siren-devs`                            |
| entity       | text                     | cortex tenant name                                                             | `odpf`                                  |

**Pagerduty Credentials**

| Column      | Type                     | Description                                     | Example                            |
| ----------- | ------------------------ | ----------------------------------------------- | ---------------------------------- |
| id          | bigint                   | Primary key                                     | 1                                  |
| created_at  | timestamp with time zone | Creation timestamp                              | `2021-03-05 12:37:56.905618+05:30` |
| updated_at  | timestamp with time zone | Last update timestamp                           | `2021-03-05 12:37:56.905618+05:30` |
| deleted_at  | timestamp with time zone | Deletion time stamp                             | `2021-03-05 12:37:56.905618+05:30` |
| service_key | text                     | pagerduty service key                           | `a7se12b1iasd7da`                  |
| team_name   | text                     | name of the team who owns this alert credential | `siren-devs`                       |
| entity      | text                     | cortex tenant name                              | `odpf`                             |
