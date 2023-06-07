# \RateLimitingApi

All URIs are relative to *http://localhost:4578*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RateLimitByAccountGet**](RateLimitingApi.md#RateLimitByAccountGet) | **Get** /rate-limit/by-account | An API with an aggressive rate limit by account
[**RateLimitByIpGet**](RateLimitingApi.md#RateLimitByIpGet) | **Get** /rate-limit/by-ip | An API with an aggressive rate limit by IP
[**RateLimitExponentialBackoffGet**](RateLimitingApi.md#RateLimitExponentialBackoffGet) | **Get** /rate-limit/exponential-backoff | An API with an aggressive rate limit with exponential backoff.



## RateLimitByAccountGet

> SuccessfulResponse RateLimitByAccountGet(ctx).AccountId(accountId).Execute()

An API with an aggressive rate limit by account

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
    accountId := "accountId_example" // string | Account ID (email)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RateLimitingApi.RateLimitByAccountGet(context.Background()).AccountId(accountId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RateLimitingApi.RateLimitByAccountGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RateLimitByAccountGet`: SuccessfulResponse
    fmt.Fprintf(os.Stdout, "Response from `RateLimitingApi.RateLimitByAccountGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRateLimitByAccountGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountId** | **string** | Account ID (email) |

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


## RateLimitByIpGet

> SuccessfulResponse RateLimitByIpGet(ctx).Execute()

An API with an aggressive rate limit by IP

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

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RateLimitingApi.RateLimitByIpGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RateLimitingApi.RateLimitByIpGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RateLimitByIpGet`: SuccessfulResponse
    fmt.Fprintf(os.Stdout, "Response from `RateLimitingApi.RateLimitByIpGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiRateLimitByIpGetRequest struct via the builder pattern


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


## RateLimitExponentialBackoffGet

> SuccessfulResponse RateLimitExponentialBackoffGet(ctx).Execute()

An API with an aggressive rate limit with exponential backoff.

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

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RateLimitingApi.RateLimitExponentialBackoffGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RateLimitingApi.RateLimitExponentialBackoffGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RateLimitExponentialBackoffGet`: SuccessfulResponse
    fmt.Fprintf(os.Stdout, "Response from `RateLimitingApi.RateLimitExponentialBackoffGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiRateLimitExponentialBackoffGetRequest struct via the builder pattern


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
