package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	service "github.com/tdmdh/fit-up-server/internal/auth/services"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/auth/utils"
)

func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	log.Printf("Register request received from %s", r.RemoteAddr)

	var payload types.RegisterRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Parsed registration payload for email: %s", payload.Email)

	if err := utils.Validate.Struct(payload); err != nil {
		log.Printf("Validation error: %v", err)
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, errors)
		return
	}

	log.Printf("Validation passed, checking for existing user")

	existingUser, err := h.store.GetUserByEmail(r.Context(), payload.Email)
	if err != nil && !errors.Is(err, types.ErrUserNotFound) {
		log.Printf("Database error while checking existing email: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if existingUser != nil {
		log.Printf("User with email %s already exists", payload.Email)
		utils.WriteError(w, http.StatusConflict, types.ErrUserAlreadyExists)
		return
	}

	existingUserByUsername, err := h.store.GetUserByUsername(r.Context(), payload.Username)
	if err != nil && !errors.Is(err, types.ErrUserNotFound) {
		log.Printf("Database error while checking existing username: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if existingUserByUsername != nil {
		log.Printf("User with username %s already exists", payload.Username)
		utils.WriteError(w, http.StatusConflict, types.ErrUsernameAlreadyExists)
		return
	}

	log.Printf("No existing user found, hashing password")

	hashedPassword, err := service.HashPassword(payload.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Password hashed successfully, creating user object")

	role := types.RoleUser
	if payload.Role != "" {
		role = payload.Role
	}

	log.Printf("Assigned role: %s", role)

	user := &types.User{
		ID:                 uuid.New().String(),
		Username:           payload.Username,
		Email:              payload.Email,
		Name:               payload.Name,
		PasswordHash:       hashedPassword,
		IsTwoFactorEnabled: false,
		Role:               role,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	log.Printf("Creating user in database for email: %s with role: %s", user.Email, user.Role)

	if err := h.store.CreateUser(r.Context(), user); err != nil {
		log.Printf("Error creating user in database: %v", err)

		if strings.Contains(err.Error(), "users_email_key") {
			utils.WriteError(w, http.StatusConflict, types.ErrUserAlreadyExists)
			return
		}
		if strings.Contains(err.Error(), "users_username_key") {
			utils.WriteError(w, http.StatusConflict, types.ErrUsernameAlreadyExists)
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("User created successfully")

	if err := h.authService.InitiateEmailVerification(r.Context(), user); err != nil && err != types.ErrEmailAlreadyVerified {
		log.Printf("Failed to initiate email verification for %s: %v", user.Email, err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Registration completed successfully for user: %s", user.Email)

	response := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"name":     user.Name,
		"message":  "User registered successfully. Please verify your email to continue.",
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}
