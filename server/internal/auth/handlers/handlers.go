package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/auth/middleware"
	"github.com/tdmdh/fit-up-server/internal/auth/repository"
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
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Auth route: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	})

	router.With(middleware.LoginRateLimit()).Post("/login", h.handleLogin)
	router.With(middleware.RegisterRateLimit()).Post("/register", h.handleRegister)
	router.Route("/oauth", func(r chi.Router) {
		r.Post("/mobile/{provider}/callback", h.handleOAuthMobileCallback)
		r.Post("/{provider}", h.handleOAuthAuthorize)
		r.Get("/callback/{provider}", h.handleOAuthCallback)
	})
	router.With(middleware.PasswordResetRateLimit()).Post("/forgot-password", h.handleForgotPassword)
	router.With(middleware.PasswordResetRateLimit()).Post("/reset-password", h.handleResetPassword)
	router.Get("/{username}", h.handleGetUser)

	router.Get("/validate-token", h.handleValidateToken)
	router.With(middleware.TokenRefreshRateLimit()).Post("/refresh-token", h.handleRefreshToken)
	router.Post("/logout", h.handleLogout)
	router.Post("/verify-email", h.handleVerifyEmail)
	router.With(middleware.EmailVerificationRateLimit()).Post("/verify-email/resend", h.handleResendVerificationEmail)

	router.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(h.store))
		r.Post("/change-password", h.handleChangePassword)
		r.Put("/update-role", h.handleUpdateRole)
		r.Put("/profile", h.handleUpdateProfile)
		r.Get("/stats", h.handleGetUserStats)
		r.Get("/today-workout", h.handleGetTodayWorkout)
		r.Post("/workout-complete", h.handleWorkoutCompletion)
		r.Get("/activity-feed", h.handleGetActivityFeed)
		r.Get("/workout-history", h.handleGetWorkoutHistory)
		r.Get("/exercise-progress", h.handleGetExerciseProgress)
		r.Get("/achievements", h.handleGetAchievements)
		r.Get("/achievement-stats", h.handleGetAchievementStats)
		r.Post("/check-achievements", h.handleCheckAchievements)

		// Workout Templates
		r.Get("/templates", h.handleGetUserTemplates)
		r.Get("/templates/public", h.handleGetPublicTemplates)
		r.Get("/templates/{templateId}", h.handleGetTemplateByID)
		r.Post("/templates", h.handleCreateTemplate)
		r.Put("/templates/{templateId}", h.handleUpdateTemplate)
		r.Delete("/templates/{templateId}", h.handleDeleteTemplate)

		r.Post("/link/{provider}", h.handleLinkAccount)
		r.Delete("/unlink/{provider}", h.handleUnlinkAccount)
		r.Get("/linked-accounts", h.handleGetLinkedAccounts)
	})
}
