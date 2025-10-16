package handlers

import (
	"net/http"

	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
)

func (h *AuthHandler) handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	var payload types.VerifyEmailRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.authService.VerifyEmail(r.Context(), payload.Token)
	if err != nil {
		switch err {
		case types.ErrVerificationTokenNotFound, types.ErrVerificationTokenExpired:
			utils.WriteError(w, http.StatusBadRequest, err)
		default:
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	tokenPair, err := h.authService.GenerateTokenPair(r.Context(), user)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Email verified successfully",
			"user":    user,
		})
		return
	}

	response := types.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         user,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) handleResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	var payload types.ResendVerificationRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.authService.ResendEmailVerification(r.Context(), payload.Email); err != nil {
		switch err {
		case types.ErrEmailNotFound:
			utils.WriteSuccess(w, http.StatusOK, "Verification email sent if the account exists")
			return
		case types.ErrEmailAlreadyVerified:
			utils.WriteSuccess(w, http.StatusOK, "Email is already verified")
			return
		default:
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}

	utils.WriteSuccess(w, http.StatusOK, "Verification email sent if the account exists")
}
