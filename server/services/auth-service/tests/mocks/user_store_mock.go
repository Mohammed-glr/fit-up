package mocks

import (
	"context"
	"time"

	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
)

// MockUserStore is a mock implementation of the UserStore interface
type MockUserStore struct {
	users               map[string]*types.User
	emails              map[string]*types.User
	usernames           map[string]*types.User
	refreshTokens       map[string]*types.RefreshToken
	passwordResetTokens map[string]*types.PasswordResetToken
	blacklistedTokens   map[string]bool
	shouldError         bool
	errorType           error
}

// NewMockUserStore creates a new mock user store
func NewMockUserStore() *MockUserStore {
	return &MockUserStore{
		users:               make(map[string]*types.User),
		emails:              make(map[string]*types.User),
		usernames:           make(map[string]*types.User),
		refreshTokens:       make(map[string]*types.RefreshToken),
		passwordResetTokens: make(map[string]*types.PasswordResetToken),
		blacklistedTokens:   make(map[string]bool),
	}
}

// SetError configures the mock to return an error
func (m *MockUserStore) SetError(err error) {
	m.shouldError = true
	m.errorType = err
}

// ClearError clears the error state
func (m *MockUserStore) ClearError() {
	m.shouldError = false
	m.errorType = nil
}

// AddUser adds a user to the mock store
func (m *MockUserStore) AddUser(user *types.User) {
	m.users[user.ID] = user
	m.emails[user.Email] = user
	m.usernames[user.Username] = user
}

// GetUserByID retrieves a user by ID
func (m *MockUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	if m.shouldError {
		return nil, m.errorType
	}
	user, exists := m.users[id]
	if !exists {
		return nil, types.ErrUserNotFound
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	if m.shouldError {
		return nil, m.errorType
	}
	user, exists := m.emails[email]
	if !exists {
		return nil, types.ErrUserNotFound
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (m *MockUserStore) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	if m.shouldError {
		return nil, m.errorType
	}
	user, exists := m.usernames[username]
	if !exists {
		return nil, types.ErrUserNotFound
	}
	return user, nil
}

// CreateUser creates a new user
func (m *MockUserStore) CreateUser(ctx context.Context, user *types.User) error {
	if m.shouldError {
		return m.errorType
	}
	if _, exists := m.emails[user.Email]; exists {
		return types.ErrUserAlreadyExists
	}
	m.users[user.ID] = user
	m.emails[user.Email] = user
	m.usernames[user.Username] = user
	return nil
}

// UpdateUser updates an existing user
func (m *MockUserStore) UpdateUser(ctx context.Context, id string, updates *types.UpdateUserRequest) error {
	if m.shouldError {
		return m.errorType
	}
	user, exists := m.users[id]
	if !exists {
		return types.ErrUserNotFound
	}
	// Apply updates (simplified for mock)
	if updates.Username != nil {
		user.Username = *updates.Username
	}
	if updates.Name != nil {
		user.Name = *updates.Name
	}
	if updates.Bio != nil {
		user.Bio = *updates.Bio
	}
	if updates.Image != nil {
		user.Image = *updates.Image
	}
	return nil
}

// UpdateUserPassword updates a user's password
func (m *MockUserStore) UpdateUserPassword(ctx context.Context, userID string, hashedPassword string) error {
	if m.shouldError {
		return m.errorType
	}
	user, exists := m.users[userID]
	if !exists {
		return types.ErrUserNotFound
	}
	user.PasswordHash = hashedPassword
	return nil
}

// CreateRefreshToken creates a refresh token
func (m *MockUserStore) CreateRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time, accessTokenJTI string) error {
	if m.shouldError {
		return m.errorType
	}
	refreshToken := &types.RefreshToken{
		ID:             token,
		UserID:         userID,
		TokenHash:      token,
		AccessTokenJTI: accessTokenJTI,
		ExpiresAt:      expiresAt,
		CreatedAt:      time.Now(),
		LastUsedAt:     time.Now(),
		IsRevoked:      false,
	}
	m.refreshTokens[token] = refreshToken
	return nil
}

// GetRefreshToken retrieves a refresh token
func (m *MockUserStore) GetRefreshToken(ctx context.Context, token string) (*types.RefreshToken, error) {
	if m.shouldError {
		return nil, m.errorType
	}
	refreshToken, exists := m.refreshTokens[token]
	if !exists {
		return nil, types.ErrRefreshTokenNotFound
	}
	return refreshToken, nil
}

// DeleteRefreshToken deletes a refresh token
func (m *MockUserStore) DeleteRefreshToken(ctx context.Context, token string) error {
	if m.shouldError {
		return m.errorType
	}
	delete(m.refreshTokens, token)
	return nil
}

// CleanupExpiredRefreshTokens cleans up expired refresh tokens
func (m *MockUserStore) CleanupExpiredRefreshTokens(ctx context.Context) error {
	if m.shouldError {
		return m.errorType
	}
	return nil
}

// RevokeRefreshToken revokes a refresh token
func (m *MockUserStore) RevokeRefreshToken(ctx context.Context, token string) error {
	if m.shouldError {
		return m.errorType
	}
	refreshToken, exists := m.refreshTokens[token]
	if !exists {
		return types.ErrRefreshTokenNotFound
	}
	refreshToken.IsRevoked = true
	now := time.Now()
	refreshToken.RevokedAt = &now
	m.refreshTokens[token] = refreshToken
	return nil
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a user
func (m *MockUserStore) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	if m.shouldError {
		return m.errorType
	}
	for token, refreshToken := range m.refreshTokens {
		if refreshToken.UserID == userID {
			refreshToken.IsRevoked = true
			now := time.Now()
			refreshToken.RevokedAt = &now
			m.refreshTokens[token] = refreshToken
		}
	}
	return nil
}

// UpdateRefreshTokenLastUsed updates the last used time of a refresh token
func (m *MockUserStore) UpdateRefreshTokenLastUsed(ctx context.Context, token string) error {
	if m.shouldError {
		return m.errorType
	}
	refreshToken, exists := m.refreshTokens[token]
	if !exists {
		return types.ErrRefreshTokenNotFound
	}
	refreshToken.LastUsedAt = time.Now()
	m.refreshTokens[token] = refreshToken
	return nil
}

// CreatePasswordResetToken creates a password reset token
func (m *MockUserStore) CreatePasswordResetToken(ctx context.Context, email string, token string, expiresAt time.Time) error {
	if m.shouldError {
		return m.errorType
	}
	resetToken := &types.PasswordResetToken{
		ID:      token,
		Email:   email,
		Token:   token,
		Expires: expiresAt,
	}
	m.passwordResetTokens[token] = resetToken
	return nil
}

// GetPasswordResetToken retrieves a password reset token
func (m *MockUserStore) GetPasswordResetToken(ctx context.Context, token string) (*types.PasswordResetToken, error) {
	if m.shouldError {
		return nil, m.errorType
	}
	resetToken, exists := m.passwordResetTokens[token]
	if !exists {
		return nil, types.ErrPasswordResetTokenNotFound
	}
	return resetToken, nil
}

// GetUserByPasswordResetToken retrieves a user by password reset token
func (m *MockUserStore) GetUserByPasswordResetToken(ctx context.Context, token string) (*types.User, error) {
	if m.shouldError {
		return nil, m.errorType
	}
	resetToken, exists := m.passwordResetTokens[token]
	if !exists {
		return nil, types.ErrPasswordResetTokenNotFound
	}
	user, exists := m.emails[resetToken.Email]
	if !exists {
		return nil, types.ErrUserNotFound
	}
	return user, nil
}

// DeletePasswordResetToken deletes a password reset token
func (m *MockUserStore) DeletePasswordResetToken(ctx context.Context, token string) error {
	if m.shouldError {
		return m.errorType
	}
	delete(m.passwordResetTokens, token)
	return nil
}

// MarkPasswordResetTokenAsUsed marks a password reset token as used
func (m *MockUserStore) MarkPasswordResetTokenAsUsed(ctx context.Context, token string) error {
	if m.shouldError {
		return m.errorType
	}
	return nil
}

// BlacklistToken blacklists a token
func (m *MockUserStore) BlacklistToken(ctx context.Context, jti, userID, reason string, expiresAt time.Time) error {
	if m.shouldError {
		return m.errorType
	}
	m.blacklistedTokens[jti] = true
	return nil
}

// IsTokenBlacklisted checks if a token is blacklisted
func (m *MockUserStore) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	if m.shouldError {
		return false, m.errorType
	}
	return m.blacklistedTokens[jti], nil
}

// CleanupExpiredTokens cleans up expired tokens
func (m *MockUserStore) CleanupExpiredTokens(ctx context.Context) error {
	if m.shouldError {
		return m.errorType
	}
	return nil
}
