package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tdmdh/lornian-backend/services/auth-service/internal/interfaces"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
	"github.com/tdmdh/lornian-backend/shared/config"
)

type OAuthService struct {
	store     interfaces.UserStore
	config    *config.Config
	providers map[string]*types.OAuthProvider
}

func NewOAuthService(store interfaces.UserStore, cfg *config.Config) *OAuthService {
	providers := map[string]*types.OAuthProvider{
		"google": {
			Name:         "google",
			ClientID:     cfg.OAuthConfig.GoogleClientID,
			ClientSecret: cfg.OAuthConfig.GoogleClientSecret,
			RedirectURI:  cfg.OAuthConfig.GoogleRedirectURI,
			AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:     "https://oauth2.googleapis.com/token",
			UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
			Scopes:       []string{"openid", "email", "profile"},
		},
		"github": {
			Name:         "github",
			ClientID:     cfg.OAuthConfig.GitHubClientID,
			ClientSecret: cfg.OAuthConfig.GitHubClientSecret,
			RedirectURI:  cfg.OAuthConfig.GitHubRedirectURI,
			AuthURL:      "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
			UserInfoURL:  "https://api.github.com/user",
			Scopes:       []string{"user:email"},
		},

	}

	return &OAuthService{
		store:     store,
		config:    cfg,
		providers: providers,
	}
}

func (s *OAuthService) GetAuthorizationURL(ctx context.Context, provider, redirectURL string) (string, error) {
	oauthProvider, exists := s.providers[provider]
	if !exists {
		return "", fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	if redirectURL == "" {
		return "", fmt.Errorf("redirect URL is required")
	}

	fmt.Printf("DEBUG - Provider: %s\n", provider)
	fmt.Printf("DEBUG - Client ID: %s\n", oauthProvider.ClientID)
	fmt.Printf("DEBUG - Redirect URL: %s\n", redirectURL)
	fmt.Printf("DEBUG - Auth URL: %s\n", oauthProvider.AuthURL)

	state, err := s.generateState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	oauthState := &types.OAuthState{
		State:       state,
		Provider:    provider,
		RedirectURL: redirectURL,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
		CreatedAt:   time.Now(),
	}

	if oauthStore, ok := s.store.(interfaces.OAuthStore); ok {
		err = oauthStore.CreateOAuthState(ctx, oauthState)
		if err != nil {
			return "", fmt.Errorf("failed to store OAuth state: %w", err)
		}
	}

	params := url.Values{}
	params.Add("client_id", oauthProvider.ClientID)
	params.Add("redirect_uri", redirectURL)
	params.Add("scope", strings.Join(oauthProvider.Scopes, " "))
	params.Add("response_type", "code")
	params.Add("state", state)

	switch provider {
	case "google":
		params.Add("access_type", "offline")
		params.Add("prompt", "consent")
	case "github":
		params.Add("allow_signup", "true")
	}

	authURL := fmt.Sprintf("%s?%s", oauthProvider.AuthURL, params.Encode())
	
	fmt.Printf("Generated OAuth URL for %s: %s\n", provider, authURL)
	fmt.Printf("All Parameters:\n")
	for key, values := range params {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
	
	return authURL, nil
}

func (s *OAuthService) HandleCallback(ctx context.Context, provider, code, state string) (*types.OAuthUserInfo, error) {
	if code == "" {
		return nil, fmt.Errorf("authorization code is required")
	}
	if state == "" {
		return nil, fmt.Errorf("state parameter is required")
	}

	if oauthStore, ok := s.store.(interfaces.OAuthStore); ok {
		storedState, err := oauthStore.GetOAuthState(ctx, state)
		if err != nil {
			return nil, fmt.Errorf("invalid or expired state parameter")
		}

		if storedState.Provider != provider {
			return nil, fmt.Errorf("state provider mismatch")
		}

		_ = oauthStore.DeleteOAuthState(ctx, state)
	}

	oauthProvider, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	accessToken, err := s.exchangeCodeForToken(ctx, oauthProvider, code, state)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	userInfo, err := s.getUserInfo(oauthProvider, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return userInfo, nil
}

func (s *OAuthService) LinkAccount(ctx context.Context, userID, provider string, userInfo *types.OAuthUserInfo) error {
	if oauthStore, ok := s.store.(interfaces.OAuthStore); ok {
		account := &types.Account{
			UserID:            userID,
			Type:              "oauth",
			Provider:          provider,
			ProviderAccountID: userInfo.ID,
		}

		return oauthStore.CreateAccount(ctx, account)
	}
	return fmt.Errorf("store does not support OAuth operations")
}

func (s *OAuthService) UnlinkAccount(ctx context.Context, userID, provider string) error {
	if oauthStore, ok := s.store.(interfaces.OAuthStore); ok {
		return oauthStore.DeleteAccount(ctx, userID, provider)
	}
	return fmt.Errorf("store does not support OAuth operations")
}

func (s *OAuthService) GetLinkedAccounts(ctx context.Context, userID string) ([]*types.Account, error) {
	if oauthStore, ok := s.store.(interfaces.OAuthStore); ok {
		return oauthStore.GetAccountsByUserID(ctx, userID)
	}
	return nil, fmt.Errorf("store does not support OAuth operations")
}

func (s *OAuthService) generateState() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *OAuthService) exchangeCodeForToken(ctx context.Context, provider *types.OAuthProvider, code, state string) (string, error) {
	var redirectURI string
	if oauthStore, ok := s.store.(interfaces.OAuthStore); ok {
		storedState, err := oauthStore.GetOAuthState(ctx, state)
		if err == nil {
			redirectURI = storedState.RedirectURL
		}
	}
	
	if redirectURI == "" {
		redirectURI = provider.RedirectURI
	}

	data := url.Values{}
	data.Set("client_id", provider.ClientID)
	data.Set("client_secret", provider.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", provider.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := json.Marshal(resp.Body)
		return "", fmt.Errorf("token exchange failed with status: %d, body: %s", resp.StatusCode, body)
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Error        string `json:"error"`
		ErrorDesc    string `json:"error_description"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	if tokenResponse.Error != "" {
		return "", fmt.Errorf("OAuth error: %s - %s", tokenResponse.Error, tokenResponse.ErrorDesc)
	}

	return tokenResponse.AccessToken, nil
}

func (s *OAuthService) getUserInfo(provider *types.OAuthProvider, accessToken string) (*types.OAuthUserInfo, error) {
	req, err := http.NewRequest("GET", provider.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info with status: %d", resp.StatusCode)
	}

	var userInfo types.OAuthUserInfo
	switch provider.Name {
	case "google":
		var googleUser struct {
			ID            string `json:"sub"`
			Email         string `json:"email"`
			Name          string `json:"name"`
			Picture       string `json:"picture"`
			EmailVerified bool   `json:"email_verified"`
		}
		err = json.NewDecoder(resp.Body).Decode(&googleUser)
		if err != nil {
			return nil, err
		}
		userInfo = types.OAuthUserInfo{
			ID:            googleUser.ID,
			Email:         googleUser.Email,
			Name:          googleUser.Name,
			AvatarURL:     googleUser.Picture,
			EmailVerified: googleUser.EmailVerified,
		}

	case "github":
		var githubUser struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
		}
		err = json.NewDecoder(resp.Body).Decode(&githubUser)
		if err != nil {
			return nil, err
		}
		userInfo = types.OAuthUserInfo{
			ID:            fmt.Sprintf("%d", githubUser.ID),
			Email:         githubUser.Email,
			Name:          githubUser.Name,
			Username:      githubUser.Login,
			AvatarURL:     githubUser.AvatarURL,
			EmailVerified: true,
		}

	case "facebook":
		var facebookUser struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Email   string `json:"email"`
			Picture struct {
				Data struct {
					URL string `json:"url"`
				} `json:"data"`
			} `json:"picture"`
		}
		err = json.NewDecoder(resp.Body).Decode(&facebookUser)
		if err != nil {
			return nil, err
		}
		userInfo = types.OAuthUserInfo{
			ID:            facebookUser.ID,
			Email:         facebookUser.Email,
			Name:          facebookUser.Name,
			AvatarURL:     facebookUser.Picture.Data.URL,
			EmailVerified: true,
		}

	default:
		return nil, fmt.Errorf("unsupported provider for user info parsing: %s", provider.Name)
	}

	return &userInfo, nil
}