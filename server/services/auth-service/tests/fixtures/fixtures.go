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
		PasswordHash:       "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password" bcrypt hashed
		Role:               types.RoleUser,
		IsTwoFactorEnabled: false,
		CreatedAt:          now,
		UpdatedAt:          now,
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

// TestUnverifiedUser creates a test user (email verification removed)
func TestUnverifiedUser() *types.User {
	user := TestUser()
	user.ID = "test-unverified-id-123"
	user.Email = "unverified@example.com"
	user.Username = "testunverified"
	user.Name = "Test Unverified"
	return user
}

// TestRefreshToken creates a test refresh token
func TestRefreshToken() *types.RefreshToken {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	return &types.RefreshToken{
		ID:             "test-refresh-token-id-123",
		UserID:         "test-user-id-123",
		TokenHash:      "hashed-token-value",
		AccessTokenJTI: "access-token-jti-123",
		ExpiresAt:      expiresAt,
		CreatedAt:      now,
		LastUsedAt:     now,
		IsRevoked:      false,
		UserAgent:      "Test User Agent",
		IPAddress:      "127.0.0.1",
	}
}

// TestPasswordResetToken creates a test password reset token
func TestPasswordResetToken() *types.PasswordResetToken {
	now := time.Now()
	expiresAt := now.Add(1 * time.Hour)

	return &types.PasswordResetToken{
		ID:      "reset-token-id-123",
		Email:   "test@example.com",
		Token:   "reset-token-123",
		Expires: expiresAt,
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
		Name:     "New User",
		Password: "password",
	}
}

// TestChangePasswordRequest creates a test change password request
func TestChangePasswordRequest() *types.ChangePasswordRequest {
	return &types.ChangePasswordRequest{
		CurrentPassword: "password",
		NewPassword:     "newpassword",
	}
}

// TestUpdateUserRequest creates a test update user request
func TestUpdateUserRequest() *types.UpdateUserRequest {
	name := "Updated Name"
	bio := "Updated bio"
	return &types.UpdateUserRequest{
		Name: &name,
		Bio:  &bio,
	}
}
