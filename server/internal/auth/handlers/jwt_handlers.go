package handlers

import (
	"net/http"

	"github.com/tdmdh/fit-up-server/internal/auth/services"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
	"github.com/tdmdh/fit-up-server/shared/config"
)

func (h *AuthHandler) handleValidateToken(w http.ResponseWriter, r *http.Request) {
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
		if err == types.ErrTokenExpired {
			utils.WriteError(w, http.StatusUnauthorized, types.ErrTokenExpired)
		} else {
			utils.WriteError(w, http.StatusUnauthorized, types.ErrInvalidToken)
		}
		return
	}

	user, err := h.store.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, types.ErrUserNotFound)
		return
	}

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
