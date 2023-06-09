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
	"context"
	"fmt"
	"math/rand"
	"net/http"
)

// ErrorsApiService is a service that implements the logic for the ErrorsApiServicer
// This service should implement the business logic for every endpoint for the ErrorsApi API.
// Include any external packages or services that will be required by this service.
type ErrorsApiService struct{}

// NewErrorsApiService creates a default api service
func NewErrorsApiService() ErrorsApiServicer {
	return &ErrorsApiService{}
}

// ErrorsPercentGet - An API that will return an error \&quot;error_percent\&quot; percent of the time
func (s *ErrorsApiService) ErrorsPercentGet(ctx context.Context, errorPercent int32) (ImplResponse, error) {
	errorFraction := float32(errorPercent) / 100.0
	if errorFraction < 0.0 || errorFraction > 1.0 {
		return Response(http.StatusBadRequest, "error_percent must be between 0 and 100"), nil
	}

	// Get a weak random number between 0 and 1
	//nolint:gosec // we're not using this for anything security related
	randomFraction := float32(0.0 + (1.0-0.0)*rand.Float64())
	if randomFraction < errorFraction {
		return Response(
			http.StatusInternalServerError,
			fmt.Sprintf("Rolled a %.2f which is less than %.2f", randomFraction, errorFraction),
		), nil
	}
	return Response(
		http.StatusOK,
		fmt.Sprintf("Rolled a %.2f which is greater than or equal to %.2f", randomFraction, errorFraction),
	), nil
}
