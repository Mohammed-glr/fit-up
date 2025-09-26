package fixtures

import (
	"time"

	"github.com/tdmdh/fit-up-server/services/auth-service/internal/types"
)

// TestUser creates a test user with default values
func TestUser() *types.User {
	now := time.Now()
	return &types.User{
		ID:                 "test-user-id-123",
		Email:              "test@example.com",
		Username:           "testuser",
		Name:               "Test User",
		Bio:                "Test user bio",
		Password:           "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password" bcrypt hashed
		Role:               types.RoleUser,
		IsTwoFactorEnabled: false,
		SubroleID:          1,
		CreatedAt:          now,
		UpdatedAt:          now,
		EmailVerified:      &now,
	}
}

// TestAdminUser creates a test admin user
func TestAdminUser() *types.User {
	user := TestUser()
	user.ID = "test-admin-id-123"
	user.Email = "admin@example.com"
	user.Username = "testadmin"
	user.Name = "Test Admin"
	user.Role = types.RoleAdmin
	return user
}

// TestUnverifiedUser creates a test unverified user
func TestUnverifiedUser() *types.User {
	user := TestUser()
	user.ID = "test-unverified-id-123"
	user.Email = "unverified@example.com"
	user.Username = "testunverified"
	user.Name = "Test Unverified"
	user.EmailVerified = nil // nil indicates not verified
	return user
}

// TestRefreshToken creates a test refresh token
func TestRefreshToken() *types.RefreshToken {
	expiresAt := time.Now().Add(24 * time.Hour)
	createdAt := time.Now()
	return &types.RefreshToken{
		ID:             "test-refresh-token-123",
		UserID:         "test-user-id-123",
		TokenHash:      "hashed_token_value",
		AccessTokenJTI: "test-jti-123",
		ExpiresAt:      expiresAt,
		CreatedAt:      createdAt,
		LastUsedAt:     createdAt,
		IsRevoked:      false,
		UserAgent:      "Test User Agent",
		IPAddress:      "127.0.0.1",
	}
}

// TestSession creates a test session
func TestSession() *types.Session {
	expiresAt := time.Now().Add(24 * time.Hour)
	createdAt := time.Now()
	return &types.Session{
		ID:           "test-session-id-123",
		UserID:       "test-user-id-123",
		SessionToken: "test-session-token-123",
		Expires:      expiresAt,
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}
}

// TestAuditEvent creates a test audit event
func TestAuditEvent() types.AuthAuditEvent {
	return types.AuthAuditEvent{
		ID:        "test-audit-id-123",
		UserID:    "test-user-id-123",
		Action:    "login",
		Success:   true,
		IPAddress: "127.0.0.1",
		UserAgent: "Test User Agent",
		Details:   map[string]interface{}{"method": "email"},
		CreatedAt: time.Now(),
	}
}

// TestLoginRequest creates a test login request with email
func TestLoginRequest() *types.LoginRequest {
	return &types.LoginRequest{
		Identifier: "test@example.com",
		Password:   "password",
	}
}

// TestLoginRequestWithUsername creates a test login request with username
func TestLoginRequestWithUsername() *types.LoginRequest {
	return &types.LoginRequest{
		Identifier: "testuser",
		Password:   "password",
	}
}

// TestRegisterRequest creates a test register request
func TestRegisterRequest() *types.RegisterRequest {
	return &types.RegisterRequest{
		Email:    "newuser@example.com",
		Username: "newuser",
		Password: "password123",
		Name:     "New User",
	}
}

// TestChangePasswordRequest creates a test change password request
func TestChangePasswordRequest() *types.ChangePasswordRequest {
	return &types.ChangePasswordRequest{
		CurrentPassword: "password",
		NewPassword:     "newpassword123",
	}
}

// TestForgotPasswordRequest creates a test forgot password request
func TestForgotPasswordRequest() *types.ForgotPasswordRequest {
	return &types.ForgotPasswordRequest{
		Email: "test@example.com",
	}
}

// TestResetPasswordRequest creates a test reset password request
func TestResetPasswordRequest() *types.ResetPasswordRequest {
	return &types.ResetPasswordRequest{
		Token:       "reset-token-123",
		NewPassword: "newpassword123",
	}
}

// TestTokenClaims creates test JWT token claims
func TestTokenClaims() *types.TokenClaims {
	now := time.Now().Unix()
	return &types.TokenClaims{
		UserID:    "test-user-id-123",
		Email:     "test@example.com",
		Role:      types.RoleUser,
		JTI:       "test-jti-123",
		Issuer:    "leornian-auth",
		Subject:   "test-user-id-123",
		Audience:  "leornian-app",
		ExpiresAt: now + 3600,
		IssuedAt:  now,
		NotBefore: now,
	}
}

// TestOAuthUserInfo creates test OAuth user info
func TestOAuthUserInfo() *types.OAuthUserInfo {
	return &types.OAuthUserInfo{
		ID:            "oauth-user-123",
		Email:         "oauth@example.com",
		Name:          "OAuth User",
		Username:      "oauthuser",
		AvatarURL:     "https://example.com/avatar.png",
		EmailVerified: true,
	}
}
