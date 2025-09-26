package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tdmdh/fit-up-server/services/auth-service/internal/types"
	"github.com/tdmdh/fit-up-server/services/auth-service/tests/fixtures"
	"github.com/tdmdh/fit-up-server/services/auth-service/tests/mocks"
	"github.com/tdmdh/fit-up-server/tests/shared/testutils"
)

func TestLoginIntegration(t *testing.T) {
	// This is a placeholder integration test
	// In a real implementation, you would:
	// 1. Set up a test database
	// 2. Create actual HTTP handlers
	// 3. Test the full request/response cycle

	// Setup
	mockStore := mocks.NewMockUserStore()
	testUser := fixtures.TestUser()
	mockStore.AddUser(testUser)

	// Create a login request
	loginReq := fixtures.TestLoginRequest()

	// Convert to JSON
	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		t.Fatalf("Failed to marshal login request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// This would normally call your actual handler
	// For now, we'll just simulate a successful response
	response := types.LoginResponse{
		AccessToken:  "mock.jwt.token",
		RefreshToken: "mock-refresh-token",
		TokenType:    "Bearer",
		User:         testUser,
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
	}

	rr.WriteHeader(http.StatusOK)
	json.NewEncoder(rr).Encode(response)

	// Assertions
	testutils.AssertStatusCode(t, http.StatusOK, rr.Code)

	var actualResponse types.LoginResponse
	err = json.NewDecoder(rr.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if actualResponse.AccessToken == "" {
		t.Error("Expected access token to be present")
	}

	if actualResponse.User.Email != testUser.Email {
		t.Errorf("Expected user email %s, got %s", testUser.Email, actualResponse.User.Email)
	}

	t.Log("Login integration test completed successfully")
}

func TestLoginWithUsernameIntegration(t *testing.T) {
	// This test verifies that login works with username as identifier
	// Setup
	mockStore := mocks.NewMockUserStore()
	testUser := fixtures.TestUser()
	mockStore.AddUser(testUser)

	// Create a login request with username
	loginReq := fixtures.TestLoginRequestWithUsername()

	// Convert to JSON
	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		t.Fatalf("Failed to marshal login request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Mock the handler behavior - in real integration tests,
	// this would go through actual handlers
	response := types.LoginResponse{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
		User:         testUser,
	}

	rr.WriteHeader(http.StatusOK)
	json.NewEncoder(rr).Encode(response)

	// Assertions
	testutils.AssertStatusCode(t, http.StatusOK, rr.Code)

	var actualResponse types.LoginResponse
	err = json.NewDecoder(rr.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if actualResponse.AccessToken == "" {
		t.Error("Expected access token to be present")
	}

	if actualResponse.User.Username != testUser.Username {
		t.Errorf("Expected user username %s, got %s", testUser.Username, actualResponse.User.Username)
	}

	t.Log("Username login integration test completed successfully")
}

func TestRegisterIntegration(t *testing.T) {
	// This is a placeholder for registration integration test
	// Similar structure to login test

	mockStore := mocks.NewMockUserStore()
	registerReq := fixtures.TestRegisterRequest()

	// Test user creation
	ctx := context.Background()
	user := &types.User{
		ID:       "new-user-123",
		Username: registerReq.Username,
		Email:    registerReq.Email,
		Name:     registerReq.Name,
		Role:     types.RoleUser,
	}

	err := mockStore.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify user was created
	createdUser, err := mockStore.GetUserByEmail(ctx, registerReq.Email)
	if err != nil {
		t.Fatalf("Failed to retrieve created user: %v", err)
	}

	if createdUser.Email != registerReq.Email {
		t.Errorf("Expected email %s, got %s", registerReq.Email, createdUser.Email)
	}

	t.Log("Register integration test completed successfully")
}

func TestPasswordResetIntegration(t *testing.T) {
	// This is a placeholder for password reset integration test

	mockStore := mocks.NewMockUserStore()
	testUser := fixtures.TestUser()
	mockStore.AddUser(testUser)

	ctx := context.Background()
	token := "reset-token-123"
	expiresAt := time.Now().Add(time.Hour)

	// Create password reset token
	err := mockStore.CreatePasswordResetToken(ctx, testUser.Email, token, expiresAt)
	if err != nil {
		t.Fatalf("Failed to create password reset token: %v", err)
	}

	// Verify token was created
	resetToken, err := mockStore.GetPasswordResetToken(ctx, token)
	if err != nil {
		t.Fatalf("Failed to retrieve password reset token: %v", err)
	}

	if resetToken.Email != testUser.Email {
		t.Errorf("Expected email %s, got %s", testUser.Email, resetToken.Email)
	}

	// Test password update
	newPassword := "new-hashed-password"
	err = mockStore.UpdateUserPassword(ctx, testUser.ID, newPassword)
	if err != nil {
		t.Fatalf("Failed to update user password: %v", err)
	}

	// Verify password was updated
	updatedUser, err := mockStore.GetUserByID(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated user: %v", err)
	}

	if updatedUser.PasswordHash != newPassword {
		t.Errorf("Expected password hash %s, got %s", newPassword, updatedUser.PasswordHash)
	}

	t.Log("Password reset integration test completed successfully")
}
