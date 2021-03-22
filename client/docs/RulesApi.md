# {{classname}}

All URIs are relative to *http://localhost:3000/*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateRuleRequest**](RulesApi.md#CreateRuleRequest) | **Put** /rules | 
[**ListRulesRequest**](RulesApi.md#ListRulesRequest) | **Get** /rules | 

# **CreateRuleRequest**
> Rule CreateRuleRequest(ctx, optional)


Upsert Rule API: This API helps in creating a new rule or update an existing one with unique combination of namespace, entity, group_name, template

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***RulesApiCreateRuleRequestOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a RulesApiCreateRuleRequestOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of Rule**](Rule.md)| Create rule request | 

### Return type

[**Rule**](Rule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListRulesRequest**
> []Template ListRulesRequest(ctx, optional)


List Rules API: This API lists all the existing rules with given filers in query params

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***RulesApiListRulesRequestOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a RulesApiListRulesRequestOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **namespace** | **optional.String**| List Rule Request | 
 **entity** | **optional.String**|  | 
 **groupName** | **optional.String**|  | 
 **status** | **optional.String**|  | 
 **template** | **optional.String**|  | 

### Return type

[**[]Template**](Template.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

