package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// IsEmailFormat checks if the identifier is in email format
func IsEmailFormat(identifier string) bool {
	// Simple email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(identifier) && strings.Contains(identifier, "@")
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is nil")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteSuccess(w http.ResponseWriter, status int, message string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]string{"message": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write success response", http.StatusInternalServerError)
	}
}



// TODO: Step 1 - Implement password utilities:
//   - HashPassword(password string) (string, error) - bcrypt hashing with salt
//   - VerifyPassword(password, hash string) bool - password verification
//   - GenerateSecurePassword(length int) string - random password generation
//   - ValidatePasswordStrength(password string) error - strength validation
// TODO: Step 2 - JWT token utilities (IMPLEMENTED ✅):
//   ✅ GenerateJWT(claims TokenClaims) (string, error) - JWT creation
//   ✅ GenerateJWTWithClaims(secret []byte, user *User) (string, error) - JWT creation with full claims
//   ✅ ValidateJWT(token string, secret []byte) (*TokenClaims, error) - JWT validation
//   ✅ ExtractClaims(token string) (*TokenClaims, error) - claims extraction
//   ✅ RefreshJWT(token string, secret []byte) (string, error) - token refresh
//   ✅ IsJWTExpired(token string) bool - expiration check
//   ✅ GetTokenRemainingTime(token string) (time.Duration, error) - remaining time
//
// JWT Implementation Details:
//   - Location: services/auth-service/internal/service/jwt.go
//   - Middleware: services/auth-service/internal/middleware/jwt_middleware.go
//   - Handlers: services/auth-service/internal/handlers/jwt_handlers.go
//   - Shared Utils: shared/utils/jwt.go
//   - Documentation: JWT_IMPLEMENTATION.md
//   - Test Script: test-jwt-implementation.sh
//
// Available JWT Endpoints:
//   - POST /validate-token - Validate JWT token and return user info
//   - POST /refresh-token - Refresh JWT token
//   - POST /logout - Logout (token invalidation placeholder)
//
// JWT Middleware:
//   - JWTAuthMiddleware() - Requires valid JWT token
//   - OptionalJWTAuthMiddleware() - Optional JWT validation
//   - RequireRoleMiddleware(role) - Role-based authorization
//   - RequireAdminMiddleware() - Admin-only access
//
// Context Utilities:
//   - GetUserIDFromContext(ctx) - Extract user ID from context
//   - GetUserClaimsFromContext(ctx) - Extract JWT claims from context
//   - GetUserFromContext(ctx) - Extract user object from context
// TODO: Step 3 - Implement security utilities:
//   - GenerateCSRFToken() string - CSRF token generation
//   - ValidateCSRFToken(token string) bool - CSRF validation
//   - GenerateOTP() string - One-time password generation
//   - ValidateOTP(otp string) bool - OTP validation
//   - SanitizeInput(input string) string - input sanitization
// TODO: Step 4 - Implement email utilities:
//   - ValidateEmail(email string) bool - email format validation
//   - GenerateVerificationToken() string - email verification token
//   - SendPasswordResetEmail(email, token string) error - password reset email
//   - SendVerificationEmail(email, token string) error - account verification email
// TODO: Step 5 - Implement rate limiting utilities:
//   - CheckRateLimit(key string, limit int, window time.Duration) bool
//   - IncrementAttempts(key string) error - failed login tracking

// Flow: service.go -> utils.go -> external libraries (bcrypt, JWT, email)
// Dependencies: bcrypt, JWT library, email service, Redis/in-memory store for rate limiting
