package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tdmdh/fit-up-server/shared/config"
)

// TestConfig returns a test configuration
func TestConfig() *config.Config {
	return &config.Config{
		PublicHost:                      "http://localhost",
		Port:                            "8080",
		DatabaseURL:                     "postgres://test_user:test_password@localhost:5433/lornian_test?sslmode=disable",
		JWTSecret:                       "test_jwt_secret_key_for_testing_purposes_only",
		JWTExpirationInSeconds:          3600,
		RefreshTokenExpirationInSeconds: 3600 * 24,
		ResendAPIKey:                    "test_resend_api_key",
		Database: config.DatabaseConfig{
			MaxConnections:    10,
			MinConnections:    1,
			MaxConnLifetime:   60,
			MaxConnIdleTime:   30,
			HealthCheckPeriod: 5,
			ConnectTimeout:    10,
		},
		OAuthConfig: config.OAuthConfig{
			GoogleClientID:       "test_google_client_id",
			GoogleClientSecret:   "test_google_client_secret",
			GoogleRedirectURI:    "http://localhost:8080/auth/google/callback",
			GitHubClientID:       "test_github_client_id",
			GitHubClientSecret:   "test_github_client_secret",
			GitHubRedirectURI:    "http://localhost:8080/auth/github/callback",
			FacebookClientID:     "test_facebook_client_id",
			FacebookClientSecret: "test_facebook_client_secret",
			FacebookRedirectURI:  "http://localhost:8080/auth/facebook/callback",
			OAuthStateSecret:     "test_oauth_state_secret",
		},
	}
}

// CreateTestRequest creates an HTTP request for testing
func CreateTestRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("Failed to encode request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateTestRequestWithAuth creates an HTTP request with authorization header
func CreateTestRequestWithAuth(t *testing.T, method, url string, body interface{}, token string) *http.Request {
	req := CreateTestRequest(t, method, url, body)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// ExecuteRequest executes an HTTP request and returns the response
func ExecuteRequest(req *http.Request, handler http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// AssertStatusCode checks if the response has the expected status code
func AssertStatusCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected status code %d, got %d", expected, actual)
	}
}

// AssertJSONResponse checks if the response body matches the expected JSON
func AssertJSONResponse(t *testing.T, expected interface{}, response *httptest.ResponseRecorder) {
	var actual interface{}
	if err := json.NewDecoder(response.Body).Decode(&actual); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)

	if string(expectedJSON) != string(actualJSON) {
		t.Errorf("Expected JSON %s, got %s", expectedJSON, actualJSON)
	}
}

// AssertContains checks if the response body contains the expected string
func AssertContains(t *testing.T, expected, actual string) {
	if !bytes.Contains([]byte(actual), []byte(expected)) {
		t.Errorf("Expected response to contain %s, got %s", expected, actual)
	}
}
