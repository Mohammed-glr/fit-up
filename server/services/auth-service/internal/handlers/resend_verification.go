package handlers

import (
	"net/http"

	"github.com/tdmdh/fit-up-server/services/auth-service/internal/types"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/utils"
)



func (h *AuthHandler) handleResendVerification(w http.ResponseWriter, r *http.Request) {
	var payload types.ResendVerificationRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.verify.ResendVerificationEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case types.ErrUserNotFound:
			utils.WriteError(w, http.StatusNotFound, types.ErrUserNotFound)
		default:
			utils.WriteError(w, http.StatusInternalServerError, types.ErrFailedToResendVerification)
		}
		return
	}

	response := map[string]string{
		"message": "Verification email resent successfully",
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
