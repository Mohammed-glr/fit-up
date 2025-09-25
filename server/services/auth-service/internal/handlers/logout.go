package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tdmdh/lornian-backend/services/auth-service/internal/middleware"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/service"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/utils"
	"github.com/tdmdh/lornian-backend/shared/config"
)

func (h *AuthHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
		return
	}

	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
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

	claims, err := service.ValidateJWT(tokenString, h.store, secret)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrInvalidToken)
		return
	}

	clientIP := middleware.GetClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	fmt.Printf("DEBUG: Starting logout for user %s with JTI %s\n", claims.UserID, claims.JTI)

	expiresAt := time.Unix(claims.ExpiresAt, 0)
	err = h.authService.LogoutWithToken(r.Context(), claims.UserID, claims.JTI, expiresAt)
	if err != nil {
		fmt.Printf("ERROR: LogoutWithToken failed for user %s: %v\n", claims.UserID, err)
		if h.auditLogger != nil {
			h.auditLogger.LogLogout(r.Context(), claims.UserID, clientIP, userAgent, false, map[string]interface{}{
				"error": err.Error(),
			})
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	fmt.Printf("DEBUG: Token blacklisted successfully, now revoking refresh tokens\n")

	err = h.store.RevokeAllUserRefreshTokens(r.Context(), claims.UserID)
	if err != nil {
		fmt.Printf("ERROR: RevokeAllUserRefreshTokens failed for user %s: %v\n", claims.UserID, err)
		if h.auditLogger != nil {
			h.auditLogger.LogLogout(r.Context(), claims.UserID, clientIP, userAgent, false, map[string]interface{}{
				"error": "failed to revoke refresh tokens",
			})
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	fmt.Printf("DEBUG: Logout completed successfully for user %s\n", claims.UserID)

	if h.auditLogger != nil {
		h.auditLogger.LogLogout(r.Context(), claims.UserID, clientIP, userAgent, true, nil)
	}

	response := map[string]interface{}{
		"message": "Logged out successfully",
		"user_id": claims.UserID,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
