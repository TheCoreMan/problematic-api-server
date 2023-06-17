//nolint:cyclop // Testing rate limiting seems complex but once you understand one case it's fine.
package ratelimiter_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/thecoreman/problematic-api-server/ratelimiter"
)

//nolint:gocognit,cyclop // Testing rate limiting seems complex but once you understand one case it's fine.
func TestShouldRateLimit(t *testing.T) {
	tests := []struct {
		name   string
		config ratelimiter.Config
		// This function returns two channels:
		// The first channel represents incoming http requests
		// The second indicates when the first one is done.
		requestsPublisher func() (<-chan *http.Request, <-chan bool)
		// These "want" slices represent what the response of the
		// rate limit function should be for each of the consumed
		// requests in the channel from the publisher IN ORDER.
		wantShouldRateLimit []bool
		wantReasonContains  []string
		wantErrors          []bool
	}{
		{
			name: "Requests to path that doesn't get rate limited, shouldn't limit",
			config: ratelimiter.Config{
				IPRateLimitWindow: time.Second,
			},
			requestsPublisher: func() (<-chan *http.Request, <-chan bool) {
				ch := make(chan *http.Request, 1)
				done := make(chan bool)
				req := httptest.NewRequest(http.MethodGet, "https://example.com/foo", nil)
				go func() {
					ch <- req
					time.Sleep(5 * time.Millisecond)
					ch <- req
					time.Sleep(5 * time.Millisecond)
					ch <- req
					time.Sleep(5 * time.Millisecond)
					ch <- req
					time.Sleep(5 * time.Millisecond)
					done <- true
				}()
				return ch, done
			},
			wantShouldRateLimit: []bool{false, false, false, false},
			wantReasonContains:  []string{"", "", "", ""},
			wantErrors:          []bool{false, false, false, false},
		},
		{
			name: "Rate Limit by IP: Two requests, same IP, shouldn't limit",
			config: ratelimiter.Config{
				IPRateLimitWindow: 500 * time.Millisecond,
			},
			requestsPublisher: func() (<-chan *http.Request, <-chan bool) {
				ch := make(chan *http.Request, 2)
				done := make(chan bool)
				req := httptest.NewRequest(http.MethodGet, "https://example.com/rate-limit/by-ip", nil)
				go func() {
					ch <- req
					time.Sleep(1 * time.Second)
					ch <- req
					time.Sleep(50 * time.Millisecond)
					done <- true
				}()
				return ch, done
			},
			wantShouldRateLimit: []bool{false, false},
			wantReasonContains:  []string{"", ""},
			wantErrors:          []bool{false, false},
		},
		{
			name: "Rate Limit by IP: Two requests, same IP, should limit",
			config: ratelimiter.Config{
				IPRateLimitWindow: time.Second,
			},
			requestsPublisher: func() (<-chan *http.Request, <-chan bool) {
				ch := make(chan *http.Request, 2)
				done := make(chan bool)
				req := httptest.NewRequest(http.MethodGet, "https://example.com/rate-limit/by-ip", nil)
				go func() {
					ch <- req
					time.Sleep(50 * time.Millisecond)
					ch <- req
					time.Sleep(50 * time.Millisecond)
					done <- true
				}()
				return ch, done
			},
			wantShouldRateLimit: []bool{false, true},
			wantReasonContains:  []string{"", "Rate limit exceeded for IP 192.0.2.1"},
			wantErrors:          []bool{false, false},
		},
		{
			name: "Rate Limit by IP: Three requests, different IPs, should limit only one",
			config: ratelimiter.Config{
				IPRateLimitWindow: time.Second,
			},
			requestsPublisher: func() (<-chan *http.Request, <-chan bool) {
				ch := make(chan *http.Request, 2)
				done := make(chan bool)
				req := httptest.NewRequest(http.MethodGet, "https://example.com/rate-limit/by-ip", nil)
				req2 := httptest.NewRequest(http.MethodGet, "https://example.com/rate-limit/by-ip", nil)
				// The default remote address is 192.0.2.1:1234
				// See https://cs.opensource.google/go/go/+/refs/tags/go1.20.5:src/net/http/httptest/httptest.go;l=75
				req2.RemoteAddr = "1.1.1.1:1234"
				go func() {
					ch <- req2
					time.Sleep(50 * time.Millisecond)
					ch <- req
					time.Sleep(50 * time.Millisecond)
					ch <- req2
					time.Sleep(50 * time.Millisecond)
					done <- true
				}()
				return ch, done
			},
			wantShouldRateLimit: []bool{false, false, true},
			wantReasonContains:  []string{"", "", "Rate limit exceeded for IP 1.1.1.1"},
			wantErrors:          []bool{false, false, false},
		},
		{
			name: "Rate Limit by IP: Invalid IP, should limit",
			config: ratelimiter.Config{
				IPRateLimitWindow: time.Second,
			},
			requestsPublisher: func() (<-chan *http.Request, <-chan bool) {
				ch := make(chan *http.Request, 2)
				done := make(chan bool)
				req := httptest.NewRequest(http.MethodGet, "https://example.com/rate-limit/by-ip", nil)
				req.RemoteAddr = "This is not an IP!"
				go func() {
					ch <- req
					done <- true
				}()
				return ch, done
			},
			wantShouldRateLimit: []bool{false},
			wantReasonContains:  []string{""},
			wantErrors:          []bool{true},
		},
		{
			name: "Rate limit by account, two requests, same account, should limit",
			config: ratelimiter.Config{
				IPRateLimitWindow:      time.Second,
				AccountRateLimitWindow: time.Minute,
			},
			requestsPublisher: func() (<-chan *http.Request, <-chan bool) {
				ch := make(chan *http.Request, 1)
				done := make(chan bool)
				req := httptest.NewRequest(http.MethodGet, "https://example.com/rate-limit/by-account", nil)
				req.Header.Set("X-Account-Id", "hello@shaynehmad.com")
				go func() {
					ch <- req
					time.Sleep(50 * time.Millisecond)
					ch <- req
					time.Sleep(50 * time.Millisecond)
					done <- true
				}()
				return ch, done
			},
			wantShouldRateLimit: []bool{false, true},
			wantReasonContains:  []string{"", "Rate limit exceeded for account hello@shaynehmad.com"},
		},
		{
			name: "Rate limit by account, invalid account, should err",
			config: ratelimiter.Config{
				IPRateLimitWindow:      time.Second,
				AccountRateLimitWindow: time.Minute,
			},
			requestsPublisher: func() (<-chan *http.Request, <-chan bool) {
				ch := make(chan *http.Request, 1)
				done := make(chan bool)
				req := httptest.NewRequest(http.MethodGet, "https://example.com/rate-limit/by-account", nil)
				req.Header.Set("X-Account-Id", "this is not an email")
				req2 := httptest.NewRequest(http.MethodGet, "https://example.com/rate-limit/by-account", nil)
				// req2 - Not adding the header
				go func() {
					ch <- req
					time.Sleep(5 * time.Millisecond)
					ch <- req2
					time.Sleep(5 * time.Millisecond)
					done <- true
				}()
				return ch, done
			},
			wantShouldRateLimit: []bool{false, false},
			wantReasonContains:  []string{"", ""},
			wantErrors:          []bool{true, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate the test data
			if len(tt.wantShouldRateLimit) != len(tt.wantReasonContains) {
				t.Fatalf("Test data is invalid: wantShouldRateLimit and wantReasonContains must be the same length")
			}

			rateLimiter := ratelimiter.NewRateLimiter(zerolog.Nop(), tt.config)
			requests, done := tt.requestsPublisher()
			rateLimitResponses := []bool{}
			rateLimitReasons := []string{}
			rateLimitErrors := []error{}
			for {
				select {
				case req := <-requests:
					t.Log("Received request")
					shouldRateLimit, reason, err := rateLimiter.ShouldRateLimit(req)
					rateLimitResponses = append(rateLimitResponses, shouldRateLimit)
					rateLimitReasons = append(rateLimitReasons, reason)
					rateLimitErrors = append(rateLimitErrors, err)
				case <-done:
					if len(rateLimitResponses) != len(tt.wantShouldRateLimit) {
						t.Fatalf("Got %d responses, want %d", len(rateLimitResponses), len(tt.wantShouldRateLimit))
					}

					for i, wantRateLimitDecision := range tt.wantShouldRateLimit {
						if rateLimitResponses[i] != wantRateLimitDecision {
							t.Errorf("Response %d: got %v, want %v", i, rateLimitResponses[i], wantRateLimitDecision)
						}
					}

					for i, wantReasonToContain := range tt.wantReasonContains {
						if !strings.Contains(rateLimitReasons[i], wantReasonToContain) {
							t.Errorf("Response %d: got %v, want to contain %v", i, rateLimitReasons[i], wantReasonToContain)
						}
					}

					for i, wantedError := range tt.wantErrors {
						gotAnError := rateLimitErrors[i] != nil
						if gotAnError != wantedError {
							t.Errorf("Response %d: got error %v, wantErr %v", i, rateLimitErrors[i], wantedError)
						}
					}

					return
				}
			}
		})
	}
}
