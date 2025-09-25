package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
)

// MockUserStore for testing
type MockUserStore struct{}

func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	return nil, types.ErrUserNotFound
}

func (m *MockUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	return &types.User{
		ID:       id,
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
		Role:     types.RoleUser,
	}, nil
}

func (m *MockUserStore) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	return nil, types.ErrUserNotFound
}

func (m *MockUserStore) CreateUser(ctx context.Context, user *types.User) error {
	return nil
}

func (m *MockUserStore) UpdateUser(ctx context.Context, id string, updates *types.UpdateUserRequest) error {
	return nil
}

func (m *MockUserStore) UpdateUserPassword(ctx context.Context, userID string, hashedPassword string) error {
	return nil
}

func (m *MockUserStore) CreateRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time, accessTokenJTI string) error {
	return nil
}

func (m *MockUserStore) GetRefreshToken(ctx context.Context, token string) (*types.RefreshToken, error) {
	return &types.RefreshToken{
		ID:             "test-refresh-id",
		UserID:         "test-user-id",
		TokenHash:      "hash",
		AccessTokenJTI: "test-jti",
		ExpiresAt:      time.Now().Add(time.Hour),
		CreatedAt:      time.Now(),
		LastUsedAt:     time.Now(),
		IsRevoked:      false,
	}, nil
}

func (m *MockUserStore) DeleteRefreshToken(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) CleanupExpiredRefreshTokens(ctx context.Context) error {
	return nil
}

func (m *MockUserStore) RevokeRefreshToken(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	return nil
}

func (m *MockUserStore) UpdateRefreshTokenLastUsed(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) CreatePasswordResetToken(ctx context.Context, email string, token string, expiresAt time.Time) error {
	return nil
}

func (m *MockUserStore) GetPasswordResetToken(ctx context.Context, token string) (*types.PasswordResetToken, error) {
	return nil, types.ErrPasswordResetTokenNotFound
}

func (m *MockUserStore) GetUserByPasswordResetToken(ctx context.Context, token string) (*types.User, error) {
	return nil, types.ErrUserNotFound
}

func (m *MockUserStore) DeletePasswordResetToken(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) MarkPasswordResetTokenAsUsed(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) BlacklistToken(ctx context.Context, jti, userID, reason string, expiresAt time.Time) error {
	return nil
}

func (m *MockUserStore) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	return false, nil
}

func (m *MockUserStore) CleanupExpiredTokens(ctx context.Context) error {
	return nil
}

func (m *MockUserStore) GetUserByVerificationToken(ctx context.Context, token string) (*types.User, error) {
	return nil, types.ErrUserNotFound
}

func (m *MockUserStore) UpdateUserVerificationStatus(ctx context.Context, userID string, verified bool) error {
	return nil
}

func (m *MockUserStore) DeleteVerificationToken(ctx context.Context, userID string) error {
	return nil
}

func (m *MockUserStore) CreateVerificationToken(ctx context.Context, userID, token string, expiresAt string) error {
	return nil
}

func (m *MockUserStore) ResendVerificationEmail(ctx context.Context, email string) error {
	return nil
}

func TestTokenPairGeneration(t *testing.T) {
	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	// Create mock store
	userStore := &MockUserStore{}

	// Create auth service
	authService := NewAuthService(userStore, userStore)

	// Test user
	user := &types.User{
		ID:       "test-user-id",
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
		Role:     types.RoleUser,
	}

	// Test token pair generation
	tokenPair, err := authService.GenerateTokenPair(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	if tokenPair.AccessToken == "" {
		t.Error("Access token should not be empty")
	}

	if tokenPair.RefreshToken == "" {
		t.Error("Refresh token should not be empty")
	}

	if tokenPair.ExpiresIn <= 0 {
		t.Error("ExpiresIn should be positive")
	}
}

func TestTokenRotation(t *testing.T) {
	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	// Create mock store
	userStore := &MockUserStore{}

	// Create auth service
	authService := NewAuthService(userStore, userStore)

	// Test token rotation
	newTokenPair, err := authService.RotateTokens(context.Background(), "test-refresh-token")
	if err != nil {
		t.Fatalf("Failed to rotate tokens: %v", err)
	}

	if newTokenPair.AccessToken == "" {
		t.Error("New access token should not be empty")
	}

	if newTokenPair.RefreshToken == "" {
		t.Error("New refresh token should not be empty")
	}

	if newTokenPair.ExpiresIn <= 0 {
		t.Error("New ExpiresIn should be positive")
	}
}
