package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	authRepo "github.com/tdmdh/fit-up-server/internal/auth/repository"
	authService "github.com/tdmdh/fit-up-server/internal/auth/services"
	authTypes "github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
	"github.com/tdmdh/fit-up-server/shared/config"
)

type contextKey string

const (
	UserIDKey     contextKey = "userID"
	UserRoleKey   contextKey = "userRole"
	AuthIDKey     contextKey = "authID"
	UserClaimsKey contextKey = "userClaims"
)

type AuthMiddleware struct {
	repo      repository.SchemaRepo
	userStore authRepo.UserStore
}

func NewAuthMiddleware(repo repository.SchemaRepo, userStore authRepo.UserStore) *AuthMiddleware {
	return &AuthMiddleware{
		repo:      repo,
		userStore: userStore,
	}
}
func (am *AuthMiddleware) ExtractUserID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUserID, ok := getAuthUserID(r)
			if !ok {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: Missing user authentication",
				})
				return
			}

			ctx := context.WithValue(r.Context(), AuthIDKey, authUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (am *AuthMiddleware) RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUserID, ok := getAuthUserID(r)
			if !ok {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: Missing user authentication",
				})
				return
			}

			userRole, err := am.repo.UserRoles().GetUserRole(r.Context(), authUserID)
			if err != nil {
				log.Printf("Error fetching user role for %s: %v", authUserID, err)
				userRole = types.RoleUser
			}

			ctx := context.WithValue(r.Context(), AuthIDKey, authUserID)
			ctx = context.WithValue(ctx, UserRoleKey, userRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (am *AuthMiddleware) RequireJWTAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := extractTokenFromHeader(r)
			if err != nil {
				respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"error": "Unauthorized: Missing or invalid authorization header",
					"code":  "UNAUTHORIZED",
				})
				return
			}

			claims, err := am.validateJWTToken(tokenString)
			if err != nil {
				status := http.StatusUnauthorized
				message := "Invalid token"

				if err == authTypes.ErrTokenExpired {
					message = "Token has expired"
				} else if err == authTypes.ErrJWTSecretNotSet {
					status = http.StatusInternalServerError
					message = "Server configuration error"
					log.Printf("JWT secret not set")
				}

				respondJSON(w, status, map[string]interface{}{
					"error": message,
					"code":  "INVALID_TOKEN",
				})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserClaimsKey, claims)
			ctx = context.WithValue(ctx, AuthIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, string(claims.Role))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (am *AuthMiddleware) OptionalJWTAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := extractTokenFromHeader(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := am.validateJWTToken(tokenString)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserClaimsKey, claims)
			ctx = context.WithValue(ctx, AuthIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, string(claims.Role))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (am *AuthMiddleware) RequireCoachRole() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUserID := r.Header.Get("X-User-ID")
			if authUserID == "" {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: Missing user authentication",
				})
				return
			}

			userRole, err := am.repo.UserRoles().GetUserRole(r.Context(), authUserID)
			if err != nil {
				log.Printf("Error fetching user role for %s: %v", authUserID, err)
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "Failed to verify user role",
				})
				return
			}

			if userRole != types.RoleCoach && userRole != types.RoleAdmin {
				respondJSON(w, http.StatusForbidden, map[string]string{
					"error": "Forbidden: Coach role required",
				})
				return
			}

			ctx := context.WithValue(r.Context(), AuthIDKey, authUserID)
			ctx = context.WithValue(ctx, UserRoleKey, userRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (am *AuthMiddleware) RequireAdminRole() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUserID := r.Header.Get("X-User-ID")
			if authUserID == "" {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: Missing user authentication",
				})
				return
			}

			userRole, err := am.repo.UserRoles().GetUserRole(r.Context(), authUserID)
			if err != nil {
				log.Printf("Error fetching user role for %s: %v", authUserID, err)
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "Failed to verify user role",
				})
				return
			}

			if userRole != types.RoleAdmin {
				respondJSON(w, http.StatusForbidden, map[string]string{
					"error": "Forbidden: Admin role required",
				})
				return
			}

			ctx := context.WithValue(r.Context(), AuthIDKey, authUserID)
			ctx = context.WithValue(ctx, UserRoleKey, userRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (am *AuthMiddleware) ValidateCoachAssignment() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUserID := GetAuthIDFromContext(r.Context())
			if authUserID == "" {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized",
				})
				return
			}

			userIDStr := chi.URLParam(r, "user_id")
			if userIDStr == "" {
				respondJSON(w, http.StatusBadRequest, map[string]string{
					"error": "Missing user_id parameter",
				})
				return
			}

			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				respondJSON(w, http.StatusBadRequest, map[string]string{
					"error": "Invalid user_id format",
				})
				return
			}

			isCoach, err := am.repo.CoachAssignments().IsCoachForUser(r.Context(), authUserID, userID)
			if err != nil {
				log.Printf("Error checking coach assignment: %v", err)
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "Failed to verify coach assignment",
				})
				return
			}

			if !isCoach {
				userRole, _ := am.repo.UserRoles().GetUserRole(r.Context(), authUserID)
				if userRole != types.RoleAdmin {
					respondJSON(w, http.StatusForbidden, map[string]string{
						"error": "Forbidden: Not assigned as coach to this client",
					})
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (am *AuthMiddleware) ValidateResourceOwnership() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUserID := GetAuthIDFromContext(r.Context())
			if authUserID == "" {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized",
				})
				return
			}

			userIDStr := chi.URLParam(r, "user_id")
			if userIDStr == "" {
				userID := GetUserIDFromContext(r.Context())
				if userID == 0 {
					respondJSON(w, http.StatusBadRequest, map[string]string{
						"error": "Missing user_id",
					})
					return
				}
			}

			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				respondJSON(w, http.StatusBadRequest, map[string]string{
					"error": "Invalid user_id format",
				})
				return
			}

			userRole, err := am.repo.UserRoles().GetUserRole(r.Context(), authUserID)
			if err != nil {
				log.Printf("Error fetching user role: %v", err)
				userRole = types.RoleUser
			}

			if userRole == types.RoleAdmin {
				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if userRole == types.RoleCoach {
				isCoach, err := am.repo.CoachAssignments().IsCoachForUser(r.Context(), authUserID, userID)
				if err == nil && isCoach {
					ctx := context.WithValue(r.Context(), UserIDKey, userID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			// Check of gebruiker eigenaar is (auth_user_id match)
			// TODO: Vergelijk authUserID met resource owner
			// Voor nu: allow als het de user's eigen resource is

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthIDFromContext(ctx context.Context) string {
	if authID, ok := ctx.Value(AuthIDKey).(string); ok {
		return authID
	}
	return ""
}

func GetUserRoleFromContext(ctx context.Context) types.UserRole {
	if role, ok := ctx.Value(UserRoleKey).(types.UserRole); ok {
		return role
	}
	return types.RoleUser
}

func GetUserIDFromContext(ctx context.Context) int {
	if userID, ok := ctx.Value(UserIDKey).(int); ok {
		return userID
	}
	return 0
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func getAuthUserID(r *http.Request) (string, bool) {
	id := r.Header.Get("X-User-ID")
	return id, id != ""
}

func extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", authTypes.ErrUnauthorized
	}

	if len(authHeader) > 7 && strings.ToUpper(authHeader[:7]) == "BEARER " {
		return authHeader[7:], nil
	}

	return "", authTypes.ErrInvalidToken
}

func (am *AuthMiddleware) validateJWTToken(tokenString string) (*authTypes.TokenClaims, error) {
	secret := []byte(config.NewConfig().JWTSecret)
	if len(secret) == 0 {
		return nil, authTypes.ErrJWTSecretNotSet
	}

	claims, err := authService.ValidateJWT(tokenString, am.userStore, secret)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func GetUserClaimsFromContext(ctx context.Context) (*authTypes.TokenClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*authTypes.TokenClaims)
	return claims, ok
}

func GetAuthUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(AuthIDKey).(string)
	return userID, ok
}
