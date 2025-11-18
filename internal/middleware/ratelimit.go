package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	requests map[string]*clientLimiter
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
	cleanup  *time.Ticker
	logger   *logrus.Logger
}

type clientLimiter struct {
	tokens     int
	lastUpdate time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed
// window: time window for the rate limit
func NewRateLimiter(rate int, window time.Duration, logger *logrus.Logger) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*clientLimiter),
		rate:     rate,
		window:   window,
		logger:   logger,
	}

	// Cleanup old entries every minute
	rl.cleanup = time.NewTicker(1 * time.Minute)
	go rl.cleanupOldEntries()

	return rl
}

// Limit wraps an HTTP handler with rate limiting
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		if !rl.allow(clientIP) {
			rl.logger.WithFields(logrus.Fields{
				"client_ip": clientIP,
				"path":      r.URL.Path,
			}).Warn("Rate limit exceeded")

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", rl.window.String())
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"Rate limit exceeded. Please try again later."}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) allow(clientIP string) bool {
	rl.mu.Lock()
	limiter, exists := rl.requests[clientIP]
	if !exists {
		limiter = &clientLimiter{
			tokens:     rl.rate,
			lastUpdate: time.Now(),
		}
		rl.requests[clientIP] = limiter
	}
	rl.mu.Unlock()

	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(limiter.lastUpdate)
	tokensToAdd := int(elapsed / (rl.window / time.Duration(rl.rate)))

	if tokensToAdd > 0 {
		limiter.tokens = min(limiter.tokens+tokensToAdd, rl.rate)
		limiter.lastUpdate = now
	}

	if limiter.tokens > 0 {
		limiter.tokens--
		return true
	}

	return false
}

// cleanupOldEntries removes old entries that haven't been used recently
func (rl *RateLimiter) cleanupOldEntries() {
	for range rl.cleanup.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, limiter := range rl.requests {
			limiter.mu.Lock()
			// Remove entries that haven't been used in the last 10 minutes
			if now.Sub(limiter.lastUpdate) > 10*time.Minute {
				delete(rl.requests, ip)
			}
			limiter.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// Stop stops the rate limiter cleanup goroutine
func (rl *RateLimiter) Stop() {
	if rl.cleanup != nil {
		rl.cleanup.Stop()
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if ip != "" {
		return ip
	}
	return "unknown"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
