# Siren APIs

Documentation of our Siren API with gRPC and gRPC-Gateway.

## Version: 0.3.0

### /ping

#### GET

##### Summary

ping

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1PingResponse](#v1beta1pingresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/alerts/cortex/{providerId}

#### POST

##### Summary

create cortex alerts

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| providerId | path |  | Yes | string (uint64) |
| body | body |  | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Alerts](#v1beta1alerts) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/alerts/{providerName}/{providerId}

#### GET

##### Summary

list alerts

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| providerName | path |  | Yes | string |
| providerId | path |  | Yes | string (uint64) |
| resourceName | query |  | No | string |
| startTime | query |  | No | string (uint64) |
| endTime | query |  | No | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Alerts](#v1beta1alerts) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/namespaces

#### GET

##### Summary

list namespaces

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1ListNamespacesResponse](#v1beta1listnamespacesresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### POST

##### Summary

create a namespace

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [v1beta1CreateNamespaceRequest](#v1beta1createnamespacerequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Namespace](#v1beta1namespace) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/namespaces/{id}

#### GET

##### Summary

get a namespace

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Namespace](#v1beta1namespace) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### DELETE

##### Summary

delete a namespace

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. |  |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### PUT

##### Summary

update a namespace

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |
| body | body |  | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Namespace](#v1beta1namespace) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/providers

#### GET

##### Summary

list providers

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| urn | query |  | No | string |
| type | query |  | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1ListProvidersResponse](#v1beta1listprovidersresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### POST

##### Summary

create a provider

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [v1beta1CreateProviderRequest](#v1beta1createproviderrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Provider](#v1beta1provider) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/providers/{id}

#### GET

##### Summary

get a provider

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Provider](#v1beta1provider) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### DELETE

##### Summary

delete a provider

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. |  |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### PUT

##### Summary

update a provider

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |
| body | body |  | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Provider](#v1beta1provider) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/receivers

#### GET

##### Summary

list receivers

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1ListReceiversResponse](#v1beta1listreceiversresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### POST

##### Summary

create a receiver

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [v1beta1CreateReceiverRequest](#v1beta1createreceiverrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Receiver](#v1beta1receiver) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/receivers/{id}

#### GET

##### Summary

get a receiver

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Receiver](#v1beta1receiver) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### DELETE

##### Summary

delete a receiver

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. |  |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### PUT

##### Summary

update a receiver

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |
| body | body |  | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Receiver](#v1beta1receiver) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/receivers/{id}/send

#### POST

##### Summary

send notification to receiver

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |
| body | body |  | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1SendReceiverNotificationResponse](#v1beta1sendreceivernotificationresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/rules

#### GET

##### Summary

list rules

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | query |  | No | string |
| namespace | query |  | No | string |
| groupName | query |  | No | string |
| template | query |  | No | string |
| providerNamespace | query |  | No | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1ListRulesResponse](#v1beta1listrulesresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### PUT

##### Summary

add/update a rule

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [v1beta1UpdateRuleRequest](#v1beta1updaterulerequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1UpdateRuleResponse](#v1beta1updateruleresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/subscriptions

#### GET

##### Summary

List subscriptions

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1ListSubscriptionsResponse](#v1beta1listsubscriptionsresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### POST

##### Summary

Create a subscription

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [v1beta1CreateSubscriptionRequest](#v1beta1createsubscriptionrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Subscription](#v1beta1subscription) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/subscriptions/{id}

#### GET

##### Summary

Get a subscription

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Subscription](#v1beta1subscription) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### DELETE

##### Summary

Delete a subscription

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. |  |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### PUT

##### Summary

Update a subscription

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string (uint64) |
| body | body |  | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1Subscription](#v1beta1subscription) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/templates

#### GET

##### Summary

list templates

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| tag | query |  | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1ListTemplatesResponse](#v1beta1listtemplatesresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### PUT

##### Summary

add/update a template

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [v1beta1UpsertTemplateRequest](#v1beta1upserttemplaterequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1TemplateResponse](#v1beta1templateresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/templates/{name}

#### GET

##### Summary

get a template

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path |  | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1TemplateResponse](#v1beta1templateresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### DELETE

##### Summary

delete a template

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path |  | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1DeleteTemplateResponse](#v1beta1deletetemplateresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/templates/{name}/render

#### POST

##### Summary

render a template

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | path |  | Yes | string |
| body | body |  | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1RenderTemplateResponse](#v1beta1rendertemplateresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### Models

#### SendReceiverNotificationRequestSlackPayload

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |
| receiverName | string |  | No |
| receiverType | string |  | No |
| blocks | [ object ] |  | No |

#### protobufAny

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| @type | string |  | No |

#### protobufNullValue

`NullValue` is a singleton enumeration to represent the null value for the
`Value` type union.

The JSON representation for `NullValue` is JSON `null`.

- NULL_VALUE: Null value.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| protobufNullValue | string | `NullValue` is a singleton enumeration to represent the null value for the `Value` type union. The JSON representation for `NullValue` is JSON `null`. - NULL_VALUE: Null value. |  |

#### rpcStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |
| message | string |  | No |
| details | [ [protobufAny](#protobufany) ] |  | No |

#### v1beta1Alert

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |
| providerId | string (uint64) |  | No |
| resourceName | string |  | No |
| metricName | string |  | No |
| metricValue | string |  | No |
| severity | string |  | No |
| rule | string |  | No |
| triggeredAt | dateTime |  | No |

#### v1beta1Alerts

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| alerts | [ [v1beta1Alert](#v1beta1alert) ] |  | No |

#### v1beta1Annotations

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| metricName | string |  | No |
| metricValue | string |  | No |
| resource | string |  | No |
| template | string |  | No |

#### v1beta1CortexAlert

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| annotations | [v1beta1Annotations](#v1beta1annotations) |  | No |
| labels | [v1beta1Labels](#v1beta1labels) |  | No |
| status | string |  | No |
| startsAt | dateTime |  | No |

#### v1beta1CreateNamespaceRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |
| urn | string |  | No |
| provider | string (uint64) |  | No |
| credentials | object |  | No |
| labels | object |  | No |
| createdAt | dateTime |  | No |
| updatedAt | dateTime |  | No |

#### v1beta1CreateProviderRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| host | string |  | No |
| urn | string |  | No |
| name | string |  | No |
| type | string |  | No |
| credentials | object |  | No |
| labels | object |  | No |

#### v1beta1CreateReceiverRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |
| type | string |  | No |
| labels | object |  | No |
| configurations | object |  | No |

#### v1beta1CreateSubscriptionRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| urn | string |  | No |
| namespace | string (uint64) |  | No |
| receivers | [ [v1beta1ReceiverMetadata](#v1beta1receivermetadata) ] |  | No |
| match | object |  | No |

#### v1beta1DeleteTemplateResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| v1beta1DeleteTemplateResponse | object |  |  |

#### v1beta1Labels

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| severity | string |  | No |

#### v1beta1ListNamespacesResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| namespaces | [ [v1beta1Namespace](#v1beta1namespace) ] |  | No |

#### v1beta1ListProvidersResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| providers | [ [v1beta1Provider](#v1beta1provider) ] |  | No |

#### v1beta1ListReceiversResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| receivers | [ [v1beta1Receiver](#v1beta1receiver) ] |  | No |

#### v1beta1ListRulesResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| rules | [ [v1beta1Rule](#v1beta1rule) ] |  | No |

#### v1beta1ListSubscriptionsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| subscriptions | [ [v1beta1Subscription](#v1beta1subscription) ] |  | No |

#### v1beta1ListTemplatesResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| templates | [ [v1beta1Template](#v1beta1template) ] |  | No |

#### v1beta1Namespace

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |
| urn | string |  | No |
| name | string |  | No |
| provider | string (uint64) |  | No |
| credentials | object |  | No |
| labels | object |  | No |
| createdAt | dateTime |  | No |
| updatedAt | dateTime |  | No |

#### v1beta1PingResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |

#### v1beta1Provider

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |
| host | string |  | No |
| urn | string |  | No |
| name | string |  | No |
| type | string |  | No |
| credentials | object |  | No |
| labels | object |  | No |
| createdAt | dateTime |  | No |
| updatedAt | dateTime |  | No |

#### v1beta1Receiver

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |
| name | string |  | No |
| type | string |  | No |
| labels | object |  | No |
| configurations | object |  | No |
| data | object |  | No |
| createdAt | dateTime |  | No |
| updatedAt | dateTime |  | No |

#### v1beta1ReceiverMetadata

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |
| configuration | object |  | No |

#### v1beta1RenderTemplateResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| body | string |  | No |

#### v1beta1Rule

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |
| name | string |  | No |
| enabled | boolean |  | No |
| groupName | string |  | No |
| namespace | string |  | No |
| template | string |  | No |
| variables | [ [v1beta1Variables](#v1beta1variables) ] |  | No |
| createdAt | dateTime |  | No |
| updatedAt | dateTime |  | No |
| providerNamespace | string (uint64) |  | No |

#### v1beta1SendReceiverNotificationResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| ok | boolean |  | No |

#### v1beta1Subscription

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |
| urn | string |  | No |
| namespace | string (uint64) |  | No |
| receivers | [ [v1beta1ReceiverMetadata](#v1beta1receivermetadata) ] |  | No |
| match | object |  | No |
| createdAt | dateTime |  | No |
| updatedAt | dateTime |  | No |

#### v1beta1Template

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |
| name | string |  | No |
| body | string |  | No |
| tags | [ string ] |  | No |
| createdAt | dateTime |  | No |
| updatedAt | dateTime |  | No |
| variables | [ [v1beta1TemplateVariables](#v1beta1templatevariables) ] |  | No |

#### v1beta1TemplateResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| template | [v1beta1Template](#v1beta1template) |  | No |

#### v1beta1TemplateVariables

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |
| type | string |  | No |
| default | string |  | No |
| description | string |  | No |

#### v1beta1UpdateRuleRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| enabled | boolean |  | No |
| groupName | string |  | No |
| namespace | string |  | No |
| template | string |  | No |
| variables | [ [v1beta1Variables](#v1beta1variables) ] |  | No |
| providerNamespace | string (uint64) |  | No |

#### v1beta1UpdateRuleResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| rule | [v1beta1Rule](#v1beta1rule) |  | No |

#### v1beta1UpsertTemplateRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |
| name | string |  | No |
| body | string |  | No |
| tags | [ string ] |  | No |
| variables | [ [v1beta1TemplateVariables](#v1beta1templatevariables) ] |  | No |

#### v1beta1Variables

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |
| value | string |  | No |
| type | string |  | No |
| description | string |  | No |
