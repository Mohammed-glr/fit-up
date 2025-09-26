package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/service"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/types"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/utils"
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

	// Check for existing email
	existingUser, err := h.store.GetUserByEmail(r.Context(), payload.Email)
	if err != nil && err != pgx.ErrNoRows {
		log.Printf("Database error while checking existing email: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if existingUser != nil {
		log.Printf("User with email %s already exists", payload.Email)
		utils.WriteError(w, http.StatusConflict, types.ErrUserAlreadyExists)
		return
	}

	// Check for existing username
	existingUserByUsername, err := h.store.GetUserByUsername(r.Context(), payload.Username)
	if err != nil && err != pgx.ErrNoRows {
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

	user := &types.User{
		ID:                 uuid.New().String(),
		Username:           payload.Username,
		Email:              payload.Email,
		Name:               payload.Name,
		PasswordHash:       hashedPassword,
		IsTwoFactorEnabled: false,
		Role:               types.RoleUser,
		EmailVerified:      nil, // Will be set when user verifies email
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	log.Printf("Creating user in database for email: %s", user.Email)

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

	log.Printf("User created successfully, generating verification token")

	token := uuid.New().String()

	log.Printf("Saving verification token to database for user: %s", user.ID)

	expiresAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	err = h.verify.CreateVerificationToken(r.Context(), user.ID, token, expiresAt)
	if err != nil {
		log.Printf("Error saving verification token: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Sending verification email to: %s", user.Email)

	if err := service.SendVerificationEmail(user.Email, token); err != nil {
		log.Printf("Warning: Failed to send verification email to %s: %v", user.Email, err)
	}

	log.Printf("Registration completed successfully for user: %s", user.Email)

	response := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"name":     user.Name,
		"message":  "User registered successfully",
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}
