package handlers

import (
	"log"
	"net/http"

	service "github.com/tdmdh/fit-up-server/internal/auth/services"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
	"github.com/tdmdh/fit-up-server/shared/config"
)

func (h *AuthHandler) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("Validate token: Handler called from %s", r.RemoteAddr)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Printf("Validate token: No authorization header")
		utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
		return
	}

	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		log.Printf("Validate token: Invalid authorization format")
		utils.WriteError(w, http.StatusUnauthorized, types.ErrInvalidToken)
		return
	}

	secret := []byte(config.NewConfig().JWTSecret)
	if len(secret) == 0 {
		log.Printf("Validate token: JWT secret not set")
		utils.WriteError(w, http.StatusInternalServerError, types.ErrJWTSecretNotSet)
		return
	}

	claims, err := service.ValidateJWT(tokenString, h.store, secret)
	if err != nil {
		log.Printf("Validate token: JWT validation error: %v", err)
		// Return error as-is with appropriate status code
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	log.Printf("Validate token: Claims validated for user ID: %s", claims.UserID)

	user, err := h.store.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		log.Printf("Validate token: Error fetching user by ID %s: %v", claims.UserID, err)
		// Return error as-is
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	log.Printf("Validate token: Successfully validated token for user: %s", user.Email)

	response := map[string]interface{}{
		"valid":   true,
		"user":    user,
		"claims":  claims,
		"message": "Token is valid",
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var payload types.RefreshTokenRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	tokenPair, err := h.authService.RotateTokens(r.Context(), payload.RefreshToken)
	if err != nil {
		if err == types.ErrRefreshTokenNotFound || err == types.ErrRefreshTokenExpired {
			utils.WriteError(w, http.StatusUnauthorized, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, tokenPair)
}
