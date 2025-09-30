package handlers

import (
	"net/http"

	"github.com/tdmdh/fit-up-server/services/auth-service/internal/middleware"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/types"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/utils"
)

func (h *AuthHandler) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	var payload types.ChangePasswordRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
		return
	}

	err := h.authService.ChangePassword(r.Context(), userID, payload.CurrentPassword, payload.NewPassword)
	if err != nil {
		switch err {
		case types.ErrInvalidCredentials:
			utils.WriteError(w, http.StatusUnauthorized, types.ErrInvalidCredentials)
		case types.ErrSamePassword:
			utils.WriteError(w, http.StatusBadRequest, types.ErrSamePassword)
		case types.ErrPasswordTooWeak:
			utils.WriteError(w, http.StatusBadRequest, types.ErrPasswordTooWeak)
		case types.ErrUserNotFound:
			utils.WriteError(w, http.StatusNotFound, types.ErrUserNotFound)
		default:
			utils.WriteError(w, http.StatusInternalServerError, types.ErrInternalServerError)
		}
		return
	}

	response := map[string]any{
		"message": "Password changed successfully",
		"user_id": userID,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
