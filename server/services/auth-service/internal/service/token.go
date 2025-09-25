package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
)

func GenerateVerificationToken(userID string) (string, error) {
	token := uuid.New().String()
	timestamp := time.Now().Unix()
	combined := fmt.Sprintf("%s-%s-%d", token, userID, timestamp)
	encoded := base64.URLEncoding.EncodeToString([]byte(combined))
	return encoded, nil
}

func CreateVerificationToken(email string) (*types.VerificationToken, error) {
	token, err := GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}
	return &types.VerificationToken{
		ID:      uuid.New().String(),
		Email:   email,
		Token:   token,
		Expires: time.Now().Add(24 * time.Hour),
	}, nil
}

func CreatePasswordResetToken(email string) (*types.PasswordResetToken, error) {
	token, err := GenerateSecureToken(32)
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

func ValidateVerificationToken(token *types.VerificationToken) bool {
	if token == nil {
		return false
	}
	return time.Now().Before(token.Expires) && token.Token != ""
}

func ValidatePasswordResetToken(token *types.PasswordResetToken) bool {
	if token == nil {
		return false
	}
	return time.Now().Before(token.Expires) && token.Token != ""
}

func GenerateSecureToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("token length must be positive")
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateUUIDToken() string {
	return uuid.New().String()
}

func GenerateSessionToken(userID string, duration time.Duration) (*types.Session, error) {
	tokenValue, err := GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}
	return &types.Session{
		ID:           uuid.New().String(),
		UserID:       userID,
		SessionToken: tokenValue,
		Expires:      time.Now().Add(duration),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func ValidateSession(session *types.Session) bool {
	if session == nil {
		return false
	}
	if time.Now().After(session.Expires) {
		return false
	}
	if session.SessionToken == "" {
		return false
	}
	return true
}

func ParseVerificationToken(tokenStr string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(tokenStr)
	if err != nil {
		return "", fmt.Errorf("invalid token format: %w", err)
	}
	if len(decoded) < 10 {
		return "", fmt.Errorf("invalid token length")
	}
	return string(decoded), nil
}

func IsTokenExpired(createdAt time.Time, duration time.Duration) bool {
	return time.Now().After(createdAt.Add(duration))
}
