package handlers

import (
	"net/http"

	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
)

func (h *AuthHandler) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload types.ResetPasswordRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.authService.ResetPassword(r.Context(), payload); err != nil {
		switch err {
		case types.ErrInvalidToken, types.ErrTokenExpired:
			utils.WriteError(w, http.StatusBadRequest, err)
		case types.ErrUserNotFound:
			utils.WriteError(w, http.StatusNotFound, err)
		case types.ErrPasswordTooWeak:
			utils.WriteError(w, http.StatusBadRequest, err)
		default:
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "Password reset successfully")
}
