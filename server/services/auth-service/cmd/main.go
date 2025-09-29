package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	"github.com/joho/godotenv"
// 	"github.com/tdmdh/fit-up-server/services/auth-service/internal/handlers"
// 	authMiddleware "github.com/tdmdh/fit-up-server/services/auth-service/internal/middleware"
// 	"github.com/tdmdh/fit-up-server/services/auth-service/internal/repository/user"
// 	"github.com/tdmdh/fit-up-server/services/auth-service/internal/service"
// 	"github.com/tdmdh/fit-up-server/shared/config"
// 	"github.com/tdmdh/fit-up-server/shared/database/postgres"
// )

// func main() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Printf("Warning: Error loading .env file: %v", err)
// 	}

// 	cfg := config.LoadConfig()

// 	if cfg.DatabaseURL == "" {
// 		log.Fatal("DATABASE_URL environment variable is required")
// 	}
// 	if cfg.JWTSecret == "" {
// 		log.Fatal("JWT_SECRET environment variable is required")
// 	}

// 	log.Printf("Starting Auth Service on port %s", cfg.Port)
// 	log.Printf("Connecting to database...")

// 	ctx := context.Background()
// 	db, err := postgres.ConnectDB(ctx, cfg.DatabaseURL)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}
// 	defer db.Close()

// 	log.Println("Successfully connected to database!")

// 	userStore := user.NewStore(db.Pool)

// 	authService := service.NewAuthService(userStore)
// 	oauthService := service.NewOAuthService(userStore, &cfg)

// 	authHandler := handlers.NewAuthHandler(userStore, authService, oauthService)

// 	router := chi.NewRouter()
// 	router.Use(authMiddleware.CORS())
// 	router.Use(middleware.Logger)
// 	router.Use(middleware.Recoverer)

// 	// Health check endpoint
// 	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Auth Service is healthy"))
// 	})

// 	// Register auth routes directly (no /auth prefix since API Gateway adds it)
// 	authHandler.RegisterRoutes(router)

// 	addr := fmt.Sprintf(":%s", cfg.Port)
// 	log.Printf("Auth Service starting on %s...", addr)

// 	if err := http.ListenAndServe(addr, router); err != nil {
// 		log.Fatalf("Server failed to start: %v", err)
// 	}
// }
