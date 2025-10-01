package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS middleware to handle cross-origin requests
func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get allowed origins from environment variable
			allowedOrigins := os.Getenv("CORS_ORIGINS")
			if allowedOrigins == "" {			// Default origins for development
			allowedOrigins = "http://localhost:3000,https://localhost:3000,http://127.0.0.1:3000,https://127.0.0.1:3000,https://*.github.dev,https://*.githubpreview.dev"
			}

			origin := r.Header.Get("Origin")
			
			// Check if the origin is allowed
			if isOriginAllowed(origin, allowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Requested-With")
			w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-Response-Time, X-RateLimit-Remaining")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

			// Handle preflight OPTIONS request
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isOriginAllowed checks if the origin is in the allowed origins list
func isOriginAllowed(origin, allowedOrigins string) bool {
	if origin == "" {
		return false
	}

	origins := strings.Split(allowedOrigins, ",")
	for _, allowedOrigin := range origins {
		allowedOrigin = strings.TrimSpace(allowedOrigin)
		
		// Exact match
		if allowedOrigin == origin {
			return true
		}
		
		// Wildcard subdomain match (*.example.com)
		if strings.HasPrefix(allowedOrigin, "*.") {
			domain := allowedOrigin[2:] // Remove "*."
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
		
		// GitHub Codespaces pattern match
		if strings.Contains(allowedOrigin, "*.github.dev") && strings.Contains(origin, ".github.dev") {
			return true
		}
		if strings.Contains(allowedOrigin, "*.githubpreview.dev") && strings.Contains(origin, ".githubpreview.dev") {
			return true
		}
	}
	
	return false
}
