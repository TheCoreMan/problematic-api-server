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

// ErrorsApiService ErrorsApi service
type ErrorsApiService service

type ApiErrorsPercentGetRequest struct {
	ctx          context.Context
	ApiService   *ErrorsApiService
	errorPercent *int32
}

// Percentage of requests that result in an error (0-100)
func (r ApiErrorsPercentGetRequest) ErrorPercent(errorPercent int32) ApiErrorsPercentGetRequest {
	r.errorPercent = &errorPercent
	return r
}

func (r ApiErrorsPercentGetRequest) Execute() (*http.Response, error) {
	return r.ApiService.ErrorsPercentGetExecute(r)
}

/*
ErrorsPercentGet An API that will return an error \"error_percent\" percent of the time

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiErrorsPercentGetRequest
*/
func (a *ErrorsApiService) ErrorsPercentGet(ctx context.Context) ApiErrorsPercentGetRequest {
	return ApiErrorsPercentGetRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
func (a *ErrorsApiService) ErrorsPercentGetExecute(r ApiErrorsPercentGetRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodGet
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ErrorsApiService.ErrorsPercentGet")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/errors/percent"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.errorPercent == nil {
		return nil, reportError("errorPercent is required and must be specified")
	}
	if *r.errorPercent < 0 {
		return nil, reportError("errorPercent must be greater than 0")
	}
	if *r.errorPercent > 100 {
		return nil, reportError("errorPercent must be less than 100")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "error_percent", r.errorPercent, "")
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}
