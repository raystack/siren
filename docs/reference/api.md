# Siren
Documentation of our Siren API.

## Version: 1.0.0

### Security
**basic**

|basic|*Basic*|
|---|---|

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

### /notifications

#### POST

##### Description

POST Notifications API This API sends notifications to configured channel

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| provider | query |  | No | string |
| Body | body |  | No | [SlackMessage](#slackmessage) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | POST codeExchange response | [OAuthExchangeResponse](#oauthexchangeresponse) |

### /oauth/slack/token

#### POST

##### Description

POST Code Exchange API This API exchanges oauth code with access token from slack server. client_id and client_secret
are read from Siren ENV vars.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Body | body |  | No | [OAuthPayload](#oauthpayload) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | POST codeExchange response | [OAuthExchangeResponse](#oauthexchangeresponse) |

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

### /teams/{teamName}/credentials

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

### /workspaces/{workspaceName}/channels

#### GET

##### Description

Get Channels API: This API gets the list of joined channels within a slack workspace

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| workspaceName | path | name of the workspace | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  | [ integer ] |

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
| created_at | dateTime |  | No |
| id | integer (uint64) |  | No |
| level | string |  | No |
| metric_name | string |  | No |
| metric_value | string |  | No |
| name | string |  | No |
| template_id | string |  | No |
| updated_at | dateTime |  | No |

#### Alerts

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| alerts | [ [Alert](#alert) ] |  | No |

#### Annotations

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| metric_name | string |  | No |
| metric_value | string |  | No |
| resource | string |  | No |
| template | string |  | No |

#### Block

Block defines an interface all block types should implement to ensure consistency between blocks.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| BlockType | [MessageBlockType](#messageblocktype) |  | No |

#### Labels

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| severity | string |  | No |

#### MessageBlockType

MessageBlockType defines a named string type to define each block type as a constant for use within the package.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| MessageBlockType | string | MessageBlockType defines a named string type to define each block type as a constant for use within the package. |  |

#### OAuthExchangeResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| ok | boolean |  | No |

#### OAuthPayload

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | string |  | No |
| workspace | string |  | No |

#### Rule

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created_at | dateTime |  | No |
| entity | string |  | No |
| group_name | string |  | No |
| id | integer (uint64) |  | No |
| name | string |  | No |
| namespace | string |  | No |
| status | string |  | No |
| template | string |  | No |
| updated_at | dateTime |  | No |
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

#### SlackMessage

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| blocks | [ [Block](#block) ] |  | No |
| entity | string |  | No |
| message | string |  | No |
| receiver_name | string |  | No |
| receiver_type | string |  | No |

#### SlackMessageSendResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| ok | boolean |  | No |

#### Template

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| body | string |  | No |
| created_at | dateTime |  | No |
| id | integer (uint64) |  | No |
| name | string |  | No |
| tags | [ string ] |  | No |
| updated_at | dateTime |  | No |
| variables | [ [Variable](#variable) ] |  | No |

#### Variable

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| default | string |  | No |
| description | string |  | No |
| name | string |  | No |
| type | string |  | No |
