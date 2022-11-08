# Siren APIs
Documentation of our Siren API with gRPC and
gRPC-Gateway.

## Version: 0.4.0

### /v1beta1/alerts/{provider_type}/{provider_id}

#### GET
##### Summary

list alerts

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| provider_type | path |  | Yes | string |
| provider_id | path |  | Yes | string (uint64) |
| resource_name | query |  | No | string |
| start_time | query |  | No | string (uint64) |
| end_time | query |  | No | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [ListAlertsResponse](#listalertsresponse) |
| default | An unexpected error response. | [Status](#status) |

#### POST
##### Summary

create alerts

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| provider_type | path |  | Yes | string |
| provider_id | path |  | Yes | string (uint64) |
| body | body |  | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [CreateAlertsResponse](#createalertsresponse) |
| default | An unexpected error response. | [Status](#status) |

### /v1beta1/namespaces

#### GET
##### Summary

list namespaces

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [ListNamespacesResponse](#listnamespacesresponse) |
| default | An unexpected error response. | [Status](#status) |

#### POST
##### Summary

create a namespace

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [CreateNamespaceRequest](#createnamespacerequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [CreateNamespaceResponse](#createnamespaceresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [GetNamespaceResponse](#getnamespaceresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [DeleteNamespaceResponse](#deletenamespaceresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [UpdateNamespaceResponse](#updatenamespaceresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [ListProvidersResponse](#listprovidersresponse) |
| default | An unexpected error response. | [Status](#status) |

#### POST
##### Summary

create a provider

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [CreateProviderRequest](#createproviderrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [CreateProviderResponse](#createproviderresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [GetProviderResponse](#getproviderresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [DeleteProviderResponse](#deleteproviderresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [UpdateProviderResponse](#updateproviderresponse) |
| default | An unexpected error response. | [Status](#status) |

### /v1beta1/receivers

#### GET
##### Summary

list receivers

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [ListReceiversResponse](#listreceiversresponse) |
| default | An unexpected error response. | [Status](#status) |

#### POST
##### Summary

create a receiver

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [CreateReceiverRequest](#createreceiverrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [CreateReceiverResponse](#createreceiverresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [GetReceiverResponse](#getreceiverresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [DeleteReceiverResponse](#deletereceiverresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [UpdateReceiverResponse](#updatereceiverresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [NotifyReceiverResponse](#notifyreceiverresponse) |
| default | An unexpected error response. | [Status](#status) |

### /v1beta1/rules

#### GET
##### Summary

list rules

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | query |  | No | string |
| namespace | query |  | No | string |
| group_name | query |  | No | string |
| template | query |  | No | string |
| provider_namespace | query |  | No | string (uint64) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [ListRulesResponse](#listrulesresponse) |
| default | An unexpected error response. | [Status](#status) |

#### PUT
##### Summary

add/update a rule

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [UpdateRuleRequest](#updaterulerequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [UpdateRuleResponse](#updateruleresponse) |
| default | An unexpected error response. | [Status](#status) |

### /v1beta1/subscriptions

#### GET
##### Summary

List subscriptions

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [ListSubscriptionsResponse](#listsubscriptionsresponse) |
| default | An unexpected error response. | [Status](#status) |

#### POST
##### Summary

Create a subscription

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [CreateSubscriptionRequest](#createsubscriptionrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [CreateSubscriptionResponse](#createsubscriptionresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [GetSubscriptionResponse](#getsubscriptionresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [DeleteSubscriptionResponse](#deletesubscriptionresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [UpdateSubscriptionResponse](#updatesubscriptionresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [ListTemplatesResponse](#listtemplatesresponse) |
| default | An unexpected error response. | [Status](#status) |

#### PUT
##### Summary

add/update a template

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [UpsertTemplateRequest](#upserttemplaterequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [UpsertTemplateResponse](#upserttemplateresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [GetTemplateResponse](#gettemplateresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [DeleteTemplateResponse](#deletetemplateresponse) |
| default | An unexpected error response. | [Status](#status) |

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
| 200 | A successful response. | [RenderTemplateResponse](#rendertemplateresponse) |
| default | An unexpected error response. | [Status](#status) |

### Models

#### Alert

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |
| metric_name | string |  | No |
| metric_value | string |  | No |
| provider_id | string (uint64) |  | No |
| resource_name | string |  | No |
| rule | string |  | No |
| severity | string |  | No |
| triggered_at | dateTime |  | No |

#### Any

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| @type | string |  | No |

#### CreateAlertsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| alerts | [ [Alert](#alert) ] |  | No |

#### CreateNamespaceRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created_at | dateTime |  | No |
| credentials | object |  | No |
| labels | object |  | No |
| name | string |  | No |
| provider | string (uint64) |  | No |
| updated_at | dateTime |  | No |
| urn | string |  | No |

#### CreateNamespaceResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |

#### CreateProviderRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| credentials | object |  | No |
| host | string |  | No |
| labels | object |  | No |
| name | string |  | No |
| type | string |  | No |
| urn | string |  | No |

#### CreateProviderResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |

#### CreateReceiverRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| configurations | object |  | No |
| labels | object |  | No |
| name | string |  | No |
| type | string |  | No |

#### CreateReceiverResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |

#### CreateSubscriptionRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| match | object |  | No |
| namespace | string (uint64) |  | No |
| receivers | [ [ReceiverMetadata](#receivermetadata) ] |  | No |
| urn | string |  | No |

#### CreateSubscriptionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |

#### DeleteNamespaceResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| DeleteNamespaceResponse | object |  |  |

#### DeleteProviderResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| DeleteProviderResponse | object |  |  |

#### DeleteReceiverResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| DeleteReceiverResponse | object |  |  |

#### DeleteSubscriptionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| DeleteSubscriptionResponse | object |  |  |

#### DeleteTemplateResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| DeleteTemplateResponse | object |  |  |

#### GetNamespaceResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| namespace | [Namespace](#namespace) |  | No |

#### GetProviderResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| provider | [Provider](#provider) |  | No |

#### GetReceiverResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| receiver | [Receiver](#receiver) |  | No |

#### GetSubscriptionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| subscription | [Subscription](#subscription) |  | No |

#### GetTemplateResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| template | [Template](#template) |  | No |

#### ListAlertsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| alerts | [ [Alert](#alert) ] |  | No |

#### ListNamespacesResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| namespaces | [ [Namespace](#namespace) ] |  | No |

#### ListProvidersResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| providers | [ [Provider](#provider) ] |  | No |

#### ListReceiversResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| receivers | [ [Receiver](#receiver) ] |  | No |

#### ListRulesResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| rules | [ [Rule](#rule) ] |  | No |

#### ListSubscriptionsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| subscriptions | [ [Subscription](#subscription) ] |  | No |

#### ListTemplatesResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| templates | [ [Template](#template) ] |  | No |

#### Namespace

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created_at | dateTime |  | No |
| credentials | object |  | No |
| id | string (uint64) |  | No |
| labels | object |  | No |
| name | string |  | No |
| provider | string (uint64) |  | No |
| updated_at | dateTime |  | No |
| urn | string |  | No |

#### NotifyReceiverResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| NotifyReceiverResponse | object |  |  |

#### NullValue

`NullValue` is a singleton enumeration to represent the null value for the
`Value` type union.

 The JSON representation for `NullValue` is JSON `null`.

- NULL_VALUE: Null value.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| NullValue | string | `NullValue` is a singleton enumeration to represent the null value for the `Value` type union.   The JSON representation for `NullValue` is JSON `null`.   - NULL_VALUE: Null value. |  |

#### Provider

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created_at | dateTime |  | No |
| credentials | object |  | No |
| host | string |  | No |
| id | string (uint64) |  | No |
| labels | object |  | No |
| name | string |  | No |
| type | string |  | No |
| updated_at | dateTime |  | No |
| urn | string |  | No |

#### Receiver

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| configurations | object |  | No |
| created_at | dateTime |  | No |
| data | object |  | No |
| id | string (uint64) |  | No |
| labels | object |  | No |
| name | string |  | No |
| type | string |  | No |
| updated_at | dateTime |  | No |

#### ReceiverMetadata

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| configuration | object |  | No |
| id | string (uint64) |  | No |

#### RenderTemplateResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| body | string |  | No |

#### Rule

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created_at | dateTime |  | No |
| enabled | boolean |  | No |
| group_name | string |  | No |
| id | string (uint64) |  | No |
| name | string |  | No |
| namespace | string |  | No |
| provider_namespace | string (uint64) |  | No |
| template | string |  | No |
| updated_at | dateTime |  | No |
| variables | [ [Variables](#variables) ] |  | No |

#### Status

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |
| details | [ [Any](#any) ] |  | No |
| message | string |  | No |

#### Subscription

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created_at | dateTime |  | No |
| id | string (uint64) |  | No |
| match | object |  | No |
| namespace | string (uint64) |  | No |
| receivers | [ [ReceiverMetadata](#receivermetadata) ] |  | No |
| updated_at | dateTime |  | No |
| urn | string |  | No |

#### Template

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| body | string |  | No |
| created_at | dateTime |  | No |
| id | string (uint64) |  | No |
| name | string |  | No |
| tags | [ string ] |  | No |
| updated_at | dateTime |  | No |
| variables | [ [TemplateVariables](#templatevariables) ] |  | No |

#### TemplateVariables

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| default | string |  | No |
| description | string |  | No |
| name | string |  | No |
| type | string |  | No |

#### UpdateNamespaceResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |

#### UpdateProviderResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |

#### UpdateReceiverResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |

#### UpdateRuleRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| enabled | boolean |  | No |
| group_name | string |  | No |
| namespace | string |  | No |
| provider_namespace | string (uint64) |  | No |
| template | string |  | No |
| variables | [ [Variables](#variables) ] |  | No |

#### UpdateRuleResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |

#### UpdateSubscriptionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |

#### UpsertTemplateRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| body | string |  | No |
| id | string (uint64) |  | No |
| name | string |  | No |
| tags | [ string ] |  | No |
| variables | [ [TemplateVariables](#templatevariables) ] |  | No |

#### UpsertTemplateResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string (uint64) |  | No |

#### Variables

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | No |
| name | string |  | No |
| type | string |  | No |
| value | string |  | No |
