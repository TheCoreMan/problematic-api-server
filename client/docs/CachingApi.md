# \CachingApi

All URIs are relative to *http://localhost:4578*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CacheableGet**](CachingApi.md#CacheableGet) | **Get** /cacheable | Get a cacheable response.



## CacheableGet

> SuccessfulResponse CacheableGet(ctx).BookTitle(bookTitle).LineNumber(lineNumber).WithControl(withControl).Execute()

Get a cacheable response.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    bookTitle := "bookTitle_example" // string | Title of the book
    lineNumber := int32(56) // int32 | Line number
    withControl := true // bool | Include the Cache-Control header in the response? (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.CachingApi.CacheableGet(context.Background()).BookTitle(bookTitle).LineNumber(lineNumber).WithControl(withControl).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `CachingApi.CacheableGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CacheableGet`: SuccessfulResponse
    fmt.Fprintf(os.Stdout, "Response from `CachingApi.CacheableGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCacheableGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **bookTitle** | **string** | Title of the book |
 **lineNumber** | **int32** | Line number |
 **withControl** | **bool** | Include the Cache-Control header in the response? |

### Return type

[**SuccessfulResponse**](SuccessfulResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)
