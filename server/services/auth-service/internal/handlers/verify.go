package handlers

import (
	"net/http"

	"github.com/tdmdh/fit-up-server/services/auth-service/internal/types"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/utils"
)

func (h *AuthHandler) handleVerify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.WriteError(w, http.StatusBadRequest, types.ErrVerificationTokenNotFound)
		return
	}

	user, err := h.verify.GetUserByVerificationToken(r.Context(), token)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, types.ErrVerificationTokenExpired)
		return
	}

	err = h.verify.UpdateUserVerificationStatus(r.Context(), user.ID, true)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, types.ErrFailedToVerifyEmail)
		return
	}

	response := types.VerifyEmailResponse{
		Message: "Email verified successfully",
		User:    user,
	}
	err = h.verify.DeleteVerificationToken(r.Context(), user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, types.ErrFailedToDeleteVerificationT)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
