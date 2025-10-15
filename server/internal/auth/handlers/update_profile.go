package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tdmdh/fit-up-server/internal/auth/middleware"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
)

func (h *AuthHandler) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
		return
	}

	var req types.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Failed to decode JSON body: %v", err)
		utils.WriteError(w, http.StatusBadRequest, types.ErrInvalidInput)
		return
	}

	log.Printf("📝 Received update request for user %s", userID)
	if req.Name != nil {
		log.Printf("✏️ Name: %s", *req.Name)
	}
	if req.Bio != nil {
		log.Printf("📄 Bio: %s", *req.Bio)
	}
	if req.Image != nil {
		log.Printf("📸 Image: base64 string (length: %d)", len(*req.Image))
	}

	if err := h.store.UpdateUser(r.Context(), userID, &req); err != nil {
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
