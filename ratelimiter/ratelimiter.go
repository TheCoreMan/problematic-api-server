package ratelimiter

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type RateLimiter interface {
	Middleware(next http.Handler) http.Handler
	// This is an internal func, exporting just to make testing easy.
	ShouldRateLimit(r *http.Request) (bool, string, error)
}

type Config struct {
	IPRateLimitWindow          time.Duration
	AccountRateLimitWindow     time.Duration
	BackoffRateLimitWindow     time.Duration
	BackoffRateLimitMultiplier int
}

// Using sync.Map for safe concurrent access.
type Impl struct {
	logger                     zerolog.Logger
	ipRateLimitWindow          time.Duration
	ipRateLimitMap             sync.Map
	accountRateLimitWindow     time.Duration
	accountRateLimitMap        sync.Map
	backoffRateLimitWindow     time.Duration
	backoffRateLimitMultiplier int
	backoffRateLimitMap        sync.Map
}

func NewRateLimiter(
	logger zerolog.Logger,
	config Config,
) RateLimiter {
	logger = logger.With().Str("component", "rate-limiter").Logger()
	return &Impl{
		logger:                     logger,
		ipRateLimitWindow:          config.IPRateLimitWindow,
		ipRateLimitMap:             sync.Map{},
		accountRateLimitWindow:     config.AccountRateLimitWindow,
		accountRateLimitMap:        sync.Map{},
		backoffRateLimitWindow:     config.BackoffRateLimitWindow,
		backoffRateLimitMultiplier: config.BackoffRateLimitMultiplier,
		backoffRateLimitMap:        sync.Map{},
	}
}

func (rl *Impl) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldRateLimit, reason, err := rl.ShouldRateLimit(r)
		if err != nil {
			rl.logger.Error().Err(err).Msg("Error in rate limiter")
			w.WriteHeader(http.StatusInternalServerError)
			_, writeErr := w.Write([]byte("Internal server error"))
			if writeErr != nil {
				rl.logger.Error().Err(writeErr).Msg("Error writing response")
			}
			return
		}
		if shouldRateLimit {
			w.WriteHeader(http.StatusTooManyRequests)
			_, writeErr := w.Write([]byte("Rate limit exceeded. reason:" + reason))
			if writeErr != nil {
				rl.logger.Error().Err(writeErr).Msg("Error writing response")
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
func (rl *Impl) ShouldRateLimit(r *http.Request) (bool, string, error) {
	switch r.URL.Path {
	case "/rate-limit/by-ip":
		if ip := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]); ip == nil {
			return false, "", errors.New("invalid IP address")
		}
		should, reason := rl.shouldRateLimitByIP(r)
		return should, reason, nil
	case "/rate-limit/by-account":
		accounts, ok := r.Header["X-Account-Id"]
		if !ok || len(accounts) != 1 {
			return false, "", errors.New("missing or Malformed X-Account-Id header")
		}

		account := accounts[0]
		_, emailParseErr := mail.ParseAddress(account)
		if emailParseErr != nil {
			return false, "",
				fmt.Errorf(
					"malformed account - must be valid email. Details: %w",
					emailParseErr,
				)
		}
		should, reason := rl.shouldRateLimitByAccount(r)
		return should, reason, nil
	case "/rate-limit/exponential-backoff":
		should, reason := rl.shouldRateLimitExponentialBackoff(r)
		return should, reason, nil
	default:
		return false, "", nil
	}
}

func (rl *Impl) shouldRateLimitByIP(r *http.Request) (bool, string) {
	// ip validated outside the function, safe to use here
	ip := strings.Split(r.RemoteAddr, ":")[0]

	now := time.Now()
	// Note: First load, then store.
	lastRequestTime, ok := rl.ipRateLimitMap.Load(ip)
	rl.ipRateLimitMap.Store(ip, now)
	if ok {
		// Check if the last request was within the window
		windowEnd := lastRequestTime.(time.Time).Add(rl.ipRateLimitWindow)
		if windowEnd.After(now) {
			return true, "Rate limit exceeded for IP " + ip + " by " + windowEnd.Sub(now).String()
		}
	}

	return false, ""
}

func (rl *Impl) shouldRateLimitByAccount(r *http.Request) (bool, string) {
	// Header is validated outside the function, safe to use here
	accounts := r.Header["X-Account-Id"]
	account := accounts[0]

	now := time.Now()
	lastRequestTime, ok := rl.accountRateLimitMap.Load(account)
	rl.accountRateLimitMap.Store(account, now)
	if ok {
		windowEnd := lastRequestTime.(time.Time).Add(rl.accountRateLimitWindow)
		if windowEnd.After(now) {
			return true, "Rate limit exceeded for account " + account + " by " + windowEnd.Sub(now).String()
		}
	}

	return false, ""
}

type exponentialRateLimitEntry struct {
	lastRequestTime time.Time
	violations      int
}

func (rl *Impl) shouldRateLimitExponentialBackoff(r *http.Request) (bool, string) {
	// ip validated outside the function, safe to use here
	ip := strings.Split(r.RemoteAddr, ":")[0]

	now := time.Now()
	// Note: First load, then store.
	existingEntry, ok := rl.backoffRateLimitMap.Load(ip)
	if !ok {
		// This is the first violation
		newEntry := exponentialRateLimitEntry{
			lastRequestTime: now,
			violations:      1,
		}
		rl.backoffRateLimitMap.Store(ip, newEntry)
	} else {
		lastRequestTime := existingEntry.(exponentialRateLimitEntry).lastRequestTime
		violations := existingEntry.(exponentialRateLimitEntry).violations
		punishmentCount := violations * rl.backoffRateLimitMultiplier
		punishmentTime := time.Duration(punishmentCount) * rl.backoffRateLimitWindow
		windowEnd := lastRequestTime.Add(punishmentTime)
		if windowEnd.After(now) {
			reasonMessage := fmt.Sprintf(
				"Rate limit exceeded for IP %s by %s. violated %d times",
				ip, windowEnd.Sub(now).String(), violations,
			)
			newEntry := exponentialRateLimitEntry{
				lastRequestTime: now,
				violations:      violations + 1,
			}
			rl.backoffRateLimitMap.Store(ip, newEntry)
			return true, reasonMessage
		}

		newEntry := exponentialRateLimitEntry{
			lastRequestTime: now,
			violations:      violations,
		}
		rl.backoffRateLimitMap.Store(ip, newEntry)
	}

	return false, ""
}
