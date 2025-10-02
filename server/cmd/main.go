package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tdmdh/fit-up-server/internal/auth/handlers"
	authMiddleware "github.com/tdmdh/fit-up-server/internal/auth/middleware"
	authRepo "github.com/tdmdh/fit-up-server/internal/auth/repository"
	authService "github.com/tdmdh/fit-up-server/internal/auth/services"
	schemaRepo "github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/shared/config"
	"github.com/tdmdh/fit-up-server/shared/database"
)

func main() {
	log.Println("🚀 Starting Fit-Up API Server...")
	cfg := config.LoadConfig()

	if cfg.DatabaseURL == "" {
		log.Fatal("❌ DATABASE_URL environment variable is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("❌ JWT_SECRET environment variable is required")
	}

	ctx := context.Background()
	log.Println("📦 Connecting to database...")
	db, err := database.ConnectDB(ctx, cfg.DatabaseURL, cfg.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	log.Println("🔐 Initializing authentication module...")
	userStore := authRepo.NewStore(db)
	authSvc := authService.NewAuthService(userStore)
	oauthService := authService.NewOAuthService(userStore, &cfg)
	authHandler := handlers.NewAuthHandler(userStore, authSvc, oauthService)

	log.Println("💪 Initializing workout/fitness module...")
	schemaStore := schemaRepo.NewStore(db)
	_ = schemaStore

	r := chi.NewRouter()

	r.Use(authMiddleware.CORS())
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"fit-up-api","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			authHandler.RegisterRoutes(r)
		})

		// Workout/Schema routes can be added here as handlers are implemented
		// Example:
		// r.Route("/workouts", func(r chi.Router) {
		// 	r.Use(authMiddleware.JWTAuthMiddleware(userStore))
		// 	// Add workout handlers here
		// })
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("================================================================================")
		log.Printf("✅ Fit-Up API Server is running")
		log.Printf("📍 Address: http://localhost%s", addr)
		log.Printf("📍 Health: http://localhost%s/health", addr)
		log.Printf("📍 Auth API: http://localhost%s/api/v1/auth/*", addr)
		log.Println("================================================================================")
		log.Println("Press Ctrl+C to stop the server")
		log.Println()

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println()
	log.Println("🛑 Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}
