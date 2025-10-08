package handlers

import (

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/auth/repository"
	"github.com/tdmdh/fit-up-server/internal/auth/middleware"
)

type AuthHandler struct {
	store        repository.UserStore
	authService  repository.AuthService
	oauthService repository.OAuthService
}

func NewAuthHandler(store repository.UserStore, authService repository.AuthService, oauthService repository.OAuthService) *AuthHandler {
	return &AuthHandler{
		store:        store,
		authService:  authService,
		oauthService: oauthService,
	}
}

func (h *AuthHandler) RegisterRoutes(router chi.Router) {
	router.With(middleware.LoginRateLimit()).Post("/login", h.handleLogin)
	router.With(middleware.RegisterRateLimit()).Post("/register", h.handleRegister)
	router.Route("/oauth", func(r chi.Router) {
		r.Post("/{provider}", h.handleOAuthAuthorize)
		r.Get("/callback/{provider}", h.handleOAuthCallback)
	})
	router.With(middleware.PasswordResetRateLimit()).Post("/forgot-password", h.handleForgotPassword)
	router.With(middleware.PasswordResetRateLimit()).Post("/reset-password", h.handleResetPassword)
	router.Get("/{username}", h.handleGetUser)

	router.Post("/validate-token", h.handleValidateToken)
	router.With(middleware.TokenRefreshRateLimit()).Post("/refresh-token", h.handleRefreshToken)
	router.Post("/logout", h.handleLogout)

	router.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(h.store))
		r.Post("/change-password", h.handleChangePassword)

		r.Post("/link/{provider}", h.handleLinkAccount)
		r.Delete("/unlink/{provider}", h.handleUnlinkAccount)
		r.Get("/linked-accounts", h.handleGetLinkedAccounts)
	})
}
