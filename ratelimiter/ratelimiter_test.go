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
		name                string
		config              ratelimiter.Config
		requestsPublisher   func() (<-chan *http.Request, <-chan bool)
		wantShouldRateLimit []bool
		wantReasonContains  []string
	}{
		{
			name: "Single request",
			config: ratelimiter.Config{
				IPRateLimitWindow: time.Second,
			},
			requestsPublisher: func() (<-chan *http.Request, <-chan bool) {
				ch := make(chan *http.Request, 1)
				done := make(chan bool)
				req := httptest.NewRequest(http.MethodGet, "https://example.com/foo", nil)
				go func() {
					ch <- req
					time.Sleep(50 * time.Millisecond)
					done <- true
				}()
				return ch, done
			},
			wantShouldRateLimit: []bool{false},
			wantReasonContains:  []string{""},
		},
		{
			name: "Two requests, same IP, IP rate limit, shouldn't limit",
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
		},
		{
			name: "Two requests, same IP, IP rate limit, should limit",
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
			wantReasonContains:  []string{"", "Rate limit exceeded for IP"},
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
			for {
				select {
				case req := <-requests:
					t.Log("Received request")
					shouldRateLimit, reason := rateLimiter.ShouldRateLimit(req)
					rateLimitResponses = append(rateLimitResponses, shouldRateLimit)
					rateLimitReasons = append(rateLimitReasons, reason)
				case <-done:
					if len(rateLimitResponses) != len(tt.wantShouldRateLimit) {
						t.Fatalf("Got %d responses, want %d", len(rateLimitResponses), len(tt.wantShouldRateLimit))
					}

					for i, want := range tt.wantShouldRateLimit {
						if rateLimitResponses[i] != want {
							t.Errorf("Response %d: got %v, want %v", i, rateLimitResponses[i], want)
						}
					}

					for i, want := range tt.wantReasonContains {
						if !strings.Contains(rateLimitReasons[i], want) {
							t.Errorf("Response %d: got %v, want to contain %v", i, rateLimitReasons[i], want)
						}
					}

					return
				}
			}
		})
	}
}
