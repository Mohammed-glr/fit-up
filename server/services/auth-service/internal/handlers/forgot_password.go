package handlers

import (
	"net/http"
	"time"

	"github.com/tdmdh/lornian-backend/services/auth-service/internal/service"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/utils"
)

func (h *AuthHandler) handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var payload types.ForgotPasswordRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.store.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {

		utils.WriteSuccess(w, http.StatusOK, "If an account with that email exists, we have sent a password reset link")
		return
	}

	resetToken, err := service.CreatePasswordResetToken(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, types.ErrInternalServerError)
		return
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	if err := h.store.CreatePasswordResetToken(r.Context(), payload.Email, resetToken.Token, expiresAt); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, types.ErrInternalServerError)
		return
	}

	if err := service.SendPasswordResetEmail(user.Email, resetToken.Token); err != nil {
		utils.WriteSuccess(w, http.StatusOK, "If an account with that email exists, we have sent a password reset link")
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "If an account with that email exists, we have sent a password reset link")
}
