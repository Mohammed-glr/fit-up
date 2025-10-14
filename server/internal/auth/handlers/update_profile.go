package handlers

import (
	"net/http"

	"github.com/tdmdh/fit-up-server/internal/auth/middleware"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
	sharedUtils "github.com/tdmdh/fit-up-server/shared/utils"
)

func (h *AuthHandler) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, types.ErrInvalidInput)
		return
	}

	var req types.UpdateUserRequest

	if name := r.FormValue("name"); name != "" {
		req.Name = &name
	}
	if bio := r.FormValue("bio"); bio != "" {
		req.Bio = &bio
	}

	file, header, err := sharedUtils.GetFormFile(r, "image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if file != nil {
		defer file.Close()

		imageURL, err := sharedUtils.HandleFileUpload(file, header, sharedUtils.FileUploadConfig{
			UploadDir:      "./uploads/profiles",
			FileNamePrefix: userID,
		})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		req.Image = &imageURL

		user, _ := h.authService.GetUser(r.Context(), userID)
		if user != nil && user.Image != "" {
			sharedUtils.DeleteFile(user.Image)
		}
	}

	if err = h.store.UpdateUser(r.Context(), userID, &req); err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := h.authService.GetUser(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var imagePtr *string
	if user.Image != "" {
		imagePtr = &user.Image
	}

	response := types.UserResponse{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Bio:                user.Bio,
		Email:              user.Email,
		Image:              imagePtr,
		Role:               user.Role,
		IsTwoFactorEnabled: user.IsTwoFactorEnabled,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Profile updated successfully",
		"user":    response,
	})
}
