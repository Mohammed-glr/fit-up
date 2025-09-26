package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/tdmdh/fit-up-server/services/auth-service/internal/interfaces"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/service"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/types"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/utils"
	"github.com/tdmdh/fit-up-server/shared/config"
)

type ContextKey string

const (
	UserIDKey     ContextKey = "user_id"
	UserClaimsKey ContextKey = "user_claims"
	UserKey       ContextKey = "user"
)

func JWTAuthMiddleware(store interfaces.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
				return
			}
			tokenString := ""
			if len(authHeader) > 7 && strings.ToUpper(authHeader[:7]) == "BEARER " {
				tokenString = authHeader[7:]
			} else {
				utils.WriteError(w, http.StatusUnauthorized, types.ErrInvalidToken)
				return
			}
			secret := []byte(config.NewConfig().JWTSecret)
			if len(secret) == 0 {
				utils.WriteError(w, http.StatusInternalServerError, types.ErrJWTSecretNotSet)
				return
			}
			claims, err := service.ValidateJWT(tokenString, store, secret)
			if err != nil {
				if err == types.ErrTokenExpired {
					utils.WriteError(w, http.StatusUnauthorized, types.ErrTokenExpired)
				} else {
					utils.WriteError(w, http.StatusUnauthorized, types.ErrInvalidToken)
				}
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalJWTAuthMiddleware(store interfaces.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}
			tokenString := ""
			if len(authHeader) > 7 && strings.ToUpper(authHeader[:7]) == "BEARER " {
				tokenString = authHeader[7:]
			} else {
				next.ServeHTTP(w, r)
				return
			}
			secret := []byte(config.NewConfig().JWTSecret)
			if len(secret) == 0 {
				next.ServeHTTP(w, r)
				return
			}
			claims, err := service.ValidateJWT(tokenString, store, secret)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRoleMiddleware(requiredRole types.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserClaimsKey).(*types.TokenClaims)
			if !ok {
				utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
				return
			}
			if claims.Role != requiredRole {
				utils.WriteError(w, http.StatusForbidden, types.ErrInsufficientPermissions)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAdminMiddleware() func(http.Handler) http.Handler {
	return RequireRoleMiddleware(types.RoleAdmin)
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

func GetUserClaimsFromContext(ctx context.Context) (*types.TokenClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*types.TokenClaims)
	return claims, ok
}

func GetUserFromContext(ctx context.Context) (*types.User, bool) {
	user, ok := ctx.Value(UserKey).(*types.User)
	return user, ok
}
