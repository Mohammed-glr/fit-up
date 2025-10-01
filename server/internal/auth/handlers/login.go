package handlers

import (
	"net/http"

	"github.com/tdmdh/fit-up-server/internal/auth/services"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
)

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var u *types.User
	var err error

	if utils.IsEmailFormat(payload.Identifier) {
		u, err = h.store.GetUserByEmail(r.Context(), payload.Identifier)
	} else {
		u, err = h.store.GetUserByUsername(r.Context(), payload.Identifier)
	}

	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrInvalidCredentials)
		return
	}

	if !service.ComparePasswords(u.PasswordHash, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrInvalidCredentials)
		return
	}

	tokenPair, err := h.authService.GenerateTokenPair(r.Context(), u)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
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
