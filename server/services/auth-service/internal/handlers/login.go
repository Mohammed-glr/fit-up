package handlers

import (
	"net/http"

	"github.com/tdmdh/lornian-backend/services/auth-service/internal/middleware"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/service"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/utils"
)

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	clientIP := middleware.GetClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	var u *types.User
	var err error
	
	if utils.IsEmailFormat(payload.Identifier) {
		u, err = h.store.GetUserByEmail(r.Context(), payload.Identifier)
	} else {
		u, err = h.store.GetUserByUsername(r.Context(), payload.Identifier)
	}
	
	if err != nil {
		if h.auditLogger != nil {
			h.auditLogger.LogFailedLoginAttempt(r.Context(), payload.Identifier, clientIP, userAgent, "user_not_found")
		}
		utils.WriteError(w, http.StatusUnauthorized, types.ErrInvalidCredentials)
		return
	}

	if !service.ComparePasswords(u.PasswordHash, []byte(payload.Password)) {
		if h.auditLogger != nil {
			h.auditLogger.LogFailedLoginAttempt(r.Context(), payload.Identifier, clientIP, userAgent, "invalid_password")
		}
		utils.WriteError(w, http.StatusUnauthorized, types.ErrInvalidCredentials)
		return
	}

	tokenPair, err := h.authService.GenerateTokenPair(r.Context(), u)
	if err != nil {
		if h.auditLogger != nil {
			h.auditLogger.LogLogin(r.Context(), u.ID, u.Email, clientIP, userAgent, false, map[string]interface{}{
				"error": err.Error(),
			})
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if h.auditLogger != nil {
		h.auditLogger.LogLogin(r.Context(), u.ID, u.Email, clientIP, userAgent, true, nil)
	}

	response := types.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    tokenPair.ExpiresIn,
		User:         u,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
