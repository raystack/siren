# {{classname}}

All URIs are relative to *http://localhost:3000/*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateAlertCredentialRequest**](AlertcredentialApi.md#CreateAlertCredentialRequest) | **Put** /alertingCredentials/teams/{teamName} | 
[**GetAlertCredentialRequest**](AlertcredentialApi.md#GetAlertCredentialRequest) | **Get** /alertingCredentials/teams/{teamName} | 

# **CreateAlertCredentialRequest**
> CreateAlertCredentialRequest(ctx, teamName, optional)


Upsert AlertCredentials API: This API helps in creating or updating the teams slack and pagerduty credentials

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamName** | **string**| name of the team | 
 **optional** | ***AlertcredentialApiCreateAlertCredentialRequestOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AlertcredentialApiCreateAlertCredentialRequestOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of AlertCredentialResponse**](AlertCredentialResponse.md)| Create AlertCredential request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAlertCredentialRequest**
> AlertCredentialResponse GetAlertCredentialRequest(ctx, teamName)


Get AlertCredentials API: This API helps in getting the teams slack and pagerduty credentials

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamName** | **string**| name of the team | 

### Return type

[**AlertCredentialResponse**](AlertCredentialResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

