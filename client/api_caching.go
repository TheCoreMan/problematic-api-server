/*
Rate Limit API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

// CachingApiService CachingApi service
type CachingApiService service

type ApiCacheableGetRequest struct {
	ctx         context.Context
	ApiService  *CachingApiService
	bookTitle   *string
	lineNumber  *int32
	withControl *bool
}

// Title of the book
func (r ApiCacheableGetRequest) BookTitle(bookTitle string) ApiCacheableGetRequest {
	r.bookTitle = &bookTitle
	return r
}

// Line number
func (r ApiCacheableGetRequest) LineNumber(lineNumber int32) ApiCacheableGetRequest {
	r.lineNumber = &lineNumber
	return r
}

// Include the Cache-Control header in the response?
func (r ApiCacheableGetRequest) WithControl(withControl bool) ApiCacheableGetRequest {
	r.withControl = &withControl
	return r
}

func (r ApiCacheableGetRequest) Execute() (*SuccessfulResponse, *http.Response, error) {
	return r.ApiService.CacheableGetExecute(r)
}

/*
CacheableGet Get a cacheable response.

Returns a cacheable response for the given book title and line number.
If `with-control` is set to true, includes Cache-Control header for
controlling caching behavior, such as max-age.

Generally, the response should be cached since it's always
the same for the same parameters.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiCacheableGetRequest
*/
func (a *CachingApiService) CacheableGet(ctx context.Context) ApiCacheableGetRequest {
	return ApiCacheableGetRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return SuccessfulResponse
func (a *CachingApiService) CacheableGetExecute(r ApiCacheableGetRequest) (*SuccessfulResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SuccessfulResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "CachingApiService.CacheableGet")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/cacheable"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.bookTitle == nil {
		return localVarReturnValue, nil, reportError("bookTitle is required and must be specified")
	}
	if r.lineNumber == nil {
		return localVarReturnValue, nil, reportError("lineNumber is required and must be specified")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "book-title", r.bookTitle, "")
	parameterAddToHeaderOrQuery(localVarQueryParams, "line-number", r.lineNumber, "")
	if r.withControl != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "with-control", r.withControl, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
