package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/utils"
)

func (h *AuthHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		utils.WriteError(w, http.StatusBadRequest, types.ErrInvalidInput)
		return
	}

	user, err := h.authService.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == pgx.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, types.ErrUserNotFound)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var imagePtr *string
	if user.Image != "" {
		imagePtr = &user.Image
	}

	publicUser := types.PublicUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Bio:       user.Bio,
		Image:     imagePtr,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, publicUser)
}
