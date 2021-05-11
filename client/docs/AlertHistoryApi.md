# {{classname}}

All URIs are relative to *http://localhost:3000/*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateAlertHistoryRequest**](AlertHistoryApi.md#CreateAlertHistoryRequest) | **Post** /history | 
[**GetAlertHistoryRequest**](AlertHistoryApi.md#GetAlertHistoryRequest) | **Get** /history | 

# **CreateAlertHistoryRequest**
> CreateAlertHistoryRequest(ctx, optional)


Create Alert History API: This API create alert history

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AlertHistoryApiCreateAlertHistoryRequestOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AlertHistoryApiCreateAlertHistoryRequestOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of []Alerts**](Alerts.md)|  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAlertHistoryRequest**
> []AlertHistoryObject GetAlertHistoryRequest(ctx, optional)


GET Alert History API: This API lists stored alert history for given filers in query params

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AlertHistoryApiGetAlertHistoryRequestOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AlertHistoryApiGetAlertHistoryRequestOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **resource** | **optional.String**|  | 
 **startTime** | **optional.Int32**|  | 
 **endTime** | **optional.Int32**|  | 

### Return type

[**[]AlertHistoryObject**](AlertHistoryObject.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

