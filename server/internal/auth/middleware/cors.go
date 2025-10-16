package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := os.Getenv("CORS_ORIGINS")
			if allowedOrigins == "" {
			allowedOrigins = "http://localhost:3000,https://localhost:3000,http://127.0.0.1:3000,https://127.0.0.1:3000,https://*.github.dev,https://*.githubpreview.dev"
			}

			origin := r.Header.Get("Origin")
			
			if isOriginAllowed(origin, allowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Requested-With")
			w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-Response-Time, X-RateLimit-Remaining")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(origin, allowedOrigins string) bool {
	if origin == "" {
		return false
	}

	origins := strings.Split(allowedOrigins, ",")
	for _, allowedOrigin := range origins {
		allowedOrigin = strings.TrimSpace(allowedOrigin)
		
		if allowedOrigin == origin {
			return true
		}
		

		if strings.HasPrefix(allowedOrigin, "*.") {
			domain := allowedOrigin[2:] // Remove "*.
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
		
		if strings.Contains(allowedOrigin, "*.github.dev") && strings.Contains(origin, ".github.dev") {
			return true
		}
		if strings.Contains(allowedOrigin, "*.githubpreview.dev") && strings.Contains(origin, ".githubpreview.dev") {
			return true
		}
	}
	
	return false
}
