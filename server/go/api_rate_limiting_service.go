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
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/thecoreman/problematic-api-server/logic"
)

// RateLimitingApiService is a service that implements the logic for the RateLimitingApiServicer
// This service should implement the business logic for every endpoint for the RateLimitingApi API.
// Include any external packages or services that will be required by this service.
type RateLimitingApiService struct {
	logger zerolog.Logger
}

// NewRateLimitingApiService creates a default api service
func NewRateLimitingApiService(logger zerolog.Logger) RateLimitingApiServicer {
	logger = logger.With().Str("component", "rate-limiting-api").Logger()
	return &RateLimitingApiService{
		logger: logger,
	}
}

// RateLimitByAccountGet - An API with an aggressive rate limit by account
func (s *RateLimitingApiService) RateLimitByAccountGet(ctx context.Context, accountId string) (ImplResponse, error) {
	// TODO - update RateLimitByAccountGet with the required logic for this service method.
	// Add api_rate_limiting_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, SuccessfulResponse{}) or use other options such as http.Ok ...
	// return Response(200, SuccessfulResponse{}), nil

	// TODO: Uncomment the next line to return response Response(429, {}) or use other options such as http.Ok ...
	// return Response(429, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("RateLimitByAccountGet method not implemented")
}

// RateLimitByIpGet - An API with an aggressive rate limit by IP
func (s *RateLimitingApiService) RateLimitByIpGet(ctx context.Context) (ImplResponse, error) {
	s.logger.Info().Msg("Got into RateLimitByIpGet")
	line, err := logic.ReadRandomLineFromFile("", -1)
	if err != nil {
		return Response(http.StatusInternalServerError, err.Error()), nil
	}

	return Response(200, SuccessfulResponse{
		BookName:   "Random",
		LineNumber: 0,
		Text:       line,
	}), nil
}

// RateLimitExponentialBackoffGet - An API with an aggressive rate limit with exponential backoff.
func (s *RateLimitingApiService) RateLimitExponentialBackoffGet(ctx context.Context) (ImplResponse, error) {
	// TODO - update RateLimitExponentialBackoffGet with the required logic for this service method.
	// Add api_rate_limiting_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, SuccessfulResponse{}) or use other options such as http.Ok ...
	// return Response(200, SuccessfulResponse{}), nil

	// TODO: Uncomment the next line to return response Response(429, CooldownResponse{}) or use other options such as http.Ok ...
	// return Response(429, CooldownResponse{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("RateLimitExponentialBackoffGet method not implemented")
}
