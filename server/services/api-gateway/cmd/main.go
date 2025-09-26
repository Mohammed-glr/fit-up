package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	gatewayMiddleware "github.com/tdmdh/fit-up-server/services/api-gateway/internal/middleware"
)

func main() {
	log.Println("API Gateway starting...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(gatewayMiddleware.CORS())
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API Gateway is healthy"))
	})

	// Basic routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Lornian Backend API Gateway"))
	})

	// Service proxying
	authBaseURL := os.Getenv("AUTH_SERVICE_URL")
	if authBaseURL == "" {
		authBaseURL = "http://lornian-auth-service:8081"
	}

	userBaseURL := os.Getenv("USER_SERVICE_URL")
	if userBaseURL == "" {
		userBaseURL = "http://lornian-user-service:8082"
	}

	aiBaseURL := os.Getenv("AI_SERVICE_URL")
	if aiBaseURL == "" {
		aiBaseURL = "http://lornian-ai-service:8083"
	}

	// Auth service routes
	r.Route("/auth", func(r chi.Router) {
		r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			proxyRequestWithPrefix(w, r, authBaseURL, "/auth")
		})
	})

	// User service routes
	r.Route("/user", func(r chi.Router) {
		r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			proxyRequestWithPrefix(w, r, userBaseURL, "/user")
		})
	})

	// AI service routes
	r.Route("/ai", func(r chi.Router) {
		r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			proxyRequestWithPrefix(w, r, aiBaseURL, "/ai")
		})
	})

	log.Printf("Starting API Gateway on port %s", port)
	log.Printf("Auth Service: %s", authBaseURL)
	log.Printf("User Service: %s", userBaseURL)
	log.Printf("AI Service: %s", aiBaseURL)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Failed to start API Gateway:", err)
	}
}

func proxyRequestWithPrefix(w http.ResponseWriter, r *http.Request, baseURL string, prefix string) {
	// Strip the prefix from the path
	path := r.URL.Path
	if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
		path = path[len(prefix):]
	}

	// Ensure path starts with /
	if path == "" || path[0] != '/' {
		path = "/" + path
	}

	// Create the target URL
	targetURL := baseURL + path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	log.Printf("Proxying %s %s to %s (stripped prefix %s)", r.Method, r.URL.Path, targetURL, prefix)

	// Create new request
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make the request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error proxying request to %s: %v", targetURL, err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers (skip CORS headers as API Gateway handles them)
	for key, values := range resp.Header {
		// Skip CORS headers to avoid conflicts
		if strings.HasPrefix(strings.ToLower(key), "access-control-") {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	log.Printf("Proxying response: status=%d, content-type=%s", resp.StatusCode, resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
