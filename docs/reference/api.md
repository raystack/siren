# Siren
Documentation of our Siren API.

## Version: 1.0.0

### Security
**basic**  

|basic|*Basic*|
|---|---|

### /alertingCredentials/teams/{teamName}

#### GET
##### Description

Get AlertCredentials API: This API helps in getting the teams slack and pagerduty credentials

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| teamName | path | name of the team | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | AlertCredentialResponse | [AlertCredentialResponse](#alertcredentialresponse) |

#### PUT
##### Description

Upsert AlertCredentials API: This API helps in creating or updating the teams slack and pagerduty credentials

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Body | body | Create AlertCredential request | No | [AlertCredentialResponse](#alertcredentialresponse) |
| teamName | path | name of the team | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 |  |

### /history

#### GET
##### Description

GET Alert History API: This API lists stored alert history for given filers in query params

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| resource | query |  | No | string |
| startTime | query |  | No | integer (uint32) |
| endTime | query |  | No | integer (uint32) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Get alertHistory response | [ [AlertHistoryObject](#alerthistoryobject) ] |

#### POST
##### Description

Create Alert History API: This API create alert history

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Body | body |  | No | [ [Alerts](#alerts) ] |

### /ping

#### GET
##### Description

Ping call

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | Response body for Ping. |

### /rules

#### GET
##### Description

List Rules API: This API lists all the existing rules with given filers in query params

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| namespace | query | List Rule Request | No | string |
| entity | query |  | No | string |
| group_name | query |  | No | string |
| status | query |  | No | string |
| template | query |  | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List rules response | [ [Rule](#rule) ] |

#### PUT
##### Description

Upsert Rule API: This API helps in creating a new rule or update an existing one with unique combination of namespace, entity, group_name, template

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Body | body | Create rule request | No | [Rule](#rule) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  | [Rule](#rule) |

### /templates

#### GET
##### Description

List Templates API: This API lists all the existing templates with given filers in query params

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| tag | query | List Template Request | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List templates response | [ [Template](#template) ] |

#### PUT
##### Description

Upsert Templates API: This API helps in creating or updating a template with unique name

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Body | body | Create template request | No | [Template](#template) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  | [Template](#template) |

### /templates/{name}

#### DELETE
##### Description

Delete Template API: This API deletes a template given the template name

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path | Delete Template Request | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  | [Template](#template) |

#### GET
##### Description

Get Template API: This API gets a template given the template name

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path | Get Template Request | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  | [Template](#template) |

### /templates/{name}/render

#### POST
##### Description

Render Template API: This API renders the given template with given values

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path | Render Template Request | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 |  |

### Models

#### Alert

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| annotations | [Annotations](#annotations) |  | No |
| labels | [Labels](#labels) |  | No |
| status | string |  | No |

#### AlertCredentialResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| entity | string |  | No |
| pagerduty_credentials | string |  | No |
| slack_config | [SlackConfig](#slackconfig) |  | No |

#### AlertHistoryObject

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created | dateTime |  | No |
| id | integer (uint64) |  | No |
| level | string |  | No |
| metric_name | string |  | No |
| metric_value | string |  | No |
| name | string |  | No |
| template_id | string |  | No |
| updated | dateTime |  | No |

#### Alerts

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| alerts | [ [Alert](#alert) ] |  | No |

#### Annotations

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| metricName | string |  | No |
| metricValue | string |  | No |
| resource | string |  | No |
| template | string |  | No |

#### Labels

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| severity | string |  | No |

#### Rule

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| CreatedAt | dateTime |  | No |
| UpdatedAt | dateTime |  | No |
| entity | string |  | No |
| group_name | string |  | No |
| id | integer (uint64) |  | No |
| name | string |  | No |
| namespace | string |  | No |
| status | string |  | No |
| template | string |  | No |
| variables | [ [RuleVariable](#rulevariable) ] |  | No |

#### RuleVariable

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | No |
| name | string |  | No |
| type | string |  | No |
| value | string |  | No |

#### SlackConfig

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| critical | [SlackCredential](#slackcredential) |  | No |
| warning | [SlackCredential](#slackcredential) |  | No |

#### SlackCredential

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| channel | string |  | No |
| username | string |  | No |
| webhook | string |  | No |

#### Template

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| CreatedAt | dateTime |  | No |
| UpdatedAt | dateTime |  | No |
| body | string |  | No |
| id | integer (uint64) |  | No |
| name | string |  | No |
| tags | [ string ] |  | No |
| variables | [ [Variable](#variable) ] |  | No |

#### Variable

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| default | string |  | No |
| description | string |  | No |
| name | string |  | No |
| type | string |  | No |
