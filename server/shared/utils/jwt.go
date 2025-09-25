package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Audience  string `json:"aud"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	NotBefore int64  `json:"nbf"`
}

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrJWTSecretNotSet = errors.New("JWT secret not set")
)

func ValidateJWTToken(tokenString string, secret []byte) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	if len(secret) == 0 {
		return nil, ErrJWTSecretNotSet
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return extractJWTClaims(claims)
}

func ExtractJWTClaims(tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return extractJWTClaims(claims)
}

func IsJWTTokenExpired(tokenString string) bool {
	claims, err := ExtractJWTClaims(tokenString)
	if err != nil {
		return true
	}
	return time.Now().Unix() > claims.ExpiresAt
}

func GetJWTTokenRemainingTime(tokenString string) (time.Duration, error) {
	claims, err := ExtractJWTClaims(tokenString)
	if err != nil {
		return 0, err
	}

	expTime := time.Unix(claims.ExpiresAt, 0)
	remaining := expTime.Sub(time.Now())

	if remaining < 0 {
		return 0, ErrTokenExpired
	}

	return remaining, nil
}

func extractJWTClaims(claims jwt.MapClaims) (*JWTClaims, error) {
	jwtClaims := &JWTClaims{}

	if userID, ok := claims["user_id"].(string); ok {
		jwtClaims.UserID = userID
	} else {
		return nil, ErrInvalidToken
	}

	if email, ok := claims["email"].(string); ok {
		jwtClaims.Email = email
	}

	if role, ok := claims["role"].(string); ok {
		jwtClaims.Role = role
	}

	if exp, ok := claims["exp"].(float64); ok {
		jwtClaims.ExpiresAt = int64(exp)
	} else {
		return nil, ErrInvalidToken
	}

	if iat, ok := claims["iat"].(float64); ok {
		jwtClaims.IssuedAt = int64(iat)
	}

	if nbf, ok := claims["nbf"].(float64); ok {
		jwtClaims.NotBefore = int64(nbf)
	}

	if iss, ok := claims["iss"].(string); ok {
		jwtClaims.Issuer = iss
	}

	if aud, ok := claims["aud"].(string); ok {
		jwtClaims.Audience = aud
	}

	if sub, ok := claims["sub"].(string); ok {
		jwtClaims.Subject = sub
	}

	return jwtClaims, nil
}
