package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tdmdh/fit-up-server/internal/auth/repository"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/shared/config"
)

type AuthService struct {
	userStore repository.UserStore
}

func NewAuthService(userStore repository.UserStore) *AuthService {
	return &AuthService{
		userStore: userStore,
	}
}

func (s *AuthService) ResetPassword(ctx context.Context, payload types.ResetPasswordRequest) error {
	resetToken, err := s.userStore.GetPasswordResetToken(ctx, payload.Token)
	if err != nil {
		if err == types.ErrPasswordResetTokenNotFound {
			return types.ErrInvalidToken
		}
		return err
	}

	if !ValidatePasswordResetToken(resetToken) {
		return types.ErrTokenExpired
	}

	user, err := s.userStore.GetUserByPasswordResetToken(ctx, payload.Token)
	if err != nil {
		if err == types.ErrPasswordResetTokenNotFound {
			return types.ErrUserNotFound
		}
		return err
	}

	if ComparePasswords(user.PasswordHash, []byte(payload.NewPassword)) {
		return types.ErrSamePassword
	}

	if err := ValidatePasswordStrength(payload.NewPassword); err != nil {
		return err
	}

	hashedPassword, err := HashPassword(payload.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userStore.UpdateUserPassword(ctx, user.ID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := s.userStore.DeletePasswordResetToken(ctx, payload.Token); err != nil {
		fmt.Printf("Warning: failed to delete reset token: %v\n", err)
	}

	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.userStore.GetUserByID(ctx, userID)
	if err != nil {
		return types.ErrUserNotFound
	}

	if !ComparePasswords(user.PasswordHash, []byte(oldPassword)) {
		return types.ErrInvalidCredentials
	}

	if ComparePasswords(user.PasswordHash, []byte(newPassword)) {
		return types.ErrSamePassword
	}

	if err := ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userStore.UpdateUserPassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *AuthService) GetUser(ctx context.Context, userID string) (*types.User, error) {
	return s.userStore.GetUserByID(ctx, userID)
}

func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	return s.userStore.GetUserByUsername(ctx, username)
}

func (s *AuthService) Logout(ctx context.Context, userID string) error {
	err := s.userStore.RevokeAllUserRefreshTokens(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh tokens: %w", err)
	}

	return nil
}

func (s *AuthService) GenerateTokenPair(ctx context.Context, user *types.User) (*types.TokenPair, error) {
	if user == nil {
		return nil, types.ErrUserNotFound
	}

	accessToken, jti, err := s.generateAccessTokenWithJTI(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshTokenExpiry := time.Now().Add(time.Duration(config.NewConfig().RefreshTokenExpirationInSeconds) * time.Second)
	err = s.userStore.CreateRefreshToken(ctx, user.ID, refreshToken, refreshTokenExpiry, jti)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	accessTokenExpiry := time.Duration(config.NewConfig().JWTExpirationInSeconds) * time.Second

	return &types.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) RotateTokens(ctx context.Context, refreshToken string) (*types.TokenPair, error) {
	if refreshToken == "" {
		return nil, types.ErrRefreshTokenNotFound
	}

	storedRefreshToken, err := s.userStore.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if storedRefreshToken.IsRevoked {
		return nil, types.ErrRefreshTokenNotFound
	}

	if time.Now().After(storedRefreshToken.ExpiresAt) {
		return nil, types.ErrRefreshTokenExpired
	}

	user, err := s.userStore.GetUserByID(ctx, storedRefreshToken.UserID)
	if err != nil {
		return nil, err
	}

	err = s.userStore.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}


	return s.GenerateTokenPair(ctx, user)
}

func ValidatePasswordResetToken(token *types.PasswordResetToken) bool {
	return time.Now().Before(token.Expires)
}

func CreatePasswordResetToken(email string) (*types.PasswordResetToken, error) {
	token, err := generateRandomToken(32)
	if err != nil {
		return nil, err
	}

	return &types.PasswordResetToken{
		ID:      uuid.New().String(),
		Email:   email,
		Token:   token,
		Expires: time.Now().Add(1 * time.Hour),
	}, nil
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *AuthService) generateAccessTokenWithJTI(user *types.User) (string, string, error) {
	secret := []byte(config.NewConfig().JWTSecret)
	if len(secret) == 0 {
		return "", "", types.ErrJWTSecretNotSet
	}

	expiration := time.Duration(config.NewConfig().JWTExpirationInSeconds) * time.Second
	now := time.Now()
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    string(user.Role),
		"jti":     jti,
		"iat":     now.Unix(),
		"exp":     now.Add(expiration).Unix(),
		"nbf":     now.Unix(),
		"iss":     "leornian-auth-service",
		"aud":     "leornian-api",
		"sub":     user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, jti, nil
}

func (s *AuthService) generateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}
