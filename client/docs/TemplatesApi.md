# {{classname}}

All URIs are relative to *http://localhost:3000/*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateTemplateRequest**](TemplatesApi.md#CreateTemplateRequest) | **Put** /templates | 
[**DeleteTemplatesRequest**](TemplatesApi.md#DeleteTemplatesRequest) | **Delete** /templates/{name} | 
[**GetTemplatesRequest**](TemplatesApi.md#GetTemplatesRequest) | **Get** /templates/{name} | 
[**ListTemplatesRequest**](TemplatesApi.md#ListTemplatesRequest) | **Get** /templates | 
[**RenderTemplatesRequest**](TemplatesApi.md#RenderTemplatesRequest) | **Post** /templates/{name}/render | 

# **CreateTemplateRequest**
> Template CreateTemplateRequest(ctx, optional)


Upsert Templates API: This API helps in creating or updating a template with unique name

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***TemplatesApiCreateTemplateRequestOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TemplatesApiCreateTemplateRequestOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of Template**](Template.md)| Create template request | 

### Return type

[**Template**](Template.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteTemplatesRequest**
> Template DeleteTemplatesRequest(ctx, name)


Delete Template API: This API deletes a template given the template name

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Delete Template Request | 

### Return type

[**Template**](Template.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetTemplatesRequest**
> Template GetTemplatesRequest(ctx, name)


Get Template API: This API gets a template given the template name

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Get Template Request | 

### Return type

[**Template**](Template.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListTemplatesRequest**
> []Template ListTemplatesRequest(ctx, optional)


List Templates API: This API lists all the existing templates with given filers in query params

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***TemplatesApiListTemplatesRequestOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TemplatesApiListTemplatesRequestOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tag** | **optional.String**| List Template Request | 

### Return type

[**[]Template**](Template.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RenderTemplatesRequest**
> RenderTemplatesRequest(ctx, name)


Render Template API: This API renders the given template with given values

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Render Template Request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

