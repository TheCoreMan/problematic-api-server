# \ErrorsApi

All URIs are relative to *http://localhost:4578*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ErrorsPercentGet**](ErrorsApi.md#ErrorsPercentGet) | **Get** /errors/percent | An API that will return an error \&quot;error_percent\&quot; percent of the time



## ErrorsPercentGet

> ErrorsPercentGet(ctx).ErrorPercent(errorPercent).Execute()

An API that will return an error \"error_percent\" percent of the time

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
    errorPercent := int32(56) // int32 | Percentage of requests that result in an error (0-100)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ErrorsApi.ErrorsPercentGet(context.Background()).ErrorPercent(errorPercent).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ErrorsApi.ErrorsPercentGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiErrorsPercentGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **errorPercent** | **int32** | Percentage of requests that result in an error (0-100) |

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)
