package middleware

import (
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// RateLimiter is a struct to manage rate limiting for multiple clients
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// GetLimiter retrieves or creates a rate limiter for a specific client
func (rl *RateLimiter) GetLimiter(clientID string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limiter, exists := rl.limiters[clientID]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.limiters[clientID] = limiter
	log.WithFields(log.Fields{
		"clientID": clientID,
	}).Info("Created new rate limiter")
	return limiter
}

// RateLimitMiddleware applies rate limiting to HTTP requests
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.RemoteAddr // Use the client's IP address as the identifier
		limiter := rl.GetLimiter(clientID)

		if !limiter.Allow() {
			log.WithFields(log.Fields{
				"clientID": clientID,
			}).Warn("Rate limit exceeded")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		log.WithFields(log.Fields{
			"clientID": clientID,
		}).Info("Request allowed")
		next.ServeHTTP(w, r)
	})
}
