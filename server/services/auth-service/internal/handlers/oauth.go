package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/interfaces"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/middleware"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/utils"
)

func (h *AuthHandler) handleOAuthAuthorize(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	var req struct {
		RedirectURI string `json:"redirect_uri"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON: %w", err))
		return
	}

	if req.RedirectURI == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("redirect_uri is required"))
		return
	}

	authURL, err := h.oauthService.GetAuthorizationURL(r.Context(), provider, req.RedirectURI)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"redirect_url": authURL,
	})
}

func (h *AuthHandler) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	userInfo, err := h.oauthService.HandleCallback(r.Context(), provider, code, state)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.createOrGetOAuthUser(r.Context(), userInfo, provider)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tokenPair, err := h.authService.GenerateTokenPair(r.Context(), user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    tokenPair.ExpiresIn,
		User:         user,
	})
}

func (h *AuthHandler) createOrGetOAuthUser(ctx context.Context, userInfo *types.OAuthUserInfo, provider string) (*types.User, error) {
	if oauthStore, ok := h.store.(interfaces.OAuthStore); ok {
		account, err := oauthStore.GetAccountByProvider(ctx, provider, userInfo.ID)
		if err == nil {
			return h.store.GetUserByID(ctx, account.UserID)
		}
	}

	existingUser, err := h.store.GetUserByEmail(ctx, userInfo.Email)
	if err == nil {
		if oauthStore, ok := h.store.(interfaces.OAuthStore); ok {
			account := &types.Account{
				UserID:            existingUser.ID,
				Type:              "oauth",
				Provider:          provider,
				ProviderAccountID: userInfo.ID,
			}
			_ = oauthStore.CreateAccount(ctx, account)
		}
		return existingUser, nil
	}

	now := time.Now()
	newUser := &types.User{
		Username:      userInfo.Username,
		Name:          userInfo.Name,
		Email:         userInfo.Email,
		EmailVerified: &now,
		Image:         userInfo.AvatarURL,
		Role:          types.RoleUser,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if newUser.Username == "" {
		newUser.Username = strings.Split(userInfo.Email, "@")[0]
	}

	err = h.store.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if oauthStore, ok := h.store.(interfaces.OAuthStore); ok {
		account := &types.Account{
			UserID:            newUser.ID,
			Type:              "oauth",
			Provider:          provider,
			ProviderAccountID: userInfo.ID,
		}
		_ = oauthStore.CreateAccount(ctx, account)
	}

	return newUser, nil
}

func (h *AuthHandler) handleLinkAccount(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	var payload types.LinkAccountRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	claims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
		return
	}

	userInfo, err := h.oauthService.HandleCallback(r.Context(), provider, payload.Code, payload.State)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.oauthService.LinkAccount(r.Context(), claims.UserID, provider, userInfo)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	clientIP := middleware.GetClientIP(r)
	userAgent := r.Header.Get("User-Agent")
	if h.auditLogger != nil {
		h.auditLogger.LogSecurityEvent(r.Context(), claims.UserID, "account_link", clientIP, userAgent, true, map[string]interface{}{
			"provider": provider,
			"email":    userInfo.Email,
		})
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "Account linked successfully",
		"provider": provider,
		"email":    userInfo.Email,
	})
}

func (h *AuthHandler) handleUnlinkAccount(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	claims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
		return
	}

	err := h.oauthService.UnlinkAccount(r.Context(), claims.UserID, provider)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	clientIP := middleware.GetClientIP(r)
	userAgent := r.Header.Get("User-Agent")
	if h.auditLogger != nil {
		h.auditLogger.LogSecurityEvent(r.Context(), claims.UserID, "account_unlink", clientIP, userAgent, true, map[string]interface{}{
			"provider": provider,
		})
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "Account unlinked successfully",
		"provider": provider,
	})
}

func (h *AuthHandler) handleGetLinkedAccounts(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, types.ErrUnauthorized)
		return
	}

	accounts, err := h.oauthService.GetLinkedAccounts(r.Context(), claims.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	safeAccounts := make([]map[string]interface{}, len(accounts))
	for i, account := range accounts {
		safeAccounts[i] = map[string]interface{}{
			"provider":            account.Provider,
			"provider_account_id": account.ProviderAccountID,
			"type":                account.Type,
		}
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"linked_accounts": safeAccounts,
	})
}
