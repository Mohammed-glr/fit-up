package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tdmdh/fit-up-server/internal/auth/middleware"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
)

func (h *AuthHandler) handleUpdateRole(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by JWT middleware)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
		return
	}

	// Parse request body
	var req types.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, types.ErrInvalidInput)
		return
	}

	// Validate role
	if req.Role != types.RoleUser && req.Role != types.RoleCoach {
		utils.WriteError(w, http.StatusBadRequest, types.ErrInvalidInput)
		return
	}

	// Update user role
	err := h.authService.UpdateUserRole(r.Context(), userID, req.Role)
	if err != nil {
		if err == types.ErrUserNotFound {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Fetch updated user
	user, err := h.authService.GetUser(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Prepare response
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
		"message": "Role updated successfully",
		"user":    response,
	})
}
