package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UserRoleKey contextKey = "userRole"
	AuthIDKey   contextKey = "authID"
)

type AuthMiddleware struct {
	repo repository.SchemaRepo
}

func NewAuthMiddleware(repo repository.SchemaRepo) *AuthMiddleware {
	return &AuthMiddleware{
		repo: repo,
	}
}

// ExtractUserID haalt de X-User-ID header op en voegt deze toe aan context
func (am *AuthMiddleware) ExtractUserID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUserID := r.Header.Get("X-User-ID")
			if authUserID == "" {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Missing X-User-ID header",
				})
				return
			}

			ctx := context.WithValue(r.Context(), AuthIDKey, authUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth verplicht authenticatie en laadt user role uit cache
func (am *AuthMiddleware) RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUserID := r.Header.Get("X-User-ID")
			if authUserID == "" {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: Missing user authentication",
				})
				return
			}

			// Haal user role op uit cache
			userRole, err := am.repo.UserRoles().GetUserRole(r.Context(), authUserID)
			if err != nil {
				log.Printf("Error fetching user role for %s: %v", authUserID, err)
				// Default naar user role als er geen cache is
				userRole = types.RoleUser
			}

			ctx := context.WithValue(r.Context(), AuthIDKey, authUserID)
			ctx = context.WithValue(ctx, UserRoleKey, userRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireCoachRole verplicht coach of admin role
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

			// Haal user role op uit cache
			userRole, err := am.repo.UserRoles().GetUserRole(r.Context(), authUserID)
			if err != nil {
				log.Printf("Error fetching user role for %s: %v", authUserID, err)
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "Failed to verify user role",
				})
				return
			}

			// Check of gebruiker coach of admin is
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

// RequireAdminRole verplicht admin role
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

			// Haal user role op uit cache
			userRole, err := am.repo.UserRoles().GetUserRole(r.Context(), authUserID)
			if err != nil {
				log.Printf("Error fetching user role for %s: %v", authUserID, err)
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "Failed to verify user role",
				})
				return
			}

			// Check of gebruiker admin is
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

// ValidateCoachAssignment controleert of de coach toegang heeft tot de client
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

			// Haal user_id uit URL parameter
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

			// Check of coach geassigneerd is aan deze user
			isCoach, err := am.repo.CoachAssignments().IsCoachForUser(r.Context(), authUserID, userID)
			if err != nil {
				log.Printf("Error checking coach assignment: %v", err)
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "Failed to verify coach assignment",
				})
				return
			}

			if !isCoach {
				// Check of gebruiker admin is (admins hebben altijd toegang)
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

// ValidateResourceOwnership controleert of gebruiker eigenaar is of coach/admin
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

			// Haal user_id uit URL of body
			userIDStr := chi.URLParam(r, "user_id")
			if userIDStr == "" {
				// Als niet in URL, probeer uit context
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

			// Haal user role op
			userRole, err := am.repo.UserRoles().GetUserRole(r.Context(), authUserID)
			if err != nil {
				log.Printf("Error fetching user role: %v", err)
				userRole = types.RoleUser
			}

			// Admin heeft altijd toegang
			if userRole == types.RoleAdmin {
				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Check of coach geassigneerd is
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
