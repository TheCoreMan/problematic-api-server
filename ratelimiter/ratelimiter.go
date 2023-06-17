package ratelimiter

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type RateLimiter interface {
	Middleware(next http.Handler) http.Handler
	// This is an internal func, exporting just to make testing easy.
	ShouldRateLimit(r *http.Request) (bool, string)
}

type Config struct {
	IPRateLimitWindow time.Duration
}

type Impl struct {
	logger            zerolog.Logger
	ipRateLimitWindow time.Duration
	ipRateLimitMap    sync.Map
}

func NewRateLimiter(
	logger zerolog.Logger,
	config Config,
) RateLimiter {
	logger = logger.With().Str("component", "rate-limiter").Logger()
	return &Impl{
		logger:            logger,
		ipRateLimitWindow: config.IPRateLimitWindow,
		ipRateLimitMap:    sync.Map{},
	}
}

func (rl *Impl) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldRateLimit, reason := rl.ShouldRateLimit(r)
		if shouldRateLimit {
			w.WriteHeader(http.StatusTooManyRequests)
			_, err := w.Write([]byte("Rate limit exceeded. reason:" + reason))
			if err != nil {
				rl.logger.Error().Err(err).Msg("Error writing response")
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

// shouldRateLimit returns true if the request should be rate limited, and the
// reason why. It returns false if the request should not be rate limited.
//
// You should generally use it like so:
//
//	shouldRateLimit, reason := rateLimiter.shouldRateLimit(r)
//	if (shouldRateLimit) {
//		w.WriteHeader(http.StatusTooManyRequests)
//		w.Write([]byte("Rate limit exceeded. reason:" + reason))
//		return
//	}
//	next.ServeHTTP(w, r)
//
// The rate limiting logic is different per API since this server is for
// educational purposes.
func (rl *Impl) ShouldRateLimit(r *http.Request) (bool, string) {
	switch r.URL.Path {
	case "/rate-limit/by-ip":
		return rl.shouldRateLimitByIP(r)
	case "/rate-limit/by-account":
		return rl.shouldRateLimitByAccount(r)
	case "/rate-limit/exponential-backoff":
		return rl.shouldRateLimitExponentialBackoff(r)
	default:
		return false, ""
	}
}

func (rl *Impl) shouldRateLimitByIP(r *http.Request) (bool, string) {
	ip := strings.Split(r.RemoteAddr, ":")[0]

	// Get the last request time for this IP from the map, or the zero value if
	// it doesn't exist.
	now := time.Now()
	lastRequestTime, ok := rl.ipRateLimitMap.Load(ip)
	if !ok {
		rl.ipRateLimitMap.Store(ip, now)
	} else {
		// Check if the last request was within the window
		windowEnd := lastRequestTime.(time.Time).Add(rl.ipRateLimitWindow)
		if windowEnd.After(now) {
			return true, "Rate limit exceeded for IP " + ip + " by " + windowEnd.Sub(now).String()
		}
	}

	return false, ""
}

func (rl *Impl) shouldRateLimitByAccount(_ *http.Request) (bool, string) {
	return false, ""
}

func (rl *Impl) shouldRateLimitExponentialBackoff(_ *http.Request) (bool, string) {
	return false, ""
}
