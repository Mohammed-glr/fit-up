package handlers

import (

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/interfaces"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/middleware"
)

type AuthHandler struct {
	store        interfaces.UserStore
	authService  interfaces.AuthService
	oauthService interfaces.OAuthService
}

func NewAuthHandler(store interfaces.UserStore, authService interfaces.AuthService, oauthService interfaces.OAuthService) *AuthHandler {
	return &AuthHandler{
		store:        store,
		authService:  authService,
		oauthService: oauthService,
	}
}

func (h *AuthHandler) RegisterRoutes(router chi.Router) {
	// Public routes (no authentication required) with rate limiting
	router.With(middleware.LoginRateLimit()).Post("/login", h.handleLogin)
	router.With(middleware.RegisterRateLimit()).Post("/register", h.handleRegister)
	//oauth2 routes
	router.Route("/oauth", func(r chi.Router) {
		r.Post("/{provider}", h.handleOAuthAuthorize)
		r.Get("/callback/{provider}", h.handleOAuthCallback)
	})
	router.With(middleware.PasswordResetRateLimit()).Post("/forgot-password", h.handleForgotPassword)
	router.With(middleware.PasswordResetRateLimit()).Post("/reset-password", h.handleResetPassword)
	router.Get("/{username}", h.handleGetUser)

	// Token management routes
	router.Post("/validate-token", h.handleValidateToken)
	router.With(middleware.TokenRefreshRateLimit()).Post("/refresh-token", h.handleRefreshToken)
	router.Post("/logout", h.handleLogout)

	// Protected routes (require authentication)
	router.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(h.store))
		r.Post("/change-password", h.handleChangePassword)

		// OAuth account management
		r.Post("/link/{provider}", h.handleLinkAccount)
		r.Delete("/unlink/{provider}", h.handleUnlinkAccount)
		r.Get("/linked-accounts", h.handleGetLinkedAccounts)
	})
}
