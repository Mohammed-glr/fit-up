package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tdmdh/fit-up-server/internal/auth/repository"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/shared/config"
)

func GenerateJWT(secret []byte, userID string) (string, error) {
	expiration := time.Duration(config.NewConfig().JWTExpirationInSeconds) * time.Second
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"iat":     now.Unix(),
		"exp":     now.Add(expiration).Unix(),
		"nbf":     now.Unix(),
		"iss":     "leornian-auth-service",
		"aud":     "leornian-api",
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func GenerateJWTWithClaims(secret []byte, user *types.User) (string, error) {
	if user == nil {
		return "", errors.New("user cannot be nil")
	}

	expiration := time.Duration(config.NewConfig().JWTExpirationInSeconds) * time.Second
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    string(user.Role),
		"jti":     uuid.New().String(),
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
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string, store repository.UserStore, secret []byte) (*types.TokenClaims, error) {
	if tokenString == "" {
		return nil, types.ErrInvalidToken
	}

	if len(secret) == 0 {
		return nil, types.ErrJWTSecretNotSet
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, types.ErrTokenExpired
		}
		return nil, types.ErrInvalidToken
	}

	if !token.Valid {
		return nil, types.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, types.ErrInvalidToken
	}

	tokenClaims, err := extractTokenClaims(claims)
	if err != nil {
		return nil, err
	}

	if time.Now().Unix() > tokenClaims.ExpiresAt {
		return nil, types.ErrTokenExpired
	}

	if time.Now().Unix() < tokenClaims.NotBefore {
		return nil, types.ErrInvalidToken
	}


	return tokenClaims, nil
}

func ExtractClaims(tokenString string) (*types.TokenClaims, error) {
	if tokenString == "" {
		return nil, types.ErrInvalidToken
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, types.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, types.ErrInvalidToken
	}

	return extractTokenClaims(claims)
}

func RefreshJWT(tokenString string, secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", types.ErrJWTSecretNotSet
	}

	claims, err := ExtractClaims(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to extract claims from token: %w", err)
	}

	maxRefreshAge := 30 * 24 * time.Hour
	if time.Now().Unix()-claims.IssuedAt > int64(maxRefreshAge.Seconds()) {
		return "", types.ErrTokenExpired
	}

	now := time.Now()
	expiration := time.Duration(config.NewConfig().JWTExpirationInSeconds) * time.Second

	newClaims := jwt.MapClaims{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    string(claims.Role),
		"jti":     uuid.New().String(),
		"iat":     now.Unix(),
		"exp":     now.Add(expiration).Unix(),
		"nbf":     now.Unix(),
		"iss":     claims.Issuer,
		"aud":     claims.Audience,
		"sub":     claims.Subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	newTokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign refreshed token: %w", err)
	}

	return newTokenString, nil
}

func extractTokenClaims(claims jwt.MapClaims) (*types.TokenClaims, error) {
	tokenClaims := &types.TokenClaims{}

	if userID, ok := claims["user_id"].(string); ok {
		tokenClaims.UserID = userID
	} else {
		return nil, types.ErrInvalidToken
	}

	if email, ok := claims["email"].(string); ok {
		tokenClaims.Email = email
	}

	if roleStr, ok := claims["role"].(string); ok {
		tokenClaims.Role = types.UserRole(roleStr)
	} else {
		tokenClaims.Role = types.RoleUser
	}

	if jti, ok := claims["jti"].(string); ok {
		tokenClaims.JTI = jti
	}

	if exp, ok := claims["exp"].(float64); ok {
		tokenClaims.ExpiresAt = int64(exp)
	} else {
		return nil, types.ErrInvalidToken
	}

	if iat, ok := claims["iat"].(float64); ok {
		tokenClaims.IssuedAt = int64(iat)
	}

	if nbf, ok := claims["nbf"].(float64); ok {
		tokenClaims.NotBefore = int64(nbf)
	}

	if iss, ok := claims["iss"].(string); ok {
		tokenClaims.Issuer = iss
	}

	if aud, ok := claims["aud"].(string); ok {
		tokenClaims.Audience = aud
	}

	if sub, ok := claims["sub"].(string); ok {
		tokenClaims.Subject = sub
	}

	return tokenClaims, nil
}

func IsJWTExpired(tokenString string) bool {
	claims, err := ExtractClaims(tokenString)
	if err != nil {
		return true
	}
	return time.Now().Unix() > claims.ExpiresAt
}

func GetTokenRemainingTime(tokenString string) (time.Duration, error) {
	claims, err := ExtractClaims(tokenString)
	if err != nil {
		return 0, err
	}

	expTime := time.Unix(claims.ExpiresAt, 0)
	remaining := expTime.Sub(time.Now())

	if remaining < 0 {
		return 0, types.ErrTokenExpired
	}

	return remaining, nil
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func HashRefreshToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateTokenPair(ctx context.Context, user *types.User, store repository.RefreshTokenStore, secret []byte) (*types.TokenPair, error) {
	accessToken, err := GenerateJWTWithClaims(secret, user)
	if err != nil {
		return nil, err
	}

	claims, err := ExtractClaims(accessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	tokenHash := HashRefreshToken(refreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = store.CreateRefreshToken(ctx, user.ID, tokenHash, expiresAt, claims.JTI)
	if err != nil {
		return nil, err
	}

	expiration := time.Duration(config.NewConfig().JWTExpirationInSeconds) * time.Second
	return &types.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(expiration.Seconds()),
	}, nil
}

func RotateTokens(ctx context.Context, refreshToken string, store repository.RefreshTokenStore, userStore repository.UserStore, secret []byte) (*types.TokenPair, error) {
	tokenHash := HashRefreshToken(refreshToken)
	storedToken, err := store.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, types.ErrInvalidRefreshToken
	}

	if storedToken.IsRevoked || time.Now().After(storedToken.ExpiresAt) {
		return nil, types.ErrRefreshTokenExpired
	}

	user, err := userStore.GetUserByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, types.ErrUserNotFound
	}

	err = store.RevokeRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	return GenerateTokenPair(ctx, user, store, secret)
}
