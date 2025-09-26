package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/interfaces"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/middleware"
	"github.com/tdmdh/fit-up-server/services/auth-service/internal/service"
)

type AuthHandler struct {
	store        interfaces.UserStore
	verify       interfaces.VerificationStore
	authService  interfaces.AuthService
	auditLogger  *service.AuditLogger
	oauthService interfaces.OAuthService
}

func NewAuthHandler(store interfaces.UserStore, verify interfaces.VerificationStore, authService interfaces.AuthService, auditLogger *service.AuditLogger, oauthService interfaces.OAuthService) *AuthHandler {
	return &AuthHandler{
		store:        store,
		verify:       verify,
		authService:  authService,
		auditLogger:  auditLogger,
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
	router.Get("/verify", h.handleVerify)
	router.With(middleware.PasswordResetRateLimit()).Post("/forgot-password", h.handleForgotPassword)
	router.With(middleware.PasswordResetRateLimit()).Post("/reset-password", h.handleResetPassword)
	router.With(middleware.EmailVerificationRateLimit()).Post("/resend-verification", h.handleResendVerification)
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
