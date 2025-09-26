package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tdmdh/fit-up-server/services/auth-service/internal/types"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/utils"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if now.Sub(reqTime) < rl.window {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[key] = validRequests
	}

	if len(rl.requests[key]) >= rl.limit {
		return false
	}

	rl.requests[key] = append(rl.requests[key], now)
	return true
}

func GetClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}

func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := GetClientIP(r)

			if !limiter.Allow(clientIP) {
				utils.WriteError(w, http.StatusTooManyRequests, types.AuthError{
					Code:    "RATE_LIMIT_EXCEEDED",
					Message: fmt.Sprintf("Rate limit exceeded. Too many requests from %s", clientIP),
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func UserRateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if userID, exists := GetUserIDFromContext(r.Context()); exists {
				if !limiter.Allow(fmt.Sprintf("user:%s", userID)) {
					utils.WriteError(w, http.StatusTooManyRequests, types.AuthError{
						Code:    "USER_RATE_LIMIT_EXCEEDED",
						Message: "Rate limit exceeded for this user",
					})
					return
				}
			} else {
				clientIP := GetClientIP(r)
				if !limiter.Allow(clientIP) {
					utils.WriteError(w, http.StatusTooManyRequests, types.AuthError{
						Code:    "RATE_LIMIT_EXCEEDED",
						Message: "Rate limit exceeded",
					})
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

var (
	LoginRateLimiter            = NewRateLimiter(5, 15*time.Minute)
	RegisterRateLimiter         = NewRateLimiter(3, time.Hour)
	PasswordResetRateLimiter    = NewRateLimiter(3, time.Hour)
	TokenRefreshRateLimiter     = NewRateLimiter(10, time.Minute)
	EmailVerificationRateLimiter = NewRateLimiter(3, time.Hour)
)

func LoginRateLimit() func(http.Handler) http.Handler {
	return RateLimitMiddleware(LoginRateLimiter)
}

func RegisterRateLimit() func(http.Handler) http.Handler {
	return RateLimitMiddleware(RegisterRateLimiter)
}

func PasswordResetRateLimit() func(http.Handler) http.Handler {
	return RateLimitMiddleware(PasswordResetRateLimiter)
}

func TokenRefreshRateLimit() func(http.Handler) http.Handler {
	return UserRateLimitMiddleware(TokenRefreshRateLimiter)
}

func EmailVerificationRateLimit() func(http.Handler) http.Handler {
	return RateLimitMiddleware(EmailVerificationRateLimiter)
}
