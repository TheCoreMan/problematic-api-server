/*
 * Rate Limit API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package server

import (
	"net/http"
	"strings"
)

// ErrorsApiController binds http requests to an api service and writes the service results to the http response
type ErrorsApiController struct {
	service      ErrorsApiServicer
	errorHandler ErrorHandler
}

// ErrorsApiOption for how the controller is set up.
type ErrorsApiOption func(*ErrorsApiController)

// WithErrorsApiErrorHandler inject ErrorHandler into controller
func WithErrorsApiErrorHandler(h ErrorHandler) ErrorsApiOption {
	return func(c *ErrorsApiController) {
		c.errorHandler = h
	}
}

// NewErrorsApiController creates a default api controller
func NewErrorsApiController(s ErrorsApiServicer, opts ...ErrorsApiOption) Router {
	controller := &ErrorsApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the ErrorsApiController
func (c *ErrorsApiController) Routes() Routes {
	return Routes{
		{
			"ErrorsPercentGet",
			strings.ToUpper("Get"),
			"/errors/percent",
			c.ErrorsPercentGet,
		},
	}
}

// ErrorsPercentGet - An API that will return an error \"error_percent\" percent of the time
func (c *ErrorsApiController) ErrorsPercentGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	errorPercentParam, err := parseInt32Parameter(query.Get("error_percent"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	result, err := c.service.ErrorsPercentGet(r.Context(), errorPercentParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)
}
