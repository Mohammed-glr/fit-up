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
	foodTrackerHandlers "github.com/tdmdh/fit-up-server/internal/food-tracker/handlers"
	foodTrackerRepo "github.com/tdmdh/fit-up-server/internal/food-tracker/repository"
	foodTrackerService "github.com/tdmdh/fit-up-server/internal/food-tracker/services"
	messageHandlers "github.com/tdmdh/fit-up-server/internal/message/handlers"
	"github.com/tdmdh/fit-up-server/internal/message/pool"
	messageRepo "github.com/tdmdh/fit-up-server/internal/message/repository"
	messageService "github.com/tdmdh/fit-up-server/internal/message/services"
	schemaHandlers "github.com/tdmdh/fit-up-server/internal/schema/handlers"
	schemaRepo "github.com/tdmdh/fit-up-server/internal/schema/repository"
	schemaService "github.com/tdmdh/fit-up-server/internal/schema/services"
	"github.com/tdmdh/fit-up-server/shared/config"
	"github.com/tdmdh/fit-up-server/shared/database"
	sharedMiddleware "github.com/tdmdh/fit-up-server/shared/middleware"
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

	// Initialize auth middleware after stores are ready
	authMW := sharedMiddleware.NewAuthMiddleware(schemaStore, userStore)

	exerciseService := schemaService.NewExerciseService(schemaStore)
	workoutService := schemaService.NewWorkoutService(schemaStore)
	planGenerationService := schemaService.NewPlanGenerationService(schemaStore)
	coachService := schemaService.NewCoachService(schemaStore)
	invitationService := schemaService.NewInvitationService(schemaStore.CoachInvitations())

	schemaRoutes := schemaHandlers.NewSchemaRoutes(
		schemaStore,
		userStore,
		exerciseService,
		workoutService,
		planGenerationService,
		coachService,
		invitationService,
	)

	log.Println("💬 Initializing message service with WebSocket support...")
	messageStore := messageRepo.NewStore(db)

	hub := pool.NewHub()

	hubCtx, hubCancel := context.WithCancel(ctx)
	defer hubCancel()
	go hub.Run(hubCtx)

	msgService := messageService.NewMessagesService(messageStore)

	realtimeService := messageService.NewRealtimeService(
		hub,
		msgService.Messages(),
		msgService.Conversations(),
		msgService.ReadStatus(),
	)

	msgService.SetRealtimeService(realtimeService)

	msgAuthMiddleware := sharedMiddleware.NewAuthMiddleware(schemaStore, userStore)

	messageHandler := messageHandlers.NewMessageHandler(msgService, msgAuthMiddleware)
	conversationHandler := messageHandlers.NewConversationHandler(msgService, msgAuthMiddleware)
	wsHandler := messageHandlers.NewWebSocketHandler(realtimeService, msgAuthMiddleware)

	log.Println("🍽️  Initializing food tracker service...")
	foodTrackerStore := foodTrackerRepo.NewStore(db)

	ingredientDB := foodTrackerService.NewSimpleIngredientDB()

	foodTrackerSvc := foodTrackerService.NewService(foodTrackerStore, ingredientDB)

	foodTrackerHandler := foodTrackerHandlers.NewFoodTrackerHandler(foodTrackerSvc, schemaStore, userStore)

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

		// Register template routes
		authHandler.RegisterTemplateRoutes(r, authMW)

		schemaRoutes.RegisterRoutes(r)

		messageHandlers.SetupMessageRoutes(r, messageHandler, conversationHandler, msgAuthMiddleware)

		foodTrackerHandler.RegisterRoutes(r)
	})

	messageHandlers.SetupWebSocketRoutes(r, wsHandler)

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
		log.Printf("📍 Exercises: http://localhost%s/api/v1/exercises/*", addr)
		log.Printf("📍 Workouts: http://localhost%s/api/v1/workouts/*", addr)
		log.Printf("📍 Sessions: http://localhost%s/api/v1/workout-sessions/*", addr)
		log.Printf("📍 Fitness: http://localhost%s/api/v1/fitness-profile/*", addr)
		log.Printf("📍 Plans: http://localhost%s/api/v1/plans/*", addr)
		log.Printf("📍 Coach: http://localhost%s/api/v1/coach/*", addr)
		log.Printf("📍 Templates: http://localhost%s/api/v1/templates/*", addr)
		log.Printf("📍 Messages: http://localhost%s/api/v1/messages/*", addr)
		log.Printf("📍 Conversations: http://localhost%s/api/v1/conversations/*", addr)
		log.Printf("📍 Food Tracker: http://localhost%s/api/v1/food-tracker/*", addr)
		log.Printf("📍 WebSocket: ws://localhost%s/ws", addr)
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

	log.Println("🔌 Stopping WebSocket hub...")
	hub.Stop()
	hubCancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}
